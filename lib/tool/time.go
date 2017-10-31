package tool

import (
	"fmt"
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
	//log.Println("Tommorrow :", tomorrow)
	//log.Println("Now       :", now)
	//log.Println("Diff      :", diff)
	return diff
}

// Todaystamp ...
func Todaystamp() string {
	now := time.Now()
	return fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())
}

// Daystamp ...
func Daystamp(n int) string {
	now := time.Now().Add(-time.Hour * 24 * time.Duration(n))
	return fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())
}
