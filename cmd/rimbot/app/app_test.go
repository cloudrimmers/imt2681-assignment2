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
	expected = "{\"text\":\"" + slackUserError + "\",\"username\":\"" + BotDefaultName + "\"}"
	if string(output) != expected {

		t.Error("Default constructor did not return correct error message. \nGot:\t\t", string(output), "\nExpected:\t", expected)
	}

	input = "Test string 123"
	output = MessageSlack(input)
	expected = "{\"text\":\"" + input + "\",\"username\":\"" + BotDefaultName + "\"}"
	if string(output) != expected {

		t.Error("Default constructor did not return correct error message. \nGot:\t\t", string(output), "\nExpected:\t", expected)
	}

	input = " "
	output = MessageSlack(input)
	expected = "{\"text\":\"" + input + "\",\"username\":\"" + BotDefaultName + "\"}"
	if string(output) != expected {

		t.Error("Default constructor did not return correct error message. \nGot:\t\t", string(output), "\nExpected:\t", expected)
	}

}
