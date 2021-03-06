package boardgame

import (
	"github.com/jkomoros/boardgame/enum"
	"github.com/jkomoros/boardgame/errors"
	"reflect"
	"strconv"
	"strings"
)

const enumStructTag = "enum"
const stackStructTag = "stack"
const fixedStackStructTag = "sizedstack"
const sanitizationStructTag = "sanitize"

type autoStackConfig struct {
	deck      *Deck
	size      int
	fixedSize bool
}

type readerValidator struct {
	autoEnumFields     map[string]enum.Enum
	autoStackFields    map[string]*autoStackConfig
	sanitizationPolicy map[string]map[int]Policy
	illegalTypes       map[PropertyType]bool
}

//newReaderValidator returns a new readerValidator configured to disallow the
//given types. It will also do an expensive processing for any nil pointer-
//properties to see if they have struct tags that tell us how to inflate them.
//This processing uses reflection, but afterwards AutoInflate can run quickly.
func newReaderValidator(exampleReader PropertyReader, exampleObj interface{}, illegalTypes map[PropertyType]bool, chest *ComponentChest, isPlayerState bool) (*readerValidator, error) {
	//TODO: there's got to be a way to not need both exampleReader and exampleObj, but only one.

	if chest == nil {
		return nil, errors.New("Passed nil chest")
	}

	if illegalTypes == nil {
		illegalTypes = make(map[PropertyType]bool)
	}

	autoEnumFields := make(map[string]enum.Enum)
	autoStackFields := make(map[string]*autoStackConfig)
	sanitizationPolicy := make(map[string]map[int]Policy)

	defaultGroup := "all"
	if isPlayerState {
		defaultGroup = "other"
	}

	for propName, propType := range exampleReader.Props() {

		sanitizationPolicy[propName] = policyFromStructTag(structTagForField(exampleObj, propName, sanitizationStructTag), defaultGroup)

		switch propType {
		case TypeStack:
			stack, err := exampleReader.StackProp(propName)
			if err != nil {
				return nil, errors.New("Couldn't fetch stack prop: " + propName)
			}
			if stack != nil {
				//This stack prop is already non-nil, so we don't need to do
				//any processing to tell how to inflate it.
				continue
			}

			isFixed := false

			tag := structTagForField(exampleObj, propName, stackStructTag)

			if tag == "" {
				tag = structTagForField(exampleObj, propName, fixedStackStructTag)
				if tag != "" {
					isFixed = true
				}
			}

			if tag != "" {

				deck, size, err := unpackStackStructTag(tag, chest)

				if err != nil {
					return nil, errors.New(propName + " was a nil SizedStack and its struct tag was not valid: " + err.Error())
				}

				if isFixed && size == 0 {
					//Size for sizedstacks defaults to 1 (which can be grown
					//easily to any other size).
					size = 1
				}

				autoStackFields[propName] = &autoStackConfig{
					deck,
					size,
					isFixed,
				}
			}
		case TypeEnum:
			enumConst, err := exampleReader.EnumProp(propName)
			if err != nil {
				return nil, errors.New("Couldn't fetch enum  prop: " + propName)
			}
			if enumConst != nil {
				//This enum prop is already non-nil, so we don't need to do
				//any processing to tell how to inflate it.
				continue
			}
			if enumName := structTagForField(exampleObj, propName, enumStructTag); enumName != "" {
				theEnum := chest.Enums().Enum(enumName)
				if theEnum == nil {
					return nil, errors.New(propName + " was a nil enum.Val and the struct tag named " + enumName + " was not a valid enum.")
				}
				//Found one!
				autoEnumFields[propName] = theEnum
			}
		}

	}

	result := &readerValidator{
		autoEnumFields,
		autoStackFields,
		sanitizationPolicy,
		illegalTypes,
	}

	if err := result.VerifyNoIllegalProps(exampleReader); err != nil {
		return nil, errors.New("Example had illegal prop types: " + err.Error())
	}

	return result, nil
}

func policyFromStructTag(tag string, defaultGroup string) map[int]Policy {
	if tag == "" {
		tag = "visible"
	}

	errorMap := map[int]Policy{
		GroupAll: PolicyInvalid,
	}

	result := make(map[int]Policy)

	pieces := strings.Split(tag, ",")
	for _, piece := range pieces {
		splitPiece := strings.Split(piece, ":")
		var groupString string
		var policyString string
		if len(splitPiece) > 2 {
			return errorMap
		}
		if len(splitPiece) == 1 {
			groupString = defaultGroup
			policyString = splitPiece[0]
		} else {
			groupString = splitPiece[0]
			policyString = splitPiece[1]
		}

		group := groupFromString(groupString)
		policy := policyFromString(policyString)

		result[group] = policy

	}

	return result

}

//AutoInflate will go through and inflate fields that are nil that it knows
//how to inflate due to comments in structs detected in the constructor for
//this validator.
func (r *readerValidator) AutoInflate(readSetConfigurer PropertyReadSetConfigurer, st State) error {

	for propName, config := range r.autoStackFields {

		stack, err := readSetConfigurer.MutableStackProp(propName)
		if stack != nil {
			//Guess it was already set!
			continue
		}
		if err != nil {
			return errors.New(propName + " had error fetching stack: " + err.Error())
		}
		if config == nil {
			return errors.New("The config for " + propName + " was unexpectedly nil")
		}
		if config.deck == nil {
			return errors.New("The deck for " + propName + " was unexpectedly nil")
		}

		if config.fixedSize {
			stack = config.deck.NewSizedStack(config.size)
		} else {
			stack = config.deck.NewStack(config.size)
		}

		if err := readSetConfigurer.ConfigureMutableStackProp(propName, stack); err != nil {
			return errors.New("Couldn't set " + propName + " to stack: " + err.Error())
		}
	}

	for propName, enum := range r.autoEnumFields {
		enumConst, err := readSetConfigurer.EnumProp(propName)
		if enumConst != nil {
			//Guess it was already set!
			continue
		}
		if err != nil {
			return errors.New(propName + " had error fetching Enum: " + err.Error())
		}
		if enum == nil {
			return errors.New("The enum for " + propName + " was unexpectedly nil")
		}
		if err := readSetConfigurer.ConfigureMutableEnumProp(propName, enum.NewMutableVal()); err != nil {
			return errors.New("Couldn't set " + propName + " to NewDefaultVal: " + err.Error())
		}
	}

	for propName, propType := range readSetConfigurer.Props() {
		switch propType {
		case TypeTimer:
			timer := NewTimer()
			if err := readSetConfigurer.ConfigureMutableTimerProp(propName, timer); err != nil {
				return errors.New("Couldn't set " + propName + " to a new timer: " + err.Error())
			}
		}
	}

	if st != nil {
		if err := setReaderStatePtr(readSetConfigurer, st); err != nil {
			return errors.New("Couldn't set state ptrs: " + err.Error())
		}
	}

	//TODO: process Stack, Timer fields (convert to state pointer if non-nil)
	return nil
}

func (r *readerValidator) VerifyNoIllegalProps(reader PropertyReader) error {
	for propName, propType := range reader.Props() {
		if propType == TypeIllegal {
			return errors.New(propName + " was TypeIllegal, which is always illegal")
		}
		if _, illegal := r.illegalTypes[propType]; illegal {
			return errors.New(propName + " was the type " + propType.String() + ", which is illegal in this context")
		}
	}
	return nil
}

//Valid will return an error if the reader is not valid according to the
//configuration of this validator.
func (r *readerValidator) Valid(reader PropertyReader) error {
	if err := r.VerifyNoIllegalProps(reader); err != nil {
		return err
	}
	for propName, propType := range reader.Props() {

		policyMap := r.sanitizationPolicy[propName]

		if policyMap == nil {
			return errors.New(propName + " had no sanitization policy")
		}

		for group, policy := range policyMap {
			if policy == PolicyInvalid {
				return errors.New(propName + " had invalid policy for group " + strconv.Itoa(group))
			}
		}

		//TODO: verifyReader should be gotten rid of in favor of this
		switch propType {
		case TypeStack:
			val, err := reader.StackProp(propName)
			if val == nil {
				return errors.New("Stack Prop " + propName + " was nil")
			}
			if err != nil {
				return errors.New("Stack prop " + propName + " had unexpected error: " + err.Error())
			}
			if val.state() == nil {
				return errors.New("Stack prop " + propName + " didn't have its state set")
			}
		case TypeTimer:
			val, err := reader.TimerProp(propName)
			if val == nil {
				return errors.New("TimerProp " + propName + " was nil")
			}
			if err != nil {
				return errors.New("TimerProp " + propName + " had unexpected error: " + err.Error())
			}
			if val.state() == nil {
				return errors.New("TimerProp " + propName + " didn't have its statePtr set")
			}
		case TypeEnum:
			val, err := reader.EnumProp(propName)
			if val == nil {
				return errors.New("EnumProp " + propName + " was nil")
			}
			if err != nil {
				return errors.New("EnumProp " + propName + " had unexpected error: " + err.Error())
			}
		}

	}
	return nil
}

func setReaderStatePtr(reader PropertyReader, st State) error {

	statePtr, ok := st.(*state)
	if !ok {
		return errors.New("The provided non-nil State could not be conveted to a state ptr")
	}

	for propName, propType := range reader.Props() {
		switch propType {
		case TypeStack:
			val, err := reader.StackProp(propName)
			if val == nil {
				return errors.New("Stack Prop " + propName + " was nil")
			}
			if err != nil {
				return errors.New("Stack prop " + propName + " had unexpected error: " + err.Error())
			}
			val.setState(statePtr)
		case TypeTimer:
			val, err := reader.TimerProp(propName)
			if val == nil {
				return errors.New("TimerProp " + propName + " was nil")
			}
			if err != nil {
				return errors.New("TimerProp " + propName + " had unexpected error: " + err.Error())
			}
			val.setState(statePtr)
		}
	}
	return nil
}

//copyReader assumes input and output container are the same "shape" (that is,
//outputContainer can have all of input's properties set). It goes through
//each property, copies it if necessary, and outputs on ReadSetConfigurer.
func copyReader(input PropertyReadSetter, outputContainer PropertyReadSetConfigurer) error {

	for propName, propType := range input.Props() {
		switch propType {
		case TypeBool:
			boolVal, err := input.BoolProp(propName)
			if err != nil {
				return errors.New(propName + " did not return a bool as expected: " + err.Error())
			}
			err = outputContainer.SetBoolProp(propName, boolVal)
			if err != nil {
				return errors.New(propName + " could not be set on output: " + err.Error())
			}
		case TypeInt:
			intVal, err := input.IntProp(propName)
			if err != nil {
				return errors.New(propName + " did not return an int as expected: " + err.Error())
			}
			err = outputContainer.SetIntProp(propName, intVal)
			if err != nil {
				return errors.New(propName + " could not be set on output: " + err.Error())
			}
		case TypeString:
			stringVal, err := input.StringProp(propName)
			if err != nil {
				return errors.New(propName + " did not return a string as expected: " + err.Error())
			}
			err = outputContainer.SetStringProp(propName, stringVal)
			if err != nil {
				return errors.New(propName + " could not be set on output: " + err.Error())
			}
		case TypePlayerIndex:
			playerIndexVal, err := input.PlayerIndexProp(propName)
			if err != nil {
				return errors.New(propName + " did not return a player index as expected: " + err.Error())
			}
			err = outputContainer.SetPlayerIndexProp(propName, playerIndexVal)
			if err != nil {
				return errors.New(propName + " could not be set on output: " + err.Error())
			}
		case TypeEnum:
			enumConst, err := input.EnumProp(propName)
			if err != nil {
				return errors.New(propName + " did not return an EnumVal as expected: " + err.Error())
			}
			err = outputContainer.ConfigureMutableEnumProp(propName, enumConst.MutableCopy())
			if err != nil {
				return errors.New(propName + " could not be set on output: " + err.Error())
			}
		case TypeIntSlice:
			intSliceVal, err := input.IntSliceProp(propName)
			if err != nil {
				return errors.New(propName + " did not return an int slice as expected: " + err.Error())
			}
			err = outputContainer.SetIntSliceProp(propName, intSliceVal)
			if err != nil {
				return errors.New(propName + " could not be set on output: " + err.Error())
			}
		case TypeBoolSlice:
			boolSliceVal, err := input.BoolSliceProp(propName)
			if err != nil {
				return errors.New(propName + " did not return an bool slice as expected: " + err.Error())
			}
			err = outputContainer.SetBoolSliceProp(propName, boolSliceVal)
			if err != nil {
				return errors.New(propName + " could not be set on output: " + err.Error())
			}
		case TypeStringSlice:
			stringSliceVal, err := input.StringSliceProp(propName)
			if err != nil {
				return errors.New(propName + " did not return an string slice as expected: " + err.Error())
			}
			err = outputContainer.SetStringSliceProp(propName, stringSliceVal)
			if err != nil {
				return errors.New(propName + " could not be set on output: " + err.Error())
			}
		case TypePlayerIndexSlice:
			playerIndexSliceVal, err := input.PlayerIndexSliceProp(propName)
			if err != nil {
				return errors.New(propName + " did not return a player index slice as expected: " + err.Error())
			}
			err = outputContainer.SetPlayerIndexSliceProp(propName, playerIndexSliceVal)
			if err != nil {
				return errors.New(propName + " could not be set on output: " + err.Error())
			}
		case TypeStack:
			stackVal, err := input.MutableStackProp(propName)
			if err != nil {
				return errors.New(propName + " did not return a stack as expected: " + err.Error())
			}
			err = outputContainer.ConfigureMutableStackProp(propName, stackVal.mutableCopy())
			if err != nil {
				return errors.New(propName + " could not be set on output: " + err.Error())
			}
		case TypeTimer:
			timerVal, err := input.MutableTimerProp(propName)
			if err != nil {
				return errors.New(propName + " did not return a timer as expected: " + err.Error())
			}
			err = outputContainer.ConfigureMutableTimerProp(propName, timerVal.mutableCopy())
			if err != nil {
				return errors.New(propName + " could not be set on output: " + err.Error())
			}
		default:
			return errors.New(propName + " was an unsupported property type: " + strconv.Itoa(int(propType)))
		}
	}

	return nil

}

func unpackStackStructTag(tag string, chest *ComponentChest) (*Deck, int, error) {
	pieces := strings.Split(tag, ",")

	if len(pieces) > 2 {
		return nil, 0, errors.New("There were more fields in the struct tag than expected")
	}

	deckName := strings.TrimSpace(pieces[0])

	deck := chest.Deck(deckName)

	if deck == nil {
		return nil, 0, errors.New("The deck name " + deckName + " was not a valid deck")
	}

	size := 0

	if len(pieces) > 1 {
		var err error
		size, err = strconv.Atoi(strings.TrimSpace(pieces[1]))
		if err != nil {
			return nil, 0, errors.New("The size in the struct tag was not a valid int: " + err.Error())
		}
	}

	return deck, size, nil

}

//structTagForField will use reflection to fetch the named field from the
//object and return the value of its `enum` field. Works even if fieldName is
//in an embedded struct.
func structTagForField(obj interface{}, fieldName string, structTag string) string {

	v := reflect.Indirect(reflect.ValueOf(obj))

	t := reflect.TypeOf(v.Interface())

	field, ok := t.FieldByNameFunc(func(str string) bool {
		return str == fieldName
	})

	if !ok {
		return ""
	}

	return field.Tag.Get(structTag)

}
