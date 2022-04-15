package main

import (
	"fmt"
	"sync"

	"github.com/bdharris08/scorekeeper"
	"github.com/bdharris08/scorekeeper/store"
)

var (
	numWorkers = 5
)

func main() {
	scoreKeeper, err := scorekeeper.New(&store.MemoryStore{})
	if err != nil {
		panic(fmt.Errorf("error creating scoreKeeper: %v", err))
	}

	scoreKeeper.Start()
	defer scoreKeeper.Stop()

	actions := []string{
		`{"action":"hop", "time":100}`,
		`{"action":"skip", "time":100}`,
		`{"action":"jump", "time":100}`,
		`{"action":"hop", "time":200}`,
		`{"action":"skip", "time":200}`,
		`{"action":"jump", "time":200}`,
		`{"action":"hop", "time":1}`,
		`{"action":"hop", "time":1}`,
		`{"action":"hop", "time":1}`,
		`{"action":"skip", "time":2}`,
		`{"action":"skip", "time":2}`,
		`{"action":"skip", "time":2}`,
		`{"action":"jump", "time":3}`,
		`{"action":"jump", "time":3}`,
		`{"action":"jump", "time":3}`,
	}

	// A simple example
	for _, a := range actions {
		if err := scoreKeeper.AddAction(a); err != nil {
			fmt.Println(err)
		}
	}
	result, err := scoreKeeper.GetStats()
	if err != nil {
		fmt.Println(err)
	}
	// Do something with the result
	fmt.Println(result)

	// You can even access the scoreKeeper concurrently!
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for _, a := range actions {
				if err := scoreKeeper.AddAction(a); err != nil {
					fmt.Println(err)
				}
			}
			result, err := scoreKeeper.GetStats()
			if err != nil {
				fmt.Println(err)
			}
			// Do something with the result,
			// this will likely be an intermediate result
			// since not all workers will have finished
			fmt.Println(result)
		}()
	}

	wg.Wait()

	result, err = scoreKeeper.GetStats()
	if err != nil {
		fmt.Println(err)
	}
	// Do something with the result
	fmt.Println(result)
}
