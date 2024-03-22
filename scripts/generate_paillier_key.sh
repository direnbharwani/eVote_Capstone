#!/bin/bash

if [ -z "$1" ]; then
    echo "Missing key length! Usage: $0 <key_length>"
    exit 1
fi

# Ensure folder for storing keys exist
if [ ! -d ./keys/ ]; then
    mkdir keys
fi

# Create temp file to generate keys

echo 'package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
    "strconv"
	"time"

	paillier "github.com/direnbharwani/evote-capstone/paillier"
)

type KeyPair struct {
	Public  string
	Private string
}

func main() {

    keylen, err := strconv.Atoi(os.Getenv("KEY_LENGTH"))
    if err != nil {
        log.Fatal(err)
    }

	publicKey, privateKey, err := paillier.GenerateKeys(keylen)
	if err != nil {
		log.Fatal(err)
	}

	var (
		now      = time.Now().Unix()
		fileName = fmt.Sprintf("keys/paillierKeyPair_%d.json", now)
	)

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Serialize keys to JSON
	publicKeyBase64, err := paillier.Base64Encode(publicKey)
	if err != nil {
		log.Fatal(err)
	}
	privateKeyBase64, err := paillier.Base64Encode(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	keys := KeyPair{publicKeyBase64, privateKeyBase64}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "	")

	err = encoder.Encode(keys)
	if err != nil {
		log.Fatal(err)
	}
}' > main.go

# Generate Keys
go mod tidy

export KEY_LENGTH="$1"

go run main.go

# Clean up
unset KEY_LENGTH
rm main.go
