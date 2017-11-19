package botNameGenerator

import (
	"math/rand"
	"time"
)

var firstNames = [...]string{"Sarah", "Dave", "Joe", "Zoe", "Peter", "Liam", "Fred", "Lisa", "martha", "Eric", "Christopher",
	"Mariuz", "Simon", "Hans", "John", "Linda"}

var adjectives = [...]string{"happy", "tall", "sassy", "giggely", "jolly", "cheerful", "relaxed", "wishful", "festive", "practical"}
var animalNames = [...]string{"lion", "zebra", "lizzard", "hamster", "parrot", "bear", "penguin"}

func Generate() (name string) {
	rand.Seed(int64(time.Now().Nanosecond()))

	name = firstNames[rand.Intn(len(firstNames))]
	name += " the "
	name += adjectives[rand.Intn(len(adjectives))]
	name += " "
	name += animalNames[rand.Intn(len(animalNames))]

	return name
}
