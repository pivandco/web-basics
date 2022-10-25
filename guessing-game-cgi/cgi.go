package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/http/cgi"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	err := cgi.Serve(http.HandlerFunc(handler))
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	i, err := readInput(r)
	var out output
	if err == nil {
		out = getOutput(i)
	} else {
		out = getCheatingDetectedOutput()
	}

	_, err = fmt.Fprintln(w, generateGameHtml(out))
	if err != nil {
		log.Println(err)
	}
}

func getCheatingDetectedOutput() output {
	inp := *getGameStartInput()
	return output{
		EncryptedInput:    encryptInput(inp),
		NumAttempts:       inp.NumAttempts,
		Message:           "Обнаружена попытка взлома. Игра начата заново.",
		MinAnswer:         minNum,
		MaxAnswer:         maxNum,
		GameOver:          false,
		PrettyAttemptsLog: "",
	}
}

func getGameStartInput() *input {
	return &input{
		CorrectAnswer: rand.Intn(maxNum-minNum) + minNum,
		NumAttempts:   int(math.Ceil(math.Log2(maxNum - minNum + 1))),
	}
}

//go:embed game.html
var gameHtmlTemplate string

func generateGameHtml(o output) string {
	t := template.Must(
		template.New("game").
			Funcs(template.FuncMap{"JoinAttempts": joinAttempts}).
			Parse(gameHtmlTemplate),
	)
	var html bytes.Buffer
	if err := t.Execute(&html, o); err != nil {
		panic(err)
	}
	return html.String()
}

func joinAttempts(ns []int) string {
	var nsStrs []string
	for _, n := range ns {
		nsStrs = append(nsStrs, strconv.Itoa(n))
	}
	return strings.Join(nsStrs, ";")
}

const (
	minNum = 1
	maxNum = 100
)

func readInput(req *http.Request) (*input, error) {
	if req.Method == "GET" {
		return nil, nil
	}

	if err := req.ParseForm(); err != nil {
		panic(err)
	}

	f := req.PostForm
	encryptedInput := f.Get("gameData")
	answer := f.Get("answer")
	input, err := decryptInput(encryptedInput, answer)
	return &input, err
}

type output struct {
	EncryptedInput       string
	NumAttempts          int
	Message              string
	MinAnswer, MaxAnswer int
	GameOver             bool
	PrettyAttemptsLog    string
}

func getOutput(i *input) output {
	if i == nil {
		i := getGameStartInput()
		return output{
			EncryptedInput:    encryptInput(*i),
			NumAttempts:       i.NumAttempts,
			MinAnswer:         minNum,
			MaxAnswer:         maxNum,
			GameOver:          false,
			PrettyAttemptsLog: "",
		}
	}

	if i.Answer < minNum || i.Answer > maxNum {
		return continueOutput("Введите число от 1 до 100", false, i)
	}

	i = reduceAndLogAttempt(*i)

	if i.Answer == i.CorrectAnswer {
		return continueOutput("Вы угадали! 🎉", true, i)
	}

	if i.NumAttempts == 0 {
		return continueOutput(
			fmt.Sprintf("К сожалению, вы проиграли! Мое число было %d. Удачи в следующий раз!", i.CorrectAnswer),
			true,
			i,
		)
	}

	if i.Answer < i.CorrectAnswer {
		return continueOutput("Не угадали (мое число больше), попробуйте еще раз.", false, i)
	}

	return continueOutput("Не угадали (мое число меньше), попробуйте еще раз.", false, i)
}

func continueOutput(msg string, gameOver bool, i *input) output {
	return output{
		EncryptedInput:    encryptInput(*i),
		NumAttempts:       i.NumAttempts,
		Message:           msg,
		MinAnswer:         minNum,
		MaxAnswer:         maxNum,
		GameOver:          gameOver,
		PrettyAttemptsLog: prettifyAttemptsLog(i),
	}
}

func prettifyAttemptsLog(i *input) string {
	pretty := ""
	for _, entry := range i.AttemptsLog {
		var text string
		if entry < i.CorrectAnswer {
			text = fmt.Sprintf("%d (нужно больше)", entry)
		}
		if entry > i.CorrectAnswer {
			text = fmt.Sprintf("%d (нужно меньше)", entry)
		}
		if entry == i.CorrectAnswer {
			text = fmt.Sprintf("%d (правильно)", entry)
		}

		pretty = fmt.Sprintf("%s<li>%s</li>", pretty, text)
	}
	return pretty
}

func reduceAndLogAttempt(i input) *input {
	numAttempts := i.NumAttempts
	attemptsLog := i.AttemptsLog
	if numAttempts > 0 {
		numAttempts--
		attemptsLog = append(i.AttemptsLog, i.Answer)
	}
	return &input{
		i.CorrectAnswer,
		i.Answer,
		numAttempts,
		attemptsLog,
	}
}
