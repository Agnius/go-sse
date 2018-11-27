package main

type Channel struct {
	name    string
	clients map[chan *Message]bool
}

func CreateChannel(name string) *Channel {
	return &Channel{
		name:    name,
		clients: make(map[chan *Message]bool),
	}
}

func (c Channel) Name() string {
	return c.name
}

func (c Channel) ClientsCount() int {
	return len(c.clients)
}
