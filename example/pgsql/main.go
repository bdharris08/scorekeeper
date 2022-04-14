package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/bdharris08/scorekeeper"
	"github.com/bdharris08/scorekeeper/store"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	numWorkers = 5
)

//TODO use cohorts to fix averages, currently just grabbing all scores for each run
// memorystore handled that by being destroyed every run

var dsn = flag.String("dsn", "postgres://postgres:xxx@localhost:5432/postgres", "dsn for postgres database")

func main() {
	flag.Parse()

	db, err := sql.Open("pgx", *dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	scoreKeeper, err := scorekeeper.New(store.NewSQLStore(db, "", ""))
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
