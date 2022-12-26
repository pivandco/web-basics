package main

type highScore struct {
	Id    int    `json:"-"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type ratedHighScore struct {
	highScore
	Rating int `json:"rating"`
}
