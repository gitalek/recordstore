package main

import (
	"github.com/gomodule/redigo/redis"
	"log"
)

func FindAlbum(id string) (*Album, error) {
	conn := pool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", "album:"+id))
	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return nil, ErrNoAlbum
	}

	album := new(Album)
	err = redis.ScanStruct(values, album)
	if err != nil {
		return nil, err
	}
	return album, nil
}

func IncrementLikes(id string) error {
	conn := pool.Get()
	defer conn.Close()

	exists, err := redis.Int(conn.Do("EXISTS", "album:"+id))
	if err != nil {
		return err
	} else if exists == 0 {
		return ErrNoAlbum
	}

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}
	err = conn.Send("HINCRBY", "album:"+id, "likes", 1)
	if err != nil {
		return err
	}
	err = conn.Send("ZINCRBY", "likes", 1, id)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

func FindTopThree() ([]*Album, error) {
	conn := pool.Get()
	defer conn.Close()

	for {
		_, err := conn.Do("WATCH", "likes")
		if err != nil {
			return nil, err
		}

		ids, err := redis.Strings(conn.Do("ZREVRANGE", "likes", 0, 2))
		if err != nil {
			return nil, err
		}

		err = conn.Send("MULTI")
		if err != nil {
			return nil, err
		}

		for _, id := range ids {
			err := conn.Send("HGETALL", "album:"+id)
			if err != nil {
				return nil, err
			}
		}

		replies, err := redis.Values(conn.Do("EXEC"))
		if err == redis.ErrNil {
			log.Println("trying again")
			continue
		} else if err != nil {
			return nil, err
		}

		albums := make([]*Album, 3)

		for i, reply := range replies {
			album := new(Album)
			err = redis.ScanStruct(reply.([]interface{}), album)
			if err != nil {
				return nil, err
			}
			albums[i] = album
		}

		return albums, nil
	}
}
