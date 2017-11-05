package moves

import (
	"errors"
	"github.com/jkomoros/boardgame"
)

//RoundRobinProperties should be implemented by your GameState if you use any
//of the RoundRobin moves, including StartRoundRobin. You don't have to do
//anything we these other than store them to a property in your gameState and
//then return them via the getters.
type RoundRobinProperties interface {
	//The next player whose round robin turn it will be
	NextRoundRobinPlayer() boardgame.PlayerIndex
	//The index of the player we started the round robin on.
	RoundRobinStarterPlayer() boardgame.PlayerIndex
	//How many complete times around the round robin we've been
	RoundRobinRoundCount() int

	SetNextRoundRobinPlayer(nextPlayer boardgame.PlayerIndex)
	SetRoundRobinStarterPlayer(index boardgame.PlayerIndex)
	SetRoundRobinRoundCount(int)
}

//RoundRobinActioner should be implemented by any moves that embed a
//RoundRobin move. It's the action that will be called on the player who is
//next in the round robin.
type RoundRobinActioner interface {
	//RoundRobinAction should do the action for the round robin to the player
	//in TargetPlayerIndex.
	RoundRobinAction(state boardgame.MutableState) error
}

//StartRoundRobin is the move you should have in the progression immediately
//before a round robin starts. It sets the NextRoundRobinPlayer to the game's
//CurrentPlayer, getting ready for moves that embed RoundRobin.
type StartRoundRobin struct {
	Base
}

func (s *StartRoundRobin) ValidConfiguration(exampleState boardgame.MutableState) error {
	if _, ok := exampleState.GameState().(RoundRobinProperties); !ok {
		return errors.New("GameState does not implement RoundRobiner interface")
	}
	return nil
}

//Apply gets the game ready for a round robin by calling
//gameState.SetNextRoundRobinPlayer to CurrentPlayerIndex.
func (s *StartRoundRobin) Apply(state boardgame.MutableState) error {
	roundRobiner, ok := state.GameState().(RoundRobinProperties)

	if !ok {
		return errors.New("GameState unexpectedly did not implement RoundRobiner interface")
	}

	currentPlayer := state.Game().Manager().Delegate().CurrentPlayerIndex(state)

	roundRobiner.SetNextRoundRobinPlayer(currentPlayer)

	return nil
}

/*

RoundRobin is a type of move that goes around every player one by one and
does some action. Other moves in this package embed RoundRobin. The
embeding move should implement RoundRobinActioner.

Round Robin moves start at a given player and go around, applying the
RoundRobinAction for each player until the RoundRobinFinished() method returns
true. Various embedders of the base RoundRobin will override the default
behavior for that method.

Round Robin keeps track of various properties on the gameState by using the
RoundRobinProperties interface.

A round robin phase must be immediately preceded by StartRoundRobin, which
sets various properties the round robin needs to operate before it starts.

*/
type RoundRobin struct {
	Base
	TargetPlayerIndex boardgame.PlayerIndex
}

//AllowMultipleInProgression returns true because RoundRobins go until the end
//condition of the round robin is met.
func (r *RoundRobin) AllowMultipleInProgression() bool {
	return true
}

//DefaultsForState sets the TargetPlayerIndex to NextRoundRobinPlayer.
func (r *RoundRobin) DefaultsForState(state boardgame.State) {
	roundRobiner, ok := state.GameState().(RoundRobinProperties)

	if !ok {
		return
	}

	r.TargetPlayerIndex = roundRobiner.NextRoundRobinPlayer()
}

func (r *RoundRobin) ValidConfiguration(exampleState boardgame.MutableState) error {
	if _, ok := exampleState.GameState().(RoundRobinProperties); !ok {
		return errors.New("GameState does not implement RoundRobiner interface")
	}

	embeddingMove := r.TopLevelStruct()

	if _, ok := embeddingMove.(RoundRobinActioner); !ok {
		return errors.New("Embedding move doesn't implement RoundRobinActioner")
	}
	return nil
}

func (r *RoundRobin) Apply(state boardgame.MutableState) error {

	embeddingMove := r.TopLevelStruct()

	actioner, ok := embeddingMove.(RoundRobinActioner)

	if !ok {
		return errors.New("Embedding move doesn't implement RoundRobinActioner")
	}

	if err := actioner.RoundRobinAction(state); err != nil {
		return errors.New("RoundRobinAction returned error: " + err.Error())
	}

	roundRobiner, ok := state.GameState().(RoundRobinProperties)

	if !ok {
		return errors.New("GameState does not implement RoundRobiner interface")
	}

	roundRobiner.SetNextRoundRobinPlayer(r.TargetPlayerIndex.Next(state))

	return nil

}
