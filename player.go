package main

import (
    "errors"
)

type Player struct {
	Client     *Client
	Playername string
	Playerhash string
}

var Players []*Player

func CreatePlayer(client *Client, name string) (err error, pl *Player) {
    for _, pl = range(Players) {
        if pl.Playername == name && pl.Client != nil {
            return errors.New("Nickname already taken"), nil
        } else if pl.Playername == name && pl.Client == nil {
            pl.Client = client
            return
        }
    }

    pl = &Player {
        Client: client,
        Playername: name,
        Playerhash: RandStringBytesRmndr(),
    }

    Players = append(Players, pl)

    return
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
