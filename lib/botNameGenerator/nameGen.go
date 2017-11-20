package botNameGenerator

import (
	"math/rand"
	"time"
)

var firstNames = [...]string{
	"Sarah", "Dave", "Joe", "Zoe", "Peter", "Liam", "Fred", "Lisa", "Martha",
	"Eric", "Christopher", "Mariuz", "Simon", "Hans", "John", "Linda", "Bob",
	"Larry", "Cory", "Sam", "Lola", "Mary", "Dory", "Siri", "Karen", "Emiliy",
	"Chuck", "Tyrone", "Jon", "Womble", "Sasafras", "David", "Finn", "Jackson",
	"Conan"}

var adjectives = [...]string{
	"happy", "tall", "sassy", "giggely", "jolly", "cheerful", "relaxed",
	"wishful", "festive", "practical", "goofy", "lucky", "daft", "pretty",
	"punky", "proper", "cheesy", "educated", "marvelous", "friendly", "genius",
	"shifty-looking", "magical", "handy", "adorable", "cool", "heartwarming",
	"weird", "bizzare", "memetastic", "uncomfortable", "savage", "bloody",
	"slimy", "spooktastic", "triggered", "gender-fluid", "fully-armed", "cyborg",
	"one-eyed", "three-eyed", "eagle-eyed", "cloudy", "air-headed", "enthusiastic",
	"full-bearded", "untrustworthy", "spacy", "shady", "lazy-eyed", "greedy",
	"confused", "international", "wise", "fishy", "speedy", "smelly", "exhausted"}

var animalNames = [...]string{
	"lion", "zebra", "lizzard", "hamster", "parrot", "bear", "penguin",
	"horse", "monkey", "donkey", "owl", "dragon", "shark", "musox", "sheep",
	"crocodile", "giraffe", "barbarian", "octompus", "raven", "eagle", "dwarf",
	"fox", "wolf", "cowboy", "unicorn", "dragon", "snake", "zombie", "troll",
	"starfish", "jellyfish", "eel", "bat", "gorilla", "deer", "seal", "squirel",
	"fish", "pirate", "sloth", "otter", "hedgehog", "elephant", "skunk", "rabbit"}

func Generate() (name string) {
	rand.Seed(int64(time.Now().Nanosecond()))

	name = firstNames[rand.Intn(len(firstNames))]
	name += " the "
	name += adjectives[rand.Intn(len(adjectives))]
	name += " "
	name += animalNames[rand.Intn(len(animalNames))]

	return name
}
