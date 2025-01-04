To generate a Golang struct from the given JSON payload, you can use the `encoding/json` package. Here's how you can define the struct:

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	IsStudent bool   `json:"isStudent"`
	Address struct {
		Street string `json:"street"`
		City  string `json:"city"`
		State string `json:"state"`
	} `json:"address"`
	Hobbies []string `json:"hobbies"`
}

func main() {
	jsonPayload := `{
	"name": "Alice",
	"age": 30,
	"isStudent": false,
	"address": {
		"street": "123 Main St",
		"city": "Springfield",
		"state": "IL"
	},
	"hobbies": ["reading", "swimming"]
}`

	var person Person
	err := json.Unmarshal([]byte(jsonPayload), &person)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	fmt.Printf("Name: %s\n", person.Name)
	fmt.Printf("Age: %d\n", person.Age)
	fmt.Printf("Is Student: %v\n", person.IsStudent)
	fmt.Printf("Address: %s\n", person.Address.Street)
	fmt.Printf("City: %s\n", person.Address.City)
	fmt.Printf("State: %s\n", person.Address.State)
	fmt.Printf("Hobbies: %v\n", person.Hobbies)
}
```

### Explanation:
- **Struct Definition**: The `Person` struct is defined with fields for `Name`, `Age`, `IsStudent`, `Address`, and `Hobbies`.
- **Unmarshalling**: The `json.Unmarshal` function is used to convert the JSON string into a `Person` struct. If there is an error during unmarshalling, it will print an error message.
- **Output**: The program prints the details of the person, including their name, age, student status, address details, and hobbies.

This code will output the following:
```
Name: Alice
Age: 30
Is Student: false
Address: Street: 123 Main St, City: Springfield, State: IL
Hobbies: [reading swimming]
```
