package blackjack

import (
	"errors"
	"github.com/jkomoros/boardgame"
	"github.com/jkomoros/boardgame/components/playingcards"
	"github.com/jkomoros/boardgame/moves"
)

func init() {

	//Make sure that we get compile-time errors if our player and game state
	//don't adhere to the interfaces that moves.FinishTurn expects
	moves.VerifyFinishTurnStates(&gameState{}, &playerState{})
}

func concreteStates(state boardgame.State) (*gameState, []*playerState) {
	game := state.GameState().(*gameState)

	players := make([]*playerState, len(state.PlayerStates()))

	for i, player := range state.PlayerStates() {
		players[i] = player.(*playerState)
	}

	return game, players
}

//+autoreader
type gameState struct {
	DiscardStack  *boardgame.GrowableStack `stack:"cards" sanitize:"len"`
	DrawStack     *boardgame.GrowableStack `stack:"cards" sanitize:"len"`
	UnusedCards   *boardgame.GrowableStack `stack:"cards"`
	CurrentPlayer boardgame.PlayerIndex
}

//+autoreader
type playerState struct {
	playerIndex    boardgame.PlayerIndex
	GotInitialDeal bool
	HiddenHand     *boardgame.GrowableStack `stack:"cards,1" sanitize:"len"`
	VisibleHand    *boardgame.GrowableStack `stack:"cards"`
	Busted         bool
	Stood          bool
}

func (g *gameState) SetCurrentPlayer(currentPlayer boardgame.PlayerIndex) {
	g.CurrentPlayer = currentPlayer
}

func (p *playerState) PlayerIndex() boardgame.PlayerIndex {
	return p.playerIndex
}

func (p *playerState) TurnDone(state boardgame.State) error {
	if !p.Busted && !p.Stood {
		return errors.New("they have neither busted nor decided to stand")
	}
	return nil
}

func (p *playerState) ResetForTurnStart(state boardgame.State) error {
	p.Stood = false
	return nil
}

func (p *playerState) ResetForTurnEnd(state boardgame.State) error {
	return nil
}

func (p *playerState) EffectiveHand() []*playingcards.Card {
	return append(playingcards.ValuesToCards(p.HiddenHand.ComponentValues()), playingcards.ValuesToCards(p.VisibleHand.ComponentValues())...)
}

//HandValue returns the value of the player's hand.
func (p *playerState) HandValue() int {

	var numUnconvertedAces int
	var currentValue int

	for _, card := range p.EffectiveHand() {
		switch card.Rank.Value() {
		case playingcards.RankAce:
			numUnconvertedAces++
			//We count the ace as 1 now. Later we'll check to see if we can
			//expand any aces.
			currentValue += 1
		case playingcards.RankJack, playingcards.RankQueen, playingcards.RankKing:
			currentValue += 10
		default:
			currentValue += card.Rank.Value()
		}
	}

	for numUnconvertedAces > 0 {

		if currentValue >= (targetScore - 10) {
			break
		}

		numUnconvertedAces--
		currentValue += 10
	}

	return currentValue

}
