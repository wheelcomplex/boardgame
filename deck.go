package boardgame

import (
	"github.com/jkomoros/boardgame/errors"
	"strconv"
)

//TODO: consider making Deck be an interface again (in some cases it
//might be nice to be able to cast the Deck directly to its underlying type to
//minimize later casts)

//A Deck represents an immutable collection of a certain type of components.
//Every component lives in one deck. 1 or more Stacks index into every Deck,
//and cover every item in the deck, with no items in more than one deck. The
//zero-value of Deck is useful. The Deck will not return items until it has
//been added to a ComponentChest, which helps enforce that Decks' values never
//change. Create a new Deck with NewDeck()
type Deck struct {
	chest *ComponentChest
	//Name is only set when it's added to the component chest.
	name string
	//Components should only ever be added at initalization time. After
	//initalization, Components should be read-only.
	components            []*Component
	shadowValues          Reader
	vendedShadowComponent *Component
	//TODO: protect shadowComponents cache with mutex to make threadsafe.
}

const genericComponentSentinel = -2

func NewDeck() *Deck {
	return &Deck{}
}

//NewSizedStack returns a new default (growable Stack) with the given size
//based on this deck. You normally do this in *Constructor delegate methods,
//if you aren't using the auto-inflating struct tags to configure your stacks.
//The returned stack will allow up to maxSize items to be inserted. If you
//don't want to set a maxSize on the stack (you often don't) pass 0 for
//maxSize to allow it to grow without limit.
func (d *Deck) NewStack(maxSize int) MutableStack {
	return newGrowableStack(d, maxSize)
}

//NewSizedStack returns a new SizedStack (a stack whose FixedSize() will
//return true). Refer to the Stack interface documentation for more about the
//difference.
func (d *Deck) NewSizedStack(size int) MutableStack {
	return newSizedStack(d, size)
}

//AddComponent adds a new component with the given values to the next spot in
//the deck. If the deck has already been added to a componentchest, this will
//do nothing.
func (d *Deck) AddComponent(v Reader) {
	if d.chest != nil {
		return
	}

	c := &Component{
		Deck:      d,
		DeckIndex: len(d.components),
		Values:    v,
	}

	d.components = append(d.components, c)
}

//AddComponentMulti is like AddComponent, but creates multiple versions of the
//same component. The exact same ComponentValues will be re-used, which is
//reasonable becasue components are read-only anyway.
func (d *Deck) AddComponentMulti(v Reader, count int) {
	for i := 0; i < count; i++ {
		d.AddComponent(v)
	}
}

//Components returns a list of Components in order in this deck, but only if
//this Deck has already been added to its ComponentChest.
func (d *Deck) Components() []*Component {
	if d.chest == nil {
		return nil
	}
	return d.components
}

//Chest points back to the chest we're part of.
func (d *Deck) Chest() *ComponentChest {
	return d.chest
}

func (d *Deck) Name() string {
	return d.name
}

//ComponentAt returns the component at a given index. It handles empty indexes
//and shadow indexes correctly.
func (d *Deck) ComponentAt(index int) *Component {
	if d.chest == nil {
		return nil
	}
	if index >= len(d.components) {
		return nil
	}
	if index >= 0 {
		return d.components[index]
	}

	if index == emptyIndexSentinel {
		return nil
	}

	return d.GenericComponent()

}

//SetShadowValues sets the SubState to return for every shadow
//component that is returned. May only be set before added to a chest. Should
//generally be the same shape of componentValues as used for other components
//in the deck.
func (d *Deck) SetShadowValues(v Reader) {
	if d.chest != nil {
		return
	}
	d.shadowValues = v
}

//GenericComponent returns the component that is considereed fully generic for
//this deck. This is the component that every component will be if a Stack is
//sanitized with PolicyLen, for example. If you want to figure out if a Stack
//was sanitized according to that policy, you can compare the component to
//this.
func (d *Deck) GenericComponent() *Component {

	if d.vendedShadowComponent == nil {
		shadow := &Component{
			Deck:      d,
			DeckIndex: genericComponentSentinel,
			Values:    d.shadowValues,
		}

		d.vendedShadowComponent = shadow
	}

	return d.vendedShadowComponent
}

var illegalComponentValuesProps = map[PropertyType]bool{
	TypeStack: true,
	TypeTimer: true,
}

//finish is called when the deck is added to a component chest. It signifies that no more items may be added.
func (d *Deck) finish(chest *ComponentChest, name string) error {

	for i, c := range d.components {
		if c.Values == nil {
			continue
		}
		validator, err := newReaderValidator(c.Values.Reader(), c.Values, illegalComponentValuesProps, chest, false)
		if err != nil {
			return errors.New("Component " + strconv.Itoa(i) + "failed to validate: " + err.Error())
		}
		if err := validator.Valid(c.Values.Reader()); err != nil {
			return errors.New("Component " + strconv.Itoa(i) + " failed to validate: " + err.Error())
		}
	}

	d.chest = chest
	//If a deck has a name, it cannot receive any more items.
	d.name = name
	return nil
}
