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
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS high_scores (name TEXT, score INTEGER, PRIMARY KEY (name, score))")
	if err != nil {
		log.Fatal(err)
	}
	return func() { db.Close() }
}

func getTopRatedHighScores(limit int) []ratedHighScore {
	rows, err := db.Query("SELECT name, score FROM high_scores ORDER BY score DESC LIMIT ?", limit)
	if err != nil {
		log.Panicln("error querying high scores:", err)
	}
	defer rows.Close()

	highScores := make([]ratedHighScore, 0, limit)
	for i := 1; rows.Next(); i++ {
		var hs ratedHighScore
		hs.Rating = i
		err := rows.Scan(&hs.Name, &hs.Score)
		if err != nil {
			log.Panicln("error scanning high score:", err)
		}
		highScores = append(highScores, hs)
	}

	return highScores
}

func getTopRatedHighScoresWithPlayerIncluded(limit int, includedPlayerName string) ([]ratedHighScore, error) {
	topRatedHighScores := getTopRatedHighScores(limit)

	playerToIncludeIsAlreadyThere := false
	for i := 0; i < len(topRatedHighScores); i++ {
		if topRatedHighScores[i].Name == includedPlayerName {
			playerToIncludeIsAlreadyThere = true
		}
	}

	if !playerToIncludeIsAlreadyThere {
		includedPlayerHighScore, err := getRatedHighScoreOfPlayer(includedPlayerName)
		if err != nil {
			return nil, err
		}
		if includedPlayerHighScore != nil {
			topRatedHighScores[limit-1] = *includedPlayerHighScore
		}
	}

	return topRatedHighScores, nil
}

func getRatedHighScoreOfPlayer(playerName string) (*ratedHighScore, error) {
	var hs ratedHighScore
	hs.Name = playerName

	row := db.QueryRow("SELECT name, score, rating FROM (SELECT name, score, RANK() OVER (ORDER BY score DESC) rating FROM high_scores) WHERE name = ?", playerName)
	err := row.Scan(&hs.Name, &hs.Score, &hs.Rating)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &hs, nil
}

func recordHighScore(hs highScore) {
	removeWorseHighScoresOfPlayer(hs)
	_, err := db.Exec("INSERT INTO high_scores (name, score) VALUES (?, ?)", hs.Name, hs.Score)
	if err != nil {
		log.Panicln("error creating high score:", err)
	}
}

func removeWorseHighScoresOfPlayer(hs highScore) {
	_, err := db.Exec("DELETE FROM high_scores WHERE name = ? AND score <= ?", hs.Name, hs.Score)
	if err != nil {
		log.Panicf("error removing worse high scores of player %s: %s\n", hs.Name, err)
	}
}
