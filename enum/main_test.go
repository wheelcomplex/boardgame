package enum

import (
	"github.com/workfit/tester/assert"
	"sort"
	"testing"
)

func TestEnum(t *testing.T) {
	enums := NewSet()

	assert.For(t).ThatActual(len(enums.EnumNames())).Equals(0)

	const (
		ColorBlue Constant = iota
		ColorGreen
		ColorRed
	)

	const (
		CardSpade Constant = ColorRed + 1 + iota
		CardClub
		CardDiamond
		CardHeart
	)

	const ConstDuplicate = iota

	assert.For(t).ThatActual(enums).IsNotNil()

	colorEnum, err := enums.Add("Color", map[Constant]string{
		ColorBlue:  "Blue",
		ColorGreen: "Green",
		ColorRed:   "Red",
	})

	assert.For(t).ThatActual(colorEnum).IsNotNil()

	assert.For(t).ThatActual(err).IsNil()

	assert.For(t).ThatActual(len(enums.EnumNames())).Equals(1)

	assert.For(t).ThatActual(enums.Enum("Color")).Equals(colorEnum)

	assert.For(t).ThatActual(enums.Membership(ColorBlue)).Equals(colorEnum)

	assert.For(t).ThatActual(colorEnum.DefaultValue()).Equals(ColorBlue)

	assert.For(t).ThatActual(colorEnum.String(ColorBlue)).Equals("Blue")

	assert.For(t).ThatActual(colorEnum.String(125)).Equals("")

	_, err = enums.Add("Color", map[Constant]string{
		ColorBlue: "Blue",
	})

	assert.For(t).ThatActual(err).IsNotNil()

	cardEnum, err := enums.Add("Card", map[Constant]string{
		CardSpade:   "Spade",
		CardClub:    "Club",
		CardDiamond: "Diamond",
		CardHeart:   "Heart",
	})

	assert.For(t).ThatActual(cardEnum).IsNotNil()

	assert.For(t).ThatActual(err).IsNil()

	assert.For(t).ThatActual(enums.Membership(CardDiamond)).Equals(cardEnum)

	_, err = enums.Add("Another", map[Constant]string{
		ConstDuplicate: "Duplicate",
	})

	assert.For(t).ThatActual(err).IsNotNil()

	val := colorEnum.ValueFromString("Blue")

	assert.For(t).ThatActual(val).Equals(ColorBlue)

	val = colorEnum.ValueFromString("Turquoise")

	assert.For(t).ThatActual(val).Equals(InvalidConstant)

	eVal := colorEnum.NewEnumValue()

	assert.For(t).ThatActual(eVal.Value()).Equals(ColorBlue)

	err = eVal.SetValue(ColorGreen)

	assert.For(t).ThatActual(err).IsNil()

	assert.For(t).ThatActual(eVal.Value()).Equals(ColorGreen)

	eVal.Lock()

	err = eVal.SetValue(ColorRed)

	assert.For(t).ThatActual(err).IsNotNil()

	assert.For(t).ThatActual(eVal.Value()).Equals(ColorGreen)

	//Do a new manager to check that adding enums after finished doesn't work

	enums = NewSet()

	_, err = enums.Add("Color", map[Constant]string{
		ColorBlue:  "Blue",
		ColorGreen: "Green",
		ColorRed:   "Blue",
	})

	assert.For(t).ThatActual(err).IsNotNil()

	enums.Finish()

	_, err = enums.Add("Card", map[Constant]string{
		CardSpade: "Spade",
		CardClub:  "Club",
	})

	assert.For(t).ThatActual(err).IsNotNil()
}

func TestCombinedEnumSets(t *testing.T) {

	firstSet := NewSet()
	secondSet := NewSet()

	const (
		ColorBlue Constant = iota
		ColorGreen
		ColorRed
	)

	const (
		CardSpade Constant = ColorRed + 1 + iota
		CardClub
		CardDiamond
		CardHeart
	)

	colorEnum, err := firstSet.Add("Color", map[Constant]string{
		ColorBlue:  "Blue",
		ColorGreen: "Green",
		ColorRed:   "Red",
	})

	cardEnum, err := secondSet.Add("Card", map[Constant]string{
		CardSpade:   "Spade",
		CardClub:    "Club",
		CardDiamond: "Diamond",
		CardHeart:   "Heart",
	})

	combinedSet, err := CombineSets(firstSet, secondSet)

	assert.For(t).ThatActual(err).IsNil()

	enumNames := combinedSet.EnumNames()

	sort.Strings(enumNames)

	assert.For(t).ThatActual(enumNames).Equals([]string{"Card", "Color"})

	assert.For(t).ThatActual(combinedSet.Enum("Color")).Equals(colorEnum)
	assert.For(t).ThatActual(combinedSet.Enum("Card")).Equals(cardEnum)
}
