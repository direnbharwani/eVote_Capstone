package main

import (
	"log"
	"time"
)

func main() {

	birthdayString := "1996-08-03"
	birthday, err := time.Parse(time.DateOnly, birthdayString)
	if err != nil {
		log.Printf("Failed to parse birthday: %s\n", err)
	}

	age := time.Now().Year() - birthday.Year()
	log.Println("Age: ", age)
}
