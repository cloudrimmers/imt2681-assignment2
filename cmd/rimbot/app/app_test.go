package app

import (
	"testing"
)

func TestMessageSlack(t *testing.T) {

	var input string
	var output []byte
	var expected string

	input = ""
	output = MessageSlack(input)
	expected = "{\"text\":\"" + slackUserError + "\",\"username\":\"Rimbot\"}"

	if string(output) != expected {

		t.Error("Default constructor did not return correct error message. \nGot:\t\t", string(output), "\nExpected:\t", expected)
	}

}
