package botNameGenerator

import "math/rand"

var firstNames = map[int]string{0: "Sarah", 1: "Dave", 2: "Joe", 3: "Zoe", 4: "Peter", 5: "Liam", 6: "Fred", 7: "Lisa"}
var adjectives = map[int]string{0: "happy", 1: "tall", 2: "sassy", 3: "giggely", 4: "jolly", 5: "cheerful", 6: "relaxing", 7: "wishful", 8: "festive"}
var animalNames = map[int]string{0: "lion", 1: "zebra", 2: "lizzard", 3: "hamster", 4: "parrot", 5: "bear", 6: "penguin"}

func Generate() (name string) {

	name = firstNames[rand.Intn(len(firstNames))]
	name += " the "
	name += adjectives[rand.Intn(len(adjectives))]
	name += " "
	name += animalNames[rand.Intn(len(animalNames))]

	return name
}
