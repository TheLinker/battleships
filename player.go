package main

import (
    "fmt"
)

type Player struct {
    Client *Client
    Playername string
    Playerhash string
}

var Players []*Player

func CreatePlayer(client *Client, name string) *Player {
    pl := &Player {
        Client: client,
        Playername: name,
        Playerhash: RandStringBytesRmndr(),
    }

    Players = append(Players, pl)

    return pl
}

func DeletePlayer(player *Player) {
    for i := range Players {
        fmt.Println("%+v", i)
        if player == Players[i] {
            Players[i] = Players[len(Players)-1]
            Players[len(Players)-1] = nil
            Players = Players[:len(Players)-1]
        }
    }
}