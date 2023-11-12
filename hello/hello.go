package main

import (
	"fmt"
	"log"

	"nicograef/greetings"
)

func main() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0) // remove printing of time, source file and line number

	names := []string{"Gladys", "Samantha", "Darrin"}

	messages, err := greetings.Hellos(names)

	if err != nil {
		log.Fatal(err)
	}

	for name, message := range messages {
		fmt.Println(name, ":", message)
	}

}
