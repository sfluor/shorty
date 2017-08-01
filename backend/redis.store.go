package main

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

const (
	// Size of the tag that is supposed to be unique
	tagSize = 5

	// Tag never expires
	tagExpiration = 0
)

// Error that redis returns when it's doesn't find a value for a given key
var notFound = errors.New("redis: nil")

func shorten(client *redis.Client, url string) string {
	stored, val := isAlreadyStored(client, url)
	// If it's already stored just return the tag
	if stored {
		return val
	} else {
		// Otherwise just process normally generating a new tag
		tag := randIdGen(tagSize)
		// Set tag: value
		err := client.Set(tag, url, tagExpiration).Err()
		if err != nil {
			logrus.Errorf("Error during storing to Redis backend: %s", err)
		}
		// But also value: tag to check if url is already stored
		err = client.Set(url, tag, tagExpiration).Err()
		if err != nil {
			logrus.Errorf("Error during storing to Redis backend: %s", err)
		}
		return tag
	}
}

func unShorten(client *redis.Client, tag string) string {
	val, err := client.Get(tag).Result()
	if err != nil {
		logrus.Errorf("Error during unshortening: %s", err)
		return ""
	} else {
		logrus.Infof("Unshortening: %v, gives : %s", tag, val)
		return val
	}
}

// Check if an url has already been shortened
func isAlreadyStored(client *redis.Client, url string) (bool, string) {
	val, err := client.Get(url).Result()
	// If an error occurs and val is not empty (meaning it does exist)
	if err != nil && val != "" {
		logrus.Errorf("An error occured when querying redis store: %s", err)
	}
	// Value is already stored
	if val != "" {
		return true, val
	}
	// Else it's not
	return false, val
}
