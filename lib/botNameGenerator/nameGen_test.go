package botNameGenerator

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {

	name := Generate()

	parts := strings.Split(name, " ")

	if !(In(firstNames, parts[0]) && parts[1] == "the" &&
		In(adjectives, parts[2]) && In(animalNames, parts[3])) {

		t.Error("Bot name generator failed to generate. Got:", name, "\nExpected a name consisting of a first name, 'the', and adjective, an animal name.")
	}
}

// In checks if an element is in a string array.
func In(array map[int]string, element string) bool {

	for i := 0; i < len(array); i++ {
		if array[i] == element {

			return true
		}
	}
	return false
}
