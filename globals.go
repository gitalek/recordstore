package main

import (
	"errors"
	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

var ErrNoAlbum = errors.New("no album found")

// Album struct holds album data
type Album struct {
	Title  string  `redis:"title"`
	Artist string  `redis:"artist"`
	Price  float64 `redis:"price"`
	Likes  int     `redis:"likes"`
}
