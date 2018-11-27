package main

import (
	"reflect"
	"testing"
)

const CHANNEL_NAME string = "Test Channel"

func TestCreateChannel(t *testing.T) {

	fakeChannel := &Channel{
		name:    CHANNEL_NAME,
		clients: make(map[chan *Message]bool),
	}

	channel := CreateChannel(CHANNEL_NAME)

	if !reflect.DeepEqual(fakeChannel, channel) {
		t.Error("Failed to create channel")
	}
}

func TestName(t *testing.T) {
	channel := CreateChannel("CHANNEL_NAME")

	if channel.Name() != CHANNEL_NAME {
		t.Errorf("Test failed name property is not equel to %s", CHANNEL_NAME)
	}
}

func TestClientsCount(t *testing.T) {
	const INTERATIONS int = 10

	channel := CreateChannel(CHANNEL_NAME)

	for i := 0; i < INTERATIONS; i++ {
		messageChannel := make(chan *Message)
		channel.clients[messageChannel] = true
	}

	if channel.ClientsCount() != INTERATIONS {
		t.Errorf("Test failed clients count does not match! Expected: %d Got: %d", INTERATIONS, channel.ClientsCount())
	}
}
