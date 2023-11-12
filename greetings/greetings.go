package greetings

import (
	"errors"
	"fmt"
	"math/rand"
)

// returns a greeting for the name
func Hello(name string) (string, error) {

	if name == "" {
		return "", errors.New("empty name")
	}

	message := fmt.Sprintf(randomFormat(), name)

	return message, nil
}

// returns a map that associates each of the names with a greeting
func Hellos(names []string) (map[string]string, error) {
	messages := make(map[string]string)

	for _, name := range names {
		message, err := Hello(name)
		
		if err != nil {
			return nil, err
		}
		
		messages[name] = message
	}

	return messages, nil
}

// returns a random greeting message format
func randomFormat() string {
	formats := []string{
		"Hi, %v. Welcome!",
		"Great to see you, %v!",
		"Howdy, %v! Well met!",
	}

	return formats[rand.Intn(len(formats))]
}
