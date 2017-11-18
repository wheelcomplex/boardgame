/************************************
 *
 * This file contains auto-generated methods to help certain structs
 * implement boardgame.SubState and boardgame.MutableSubState. It was
 * generated by autoreader.
 *
 * DO NOT EDIT by hand.
 *
 ************************************/

package moves

import (
	"errors"
	"github.com/jkomoros/boardgame"
	"github.com/jkomoros/boardgame/enum"
)

// Implementation for gameState

var __gameStateReaderProps map[string]boardgame.PropertyType = map[string]boardgame.PropertyType{
	"Counter":         boardgame.TypeInt,
	"CurrentPlayer":   boardgame.TypePlayerIndex,
	"DiscardStack":    boardgame.TypeStack,
	"DrawStack":       boardgame.TypeStack,
	"Phase":           boardgame.TypeEnum,
	"RRHasStarted":    boardgame.TypeBool,
	"RRLastPlayer":    boardgame.TypePlayerIndex,
	"RRRoundCount":    boardgame.TypeInt,
	"RRStarterPlayer": boardgame.TypePlayerIndex,
}

type __gameStateReader struct {
	data *gameState
}

func (g *__gameStateReader) Props() map[string]boardgame.PropertyType {
	return __gameStateReaderProps
}

func (g *__gameStateReader) Prop(name string) (interface{}, error) {
	props := g.Props()
	propType, ok := props[name]

	if !ok {
		return nil, errors.New("No such property with that name: " + name)
	}

	switch propType {
	case boardgame.TypeBool:
		return g.BoolProp(name)
	case boardgame.TypeBoolSlice:
		return g.BoolSliceProp(name)
	case boardgame.TypeEnum:
		return g.EnumProp(name)
	case boardgame.TypeInt:
		return g.IntProp(name)
	case boardgame.TypeIntSlice:
		return g.IntSliceProp(name)
	case boardgame.TypePlayerIndex:
		return g.PlayerIndexProp(name)
	case boardgame.TypePlayerIndexSlice:
		return g.PlayerIndexSliceProp(name)
	case boardgame.TypeStack:
		return g.StackProp(name)
	case boardgame.TypeString:
		return g.StringProp(name)
	case boardgame.TypeStringSlice:
		return g.StringSliceProp(name)
	case boardgame.TypeTimer:
		return g.TimerProp(name)

	}

	return nil, errors.New("Unexpected property type: " + propType.String())
}

func (g *__gameStateReader) SetProp(name string, value interface{}) error {
	props := g.Props()
	propType, ok := props[name]

	if !ok {
		return errors.New("No such property with that name: " + name)
	}

	switch propType {
	case boardgame.TypeBool:
		val, ok := value.(bool)
		if !ok {
			return errors.New("Provided value was not of type bool")
		}
		return g.SetBoolProp(name, val)
	case boardgame.TypeBoolSlice:
		val, ok := value.([]bool)
		if !ok {
			return errors.New("Provided value was not of type []bool")
		}
		return g.SetBoolSliceProp(name, val)
	case boardgame.TypeInt:
		val, ok := value.(int)
		if !ok {
			return errors.New("Provided value was not of type int")
		}
		return g.SetIntProp(name, val)
	case boardgame.TypeIntSlice:
		val, ok := value.([]int)
		if !ok {
			return errors.New("Provided value was not of type []int")
		}
		return g.SetIntSliceProp(name, val)
	case boardgame.TypeEnum:
		return errors.New("SetProp does not allow setting mutable types. Use ConfigureProp instead.")
	case boardgame.TypeStack:
		return errors.New("SetProp does not allow setting mutable types. Use ConfigureProp instead.")
	case boardgame.TypeTimer:
		return errors.New("SetProp does not allow setting mutable types. Use ConfigureProp instead.")
	case boardgame.TypePlayerIndex:
		val, ok := value.(boardgame.PlayerIndex)
		if !ok {
			return errors.New("Provided value was not of type boardgame.PlayerIndex")
		}
		return g.SetPlayerIndexProp(name, val)
	case boardgame.TypePlayerIndexSlice:
		val, ok := value.([]boardgame.PlayerIndex)
		if !ok {
			return errors.New("Provided value was not of type []boardgame.PlayerIndex")
		}
		return g.SetPlayerIndexSliceProp(name, val)
	case boardgame.TypeString:
		val, ok := value.(string)
		if !ok {
			return errors.New("Provided value was not of type string")
		}
		return g.SetStringProp(name, val)
	case boardgame.TypeStringSlice:
		val, ok := value.([]string)
		if !ok {
			return errors.New("Provided value was not of type []string")
		}
		return g.SetStringSliceProp(name, val)

	}

	return errors.New("Unexpected property type: " + propType.String())
}

func (g *__gameStateReader) ConfigureProp(name string, value interface{}) error {
	props := g.Props()
	propType, ok := props[name]

	if !ok {
		return errors.New("No such property with that name: " + name)
	}

	switch propType {
	case boardgame.TypeBool:
		val, ok := value.(bool)
		if !ok {
			return errors.New("Provided value was not of type bool")
		}
		return g.SetBoolProp(name, val)
	case boardgame.TypeBoolSlice:
		val, ok := value.([]bool)
		if !ok {
			return errors.New("Provided value was not of type []bool")
		}
		return g.SetBoolSliceProp(name, val)
	case boardgame.TypeInt:
		val, ok := value.(int)
		if !ok {
			return errors.New("Provided value was not of type int")
		}
		return g.SetIntProp(name, val)
	case boardgame.TypeIntSlice:
		val, ok := value.([]int)
		if !ok {
			return errors.New("Provided value was not of type []int")
		}
		return g.SetIntSliceProp(name, val)
	case boardgame.TypeEnum:
		val, ok := value.(enum.MutableVal)
		if !ok {
			return errors.New("Provided value was not of type enum.MutableVal")
		}
		return g.ConfigureMutableEnumProp(name, val)
	case boardgame.TypeStack:
		val, ok := value.(boardgame.MutableStack)
		if !ok {
			return errors.New("Provided value was not of type boardgame.MutableStack")
		}
		return g.ConfigureMutableStackProp(name, val)
	case boardgame.TypeTimer:
		val, ok := value.(boardgame.MutableTimer)
		if !ok {
			return errors.New("Provided value was not of type boardgame.MutableTimer")
		}
		return g.ConfigureMutableTimerProp(name, val)
	case boardgame.TypePlayerIndex:
		val, ok := value.(boardgame.PlayerIndex)
		if !ok {
			return errors.New("Provided value was not of type boardgame.PlayerIndex")
		}
		return g.SetPlayerIndexProp(name, val)
	case boardgame.TypePlayerIndexSlice:
		val, ok := value.([]boardgame.PlayerIndex)
		if !ok {
			return errors.New("Provided value was not of type []boardgame.PlayerIndex")
		}
		return g.SetPlayerIndexSliceProp(name, val)
	case boardgame.TypeString:
		val, ok := value.(string)
		if !ok {
			return errors.New("Provided value was not of type string")
		}
		return g.SetStringProp(name, val)
	case boardgame.TypeStringSlice:
		val, ok := value.([]string)
		if !ok {
			return errors.New("Provided value was not of type []string")
		}
		return g.SetStringSliceProp(name, val)

	}

	return errors.New("Unexpected property type: " + propType.String())
}

func (g *__gameStateReader) BoolProp(name string) (bool, error) {

	switch name {
	case "RRHasStarted":
		return g.data.RRHasStarted, nil

	}

	return false, errors.New("No such Bool prop: " + name)

}

func (g *__gameStateReader) SetBoolProp(name string, value bool) error {

	switch name {
	case "RRHasStarted":
		g.data.RRHasStarted = value
		return nil

	}

	return errors.New("No such Bool prop: " + name)

}

func (g *__gameStateReader) BoolSliceProp(name string) ([]bool, error) {

	return []bool{}, errors.New("No such BoolSlice prop: " + name)

}

func (g *__gameStateReader) SetBoolSliceProp(name string, value []bool) error {

	return errors.New("No such BoolSlice prop: " + name)

}

func (g *__gameStateReader) EnumProp(name string) (enum.Val, error) {

	switch name {
	case "Phase":
		return g.data.Phase, nil

	}

	return nil, errors.New("No such Enum prop: " + name)

}

func (g *__gameStateReader) ConfigureMutableEnumProp(name string, value enum.MutableVal) error {

	switch name {
	case "Phase":
		g.data.Phase = value
		return nil

	}

	return errors.New("No such MutableEnum prop: " + name)

}

func (g *__gameStateReader) MutableEnumProp(name string) (enum.MutableVal, error) {

	switch name {
	case "Phase":
		return g.data.Phase, nil

	}

	return nil, errors.New("No such Enum prop: " + name)

}

func (g *__gameStateReader) IntProp(name string) (int, error) {

	switch name {
	case "Counter":
		return g.data.Counter, nil
	case "RRRoundCount":
		return g.data.RRRoundCount, nil

	}

	return 0, errors.New("No such Int prop: " + name)

}

func (g *__gameStateReader) SetIntProp(name string, value int) error {

	switch name {
	case "Counter":
		g.data.Counter = value
		return nil
	case "RRRoundCount":
		g.data.RRRoundCount = value
		return nil

	}

	return errors.New("No such Int prop: " + name)

}

func (g *__gameStateReader) IntSliceProp(name string) ([]int, error) {

	return []int{}, errors.New("No such IntSlice prop: " + name)

}

func (g *__gameStateReader) SetIntSliceProp(name string, value []int) error {

	return errors.New("No such IntSlice prop: " + name)

}

func (g *__gameStateReader) PlayerIndexProp(name string) (boardgame.PlayerIndex, error) {

	switch name {
	case "CurrentPlayer":
		return g.data.CurrentPlayer, nil
	case "RRLastPlayer":
		return g.data.RRLastPlayer, nil
	case "RRStarterPlayer":
		return g.data.RRStarterPlayer, nil

	}

	return 0, errors.New("No such PlayerIndex prop: " + name)

}

func (g *__gameStateReader) SetPlayerIndexProp(name string, value boardgame.PlayerIndex) error {

	switch name {
	case "CurrentPlayer":
		g.data.CurrentPlayer = value
		return nil
	case "RRLastPlayer":
		g.data.RRLastPlayer = value
		return nil
	case "RRStarterPlayer":
		g.data.RRStarterPlayer = value
		return nil

	}

	return errors.New("No such PlayerIndex prop: " + name)

}

func (g *__gameStateReader) PlayerIndexSliceProp(name string) ([]boardgame.PlayerIndex, error) {

	return []boardgame.PlayerIndex{}, errors.New("No such PlayerIndexSlice prop: " + name)

}

func (g *__gameStateReader) SetPlayerIndexSliceProp(name string, value []boardgame.PlayerIndex) error {

	return errors.New("No such PlayerIndexSlice prop: " + name)

}

func (g *__gameStateReader) StackProp(name string) (boardgame.Stack, error) {

	switch name {
	case "DiscardStack":
		return g.data.DiscardStack, nil
	case "DrawStack":
		return g.data.DrawStack, nil

	}

	return nil, errors.New("No such Stack prop: " + name)

}

func (g *__gameStateReader) ConfigureMutableStackProp(name string, value boardgame.MutableStack) error {

	switch name {
	case "DiscardStack":
		g.data.DiscardStack = value
		return nil
	case "DrawStack":
		g.data.DrawStack = value
		return nil

	}

	return errors.New("No such MutableStack prop: " + name)

}

func (g *__gameStateReader) MutableStackProp(name string) (boardgame.MutableStack, error) {

	switch name {
	case "DiscardStack":
		return g.data.DiscardStack, nil
	case "DrawStack":
		return g.data.DrawStack, nil

	}

	return nil, errors.New("No such Stack prop: " + name)

}

func (g *__gameStateReader) StringProp(name string) (string, error) {

	return "", errors.New("No such String prop: " + name)

}

func (g *__gameStateReader) SetStringProp(name string, value string) error {

	return errors.New("No such String prop: " + name)

}

func (g *__gameStateReader) StringSliceProp(name string) ([]string, error) {

	return []string{}, errors.New("No such StringSlice prop: " + name)

}

func (g *__gameStateReader) SetStringSliceProp(name string, value []string) error {

	return errors.New("No such StringSlice prop: " + name)

}

func (g *__gameStateReader) TimerProp(name string) (boardgame.Timer, error) {

	return nil, errors.New("No such Timer prop: " + name)

}

func (g *__gameStateReader) ConfigureMutableTimerProp(name string, value boardgame.MutableTimer) error {

	return errors.New("No such MutableTimer prop: " + name)

}

func (g *__gameStateReader) MutableTimerProp(name string) (boardgame.MutableTimer, error) {

	return nil, errors.New("No such Timer prop: " + name)

}

func (g *gameState) Reader() boardgame.PropertyReader {
	return &__gameStateReader{g}
}

func (g *gameState) ReadSetter() boardgame.PropertyReadSetter {
	return &__gameStateReader{g}
}

func (g *gameState) ReadSetConfigurer() boardgame.PropertyReadSetConfigurer {
	return &__gameStateReader{g}
}

// Implementation for playerState

var __playerStateReaderProps map[string]boardgame.PropertyType = map[string]boardgame.PropertyType{
	"Counter": boardgame.TypeInt,
	"Hand":    boardgame.TypeStack,
}

type __playerStateReader struct {
	data *playerState
}

func (p *__playerStateReader) Props() map[string]boardgame.PropertyType {
	return __playerStateReaderProps
}

func (p *__playerStateReader) Prop(name string) (interface{}, error) {
	props := p.Props()
	propType, ok := props[name]

	if !ok {
		return nil, errors.New("No such property with that name: " + name)
	}

	switch propType {
	case boardgame.TypeBool:
		return p.BoolProp(name)
	case boardgame.TypeBoolSlice:
		return p.BoolSliceProp(name)
	case boardgame.TypeEnum:
		return p.EnumProp(name)
	case boardgame.TypeInt:
		return p.IntProp(name)
	case boardgame.TypeIntSlice:
		return p.IntSliceProp(name)
	case boardgame.TypePlayerIndex:
		return p.PlayerIndexProp(name)
	case boardgame.TypePlayerIndexSlice:
		return p.PlayerIndexSliceProp(name)
	case boardgame.TypeStack:
		return p.StackProp(name)
	case boardgame.TypeString:
		return p.StringProp(name)
	case boardgame.TypeStringSlice:
		return p.StringSliceProp(name)
	case boardgame.TypeTimer:
		return p.TimerProp(name)

	}

	return nil, errors.New("Unexpected property type: " + propType.String())
}

func (p *__playerStateReader) SetProp(name string, value interface{}) error {
	props := p.Props()
	propType, ok := props[name]

	if !ok {
		return errors.New("No such property with that name: " + name)
	}

	switch propType {
	case boardgame.TypeBool:
		val, ok := value.(bool)
		if !ok {
			return errors.New("Provided value was not of type bool")
		}
		return p.SetBoolProp(name, val)
	case boardgame.TypeBoolSlice:
		val, ok := value.([]bool)
		if !ok {
			return errors.New("Provided value was not of type []bool")
		}
		return p.SetBoolSliceProp(name, val)
	case boardgame.TypeInt:
		val, ok := value.(int)
		if !ok {
			return errors.New("Provided value was not of type int")
		}
		return p.SetIntProp(name, val)
	case boardgame.TypeIntSlice:
		val, ok := value.([]int)
		if !ok {
			return errors.New("Provided value was not of type []int")
		}
		return p.SetIntSliceProp(name, val)
	case boardgame.TypeEnum:
		return errors.New("SetProp does not allow setting mutable types. Use ConfigureProp instead.")
	case boardgame.TypeStack:
		return errors.New("SetProp does not allow setting mutable types. Use ConfigureProp instead.")
	case boardgame.TypeTimer:
		return errors.New("SetProp does not allow setting mutable types. Use ConfigureProp instead.")
	case boardgame.TypePlayerIndex:
		val, ok := value.(boardgame.PlayerIndex)
		if !ok {
			return errors.New("Provided value was not of type boardgame.PlayerIndex")
		}
		return p.SetPlayerIndexProp(name, val)
	case boardgame.TypePlayerIndexSlice:
		val, ok := value.([]boardgame.PlayerIndex)
		if !ok {
			return errors.New("Provided value was not of type []boardgame.PlayerIndex")
		}
		return p.SetPlayerIndexSliceProp(name, val)
	case boardgame.TypeString:
		val, ok := value.(string)
		if !ok {
			return errors.New("Provided value was not of type string")
		}
		return p.SetStringProp(name, val)
	case boardgame.TypeStringSlice:
		val, ok := value.([]string)
		if !ok {
			return errors.New("Provided value was not of type []string")
		}
		return p.SetStringSliceProp(name, val)

	}

	return errors.New("Unexpected property type: " + propType.String())
}

func (p *__playerStateReader) ConfigureProp(name string, value interface{}) error {
	props := p.Props()
	propType, ok := props[name]

	if !ok {
		return errors.New("No such property with that name: " + name)
	}

	switch propType {
	case boardgame.TypeBool:
		val, ok := value.(bool)
		if !ok {
			return errors.New("Provided value was not of type bool")
		}
		return p.SetBoolProp(name, val)
	case boardgame.TypeBoolSlice:
		val, ok := value.([]bool)
		if !ok {
			return errors.New("Provided value was not of type []bool")
		}
		return p.SetBoolSliceProp(name, val)
	case boardgame.TypeInt:
		val, ok := value.(int)
		if !ok {
			return errors.New("Provided value was not of type int")
		}
		return p.SetIntProp(name, val)
	case boardgame.TypeIntSlice:
		val, ok := value.([]int)
		if !ok {
			return errors.New("Provided value was not of type []int")
		}
		return p.SetIntSliceProp(name, val)
	case boardgame.TypeEnum:
		val, ok := value.(enum.MutableVal)
		if !ok {
			return errors.New("Provided value was not of type enum.MutableVal")
		}
		return p.ConfigureMutableEnumProp(name, val)
	case boardgame.TypeStack:
		val, ok := value.(boardgame.MutableStack)
		if !ok {
			return errors.New("Provided value was not of type boardgame.MutableStack")
		}
		return p.ConfigureMutableStackProp(name, val)
	case boardgame.TypeTimer:
		val, ok := value.(boardgame.MutableTimer)
		if !ok {
			return errors.New("Provided value was not of type boardgame.MutableTimer")
		}
		return p.ConfigureMutableTimerProp(name, val)
	case boardgame.TypePlayerIndex:
		val, ok := value.(boardgame.PlayerIndex)
		if !ok {
			return errors.New("Provided value was not of type boardgame.PlayerIndex")
		}
		return p.SetPlayerIndexProp(name, val)
	case boardgame.TypePlayerIndexSlice:
		val, ok := value.([]boardgame.PlayerIndex)
		if !ok {
			return errors.New("Provided value was not of type []boardgame.PlayerIndex")
		}
		return p.SetPlayerIndexSliceProp(name, val)
	case boardgame.TypeString:
		val, ok := value.(string)
		if !ok {
			return errors.New("Provided value was not of type string")
		}
		return p.SetStringProp(name, val)
	case boardgame.TypeStringSlice:
		val, ok := value.([]string)
		if !ok {
			return errors.New("Provided value was not of type []string")
		}
		return p.SetStringSliceProp(name, val)

	}

	return errors.New("Unexpected property type: " + propType.String())
}

func (p *__playerStateReader) BoolProp(name string) (bool, error) {

	return false, errors.New("No such Bool prop: " + name)

}

func (p *__playerStateReader) SetBoolProp(name string, value bool) error {

	return errors.New("No such Bool prop: " + name)

}

func (p *__playerStateReader) BoolSliceProp(name string) ([]bool, error) {

	return []bool{}, errors.New("No such BoolSlice prop: " + name)

}

func (p *__playerStateReader) SetBoolSliceProp(name string, value []bool) error {

	return errors.New("No such BoolSlice prop: " + name)

}

func (p *__playerStateReader) EnumProp(name string) (enum.Val, error) {

	return nil, errors.New("No such Enum prop: " + name)

}

func (p *__playerStateReader) ConfigureMutableEnumProp(name string, value enum.MutableVal) error {

	return errors.New("No such MutableEnum prop: " + name)

}

func (p *__playerStateReader) MutableEnumProp(name string) (enum.MutableVal, error) {

	return nil, errors.New("No such Enum prop: " + name)

}

func (p *__playerStateReader) IntProp(name string) (int, error) {

	switch name {
	case "Counter":
		return p.data.Counter, nil

	}

	return 0, errors.New("No such Int prop: " + name)

}

func (p *__playerStateReader) SetIntProp(name string, value int) error {

	switch name {
	case "Counter":
		p.data.Counter = value
		return nil

	}

	return errors.New("No such Int prop: " + name)

}

func (p *__playerStateReader) IntSliceProp(name string) ([]int, error) {

	return []int{}, errors.New("No such IntSlice prop: " + name)

}

func (p *__playerStateReader) SetIntSliceProp(name string, value []int) error {

	return errors.New("No such IntSlice prop: " + name)

}

func (p *__playerStateReader) PlayerIndexProp(name string) (boardgame.PlayerIndex, error) {

	return 0, errors.New("No such PlayerIndex prop: " + name)

}

func (p *__playerStateReader) SetPlayerIndexProp(name string, value boardgame.PlayerIndex) error {

	return errors.New("No such PlayerIndex prop: " + name)

}

func (p *__playerStateReader) PlayerIndexSliceProp(name string) ([]boardgame.PlayerIndex, error) {

	return []boardgame.PlayerIndex{}, errors.New("No such PlayerIndexSlice prop: " + name)

}

func (p *__playerStateReader) SetPlayerIndexSliceProp(name string, value []boardgame.PlayerIndex) error {

	return errors.New("No such PlayerIndexSlice prop: " + name)

}

func (p *__playerStateReader) StackProp(name string) (boardgame.Stack, error) {

	switch name {
	case "Hand":
		return p.data.Hand, nil

	}

	return nil, errors.New("No such Stack prop: " + name)

}

func (p *__playerStateReader) ConfigureMutableStackProp(name string, value boardgame.MutableStack) error {

	switch name {
	case "Hand":
		p.data.Hand = value
		return nil

	}

	return errors.New("No such MutableStack prop: " + name)

}

func (p *__playerStateReader) MutableStackProp(name string) (boardgame.MutableStack, error) {

	switch name {
	case "Hand":
		return p.data.Hand, nil

	}

	return nil, errors.New("No such Stack prop: " + name)

}

func (p *__playerStateReader) StringProp(name string) (string, error) {

	return "", errors.New("No such String prop: " + name)

}

func (p *__playerStateReader) SetStringProp(name string, value string) error {

	return errors.New("No such String prop: " + name)

}

func (p *__playerStateReader) StringSliceProp(name string) ([]string, error) {

	return []string{}, errors.New("No such StringSlice prop: " + name)

}

func (p *__playerStateReader) SetStringSliceProp(name string, value []string) error {

	return errors.New("No such StringSlice prop: " + name)

}

func (p *__playerStateReader) TimerProp(name string) (boardgame.Timer, error) {

	return nil, errors.New("No such Timer prop: " + name)

}

func (p *__playerStateReader) ConfigureMutableTimerProp(name string, value boardgame.MutableTimer) error {

	return errors.New("No such MutableTimer prop: " + name)

}

func (p *__playerStateReader) MutableTimerProp(name string) (boardgame.MutableTimer, error) {

	return nil, errors.New("No such Timer prop: " + name)

}

func (p *playerState) Reader() boardgame.PropertyReader {
	return &__playerStateReader{p}
}

func (p *playerState) ReadSetter() boardgame.PropertyReadSetter {
	return &__playerStateReader{p}
}

func (p *playerState) ReadSetConfigurer() boardgame.PropertyReadSetConfigurer {
	return &__playerStateReader{p}
}

// Implementation for StartPhase

var __StartPhaseReaderProps map[string]boardgame.PropertyType = map[string]boardgame.PropertyType{}

type __StartPhaseReader struct {
	data *StartPhase
}

func (s *__StartPhaseReader) Props() map[string]boardgame.PropertyType {
	return __StartPhaseReaderProps
}

func (s *__StartPhaseReader) Prop(name string) (interface{}, error) {
	props := s.Props()
	propType, ok := props[name]

	if !ok {
		return nil, errors.New("No such property with that name: " + name)
	}

	switch propType {
	case boardgame.TypeBool:
		return s.BoolProp(name)
	case boardgame.TypeBoolSlice:
		return s.BoolSliceProp(name)
	case boardgame.TypeEnum:
		return s.EnumProp(name)
	case boardgame.TypeInt:
		return s.IntProp(name)
	case boardgame.TypeIntSlice:
		return s.IntSliceProp(name)
	case boardgame.TypePlayerIndex:
		return s.PlayerIndexProp(name)
	case boardgame.TypePlayerIndexSlice:
		return s.PlayerIndexSliceProp(name)
	case boardgame.TypeStack:
		return s.StackProp(name)
	case boardgame.TypeString:
		return s.StringProp(name)
	case boardgame.TypeStringSlice:
		return s.StringSliceProp(name)
	case boardgame.TypeTimer:
		return s.TimerProp(name)

	}

	return nil, errors.New("Unexpected property type: " + propType.String())
}

func (s *__StartPhaseReader) SetProp(name string, value interface{}) error {
	props := s.Props()
	propType, ok := props[name]

	if !ok {
		return errors.New("No such property with that name: " + name)
	}

	switch propType {
	case boardgame.TypeBool:
		val, ok := value.(bool)
		if !ok {
			return errors.New("Provided value was not of type bool")
		}
		return s.SetBoolProp(name, val)
	case boardgame.TypeBoolSlice:
		val, ok := value.([]bool)
		if !ok {
			return errors.New("Provided value was not of type []bool")
		}
		return s.SetBoolSliceProp(name, val)
	case boardgame.TypeInt:
		val, ok := value.(int)
		if !ok {
			return errors.New("Provided value was not of type int")
		}
		return s.SetIntProp(name, val)
	case boardgame.TypeIntSlice:
		val, ok := value.([]int)
		if !ok {
			return errors.New("Provided value was not of type []int")
		}
		return s.SetIntSliceProp(name, val)
	case boardgame.TypeEnum:
		return errors.New("SetProp does not allow setting mutable types. Use ConfigureProp instead.")
	case boardgame.TypeStack:
		return errors.New("SetProp does not allow setting mutable types. Use ConfigureProp instead.")
	case boardgame.TypeTimer:
		return errors.New("SetProp does not allow setting mutable types. Use ConfigureProp instead.")
	case boardgame.TypePlayerIndex:
		val, ok := value.(boardgame.PlayerIndex)
		if !ok {
			return errors.New("Provided value was not of type boardgame.PlayerIndex")
		}
		return s.SetPlayerIndexProp(name, val)
	case boardgame.TypePlayerIndexSlice:
		val, ok := value.([]boardgame.PlayerIndex)
		if !ok {
			return errors.New("Provided value was not of type []boardgame.PlayerIndex")
		}
		return s.SetPlayerIndexSliceProp(name, val)
	case boardgame.TypeString:
		val, ok := value.(string)
		if !ok {
			return errors.New("Provided value was not of type string")
		}
		return s.SetStringProp(name, val)
	case boardgame.TypeStringSlice:
		val, ok := value.([]string)
		if !ok {
			return errors.New("Provided value was not of type []string")
		}
		return s.SetStringSliceProp(name, val)

	}

	return errors.New("Unexpected property type: " + propType.String())
}

func (s *__StartPhaseReader) ConfigureProp(name string, value interface{}) error {
	props := s.Props()
	propType, ok := props[name]

	if !ok {
		return errors.New("No such property with that name: " + name)
	}

	switch propType {
	case boardgame.TypeBool:
		val, ok := value.(bool)
		if !ok {
			return errors.New("Provided value was not of type bool")
		}
		return s.SetBoolProp(name, val)
	case boardgame.TypeBoolSlice:
		val, ok := value.([]bool)
		if !ok {
			return errors.New("Provided value was not of type []bool")
		}
		return s.SetBoolSliceProp(name, val)
	case boardgame.TypeInt:
		val, ok := value.(int)
		if !ok {
			return errors.New("Provided value was not of type int")
		}
		return s.SetIntProp(name, val)
	case boardgame.TypeIntSlice:
		val, ok := value.([]int)
		if !ok {
			return errors.New("Provided value was not of type []int")
		}
		return s.SetIntSliceProp(name, val)
	case boardgame.TypeEnum:
		val, ok := value.(enum.MutableVal)
		if !ok {
			return errors.New("Provided value was not of type enum.MutableVal")
		}
		return s.ConfigureMutableEnumProp(name, val)
	case boardgame.TypeStack:
		val, ok := value.(boardgame.MutableStack)
		if !ok {
			return errors.New("Provided value was not of type boardgame.MutableStack")
		}
		return s.ConfigureMutableStackProp(name, val)
	case boardgame.TypeTimer:
		val, ok := value.(boardgame.MutableTimer)
		if !ok {
			return errors.New("Provided value was not of type boardgame.MutableTimer")
		}
		return s.ConfigureMutableTimerProp(name, val)
	case boardgame.TypePlayerIndex:
		val, ok := value.(boardgame.PlayerIndex)
		if !ok {
			return errors.New("Provided value was not of type boardgame.PlayerIndex")
		}
		return s.SetPlayerIndexProp(name, val)
	case boardgame.TypePlayerIndexSlice:
		val, ok := value.([]boardgame.PlayerIndex)
		if !ok {
			return errors.New("Provided value was not of type []boardgame.PlayerIndex")
		}
		return s.SetPlayerIndexSliceProp(name, val)
	case boardgame.TypeString:
		val, ok := value.(string)
		if !ok {
			return errors.New("Provided value was not of type string")
		}
		return s.SetStringProp(name, val)
	case boardgame.TypeStringSlice:
		val, ok := value.([]string)
		if !ok {
			return errors.New("Provided value was not of type []string")
		}
		return s.SetStringSliceProp(name, val)

	}

	return errors.New("Unexpected property type: " + propType.String())
}

func (s *__StartPhaseReader) BoolProp(name string) (bool, error) {

	return false, errors.New("No such Bool prop: " + name)

}

func (s *__StartPhaseReader) SetBoolProp(name string, value bool) error {

	return errors.New("No such Bool prop: " + name)

}

func (s *__StartPhaseReader) BoolSliceProp(name string) ([]bool, error) {

	return []bool{}, errors.New("No such BoolSlice prop: " + name)

}

func (s *__StartPhaseReader) SetBoolSliceProp(name string, value []bool) error {

	return errors.New("No such BoolSlice prop: " + name)

}

func (s *__StartPhaseReader) EnumProp(name string) (enum.Val, error) {

	return nil, errors.New("No such Enum prop: " + name)

}

func (s *__StartPhaseReader) ConfigureMutableEnumProp(name string, value enum.MutableVal) error {

	return errors.New("No such MutableEnum prop: " + name)

}

func (s *__StartPhaseReader) MutableEnumProp(name string) (enum.MutableVal, error) {

	return nil, errors.New("No such Enum prop: " + name)

}

func (s *__StartPhaseReader) IntProp(name string) (int, error) {

	return 0, errors.New("No such Int prop: " + name)

}

func (s *__StartPhaseReader) SetIntProp(name string, value int) error {

	return errors.New("No such Int prop: " + name)

}

func (s *__StartPhaseReader) IntSliceProp(name string) ([]int, error) {

	return []int{}, errors.New("No such IntSlice prop: " + name)

}

func (s *__StartPhaseReader) SetIntSliceProp(name string, value []int) error {

	return errors.New("No such IntSlice prop: " + name)

}

func (s *__StartPhaseReader) PlayerIndexProp(name string) (boardgame.PlayerIndex, error) {

	return 0, errors.New("No such PlayerIndex prop: " + name)

}

func (s *__StartPhaseReader) SetPlayerIndexProp(name string, value boardgame.PlayerIndex) error {

	return errors.New("No such PlayerIndex prop: " + name)

}

func (s *__StartPhaseReader) PlayerIndexSliceProp(name string) ([]boardgame.PlayerIndex, error) {

	return []boardgame.PlayerIndex{}, errors.New("No such PlayerIndexSlice prop: " + name)

}

func (s *__StartPhaseReader) SetPlayerIndexSliceProp(name string, value []boardgame.PlayerIndex) error {

	return errors.New("No such PlayerIndexSlice prop: " + name)

}

func (s *__StartPhaseReader) StackProp(name string) (boardgame.Stack, error) {

	return nil, errors.New("No such Stack prop: " + name)

}

func (s *__StartPhaseReader) ConfigureMutableStackProp(name string, value boardgame.MutableStack) error {

	return errors.New("No such MutableStack prop: " + name)

}

func (s *__StartPhaseReader) MutableStackProp(name string) (boardgame.MutableStack, error) {

	return nil, errors.New("No such Stack prop: " + name)

}

func (s *__StartPhaseReader) StringProp(name string) (string, error) {

	return "", errors.New("No such String prop: " + name)

}

func (s *__StartPhaseReader) SetStringProp(name string, value string) error {

	return errors.New("No such String prop: " + name)

}

func (s *__StartPhaseReader) StringSliceProp(name string) ([]string, error) {

	return []string{}, errors.New("No such StringSlice prop: " + name)

}

func (s *__StartPhaseReader) SetStringSliceProp(name string, value []string) error {

	return errors.New("No such StringSlice prop: " + name)

}

func (s *__StartPhaseReader) TimerProp(name string) (boardgame.Timer, error) {

	return nil, errors.New("No such Timer prop: " + name)

}

func (s *__StartPhaseReader) ConfigureMutableTimerProp(name string, value boardgame.MutableTimer) error {

	return errors.New("No such MutableTimer prop: " + name)

}

func (s *__StartPhaseReader) MutableTimerProp(name string) (boardgame.MutableTimer, error) {

	return nil, errors.New("No such Timer prop: " + name)

}

func (s *StartPhase) Reader() boardgame.PropertyReader {
	return &__StartPhaseReader{s}
}

func (s *StartPhase) ReadSetter() boardgame.PropertyReadSetter {
	return &__StartPhaseReader{s}
}

func (s *StartPhase) ReadSetConfigurer() boardgame.PropertyReadSetConfigurer {
	return &__StartPhaseReader{s}
}
