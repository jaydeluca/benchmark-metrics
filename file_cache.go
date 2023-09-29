package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type SingleFileCache struct {
	location string
}

func NewSingleFileCache(location string) *SingleFileCache {
	if _, err := os.Stat(location); os.IsNotExist(err) {
		// Create an empty JSON file if it doesn't exist
		data := []byte("{}")
		err := ioutil.WriteFile(location, data, 0644)
		if err != nil {
			panic(err)
		}
	}

	return &SingleFileCache{
		location: location,
	}
}

func (c *SingleFileCache) AddToCache(key string, value string) error {
	// Read the existing cache
	data, err := ioutil.ReadFile(c.location)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data into a map
	cache := map[string]string{}
	err = json.Unmarshal(data, &cache)
	if err != nil {
		return err
	}

	// Add or update the key-value pair
	cache[key] = value

	// Marshal the updated cache and write it back to the file
	updatedData, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.location, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (c *SingleFileCache) RetrieveValue(key string) (string, error) {
	// Read the existing cache
	data, err := ioutil.ReadFile(c.location)
	if err != nil {
		return "", err
	}

	// Unmarshal the JSON data into a map
	cache := map[string]string{}
	err = json.Unmarshal(data, &cache)
	if err != nil {
		return "", err
	}

	return cache[key], nil
}

func (c *SingleFileCache) DeleteCache() error {
	if _, err := os.Stat(c.location); os.IsNotExist(err) {
		return nil // Cache file doesn't exist, nothing to delete
	}

	err := os.Remove(c.location)
	if err != nil {
		return err
	}

	return nil
}
