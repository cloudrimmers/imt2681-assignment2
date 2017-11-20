package botNameGenerator

import (
	"math/rand"
	"time"
)

var firstNames = [...]string{
	"Sarah", "Dave", "Joe", "Zoe", "Peter", "Liam", "Fred", "Lisa", "Martha",
	"Eric", "Christopher", "Mariuz", "Simon", "Hans", "John", "Linda", "Bob",
	"Larry", "Cory", "Sam", "Lola", "Mary", "Dory", "Siri", "Karen", "Emiliy"}

var adjectives = [...]string{
	"happy", "tall", "sassy", "giggely", "jolly", "cheerful", "relaxed",
	"wishful", "festive", "practical", "goofy", "lucky", "daft", "pretty",
	"punky", "proper", "cheesy"}
var animalNames = [...]string{
	"lion", "zebra", "lizzard", "hamster", "parrot", "bear", "penguin",
	"horse", "monkey", "donkey", "owl", "dragon", "shark", "musox", "sheep",
	"crocodile", "giraffe"}

func Generate() (name string) {
	rand.Seed(int64(time.Now().Nanosecond()))

	name = firstNames[rand.Intn(len(firstNames))]
	name += " the "
	name += adjectives[rand.Intn(len(adjectives))]
	name += " "
	name += animalNames[rand.Intn(len(animalNames))]

	return name
}
