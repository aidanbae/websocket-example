package main

import (
	"time"
	"fmt"
)

type World struct {
	clientMap map[*Client]bool
	ChanEnter chan *Client
	ChanLeave chan *Client
	broadcast chan []byte
}

func newWorld () *World {
	return &World{
		clientMap: make(map[*Client]bool, 5),
		broadcast: make(chan []byte),
	}
}

func (w *World) run() {
	w.ChanEnter = make(chan *Client)
	w.ChanLeave = make(chan *Client)

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case client := <- w.ChanEnter:
			fmt.Println("클라이언트 입장")
			w.clientMap[client] = true
		case client := <- w.ChanLeave:
			if _, ok := w.clientMap[client]; ok {
				delete(w.clientMap, client)
				fmt.Println("클라이언트 퇴장")
				close(client.send)
			}
		case message := <- w.broadcast:
			for client := range w.clientMap {
				client.send <- message
			}
		case tick := <- ticker.C:
			for client := range w.clientMap {
				client.send <- []byte(tick.String())
			}
			fmt.Println(tick.String())
		}
	}
}

