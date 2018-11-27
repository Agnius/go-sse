package main

type Client struct {
	channel        string
	messageChannel chan *Message
}

func CreateClient(channel string, messageChan chan *Message) *Client {
	return &Client{
		channel:        channel,
		messageChannel: messageChan,
	}
}

func (c *Client) Channel() string {
	return c.channel
}

func (c *Client) Close(clients map[chan *Message]bool) {
	delete(clients, c.messageChannel)
	close(c.messageChannel)
}
