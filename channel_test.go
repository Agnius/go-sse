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
		t.Error("Channel properties is not euqal to fakeChannel")
	}
}

func TestName(t *testing.T) {
	channel := CreateChannel(CHANNEL_NAME)

	if channel.Name() != CHANNEL_NAME {
		t.Errorf("Name property is not equal to %s", CHANNEL_NAME)
	}
}

func TestClientsCount(t *testing.T) {
	const ITERATIONS int = 10

	channel := CreateChannel(CHANNEL_NAME)

	for i := 0; i < ITERATIONS; i++ {
		messageChannel := make(chan *Message)
		channel.clients[messageChannel] = true
	}

	if channel.ClientsCount() != ITERATIONS {
		t.Errorf("Test failed clients count does not match! Expected: %d Got: %d", ITERATIONS, channel.ClientsCount())
	}
}
