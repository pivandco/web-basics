package main

type highScore struct {
	Id    int    `json:"-"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}
