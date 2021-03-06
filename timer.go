package boardgame

import (
	"container/heap"
	"time"
)

//A Timer is a type of property that can be used in states that represents a
//countdown. The MutableTimer interface includes mutator methods as well. See
//the package documentation for more on timers.
type Timer interface {
	Active() bool
	TimeLeft() time.Duration
	copy() Timer
	id() int
	state() *state
	setState(*state)
}

//MutableTimer is a Timer that also includes mutator methods.
type MutableTimer interface {
	Timer
	Start(time.Duration, Move)
	Cancel() bool
	mutableCopy() MutableTimer
}

type timer struct {
	//Id will be an opaque identifier that is used to keep track of the
	//corresponding underlying Timer object in the game engine. It is not
	//meaningful to inspect yourself and should not be modified.
	Id       int
	statePtr *state
}

func NewTimer() MutableTimer {
	return &timer{}
}

func (t *timer) id() int {
	return t.Id
}

func (t *timer) state() *state {
	return t.statePtr
}

func (t *timer) setState(state *state) {
	t.statePtr = state
}

func (t *timer) copy() Timer {
	return t.mutableCopy()
}

func (t *timer) mutableCopy() MutableTimer {
	var result timer
	result = *t
	return &result
}

func (t *timer) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{
		"Id": t.Id,
		//TimeLeft is only ever for the client (it's not read back in when
		//deserialized), so put it in the more traditional milliseconds units,
		//not nanoseconds.
		"TimeLeft": t.TimeLeft() / time.Millisecond,
	}

	return DefaultMarshalJSON(obj)
}

//Active returns true if the timer is active and counting down.
func (t *timer) Active() bool {
	return t.statePtr.game.manager.timers.TimerActive(t.Id)
}

//TimeLeft returns the number of nanoseconds left until this timer fires.
func (t *timer) TimeLeft() time.Duration {
	return t.statePtr.game.manager.timers.GetTimerRemaining(t.Id)
}

//Start starts the timer. After duration has passed, the Move will be proposed
//via proposeMove. If the timer is already active, it will be canceled before
//the new timer is configured.
func (t *timer) Start(duration time.Duration, move Move) {

	if t.Active() {
		t.Cancel()
	}

	game := t.statePtr.game
	manager := game.manager

	t.Id = manager.timers.PrepareTimer(duration, game, move)

	t.statePtr.timersToStart = append(t.statePtr.timersToStart, t.Id)
}

//Cancel cancels an active timer. If the timer is not active, it has no
//effect. Returns true if the timer was active and canceled, false if the
//timer was not active.
func (t *timer) Cancel() bool {

	wasActive := t.Active()

	manager := t.statePtr.game.manager

	manager.timers.CancelTimer(t.Id)

	//Technically there's a case where Start() was called, but the state was
	//never fully committed. However, StartTimer() on a canceled timer is a
	//no-op so it's fine.

	t.Id = 0

	return wasActive
}

type timerRecord struct {
	id    int
	index int
	//When the timer should fire. Set after the timer is fully Started().
	//Before that it is impossibly far in the future.
	fireTime time.Time
	//The duration we were configured with. Only set when the timer is
	//Prepared() but not yet Started().
	duration time.Duration
	game     *Game
	move     Move
}

func (t *timerRecord) TimeRemaining() time.Duration {

	//Before a timer is Started(), just say its duration as the time
	//remaining.
	if t.duration > 0 {
		return t.duration
	}

	duration := t.fireTime.Sub(time.Now())

	if duration < 0 {
		duration = 0
	}

	return duration
}

type timerQueue []*timerRecord

type timerManager struct {
	nextId      int
	records     timerQueue
	recordsById map[int]*timerRecord
	manager     *GameManager
}

func newTimerManager(gameManager *GameManager) *timerManager {
	return &timerManager{
		//the default id in TimerProps is 0, so we should start beyond that.
		nextId:      1,
		records:     make(timerQueue, 0),
		recordsById: make(map[int]*timerRecord),
		manager:     gameManager,
	}
}

//PrepareTimer creates a timer entry and gets it ready and an Id allocated.
//However, the timer doesn't actually start counting down until
//manager.StartTimer(id) is called.
func (t *timerManager) PrepareTimer(duration time.Duration, game *Game, move Move) int {
	record := &timerRecord{
		id:       t.nextId,
		index:    -1,
		duration: duration,
		//fireTime will be set when StartTimer is called. For now, set it to
		//something impossibly far in the future.
		fireTime: time.Now().Add(time.Hour * 100000),
		game:     game,
		move:     move,
	}
	t.nextId++

	t.recordsById[record.id] = record

	heap.Push(&t.records, record)

	return record.id
}

//StartTimer actually triggers a timer that was previously PrepareTimer'd to
//start counting down.
func (t *timerManager) StartTimer(id int) {

	if t.TimerActive(id) {
		return
	}

	record := t.recordsById[id]

	if record == nil {
		return
	}

	record.fireTime = time.Now().Add(record.duration)
	record.duration = 0

	heap.Fix(&t.records, record.index)
}

//TimerActive returns if the timer is active and counting down.
func (t *timerManager) TimerActive(id int) bool {
	record := t.recordsById[id]

	if record == nil {
		return false
	}

	if record.duration > 0 {
		return false
	}

	return true
}

func (t *timerManager) GetTimerRemaining(id int) time.Duration {
	record := t.recordsById[id]

	if record == nil {
		return 0
	}

	return record.TimeRemaining()
}

func (t *timerManager) CancelTimer(id int) {
	record := t.recordsById[id]

	if record == nil {
		return
	}

	heap.Remove(&t.records, record.index)

	record.index = -1

	delete(t.recordsById, record.id)

}

//Should be called regularly by the manager to tell this to check and see if
//any timers have fired, and execute them if so.
func (t *timerManager) Tick() {
	for t.nextTimerFired() {
		record := t.popNext()
		if record == nil {
			continue
		}

		if err := <-record.game.ProposeMove(record.move, AdminPlayerIndex); err != nil {
			//TODO: log the error or something
			t.manager.Logger().Info("When timer failed the move could not be made: ", err, record.move)
		}
	}
}

//Whether the next timer in the queue is already fired
func (t *timerManager) nextTimerFired() bool {
	if len(t.records) == 0 {
		return false
	}

	record := t.records[0]

	return record.TimeRemaining() <= 0
}

func (t *timerManager) popNext() *timerRecord {
	if !t.nextTimerFired() {
		return nil
	}

	x := heap.Pop(&t.records)

	record := x.(*timerRecord)

	delete(t.recordsById, record.id)

	return record
}

func (t timerQueue) Len() int {
	return len(t)
}

func (t timerQueue) Less(i, j int) bool {
	return t[i].fireTime.Sub(t[j].fireTime) < 0
}

func (t timerQueue) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
	t[i].index = i
	t[j].index = j
}

//DO NOT USE THIS DIRECTLY. Use heap.Push(t)
func (t *timerQueue) Push(x interface{}) {
	n := len(*t)
	item := x.(*timerRecord)
	item.index = n
	*t = append(*t, item)
}

//DO NOT USE THIS DIRECTLY. Use heap.Pop()
func (t *timerQueue) Pop() interface{} {
	old := *t
	n := len(old)
	item := old[n-1]
	item.index = -1
	*t = old[0 : n-1]
	return item
}
