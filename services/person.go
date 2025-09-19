package services

import (
	"flag"
	"fmt"
)

const ageMessage = "Should be a positive integer"
const nameMessage = "Name cannot be empty"

type Person struct {
	Age  int
	Name string
}

func NewPerson() *Person {
	return &Person{}
}

func (p *Person) ParseFlags() error {
	flag.IntVar(&p.Age, "age", -1, ageMessage)
	flag.StringVar(&p.Name, "name", "", nameMessage)

	flag.Parse()

	if err := p.UpdateAge(p.Age); err != nil {
		return err
	}

	if err := p.UpdateName(p.Name); err != nil {
		return err
	}

	fmt.Printf("%+v\n", p)
	return nil
}

func (p *Person) UpdateAge(age int) error {
	if age <= 0 {
		return fmt.Errorf("%s", ageMessage)
	}
	p.Age = age
	fmt.Printf("Age updated to %d\n", age)
	return nil
}

func (p *Person) UpdateName(name string) error {
	if name == "" {
		return fmt.Errorf("%s", nameMessage)
	}
	p.Name = name
	fmt.Printf("Name updated to %s\n", name)
	return nil
}
