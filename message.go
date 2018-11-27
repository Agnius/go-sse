package main

// It's a global state for all messages passed to the application
// TODO: Because we using channels I doubnt that I need to use
// Mutual exclusion for race conditions, although need to make research on it.
var messagesCount int = 0

type Message struct {
	id      int
	event   string
	data    string
	channel string
}

func NewMessage(event string, data string, channel string) *Message {
	messagesCount++

	return &Message{
		id:      messagesCount,
		event:   event,
		data:    data,
		channel: channel,
	}
}

func MessagesCount() int {
	return messagesCount
}
