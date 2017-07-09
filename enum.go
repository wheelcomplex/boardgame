package boardgame

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
)

//EnumSet is a set of enums. Normally you will create one in your package, add
//enums to it during initalization, and then use it for all managers you
//create.
type EnumSet struct {
	finished bool
	enums    map[string]*Enum
	//A map of which int goes to which Enum
	values map[int]*Enum
}

//Enum is a named set of values within a set. Get a new one with
//enumSet.Add().
type Enum struct {
	name         string
	set          *EnumSet
	values       map[int]string
	defaultValue int
}

//An EnumValue is an instantiation of a value that must be set to a value in
//the given enum.
type EnumValue struct {
	enum   *Enum
	locked bool
	val    int
}

//NewEnumSet returns a new EnumSet. Generally you'll call this once in a
//package and create the set during initalization.
func NewEnumSet() *EnumSet {
	return &EnumSet{
		false,
		make(map[string]*Enum),
		make(map[int]*Enum),
	}
}

//Finish finalizes an EnumSet so that no more enums may be added. After this
//is called it is safe to use this in a multi-threaded environment. Repeated
//calls do nothing. ComponenChest automatically calls Finish() on the set you
//pass it.
func (e *EnumSet) Finish() {
	e.finished = true
}

//Membership returns the enum that the given val is a member of.
func (e *EnumSet) Membership(val int) *Enum {
	return e.values[val]
}

//MustAdd is like Add, but instead of an error it will panic if the enum
//cannot be added. This is useful for defining your enums at the package level
//outside of an init().
func (e *EnumSet) MustAdd(enumName string, values map[int]string) *Enum {
	result, err := e.Add(enumName, values)

	if err != nil {
		panic("Couldn't add to enumset: " + err.Error())
	}

	return result
}

/*
Add ads an enum with the given name and values to the enum manager. Will error
if that name has already been added, or any of the int values has been used
for any other enum item already. This means that enums must be unique within a
manager. Check out the package doc for the idiomatic way to initalize enums.
*/
func (e *EnumSet) Add(enumName string, values map[int]string) (*Enum, error) {
	if e.finished {
		return nil, errors.New("The set has been finished so no more enums can be added")
	}

	if len(values) == 0 {
		return nil, errors.New("No values provided")
	}

	if _, ok := e.enums[enumName]; ok {
		return nil, errors.New("That enum name has already been provided")
	}

	enum := &Enum{
		enumName,
		e,
		make(map[int]string),
		math.MaxInt64,
	}

	seenValues := make(map[string]bool)

	for v, s := range values {
		if _, ok := e.values[v]; ok {
			//Already registered
			return nil, errors.New("Value " + strconv.Itoa(v) + " was registered twice")
		}

		e.values[v] = enum

		if seenValues[s] {
			return nil, errors.New("String " + s + " was not unique within enum " + enumName)
		}

		seenValues[s] = true

		enum.values[v] = s

		if v < enum.defaultValue {
			enum.defaultValue = v
		}

		e.enums[enumName] = enum
	}
	return enum, nil
}

//DefaultValue returns the default value for this enum (the lowest valid value
//in it).
func (e *Enum) DefaultValue() int {
	return e.defaultValue
}

//Valid returns whether the given value is a valid member of this enum.
func (e *Enum) Valid(val int) bool {
	_, ok := e.values[val]
	return ok
}

//String returns the string value associated with the given value.
func (e *Enum) String(val int) string {
	return e.values[val]
}

//ValueFromString returns the enum value that corresponds to the given string,
//or -1 if no value has that string.
func (e *Enum) ValueFromString(in string) int {
	for v, str := range e.values {
		if str == in {
			return v
		}
	}
	return -1
}

func (e *EnumValue) copy() *EnumValue {
	return &EnumValue{
		e.enum,
		e.locked,
		e.val,
	}
}

//NewEnumValue returns a new EnumValue associated with this enum, set to the
//Enum's DefaultValue to start.
func (e *Enum) NewEnumValue() *EnumValue {
	return &EnumValue{
		e,
		false,
		e.DefaultValue(),
	}
}

//The enum marshals as the string value of the enum so it's more readable.
func (e *EnumValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

//UnmarshalJSON expects the blob to be the string value. Will error if that
//doesn't correspond to a valid value for this enum.
func (e *EnumValue) UnmarshalJSON(blob []byte) error {
	var str string
	if err := json.Unmarshal(blob, &str); err != nil {
		return err
	}
	val := e.enum.ValueFromString(str)
	if val == -1 {
		return errors.New("That string value had no enum in the value")
	}
	return e.SetValue(val)
}

func (e *EnumValue) Enum() *Enum {
	return e.enum
}

func (e *EnumValue) Value() int {
	return e.val
}

func (e *EnumValue) String() string {
	return e.enum.String(e.val)
}

//SetValue changes the value. Returns true if successful. Will fail if the
//value is locked or the val you want to set is not a valid number for the
//enum this value is associated with.
func (e *EnumValue) SetValue(val int) error {
	if e.locked {
		return errors.New("Value is locked")
	}
	if !e.enum.Valid(val) {
		return errors.New("That value is not valid for this enum")
	}
	e.val = val
	return nil
}

//Lock locks in the value of the EnumValue, so that in the future all calls to
//SetValue will fail.
func (e *EnumValue) Lock() {
	e.locked = true
}
