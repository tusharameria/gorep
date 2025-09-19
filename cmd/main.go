package main

import (
	"flag"
	"fmt"
	"os"
)

const ageMessage = "Should be a positive integer"
const nameMessage = "Name cannot be empty"

type Person struct {
	Age  int
	Name string
}

func main() {
	var person Person

	flag.IntVar(&person.Age, "age", -1, ageMessage)
	flag.StringVar(&person.Name, "name", "", nameMessage)

	flag.Parse()

	if person.Age <= 0 {
		fmt.Printf("%s\n", ageMessage)
		os.Exit(1)
	}

	if person.Name == "" {
		fmt.Printf("%s\n", nameMessage)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", person)
}
