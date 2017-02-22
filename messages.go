package main


type Envelope struct {
    Type string
    Msg  interface{}
}

// C -> S
type Registration struct {
    Playername string
}

// S -> C
type RegistrationOK struct {
    Playername string
    Playerhash string
}

// C <-> S
type Chat struct {
    Lobby   string
    Message string
}

