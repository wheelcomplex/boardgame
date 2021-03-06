package boardgame

import (
	"reflect"
	"sort"
	"testing"
)

func TestComponentChest(t *testing.T) {

	chest := NewComponentChest(nil)

	if chest.DeckNames() != nil {
		t.Error("We got a deck names array before we'd added anything")
	}

	deckOne := NewDeck()

	componentOne := &testingComponent{
		"foo",
		1,
	}

	deckOne.AddComponent(componentOne)

	componentTwo := &testingComponent{
		"bar",
		2,
	}

	deckOne.AddComponent(componentTwo)

	if deckOne.Components() != nil {
		t.Error("We got non-nil components before it was added to the chest")
	}

	chest.AddDeck("test", deckOne)

	componentValues := make([]Reader, 2)

	for i, component := range deckOne.Components() {
		componentValues[i] = component.Values
	}

	if !reflect.DeepEqual(componentValues, []Reader{componentOne, componentTwo}) {
		t.Error("Deck gave back wrong items after being added to chest")
	}

	deckOne.AddComponent(&testingComponent{
		"illegal",
		-1,
	})

	componentValues = make([]Reader, 2)

	for i, component := range deckOne.Components() {
		componentValues[i] = component.Values
	}

	if !reflect.DeepEqual(componentValues, []Reader{componentOne, componentTwo}) {
		t.Error("Deck allowed itself to be mutated after it was added to chest")
	}

	if chest.DeckNames() != nil {
		t.Error("We got decknames before we called freeze")
	}

	if chest.Deck("test") != nil {
		t.Error("We got a deck back before freeze was called")
	}

	deckTwo := NewDeck()

	deckTwo.AddComponent(&testingComponent{
		"another",
		3,
	})

	chest.AddDeck("other", deckTwo)

	chest.Finish()

	chest.AddDeck("shouldfail", deckOne)

	if chest.decks["shouldfail"] != nil {
		t.Fatal("We were able to add a deck after freezing")
	}

	sortedDeckNames := chest.DeckNames()

	sort.Strings(sortedDeckNames)

	expectedDeckNames := []string{"other", "test"}

	if !reflect.DeepEqual(sortedDeckNames, expectedDeckNames) {
		t.Error("Got unexpected decknames. got", sortedDeckNames, "wanted", expectedDeckNames)
	}

	if chest.Deck("test") != deckOne {
		t.Error("Got wrong value for deck one. Got", chest.Deck("test"), "wanted", deckOne)
	}

	if chest.Deck("other") != deckTwo {
		t.Error("Got wrong value for deck two. Got", chest.Deck("other"), "wanted", deckTwo)
	}

	if deckOne.name != "test" {
		t.Error("DeckOne didn't have its name set when added to the chest. Got", deckOne.name, "wanted test")
	}

	if deckTwo.name != "other" {
		t.Error("DeckTwo didn't have its name set when added to the chest. Got", deckTwo.name, "wanted other")
	}

	for i, c := range deckOne.Components() {
		if c.Deck != deckOne {
			t.Error("At position", i, "deck name was not set correctly in component. Got", c.Deck, "wanted", deckOne)
		}
		if c.DeckIndex != i {
			t.Error("At position", i, "index was not set correctly in component. Got", c.DeckIndex, "wanted", i)
		}
	}

}
