package errors

import "log"

func Must[T any](value T, err error) T {
	if err != nil {
		log.Fatalf("Critical error: %v", err)
	}
	return value
}

func LogError[T any](value T, err error) T {
	if err != nil {
		log.Printf("Error: %v", err)
	}
	return value
}
