package main

import (
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"github.com/gtank/cryptopasta"
	"strconv"
)

type input struct {
	CorrectAnswer, Answer, NumAttempts int
	AttemptsLog                        []int
}

//go:embed key.bin
var key []byte

func getKey() *[32]byte {
	if len(key) != 32 {
		panic("Key must be 32 bytes in size")
	}
	return (*[32]byte)(key)
}

func decryptInput(encryptedJsonInBase64, answer string) (i input, err error) {
	var encryptedJson []byte
	if encryptedJson, err = base64.StdEncoding.DecodeString(encryptedJsonInBase64); err != nil {
		return
	}

	jsonBytes, err := decryptJson(encryptedJson)
	if err != nil {
		return
	}

	if err = json.Unmarshal(jsonBytes, &i); err != nil {
		return
	}

	var answerNum int
	if answerNum, err = strconv.Atoi(answer); err != nil {
		return
	}

	i.Answer = answerNum
	return
}

func decryptJson(encryptedJson []byte) (jsonBytes []byte, err error) {
	key := getKey()
	jsonBytes, err = cryptopasta.Decrypt(encryptedJson, key)
	return
}

func encryptInput(i input) string {
	jsonBytes, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}

	key := getKey()
	encryptedJson, err := cryptopasta.Encrypt(jsonBytes, (*[32]byte)(key))
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(encryptedJson)
}
