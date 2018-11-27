package main

import (
	"reflect"
	"testing"
)

const TEST_CHANNEL_NAME = "0x80808000"

func TestCreateClient(t *testing.T) {
	messageChan := make(chan *Message)

	fakeClient := &Client{
		channel:        TEST_CHANNEL_NAME,
		messageChannel: messageChan,
	}

	client := CreateClient(TEST_CHANNEL_NAME, messageChan)

	if !reflect.DeepEqual(client, fakeClient) {
		t.Error("Client properties is not euqal to fakeClient")
	}
}

func TestClientChannelName(t *testing.T) {
	client := CreateClient(TEST_CHANNEL_NAME, make(chan *Message))

	if client.Channel() != TEST_CHANNEL_NAME {
		t.Errorf("Channel property is not equal to %s", TEST_CHANNEL_NAME)

	}
}

func TestClose(t *testing.T) {
	sse := &SSE{
		channels:       make(map[string]*Channel),
		incomingClient: make(chan *RequestModel),
		closingClient:  make(chan *Client),
	}

	channel := sse.createChannelIfNotExist(TEST_CHANNEL_NAME)
	clientChan := make(chan *Message)
	client := CreateClient(TEST_CHANNEL_NAME, clientChan)
	channel.clients[clientChan] = true

	client.Close(channel.clients)

	if !isClosed(client.messageChannel) {
		t.Error("Failed to close client message channel!")
	}
}

func isClosed(ch <-chan *Message) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}
