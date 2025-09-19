package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/tusharameria/gorep/services"
)

func main() {
	person := services.NewPerson()

	if err := person.ParseFlags(); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter new name if you want to change :")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	fmt.Printf("You entered : %s", text)

	if err := person.UpdateName(text); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)

	i := 0
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			break
		}
		fmt.Printf("%d : %s\n", i, text)
		i++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	fmt.Println("Program Completed!!!")
}
