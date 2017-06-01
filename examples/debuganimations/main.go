/*

	debuganimations is a very simple debug "game" designed to allow us to
	exercise component animations very directly and purely, in order to build
	and debug that system.

*/
package debuganimations

import (
	"github.com/jkomoros/boardgame"
)

//go:generate autoreader

type gameDelegate struct {
	boardgame.DefaultGameDelegate
}

func (g *gameDelegate) Name() string {
	return "debuganimations"
}

func (g *gameDelegate) DisplayName() string {
	return "Animations Debugger"
}

func (g *gameDelegate) DefaultNumPlayeres() int {
	return 2
}

func (g *gameDelegate) LegalNumPlayers(numPlayers int) bool {
	return numPlayers == 2
}

func (g *gameDelegate) CurrentPlayerIndex(state boardgame.State) boardgame.PlayerIndex {
	game, _ := concreteStates(state)
	return game.CurrentPlayer
}

func (g *gameDelegate) EmptyGameState() boardgame.MutableSubState {

	cards := g.Manager().Chest().Deck(cardsDeckName)

	if cards == nil {
		return nil
	}

	return &gameState{
		CurrentPlayer:    0,
		DrawStack:        boardgame.NewGrowableStack(cards, 0),
		DiscardStack:     boardgame.NewGrowableStack(cards, 0),
		FirstShortStack:  boardgame.NewGrowableStack(cards, 0),
		SecondShortStack: boardgame.NewGrowableStack(cards, 0),
		HiddenCard:       boardgame.NewSizedStack(cards, 1),
		RevealedCard:     boardgame.NewSizedStack(cards, 1),
		FanStack:         boardgame.NewGrowableStack(cards, 0),
		FanDiscard:       boardgame.NewGrowableStack(cards, 0),
	}
}

func (g *gameDelegate) EmptyPlayerState(playerIndex boardgame.PlayerIndex) boardgame.MutablePlayerState {

	cards := g.Manager().Chest().Deck(cardsDeckName)

	if cards == nil {
		return nil
	}

	return &playerState{
		playerIndex: playerIndex,
		Hand:        boardgame.NewGrowableStack(cards, 0),
	}
}

func (g *gameDelegate) DistributeComponentToStarterStack(state boardgame.State, c *boardgame.Component) (boardgame.Stack, error) {
	game, _ := concreteStates(state)

	if game.FirstShortStack.NumComponents() < 1 {
		return game.FirstShortStack, nil
	}

	if game.SecondShortStack.NumComponents() < 1 {
		return game.SecondShortStack, nil
	}

	if game.DiscardStack.NumComponents() < 2 {
		return game.DiscardStack, nil
	}

	if game.HiddenCard.NumComponents() < 1 {
		return game.HiddenCard, nil
	}

	if game.FanStack.NumComponents() < 6 {
		return game.FanStack, nil
	}

	if game.FanDiscard.NumComponents() < 3 {
		return game.FanDiscard, nil
	}

	return game.DrawStack, nil

}

func (g *gameDelegate) FinishSetUp(state boardgame.MutableState) {
	game, _ := concreteStates(state)

	game.DrawStack.Shuffle()

}

func (g *gameDelegate) Diagram(state boardgame.State) string {
	return "Not implemented"
}

var policy *boardgame.StatePolicy

func (g *gameDelegate) StateSanitizationPolicy() *boardgame.StatePolicy {

	if policy == nil {
		policy = &boardgame.StatePolicy{
			Game: map[string]boardgame.GroupPolicy{
				"DrawStack": {
					boardgame.GroupAll: boardgame.PolicyOrder,
				},
				"HiddenCard": {
					boardgame.GroupAll: boardgame.PolicyOrder,
				},
				"FirstShortStack": {
					boardgame.GroupAll: boardgame.PolicyOrder,
				},
				"SecondShortStack": {
					boardgame.GroupAll: boardgame.PolicyOrder,
				},
				"FanDiscard": {
					boardgame.GroupAll: boardgame.PolicyOrder,
				},
			},
		}
	}

	return policy

}

func (g *gameDelegate) CheckGameFinished(state boardgame.State) (finished bool, winners []boardgame.PlayerIndex) {
	//This debug game is never finished
	return false, nil
}

func NewManager(storage boardgame.StorageManager) *boardgame.GameManager {
	chest := boardgame.NewComponentChest()

	cards := boardgame.NewDeck()

	for _, val := range cardNames {
		cards.AddComponentMulti(&cardValue{
			Type: val,
		}, 3)
	}

	cards.SetShadowValues(&cardValue{
		Type: "<hidden>",
	})

	chest.AddDeck(cardsDeckName, cards)

	manager := boardgame.NewGameManager(&gameDelegate{}, chest, storage)

	if manager == nil {
		panic("No manager returned")
	}

	manager.AddPlayerMoveFactory(MoveMoveCardBetweenShortStacksFactory)
	manager.AddPlayerMoveFactory(MoveMoveCardBetweenDrawAndDiscardStacksFactory)
	manager.AddPlayerMoveFactory(MoveFlipHiddenCardFactory)
	manager.AddPlayerMoveFactory(MoveMoveCardBetweenFanStacksFactory)

	manager.SetUp()

	return manager
}
