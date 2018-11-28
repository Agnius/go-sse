package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	// TIMEOUT in seconds
	TIMEOUT      = 30
	DEBUG   bool = true
)

// RequestModel is a simple Request model for better abstraction
// It receives data from Request and passes between code
type RequestModel struct {
	lastEventId string
	topic       string
	message     string
}

type SSE struct {
	channels       map[string]*Channel
	incomingClient chan *RequestModel
	closingClient  chan *Client
}

func createSSE() (sse *SSE) {
	sse = &SSE{
		channels:       make(map[string]*Channel),
		incomingClient: make(chan *RequestModel),
		closingClient:  make(chan *Client),
	}

	go sse.dispatch()

	return
}

func main() {
	sse := createSSE()
	log.Fatal(http.ListenAndServe(":8000", sse))
}

func (sse *SSE) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Flusher is not supported", http.StatusInternalServerError)
		return
	}

	topic := req.URL.Path[len("/infocenter/"):]

	// In every case we should receive topic
	if topic == "" {
		http.Error(rw, "Topic is missing", http.StatusBadRequest)
		return
	}

	if req.Method == "POST" {
		b, _ := ioutil.ReadAll(req.Body)

		go sse.addToChannel(&RequestModel{
			topic:   topic,
			message: string(b),
		})

		rw.WriteHeader(204)
	} else if req.Method == "GET" {
		// Since we can listen non existing channel, we have to create new one
		channel := sse.createChannelIfNotExist(topic)
		// Every client channel receives messages
		clientChan := make(chan *Message)
		client := CreateClient(topic, clientChan)
		channel.clients[clientChan] = true

		if DEBUG {
			fmt.Printf("New client has been registered in %s channel. Total: %d\n\n", channel.name, len(channel.clients))
		}

		rw.Header().Set("Content-Type", "text/event-stream")
		rw.Header().Set("Cache-Control", "no-cache")
		rw.Header().Set("Access-Control-Allow-Origin", "*")

		notify := rw.(http.CloseNotifier).CloseNotify()

		go func() {
			<-notify
			sse.closingClient <- client
		}()

		// After we started SSE connection between server and client
		// We have to set timeout for client
		time := time.NewTimer(time.Second * TIMEOUT)

		for {
			select {
			case msg := <-clientChan:
				if msg == nil {
					return
				}

				fmt.Fprintf(rw, "id: %d\n", msg.id)
				fmt.Fprintf(rw, "event: %s\n", msg.event)
				fmt.Fprintf(rw, "data: %s\n\n", strings.Replace(string(msg.data), "\n", "\ndata: ", -1))
				flusher.Flush()
			case <-time.C:
				fmt.Fprintf(rw, "event: %s\n", "timeout")
				fmt.Fprintf(rw, "data: %ds\n\n", TIMEOUT)
				flusher.Flush()

				client.Close(sse.channels[client.channel].clients)
				return
			}
		}
	}
}

func (sse *SSE) addToChannel(m *RequestModel) {
	sse.createChannelIfNotExist(m.topic)

	sse.incomingClient <- m
}

func (sse *SSE) createChannelIfNotExist(chanName string) *Channel {
	if !sse.doesChannelExist(chanName) {
		sse.channels[chanName] = CreateChannel(chanName)
	}

	return sse.channels[chanName]
}

func (sse *SSE) dispatch() {
	for {
		select {
		case s := <-sse.incomingClient:
			message := NewMessage("msg", s.message, s.topic)
			go addMessage(sse.channels[s.topic], message)
		case c := <-sse.closingClient:
			c.Close(sse.channels[c.Channel()].clients)

			if DEBUG {
				fmt.Printf("Removed client from %s channel. Clients left: %d", c.Channel(), sse.channels[c.Channel()].ClientsCount())
			}
			return
		}
	}
}

func addMessage(channel *Channel, message *Message) {
	// Add message for every client in this channel
	for client := range channel.clients {
		client <- message
	}
}

func (s *SSE) doesChannelExist(name string) bool {
	_, ok := s.channels[name]

	return ok
}
