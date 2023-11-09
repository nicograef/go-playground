package main

import "fmt"

type Person struct {
	name string
	age  int
}

func (person Person) sayHello() {
	fmt.Println("Hello world,")
	fmt.Println("my name is ", person.name, " and I am ", person.age, " years old.")
	fmt.Println("Alright, goodbye!")
}

func main() {

	var john Person
	john.name = "John Doe"
	john.age = 36

	alice := Person{
		name: "Alice",
		age:  25,
	}

	bob := Person{
		name: "Bob Barker",
		age:  99,
	}

	var persons = [3]Person{john, alice, bob}

	for _, person := range persons {
		person.sayHello()
		fmt.Println("")
	}

}
