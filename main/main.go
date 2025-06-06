package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"medbot_project/medbot"
)

func main() {
	user := &medbot.User{Name: "Alice"}

	db := medbot.MedDataBase{
		Doctors: []medbot.Doctor{
			{ID: "1", Name: "Dr. Smith", Special: "General", Location: "New York"},
			{ID: "2", Name: "Dr. Life", Special: "Cardiology", Location: "San Francisco"},
			{ID: "3", Name: "Dr. Bold", Special: "Pediatrics", Location: "Chicago"},
			{ID: "4", Name: "Dr. Rose", Special: "Dermatology", Location: "Los Angeles"},
			{ID: "5", Name: "Dr. Green", Special: "Neurology", Location: "Boston"},
		},
		Medicines: []string{
			"Paracetamol",
			"Amoxicillin",
			"Ibuprofen",
			"Loratadine",
			"Omeprazole",
		},
	}

	bot := medbot.Chatbot{Base: db}

	fmt.Println("You can ask me about your symptoms, or you can ask for a list of doctors. For example:\n" +
		"- Type 'list' to see a list of doctors or you can type recommend.\n" +
		"- Ask about symptoms like 'headache' or 'fever' to get advice on how to deal with them.\n" +
		"- You can add medicines using 'buy <medication> [quantity]' or remove with 'remove [quantity] <medication>'.\n" +
		"- Type 'view cart' to see items in your cart.\n" +
		"- Type 'book' to book an appointment 'book 3 14:00'.")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if strings.ToLower(input) == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		query := medbot.Query{Content: input}
		response := bot.GenerateResponse(query, user)
		fmt.Printf("[%s] Bot: %s\n", response.Time.Format("15:04:05"), response.Content)
	}
}
