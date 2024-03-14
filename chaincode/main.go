package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

func main() {

	randomUUID, err := uuid.NewV7()
	if err != nil {
		log.Fatal(err)
	}

	var ballotID string = "b-" + randomUUID.String()
	fmt.Println(ballotID)
}
