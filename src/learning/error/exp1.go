package error

import (
	"fmt"
	"time"
)

type MyError struct {
	When time.Time
	What string
}

func (err *MyError) Error() string {
	return fmt.Sprintf("at %v, %s", err.When, err.What)
}

func run() error {
	return &MyError{
		time.Now(),
		"Some error has happened",
	}
}

func Experiment1() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
