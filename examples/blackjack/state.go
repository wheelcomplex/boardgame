package blackjack

import (
	"errors"
	"github.com/jkomoros/boardgame"
	"github.com/jkomoros/boardgame/components/playingcards"
	"github.com/jkomoros/boardgame/enum"
	"github.com/jkomoros/boardgame/moves/moveinterfaces"
)

//+autoreader
const (
	PhaseInitialDeal = iota
	PhaseNormalPlay
)

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
	moveinterfaces.RoundRobinBaseGameState
	Phase         enum.MutableVal        `enum:"Phase"`
	DiscardStack  boardgame.MutableStack `stack:"cards" sanitize:"len"`
	DrawStack     boardgame.MutableStack `stack:"cards" sanitize:"len"`
	UnusedCards   boardgame.MutableStack `stack:"cards"`
	CurrentPlayer boardgame.PlayerIndex
}

//+autoreader
type playerState struct {
	boardgame.BaseSubState
	playerIndex boardgame.PlayerIndex
	HiddenHand  boardgame.MutableStack `stack:"cards,1" sanitize:"len"`
	VisibleHand boardgame.MutableStack `stack:"cards"`
	Busted      bool
	Stood       bool
}

func (g *gameState) SetCurrentPlayer(currentPlayer boardgame.PlayerIndex) {
	g.CurrentPlayer = currentPlayer
}

func (g *gameState) SetCurrentPhase(phase int) {
	g.Phase.SetValue(phase)
}

func (p *playerState) PlayerIndex() boardgame.PlayerIndex {
	return p.playerIndex
}

func (p *playerState) TurnDone() error {
	if !p.Busted && !p.Stood {
		return errors.New("they have neither busted nor decided to stand")
	}
	return nil
}

func (p *playerState) ResetForTurnStart() error {
	p.Stood = false
	return nil
}

func (p *playerState) ResetForTurnEnd() error {
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
