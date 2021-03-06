package boardgame

import (
	"encoding/json"
	"github.com/workfit/tester/assert"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestPlayerIndexNextPrevious(t *testing.T) {

	game := testGame(t)

	game.SetUp(3, nil, nil)

	state := game.CurrentState()

	tests := []struct {
		p            PlayerIndex
		expectedNext PlayerIndex
		expectedPrev PlayerIndex
	}{
		{
			0,
			1,
			2,
		},
		{
			2,
			0,
			1,
		},
		{
			AdminPlayerIndex,
			AdminPlayerIndex,
			AdminPlayerIndex,
		},
		{
			ObserverPlayerIndex,
			ObserverPlayerIndex,
			ObserverPlayerIndex,
		},
	}

	for i, test := range tests {
		result := test.p.Next(state)
		assert.For(t, "next", i).ThatActual(result).Equals(test.expectedNext)

		result = test.p.Previous(state)

		assert.For(t, "prev", i).ThatActual(result).Equals(test.expectedPrev)
	}
}

func TestPlayerIndexValid(t *testing.T) {

	gameTwoPlayers := testGame(t)

	gameTwoPlayers.SetUp(2, nil, nil)

	stateTwoPlayers := gameTwoPlayers.CurrentState()

	tests := []struct {
		p        PlayerIndex
		state    State
		expected bool
	}{
		{
			0,
			stateTwoPlayers,
			true,
		},
		{
			ObserverPlayerIndex,
			stateTwoPlayers,
			true,
		},
		{
			AdminPlayerIndex,
			stateTwoPlayers,
			true,
		},
		{
			AdminPlayerIndex - 1,
			stateTwoPlayers,
			false,
		},
		{
			3,
			stateTwoPlayers,
			false,
		},
		{
			2,
			stateTwoPlayers,
			false,
		},
	}

	for i, test := range tests {
		result := test.p.Valid(test.state)
		assert.For(t, "valid", i).ThatActual(result).Equals(test.expected)
	}
}

func TestPlayerIndexEquivalent(t *testing.T) {

	equivalentTests := []struct {
		p        PlayerIndex
		other    PlayerIndex
		expected bool
	}{
		{
			0,
			0,
			true,
		},
		{
			0,
			1,
			false,
		},
		{
			AdminPlayerIndex,
			0,
			true,
		},
		{
			AdminPlayerIndex,
			ObserverPlayerIndex,
			false,
		},
		{
			ObserverPlayerIndex,
			1,
			false,
		},
		{
			0,
			AdminPlayerIndex,
			true,
		},
		{
			AdminPlayerIndex,
			AdminPlayerIndex,
			true,
		},
		{
			ObserverPlayerIndex,
			ObserverPlayerIndex,
			false,
		},
	}

	for i, test := range equivalentTests {
		result := test.p.Equivalent(test.other)

		assert.For(t, "equivalent", i).ThatActual(result).Equals(test.expected)
	}
}

func TestSecretMoveCount(t *testing.T) {

	game := testGame(t)

	makeTestGameIdsStable(game)

	game.SetUp(0, nil, nil)

	currentState := game.CurrentState()

	assert.For(t).ThatActual(currentState.Version()).Equals(game.Version())

	gameState, _ := concreteStates(currentState)

	s := currentState.(*state)

	for i, c := range gameState.DrawDeck.Components() {
		assert.For(t, i).ThatActual(c.secretMoveCount(s)).Equals(0)
	}

	idBefore := gameState.DrawDeck.ComponentAt(0).Id(s)

	gameState.DrawDeck.ComponentAt(0).movedSecretly(s)

	assert.For(t).ThatActual(gameState.DrawDeck.ComponentAt(0).Id(s)).DoesNotEqual(idBefore)

	for i, c := range gameState.DrawDeck.Components() {
		if i == 0 {
			assert.For(t, i).ThatActual(c.secretMoveCount(s)).Equals(1)
		} else {
			assert.For(t, i).ThatActual(c.secretMoveCount(s)).Equals(0)
		}
	}

	//We're going to do a faked save to verify that these things round trip
	game.version++

	blob, err := json.MarshalIndent(s, "", "\t")

	assert.For(t).ThatActual(err).IsNil()

	game.manager.Storage().SaveGameAndCurrentState(game.StorageRecord(), blob, nil)

	//Read back in the game and verify that the secretMoveCount round-tripped.

	refriedGame := game.manager.Game(game.Id())

	refriedS := refriedGame.CurrentState().(*state)

	refriedGameState, _ := concreteStates(refriedGame.CurrentState())

	for i, c := range refriedGameState.DrawDeck.Components() {
		if i == 0 {
			assert.For(t, i).ThatActual(c.secretMoveCount(refriedS)).Equals(1)
		} else {
			assert.For(t, i).ThatActual(c.secretMoveCount(refriedS)).Equals(0)
		}
	}

}

func testSubStatesHaveStateSet(t *testing.T, state *state) {
	assert.For(t).ThatActual(state.GameState().(*testGameState).state).Equals(state)

	for i := 0; i < len(state.playerStates); i++ {
		assert.For(t, i).ThatActual(state.PlayerStates()[i].(*testPlayerState).state).Equals(state)
	}

	for _, dynamicComponents := range state.DynamicComponentValues() {
		for i, component := range dynamicComponents {
			assert.For(t, i).ThatActual(component.(*testingComponentDynamic).state).Equals(state)
		}
	}
}

func TestState(t *testing.T) {

	game := testGame(t)

	makeTestGameIdsStable(game)

	game.SetUp(0, nil, nil)

	assert.For(t).ThatActual(game.CurrentState().Version()).Equals(game.Version())

	theState := game.CurrentState().(*state)

	testSubStatesHaveStateSet(t, theState)

	testSubStatesHaveStateSet(t, theState.Copy(false).(*state))

	record, err := game.Manager().Storage().State(game.Id(), game.Version())

	if err != nil {
		t.Error("Unexpected error", err)
	}

	state, err := game.Manager().stateFromRecord(record)
	state.game = game

	if err != nil {
		t.Error("StateFromBlob err", err)
	}

	if state == nil {
		t.Error("State could not be created")
	}

	assert.For(t).ThatActual(state.Version()).Equals(game.Version())

	testSubStatesHaveStateSet(t, state)

	currentJson, _ := json.Marshal(state)
	golden := goldenJSON("basic_state.json", t)

	compareJSONObjects(currentJson, golden, "Basic state", t)

	stateCopy := state.Copy(false)

	copyJson, _ := DefaultMarshalJSON(stateCopy)

	compareJSONObjects(copyJson, currentJson, "Copy was not same", t)

	_, playerStatesCopy := concreteStates(stateCopy)

	playerStatesCopy[0].MovesLeftThisTurn = 10

	_, playerStates := concreteStates(state)

	if playerStates[0].MovesLeftThisTurn == 10 {
		t.Error("Modifying a copy change the original")
	}

	if state.Sanitized() {
		t.Error("State reported being sanitized even when it wasn't")
	}

	sanitizedStateCopy := stateCopy.Copy(true)

	if !sanitizedStateCopy.Sanitized() {
		t.Error("A copy that was told it was sanitized did not report being sanitized.")
	}

	//TODO: test that GAmeState and UserStates are also copies
}

func TestStateSerialization(t *testing.T) {

	game := testGame(t)

	game.SetUp(0, nil, nil)

	gameState, _ := concreteStates(game.CurrentState())

	if gameState.Timer.state() == nil {
		t.Error("The set up timer did no thave a stateptr")
	}

	rawMove := game.PlayerMoveByName("test")

	move := rawMove.(*testMove)

	move.AString = "bam"
	move.ScoreIncrement = 3
	move.TargetPlayerIndex = 0
	move.ABool = true

	if err := <-game.ProposeMove(move, AdminPlayerIndex); err != nil {
		t.Fatal("Couldn't make move", err)
	}

	blob, err := json.Marshal(game.CurrentState())

	if err != nil {
		t.Fatal("Couldn't serialize state:", err)
	}

	reconstitutedState, err := game.Manager().stateFromRecord(blob)

	if err != nil {
		t.Error("StateFromBlob returned unexpected err", err)
	}

	reconstitutedState.game = game

	gameState, _ = concreteStates(reconstitutedState)

	if !gameState.DrawDeck.inflated() {
		t.Error("The stack was not inflated when it came back from StateFromBlob")
	}

	if !gameState.DrawDeck.ComponentAt(0).DynamicValues(reconstitutedState).(*testingComponentDynamic).Stack.inflated() {
		t.Error("The stack on a component's dynamic value was not inflated coming back from storage.")
	}

	if gameState.Timer.state() == nil {
		t.Error("The timer did not come back inflated from storage")
	}

	//This is lame, but when you create json for a State, it touches Computed,
	//which will make it non-nil, so if you're doing direct comparison they
	//won't compare equal even though they basically are. At this point
	//CurrentState has already been touched above by creating a json blob. So
	//just touch reconstitutedState, too. ¯\_(ツ)_/¯

	_, _ = json.Marshal(reconstitutedState)

	if !reflect.DeepEqual(reconstitutedState, game.CurrentState()) {

		rStateBlob, _ := json.Marshal(reconstitutedState)
		oStateBlob, _ := json.Marshal(game.CurrentState())

		t.Error("Reconstituted state and original state were not the same. Got", string(rStateBlob), "wanted", string(oStateBlob))
	}
}

func compareJSONObjects(in []byte, golden []byte, message string, t *testing.T) {

	//recreated in server/internal/teststoragemanager

	var deserializedIn interface{}
	var deserializedGolden interface{}

	json.Unmarshal(in, &deserializedIn)
	json.Unmarshal(golden, &deserializedGolden)

	assert.For(t).ThatActual(deserializedIn).IsNotNil()

	assert.For(t).ThatActual(deserializedGolden).IsNotNil()

	assert.For(t, message).ThatActual(deserializedGolden).Equals(deserializedIn).ThenDiffOnFail()

}

func goldenJSON(fileName string, t *testing.T) []byte {
	contents, err := ioutil.ReadFile("./test/" + fileName)

	if !assert.For(t, fileName).ThatActual(err).IsNil().Passed() {
		t.FailNow()
	}

	return contents

}
