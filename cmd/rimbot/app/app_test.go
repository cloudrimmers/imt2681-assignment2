package app

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func TestMessageSlack(t *testing.T) {

	var input string
	var output []byte
	var expected string

	input = ""
	output = MessageSlack(input, false)
	expected = "{\"text\":\"" + slackUserError + "\",\"username\":\"" + BotDefaultName + "\"}"
	if string(output) != expected {

		t.Error("Default constructor did not return correct error message. \nGot:\t\t", string(output), "\nExpected:\t", expected)
	}

	input = "Test string 123"
	output = MessageSlack(input, false)
	expected = "{\"text\":\"" + input + "\",\"username\":\"" + BotDefaultName + "\"}"
	if string(output) != expected {

		t.Error("Default constructor did not return correct error message. \nGot:\t\t", string(output), "\nExpected:\t", expected)
	}

	input = " "
	output = MessageSlack(input, false)
	expected = "{\"text\":\"" + input + "\",\"username\":\"" + BotDefaultName + "\"}"
	if string(output) != expected {

		t.Error("Default constructor did not return correct error message. \nGot:\t\t", string(output), "\nExpected:\t", expected)
	}

}

func TestParseFixerResponse(t *testing.T) {

	var input io.ReadCloser
	var output float64
	var expected float64
	var err error

	input = ioutil.NopCloser(bytes.NewReader([]byte("3.14")))
	output, err = ParseFixerResponse(input)
	expected = 3.14
	if err != nil {

		t.Error("Parsing fixer data failed. Got error: ", err)

	} else if output != expected {

		t.Error("Parsing fixer data failed. Got: ", output, "\tExpected: ", expected)
	}

	input = ioutil.NopCloser(bytes.NewReader([]byte("")))
	output, err = ParseFixerResponse(input)
	expected = 0.0
	if err == nil {

		t.Error("Parsing fixer data succeeded when it should have failed.")
	}

}
