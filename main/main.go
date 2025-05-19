package main

import (
	"bufio" // For buffered input from the user
	"fmt"
	"os"      // For accessing operating system features like stdin
	"strings" // For string manipulation
	"time"    // For working with dates and times

	"medbot_project/medbot" // Importing the custom medbot package that defines the chatbot logic
)

func main() {

	user := medbot.User{
		ID:   1,
		Name: "Dmytro",
	}

	// It has a name and a base database
	bot := medbot.Chatbot{
		Name: "MedBot",
		Base: medbot.MedDataBase{}, // Provide an instance of the medical database
	}

	fmt.Println("\nWelcome to MedBot - your vitual medical assistant.")
	fmt.Println(" Type 'help' to see what i can do or 'exit' to leave.")
	fmt.Println()
	// Create a buffered reader to read user input from the terminal
	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print("You: ")

		// Read a full line of text until the user presses "Enter"
		input, err := reader.ReadString('\n') // Read input including newline
		if err != nil {
			// Handle any error that occurred while reading input
			fmt.Println("Error reading input:", err)
			continue
		}

		// Remove any leading whitespace
		input = strings.TrimSpace(input)

		if strings.ToLower(input) == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		// Pass the user's input to the chatbot and get a parsed query
		query := bot.ReceiveInput(&user, input)

		// Ask the bot to generate a response based on the query and user
		response := bot.GenerateResponse(query, &user)

		// Print the bot's response with a timestamp
		fmt.Printf("\nMedBot [%s]: %s\n", response.Time.Format(time.Kitchen), response.Content)
	}
}
