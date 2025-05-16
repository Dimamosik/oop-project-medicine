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

	// This simulates the person chatting with the bot
	user := medbot.User{
		ID:   1,       // Give the user an ID
		Name: "Alice", // Give the user a name
	}

	// It has a name and a base database
	bot := medbot.Chatbot{
		Name: "MedBot",
		Base: medbot.MedDataBase{}, // Provide an instance of the medical database
	}

	fmt.Println("\nWelcome to MedBot. Ask me a question. (type 'help' for information)")

	// Create a buffered reader to read user input from the terminal
	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print("You: ")

		// Read a full line of text until the user presses "Enter"
		input, err := reader.ReadString('\n') // Read input including newline
		if err != nil {
			// Handle any error that occurred while reading input
			fmt.Println("Error reading input:", err)
			continue // Skip this iteration and wait for new input
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
