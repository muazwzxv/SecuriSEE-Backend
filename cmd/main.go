package main

import (
	"Oracle-Hackathon-BE/config"
	"fmt"
)

func main() {
	fmt.Println("hello bitches")

	config := config.New()
	fmt.Println(config.ReadEnv("Database.Password"))
}
