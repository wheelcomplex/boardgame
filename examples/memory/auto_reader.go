/************************************
 *
 * This file contains auto-generated methods to help certain structs
 * implement boardgame.SubState and boardgame.MutableSubState. It was
 * generated by autoreader.
 *
 * DO NOT EDIT by hand.
 *
 ************************************/
package memory

import (
	"github.com/jkomoros/boardgame"
)

// Implementation for cardValue

func (c *cardValue) Reader() boardgame.PropertyReader {
	return boardgame.DefaultReader(c)
}

// Implementation for MoveAdvanceNextPlayer

func (M *MoveAdvanceNextPlayer) ReadSetter() boardgame.PropertyReadSetter {
	return boardgame.DefaultReadSetter(M)
}

// Implementation for MoveRevealCard

func (M *MoveRevealCard) ReadSetter() boardgame.PropertyReadSetter {
	return boardgame.DefaultReadSetter(M)
}

// Implementation for MoveStartHideCardsTimer

func (M *MoveStartHideCardsTimer) ReadSetter() boardgame.PropertyReadSetter {
	return boardgame.DefaultReadSetter(M)
}

// Implementation for MoveCaptureCards

func (M *MoveCaptureCards) ReadSetter() boardgame.PropertyReadSetter {
	return boardgame.DefaultReadSetter(M)
}

// Implementation for MoveHideCards

func (M *MoveHideCards) ReadSetter() boardgame.PropertyReadSetter {
	return boardgame.DefaultReadSetter(M)
}

// Implementation for gameState

func (g *gameState) Reader() boardgame.PropertyReader {
	return boardgame.DefaultReader(g)
}

func (g *gameState) ReadSetter() boardgame.PropertyReadSetter {
	return boardgame.DefaultReadSetter(g)
}

// Implementation for playerState

func (p *playerState) Reader() boardgame.PropertyReader {
	return boardgame.DefaultReader(p)
}

func (p *playerState) ReadSetter() boardgame.PropertyReadSetter {
	return boardgame.DefaultReadSetter(p)
}
