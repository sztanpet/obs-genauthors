package main

import (
	"math/rand"
	"time"
)

type Config struct {
	Addr     string `toml:"listenaddr"`
	Secret   string `toml:"secret"`
	RepoPath string `toml:"repopath"`
}

var sampleconf = `listenaddr=":80"
repopath="."
secret="`

func init() {
	// from http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
	rand.Seed(time.Now().UnixNano())

	runes := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 32)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}

	sampleconf += string(b) + "\"\n"
}
