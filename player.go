package main

import ()

type Player struct {
	Client     *Client
	Playername string
	Playerhash string
}

var Players []*Player

func CreatePlayer(client *Client, name string) *Player {
	pl := &Player{
		Client:     client,
		Playername: name,
		Playerhash: RandStringBytesRmndr(),
	}

	Players = append(Players, pl)

	return pl
}

func DeletePlayer(player *Player) bool {
	for i := range Players {
		if player == Players[i] {
			Players[i] = Players[len(Players)-1]
			Players[len(Players)-1] = nil
			Players = Players[:len(Players)-1]
			return true
		}
	}
	return false
}
