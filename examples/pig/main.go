/*

pig is a very simple game involving dice rolls.

*/
package pig

import (
	"errors"
	"github.com/jkomoros/boardgame"
	"github.com/jkomoros/boardgame/components/dice"
	"math/rand"
	"strconv"
	"strings"
)

//go:generate autoreader

const DefaultTargetScore = 100
const diceDeckName = "dice"

type gameDelegate struct {
	boardgame.DefaultGameDelegate
}

func (g *gameDelegate) Name() string {
	return "pig"
}

func (g *gameDelegate) DisplayName() string {
	return "Pig"
}

func (g *gameDelegate) Description() string {
	return "Players roll the dice, collecting points, but bust if they roll a one."
}

func (g *gameDelegate) MinNumPlayers() int {
	return 2
}

func (g *gameDelegate) MaxNumPlayers() int {
	return 6
}

func (g *gameDelegate) DefaultNumPlayers() int {
	return 2
}

func (g *gameDelegate) DistributeComponentToStarterStack(state boardgame.State, c *boardgame.Component) (boardgame.Stack, error) {
	game, _ := concreteStates(state)
	return game.Die, nil
}

func (g *gameDelegate) FinishSetUp(state boardgame.MutableState) error {

	game, _ := concreteStates(state)

	//Pick a player to start randomly.
	startingPlayer := boardgame.PlayerIndex(rand.Intn(len(state.PlayerStates())))

	game.CurrentPlayer = startingPlayer
	game.TargetScore = DefaultTargetScore

	return nil

}

func (g *gameDelegate) GameEndConditionMet(state boardgame.State) bool {
	game, players := concreteStates(state)

	for _, player := range players {
		if player.TotalScore >= game.TargetScore {
			return true
		}
	}

	return false
}

func (g *gameDelegate) PlayerScore(pState boardgame.PlayerState) int {
	return pState.(*playerState).TotalScore
}

func (g *gameDelegate) Diagram(state boardgame.State) string {
	var parts []string

	game, players := concreteStates(state)

	dieValue := game.Die.ComponentAt(0).DynamicValues(state).(*dice.DynamicValue).Value

	parts = append(parts, "Die: "+strconv.Itoa(dieValue))

	parts = append(parts, "\nPlayers")

	for i, player := range players {
		parts = append(parts, "Player "+strconv.Itoa(i)+": "+strconv.Itoa(player.RoundScore)+", "+strconv.Itoa(player.TotalScore))
	}

	return strings.Join(parts, "\n")
}

func (g *gameDelegate) GameStateConstructor() boardgame.ConfigurableSubState {
	return &gameState{
		CurrentPlayer: 0,
		TargetScore:   DefaultTargetScore,
	}
}

func (g *gameDelegate) PlayerStateConstructor(index boardgame.PlayerIndex) boardgame.ConfigurablePlayerState {
	return &playerState{
		playerIndex: index,
	}
}

func (g *gameDelegate) DynamicComponentValuesConstructor(deck *boardgame.Deck) boardgame.ConfigurableSubState {
	if deck.Name() == diceDeckName {
		return &dice.DynamicValue{
			Value: 1,
		}
	}
	return nil
}

func (g *gameDelegate) ConfigureMoves() *boardgame.MoveTypeConfigBundle {
	return boardgame.NewMoveTypeConfigBundle().AddMoves(
		&moveRollDiceConfig,
		&moveDoneTurnConfig,
		&moveCountDieConfig,
		&moveFinishTurnConfig,
	)
}

func NewManager(storage boardgame.StorageManager) (*boardgame.GameManager, error) {
	chest := boardgame.NewComponentChest(nil)

	diceDeck := boardgame.NewDeck()

	diceDeck.AddComponent(dice.DefaultDie())

	if err := chest.AddDeck(diceDeckName, diceDeck); err != nil {
		return nil, errors.New("Couldn't add deck: " + err.Error())
	}

	return boardgame.NewGameManager(&gameDelegate{}, chest, storage)

}
