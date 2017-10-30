package tool

import (
	"log"
	"time"
)

// UntilTomorrow ...
func UntilTomorrow() time.Duration {
	// @doc https://stackoverflow.com/a/36988882
	now := time.Now()
	tomorrow := now.Add(time.Hour * 24)
	tomorrow = time.Date(
		tomorrow.Year(),
		tomorrow.Month(),
		tomorrow.Day(),
		0, 0, 0, 0,
		tomorrow.Location()) // Round to 00:00:00

	diff := tomorrow.Sub(now)

	// @debug
	log.Println("Tommorrow :", tomorrow)
	log.Println("Today     :", now)
	log.Println("Diff      :", diff)
	return diff
}
