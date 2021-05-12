package common

import (
	"math/rand"
)

type HealthState int32

const (
	Init  HealthState = 0
	Alive HealthState = 1
	Die   HealthState = 2
)

var HealthState_name = map[int32]string{
	0: "Init",
	1: "Alive",
	2: "Die",
}

func RandStringBytes(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
