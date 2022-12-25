package main

import (
	"database/sql"
	"log"
)

var db *sql.DB

func initDb() func() {
	var err error
	db, err = sql.Open("sqlite3", "highscores.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS high_scores (id INTEGER PRIMARY KEY, name TEXT, score INTEGER)")
	if err != nil {
		log.Fatal(err)
	}
	return func() { db.Close() }
}

func getHighScores() []highScore {
	rows, err := db.Query("SELECT name, score FROM high_scores ORDER BY score DESC")
	if err != nil {
		log.Panicln("error querying high scores:", err)
	}
	defer rows.Close()
	highScores := make([]highScore, 0)
	for rows.Next() {
		hs := highScore{}
		err := rows.Scan(&hs.Name, &hs.Score)
		if err != nil {
			log.Panicln("error scanning high score:", err)
		}
		highScores = append(highScores, hs)
	}
	return highScores
}

func createHighScore(hs highScore) {
	_, err := db.Exec("INSERT INTO high_scores (name, score) VALUES (?, ?)", hs.Name, hs.Score)
	if err != nil {
		log.Panicln("error creating high score:", err)
	}
}
