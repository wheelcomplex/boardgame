/************************************
 *
 * This file contains auto-generated methods to help configure enums.
 * It was generated by autoreader.
 *
 * DO NOT EDIT by hand.
 *
 ************************************/

package examplepkg

import (
	"github.com/jkomoros/boardgame/enum"
)

var Enums = enum.NewSet()

var ColorEnum = Enums.MustAdd("Color", map[int]string{
	ColorBlue:    "Blue",
	ColorGreen:   "Green",
	ColorRed:     "Red",
	ColorUnknown: "Unknown",
})

var PhaseEnum = Enums.MustAdd("Phase", map[int]string{
	PhaseMultiWord:    "Multi Word",
	PhaseUnknown:      "Unknown",
	PhaseVeryLongName: "Very Long Name",
})

var FooEnum = Enums.MustAdd("Foo", map[int]string{
	FooBlue:           "Blue",
	FooOverride:       "Green",
	FooOverrideBlank:  "",
	FooOverrideQuoted: "My name is \"Blue\"",
})

var TransformExampleEnum = Enums.MustAdd("TransformExample", map[int]string{
	TransformExampleLowerCase:                 "lower case",
	TransformExampleNormalConfiguredTransform: "Normal Configured Transform",
	TransformExampleNormalTransform:           "Normal Transform",
	TransformExampleUpperCase:                 "UPPER CASE",
})

var DefaultTransformEnum = Enums.MustAdd("DefaultTransform", map[int]string{
	DefaultTransformBlue:  "BLUE",
	DefaultTransformGreen: "GREEN",
	DefaultTransformRed:   "Red",
})
