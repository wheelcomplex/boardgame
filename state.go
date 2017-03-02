package boardgame

import (
	"encoding/json"
)

//StatePayload is where the "meat" of the state goes. It is one object so that
//client games can cast it quickly to the concrete struct for their game, so
//that they can get to a type-checked world with minimal fuss inside of
//Move.Legal and move.Apply. Your underlying struct should have a Game and
//Players property, so they serialize properly to JSON. Most importantly,
//json.Marshal() should round trip through your GameDelegate.StateFromBlob()
//without modifications in order for persistence to work. Each PlayerState you
//return should be the same underlying type. This means that if different
//players have very different roles in a game, there might be many properties
//that are not in use for any given player.
type State struct {
	Game      GameState
	Players   []PlayerState
	sanitized bool
	delegate  GameDelegate
}

//Note: the MarshalJSON output of State is appropriate for sending to client
//or persisting to storage.  In the future what we marshal for storage and
//what we marshal for sending to client might be different. (e.g. all computed
//properties, which might be lazily computed server side). But we'll cross
//that bridge when we come to it.

//Copy returns a copy of the Payload. Be sure it's a deep copy that makes
//a copy of any pointer arguments. If the Copy will be used to create a
//Sanitized version of state, pass sanitized = true.
func (s *State) Copy(sanitized bool) *State {

	players := make([]PlayerState, len(s.Players))

	for i, player := range s.Players {
		players[i] = player.Copy()
	}

	result := &State{
		Game:      s.Game.Copy(),
		Players:   players,
		sanitized: sanitized,
		delegate:  s.delegate,
	}

	return result

}

//Diagram should return a basic debug rendering of state in multi-line ascii
//art. Useful for debugging. Will thunk out to Delegate.Diagram()
func (s *State) Diagram() string {
	return s.delegate.Diagram(s)
}

//Santizied will return false if this is a full-fidelity State object, or
//true if it has been sanitized, which means that some properties might be
//hidden or otherwise altered. This should return true if the object was
//created with Copy(true)
func (s *State) Sanitized() bool {
	return s.sanitized
}

//SanitizedForPlayer produces a copy state object that has been sanitized for
//the player at the given index. The state object returned will have
//Sanitized() return true. Will call GameDelegate.StateSanitizationPolicy to
//retrieve the policy in place. See the package level comment for an overview
//of how state sanitization works.
func (s *State) SanitizedForPlayer(playerIndex int) *State {

	//If the playerIndex isn't an actuall player's index, just return self.
	if playerIndex < 0 || playerIndex >= len(s.Players) {
		return s
	}

	policy := s.delegate.StateSanitizationPolicy()

	if policy == nil {
		policy = &StatePolicy{}
	}

	sanitized := s.Copy(true)

	sanitizeStateObj(sanitized.Game.Reader(), policy.Game, -1, playerIndex)

	playerStates := sanitized.Players

	for i := 0; i < len(playerStates); i++ {
		sanitizeStateObj(playerStates[i].Reader(), policy.Player, i, playerIndex)
	}

	return sanitized

}

//BaseState is the interface that all state objects--UserStates and GameStates
//--implement.
type BaseState interface {
	Reader() PropertyReadSetter
}

//PlayerState represents the state of a game associated with a specific user.
type PlayerState interface {
	//PlayerIndex encodes the index this user's state is in the containing
	//state object.
	PlayerIndex() int
	//Copy produces a copy of our current state. Be sure it's a deep copy that
	//makes a copy of any pointer arguments.
	Copy() PlayerState
	BaseState
}

//GameState represents the state of a game that is not associated with a
//particular user. For example, the draw stack of cards, who the current
//player is, and other properites.
type GameState interface {
	//Copy returns a copy of our current state. Be sure it's a deep copy that
	//makes a copy of any pointer arguments.
	Copy() GameState
	BaseState
}

//DefaultMarshalJSON is a simple wrapper around json.MarshalIndent, with the
//right defaults set. If your structs need to implement MarshaLJSON to output
//JSON, use this to encode it.
func DefaultMarshalJSON(obj interface{}) ([]byte, error) {
	return json.MarshalIndent(obj, "", "  ")
}
