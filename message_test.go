package main

import (
	"reflect"
	"testing"
)

const (
	TEST_EVENT   = "erwferg5b1vsd315"
	TEST_DATA    = "X12ds34fv46we4rfg68ewg4"
	TEST_CHANNEL = "897878refv1"
)

func TestNewMessage(t *testing.T) {

	fakeMessage := &Message{
		id:      1,
		event:   TEST_EVENT,
		data:    TEST_DATA,
		channel: TEST_CHANNEL,
	}

	message := NewMessage(TEST_EVENT, TEST_DATA, TEST_CHANNEL)

	if !reflect.DeepEqual(message, fakeMessage) {
		t.Error("Message properties is not euqal to fakeMessage")
	}
}

func TestMessagesCount(t *testing.T) {
	const ITERATIONS = 10

	for range [ITERATIONS]int{} {
		NewMessage(TEST_EVENT, TEST_DATA, TEST_CHANNEL)
	}

	if MessagesCount() != ITERATIONS {
		t.Errorf("Test failed messagesCount does not match! Expected: %d Got: %d", ITERATIONS, MessagesCount())
	}
}
