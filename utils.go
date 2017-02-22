package main

import (
	//"crypto/sha256"
	"math/rand"
)

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const n = 32

func RandStringBytesRmndr() string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

type Vector struct {
	X, Y int
}

func (v *Vector) AddVector(v2 *Vector) Vector {
	var ret Vector
	ret.X = v.X + v2.X
	ret.Y = v.Y + v2.Y
	return ret
}
