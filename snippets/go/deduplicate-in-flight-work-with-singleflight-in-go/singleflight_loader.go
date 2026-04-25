package main

import (
	"errors"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type InFlightStringLoader struct {
	group singleflight.Group
	load  func(string) (string, error)
}

func NewInFlightStringLoader(load func(string) (string, error)) (*InFlightStringLoader, error) {
	if load == nil {
		return nil, errors.New("load function is required")
	}

	return &InFlightStringLoader{load: load}, nil
}

// singleflight shares only concurrent loads for the same key and reruns the work after the in-flight call finishes.
func (loader *InFlightStringLoader) Load(key string) (string, bool, error) {
	if loader == nil {
		return "", false, errors.New("loader is required")
	}
	if key == "" {
		return "", false, errors.New("key is required")
	}

	value, err, shared := loader.group.Do(key, func() (interface{}, error) {
		return loader.load(key)
	})
	if err != nil {
		return "", shared, err
	}

	loadedValue, ok := value.(string)
	if !ok {
		return "", shared, errors.New("singleflight returned a non-string value")
	}

	return loadedValue, shared, nil
}

func main() {
	loader, err := NewInFlightStringLoader(func(key string) (string, error) {
		time.Sleep(50 * time.Millisecond)
		return "value-for-" + key, nil
	})
	if err != nil {
		log.Fatal(err)
	}

	var waitGroup sync.WaitGroup
	for range 2 {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			value, shared, loadErr := loader.Load("profile:42")
			if loadErr != nil {
				log.Printf("load error: %v", loadErr)
				return
			}
			log.Printf("value=%s shared=%t", value, shared)
		}()
	}
	waitGroup.Wait()
}
