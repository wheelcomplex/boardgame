/************************************
 *
 * This file contains auto-generated methods to help certain structs
 * implement boardgame.SubState and boardgame.MutableSubState. It was
 * generated by autoreader.
 *
 * DO NOT EDIT by hand.
 *
 ************************************/
package playingcards

import (
	"errors"
	"github.com/jkomoros/boardgame"
	"github.com/jkomoros/boardgame/enum"
)

// Implementation for Card

var __CardReaderProps map[string]boardgame.PropertyType = map[string]boardgame.PropertyType{
	"Rank": boardgame.TypeEnumConst,
	"Suit": boardgame.TypeEnumConst,
}

type __CardReader struct {
	data *Card
}

func (c *__CardReader) Props() map[string]boardgame.PropertyType {
	return __CardReaderProps
}

func (c *__CardReader) Prop(name string) (interface{}, error) {
	props := c.Props()
	propType, ok := props[name]

	if !ok {
		return nil, errors.New("No such property with that name: " + name)
	}

	switch propType {
	case boardgame.TypeBool:
		return c.BoolProp(name)
	case boardgame.TypeBoolSlice:
		return c.BoolSliceProp(name)
	case boardgame.TypeEnumConst:
		return c.EnumConstProp(name)
	case boardgame.TypeEnumVar:
		return c.EnumVarProp(name)
	case boardgame.TypeGrowableStack:
		return c.GrowableStackProp(name)
	case boardgame.TypeInt:
		return c.IntProp(name)
	case boardgame.TypeIntSlice:
		return c.IntSliceProp(name)
	case boardgame.TypePlayerIndex:
		return c.PlayerIndexProp(name)
	case boardgame.TypePlayerIndexSlice:
		return c.PlayerIndexSliceProp(name)
	case boardgame.TypeSizedStack:
		return c.SizedStackProp(name)
	case boardgame.TypeString:
		return c.StringProp(name)
	case boardgame.TypeStringSlice:
		return c.StringSliceProp(name)
	case boardgame.TypeTimer:
		return c.TimerProp(name)

	}

	return nil, errors.New("Unexpected property type: " + propType.String())
}

func (c *__CardReader) BoolProp(name string) (bool, error) {

	return false, errors.New("No such Bool prop: " + name)

}

func (c *__CardReader) BoolSliceProp(name string) ([]bool, error) {

	return []bool{}, errors.New("No such BoolSlice prop: " + name)

}

func (c *__CardReader) EnumConstProp(name string) (enum.Const, error) {

	switch name {
	case "Suit":
		return c.data.Suit, nil
	case "Rank":
		return c.data.Rank, nil

	}

	return nil, errors.New("No such EnumConst prop: " + name)

}

func (c *__CardReader) EnumVarProp(name string) (enum.Var, error) {

	return nil, errors.New("No such EnumVar prop: " + name)

}

func (c *__CardReader) GrowableStackProp(name string) (*boardgame.GrowableStack, error) {

	return nil, errors.New("No such GrowableStack prop: " + name)

}

func (c *__CardReader) IntProp(name string) (int, error) {

	return 0, errors.New("No such Int prop: " + name)

}

func (c *__CardReader) IntSliceProp(name string) ([]int, error) {

	return []int{}, errors.New("No such IntSlice prop: " + name)

}

func (c *__CardReader) PlayerIndexProp(name string) (boardgame.PlayerIndex, error) {

	return 0, errors.New("No such PlayerIndex prop: " + name)

}

func (c *__CardReader) PlayerIndexSliceProp(name string) ([]boardgame.PlayerIndex, error) {

	return []boardgame.PlayerIndex{}, errors.New("No such PlayerIndexSlice prop: " + name)

}

func (c *__CardReader) SizedStackProp(name string) (*boardgame.SizedStack, error) {

	return nil, errors.New("No such SizedStack prop: " + name)

}

func (c *__CardReader) StringProp(name string) (string, error) {

	return "", errors.New("No such String prop: " + name)

}

func (c *__CardReader) StringSliceProp(name string) ([]string, error) {

	return []string{}, errors.New("No such StringSlice prop: " + name)

}

func (c *__CardReader) TimerProp(name string) (*boardgame.Timer, error) {

	return nil, errors.New("No such Timer prop: " + name)

}

func (c *Card) Reader() boardgame.PropertyReader {
	return &__CardReader{c}
}
