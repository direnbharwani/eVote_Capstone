package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	paillier "github.com/direnbharwani/go-paillier/pkg"
)

func main() {

	publicKey, privateKey, err := paillier.GenerateKeys(128)
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
	publicKeyJSON, err := paillier.SerializeJSON(publicKey)
	if err != nil {
		log.Fatal(err)
	}
	privateKeyJSON, err := paillier.SerializeJSON(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	keys := struct {
		Public  string
		Private string
	}{
		Public:  base64.StdEncoding.EncodeToString(publicKeyJSON),
		Private: base64.StdEncoding.EncodeToString(privateKeyJSON),
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "	")

	err = encoder.Encode(keys)
	if err != nil {
		log.Fatal(err)
	}
}
