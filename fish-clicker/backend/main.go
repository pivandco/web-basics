package main

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

const port = ":8080"

func main() {
	closeDb := initDb()
	defer closeDb()
	http.HandleFunc("/high-scores", handle)
	log.Println("Listening on", port)
	log.Fatal(http.ListenAndServe(port, http.HandlerFunc(serveWithRequestLog)))
}

func serveWithRequestLog(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
	http.DefaultServeMux.ServeHTTP(w, r)
}

func handle(w http.ResponseWriter, r *http.Request) {
	defer recoverPanic(w)

	log.Println(r.Method, r.URL)
	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodPost:
		handlePost(w, r)
	case http.MethodDelete:
		handleDelete(w, r)
	default:
		handleUnknown(w, r)
	}
}

func recoverPanic(w http.ResponseWriter) {
	r := recover()
	if r == nil {
		return
	}

	log.Printf("Panic: %s\n", r)
	debug.PrintStack()

	var msg string
	switch x := r.(type) {
	case string:
		msg = x
	case error:
		msg = x.Error()
	default:
		msg = "Unknown error"
	}
	w.WriteHeader(500)
	w.Header().Add("Content-Type", "application/json")
	respBody, err := json.Marshal(map[string]string{"error": msg})
	if err != nil {
		log.Println(err.Error())
		w.Write(([]byte)(`{"error": "Unknown - failed to serialize error message to JSON"}`))
		return
	}
	w.Write(respBody)
}

const topHighScoresAmount = 10

func handleGet(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	playerNameToInclude := query.Get("myName")
	highScores, err := getTopRatedHighScoresWithPlayerIncluded(topHighScoresAmount, playerNameToInclude)
	if err != nil {
		log.Panicf("error getting top rated high scores with player %s included: %s\n", playerNameToInclude, err)
	}

	topHighScores := make([]ratedHighScore, 0, topHighScoresAmount)
	for i := 0; i < topHighScoresAmount && i < len(highScores); i++ {
		topHighScores = append(topHighScores, highScores[i])
	}

	scoresJson, err := json.Marshal(topHighScores)
	if err != nil {
		log.Panicln("error marshaling high scores:", err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(scoresJson)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var hs highScore
	err := json.NewDecoder(r.Body).Decode(&hs)
	if err != nil {
		handleBodyParseError(w, err)
		return
	}
	recordHighScore(hs)
	w.WriteHeader(http.StatusCreated)
}

//go:embed "delete-password.txt"
var deletePassword string

func handleDelete(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		handleBodyParseError(w, err)
		return
	}
	if body.Password != strings.TrimSpace(deletePassword) {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		w.Write(([]byte)(`{"error": "Wrong password"}`))
		return
	}
	_, err = db.Exec("DELETE FROM high_scores")
	if err != nil {
		log.Panicln("error deleting high scores:", err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func handleUnknown(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func handleBodyParseError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}
