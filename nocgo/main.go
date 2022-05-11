package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

var perfTune = `
pragma journal_mode = WAL;
pragma synchronous = normal;
pragma temp_store = memory;
pragma mmap_size = 30000000000;`

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}

	words, err := os.ReadFile("/usr/share/dict/words")
	if err != nil {
		panic(err)
	}

	rand.Seed(0)
	times := 10

	for i := 0; i < times; i++ {
		db.Exec("DROP TABLE IF EXISTS people")
		db.Exec(`
CREATE TABLE people (
  name TEXT,
  country TEXT,
  region TEXT,
  occupation TEXT,
  age INT,
  company TEXT,
  favorite_team TEXT,
  favorite_sport TEXT
)`)

		stmt, err := db.Prepare("INSERT INTO people VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			panic(err)
		}

		t1 := time.Now()
		for i := 0; i < len(words)*10; i++ {
			rnd_name := words[rand.Int()%len(words)]
			rnd_country := words[rand.Int()%len(words)]
			rnd_region := words[rand.Int()%len(words)]
			rnd_occupation := words[rand.Int()%len(words)]
			rnd_company := words[rand.Int()%len(words)]
			rnd_age := rand.Int() % 110
			rnd_fav_team := words[rand.Int()%len(words)]
			rnd_fav_sport := words[rand.Int()%len(words)]

			_, err := stmt.Exec(rnd_name, rnd_country, rnd_region, rnd_occupation, rnd_age, rnd_company, rnd_fav_team, rnd_fav_sport)
			if err != nil {
				panic(err)
			}

			if i % 100_000 == 0 {
				fmt.Printf("%d%%\n", int((float64(i) / float64(len(words)*10)) * 100))
			}
		}
		fmt.Println(time.Now().Sub(t1))
	}
}
