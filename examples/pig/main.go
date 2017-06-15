/*
 *
 * pig is a very simple game involving dice rolls.
 *
 */
package pig

import (
	"github.com/jkomoros/boardgame"
	"math/rand"
	"strconv"
	"strings"
)

//go:generate autoreader

//TODO: this should be configurable, and thus in the gameState.
const TargetScore = 100

type gameDelegate struct {
	boardgame.DefaultGameDelegate
}

func (g *gameDelegate) Name() string {
	return "pig"
}

func (g *gameDelegate) DisplayName() string {
	return "Pig"
}

func (g *gameDelegate) LegalNumPlayers(numPlayers int) bool {
	return numPlayers > 0 && numPlayers < 6
}

func (g *gameDelegate) DefaultNumPlayers() int {
	return 2
}

func (g *gameDelegate) DistributeComponentToStarterStack(state boardgame.State, c *boardgame.Component) (boardgame.Stack, error) {
	game, _ := concreteStates(state)
	return game.Die, nil
}

func (g *gameDelegate) FinishSetUp(state boardgame.MutableState) {

	game, _ := concreteStates(state)

	//Pick a player to start randomly.
	startingPlayer := boardgame.PlayerIndex(rand.Intn(len(state.Players())))

	game.CurrentPlayer = startingPlayer

}

func (g *gameDelegate) CheckGameFinished(state boardgame.State) (finished bool, winners []boardgame.PlayerIndex) {
	_, players := concreteStates(state)

	for i, player := range players {
		if player.TotalScore >= TargetScore {
			winners = append(winners, boardgame.PlayerIndex(i))
		}
	}

	if len(winners) > 0 {
		return true, winners
	}

	return false, nil
}

func (g *gameDelegate) Diagram(state boardgame.State) string {
	var parts []string

	game, players := concreteStates(state)

	dieValue := game.Die.ComponentAt(0).DynamicValues(state).(*dieDynamicValue).Value

	parts = append(parts, "Die: "+strconv.Itoa(dieValue))

	parts = append(parts, "\nPlayers")

	for i, player := range players {
		parts = append(parts, "Player "+strconv.Itoa(i)+": "+strconv.Itoa(player.RoundScore)+", "+strconv.Itoa(player.TotalScore))
	}

	return strings.Join(parts, "\n")
}

func (g *gameDelegate) EmptyGameState() boardgame.MutableSubState {
	dice := g.Manager().Chest().Deck(diceDeckName)

	if dice == nil {
		return nil
	}

	return &gameState{
		CurrentPlayer: 0,
		Die:           boardgame.NewSizedStack(dice, 1),
	}
}

func (g *gameDelegate) EmptyPlayerState(index boardgame.PlayerIndex) boardgame.MutablePlayerState {
	return &playerState{
		playerIndex: index,
		TotalScore:  0,
		RoundScore:  0,
		DieCounted:  true,
		Done:        false,
		Busted:      false,
	}
}

func (g *gameDelegate) EmptyDynamicComponentValues(deck *boardgame.Deck) boardgame.MutableSubState {
	if deck.Name() == diceDeckName {
		return &dieDynamicValue{
			Value: 1,
		}
	}
	return nil
}

func NewManager(storage boardgame.StorageManager) *boardgame.GameManager {
	chest := boardgame.NewComponentChest()

	dice := boardgame.NewDeck()

	dice.AddComponent(DefaultDie())

	chest.AddDeck(diceDeckName, dice)

	manager := boardgame.NewGameManager(&gameDelegate{}, chest, storage)

	if manager == nil {
		panic("No manager returned")
	}

	manager.AddPlayerMoveFactory(MoveRollDiceFactory)
	manager.AddPlayerMoveFactory(MoveDoneTurnFactory)

	manager.AddFixUpMoveFactory(MoveCountDieFactory)
	manager.AddFixUpMoveFactory(MoveAdvanceNextPlayerFactory)

	manager.SetUp()

	return manager
}