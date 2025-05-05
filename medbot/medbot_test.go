package medbot

import "testing"

func TestGenerateResponse(t *testing.T) {
	// Initialize the MedDataBase and chatbot for testing
	db := MedDataBase{}
	bot := Chatbot{
		Name: "TestBot",
		Base: db,
	}
	// Expected formatted list of doctors for the "list" command
	expectedDoctorList := "Here are the available doctors:\n" +
		"1. Dr. Smith - General (Location: New York, Rating: 4.7)\n" +
		"2. Dr. Life - Cardiology (Location: San Francisco, Rating: 4.5)\n" +
		"3. Dr. Bold - Pediatrics (Location: Chicago, Rating: 4.9)\n" +
		"4. Dr. Rose - Dermatology (Location: Los Angeles, Rating: 4.6)\n" +
		"5. Dr. Green - Neurology (Location: Boston, Rating: 4.8)\n"
		// Table of test cases to check chatbot behavior for different inputs
	testCases := []struct {
		input    string
		expected string
	}{
		{"I have a headache", "Use a hot or cold compress on your head or neck. Try gentle massage. Drink small amounts of caffeine. Take over-the-counter pain relievers like ibuprofen or aspirin."},
		{"fever", "If you have a fever rest, stay hydrated, and take fever-reducing medicine like acetaminophen or ibuprofen."},
		{"toothache", " Rinse your mouth with warm salt water. Use a cold compress on your cheek. Apply clove oil to the tooth.Take over-the-counter painkillers like ibuprofen. "},
		{"diarrhea", " Drink plenty of fluids and oral rehydration solution. Eat bland foods like bananas and rice.Avoid dairy, caffeine, and greasy food.Rest as much as possible."},
		{"chill", "Wear warm clothing and use blankets. Drink hot tea or soup. Rest and check for fever. Use a warm compress if you feel tense or achy."},
		{"runny nose", "Drink lots of water. Use a saline spray or rinse. Try a warm compress on your face.Rest and avoid allergens or irritants."},
		{"vomiting", "Sip water or electrolyte drinks slowly. Avoid solid food until vomiting stops. Eat plain food like toast or crackers after.Rest and avoid strong smells."},
		{"cut", "Wash the cut with water and mild soap. Press with a clean cloth to stop bleeding. Apply antiseptic and cover with a bandage.Change the bandage daily."},
		{"Which doctor should I see?", "Recommended doctor: Dr. Smith (General, Rating: 4.7)"},
		{"Unknown symptom", "I don't understand that. Type 'help' for instructions on how to interact with me."},
		{"list", expectedDoctorList},
	}
	// Loop through each test case and compare actual vs expected responses
	for _, tc := range testCases {
		// Simulate a user query
		query := Query{UserID: 1, Content: tc.input}

		resp := bot.GenerateResponse(query, &User{ID: 1})
		// If the response doesn't match what we expected, report an error
		if resp.Content != tc.expected {
			t.Errorf("For input '%s', expected '%s', got '%s'", tc.input, tc.expected, resp.Content)
		}
	}
}

// TestCheckPharmacy checks if pharmacy queries and purchase flows work correctly
func TestCheckPharmacy(t *testing.T) {
	// Create a user and a chatbot instance
	user := &User{ID: 1}
	db := MedDataBase{}
	bot := Chatbot{
		Base: db,
	}
	// Check the pharmacy inventory listing
	result := bot.Base.CheckPharmacy("pharmacy", user)

	expected := "Available medications in the pharmacy:\n" +
		"- Ibuprofen: $10 (Availability: 20)\n" +
		"- Paracetamol: $5 (Availability: 50)\n" +
		"- Cough Syrup: $8 (Availability: 10)\n" +
		"- Aspirin: $7 (Availability: 30)\n"

	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	// Reset the user's cart and try to buy a medicine
	user.Cart = nil
	buyResult := bot.Base.CheckPharmacy("buy Ibuprofen", user)

	expectedBuy := "Ibuprofen has been added to your cart."

	if buyResult != expectedBuy {
		t.Errorf("Expected '%s', but got '%s'", expectedBuy, buyResult)
	}
	// Verify cart contents after purchase
	cartContents := user.ViewCart()
	expectedCart := "Your cart contains the following items:\n" +
		"- Ibuprofen: $10\n"
	if cartContents != expectedCart {
		t.Errorf("Expected '%s', but got '%s'", expectedCart, cartContents)
	}
	// Test checkout and final message
	checkoutResult := user.Checkout()
	expectedCheckout := "Your total is $10. Thank you for your purchase!"
	if checkoutResult != expectedCheckout {
		t.Errorf("Expected '%s', but got '%s'", expectedCheckout, checkoutResult)
	}
}
