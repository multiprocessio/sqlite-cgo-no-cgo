package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var perfTune = `
pragma journal_mode = WAL;
pragma synchronous = normal;
pragma temp_store = memory;
pragma mmap_size = 30000000000;`

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	wordsFile, err := os.ReadFile("/usr/share/dict/words")
	if err != nil {
		panic(err)
	}
	words := strings.Split(string(wordsFile), "\n")

	rand.Seed(0)
	times := 10
	rowTests := []int{10_000, len(words), len(words) * 10}

	_, err = db.Exec(perfTune)
	if err != nil {
		panic(err)
	}

	fmt.Println("time,rows,category,version")

	for _, rows := range rowTests {
		for i := 0; i < times; i++ {
			_, err = db.Exec("DROP TABLE IF EXISTS people")
			if err != nil {
				panic(err)
			}
			_, err = db.Exec(`
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
			if err != nil {
				panic(err)
			}

			stmt, err := db.Prepare("INSERT INTO people VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				panic(err)
			}

			t1 := time.Now()
			for i := 0; i < rows; i++ {
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
			}
			fmt.Printf("%f,%d,insert,cgo\n", float64(time.Now().Sub(t1)) / 1e9, rows)

			t1 = time.Now()
			_, err = db.Query("SELECT COUNT(1), age FROM people GROUP BY age ORDER BY COUNT(1) DESC")
			fmt.Printf("%f,%d,group_by,cgo\n", float64(time.Now().Sub(t1)) / 1e9, rows)
		}
	}
}
