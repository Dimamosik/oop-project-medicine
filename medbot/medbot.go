package medbot

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// USER STRUCT AND METHODS

/*
User represents a patient using the chatbot. It contains basic personal data,
a cart for pharmacy purchases, and interaction history.
*/
type User struct {
	ID      int
	Name    string
	Cart    []Pharmacy // Composition: Cart is composed of Pharmacy items
	History []string   // Stores all user interactions
}

/*
AddToCart adds a medication to the user's cart.
Demonstrates encapsulation: Cart is modified through a method.
*/
func (user *User) AddToCart(med Pharmacy, quantity int) {
	for i, item := range user.Cart {
		if strings.EqualFold(item.Med, med.Med) {
			user.Cart[i].Quantity += quantity
			return
		}
	}
	med.Quantity = quantity
	user.Cart = append(user.Cart, med)
}

// ViewCart returns a string listing all the items in the cart.
func (user *User) ViewCart() string {
	if len(user.Cart) == 0 {
		return "Your cart is empty."
	}

	var response strings.Builder
	var total int
	response.WriteString("Your cart contains the following items:\n")
	for _, item := range user.Cart {
		subtotal := item.Price * item.Quantity
		total += subtotal
		response.WriteString(fmt.Sprintf("- %s x%d: $%d\n", item.Med, item.Quantity, subtotal))
	}
	response.WriteString(fmt.Sprintf("Total: $%d", total))
	return response.String()
}

// Checkout finalizes the purchase, returns the total cost, and clears the cart.
func (user *User) Checkout() string {
	var total int
	for _, item := range user.Cart {
		total += item.Price * item.Quantity
	}
	user.Cart = []Pharmacy{} // Clear the cart
	return fmt.Sprintf("Your total is $%d. Thank you for your purchase!", total)
}

func (user *User) RemoveFromCart(medName string) string {
	for i, item := range user.Cart {
		if strings.EqualFold(item.Med, medName) {
			user.Cart = append(user.Cart[:i], user.Cart[i+1])
			return fmt.Sprintf("%s has been removed from your cart.", item.Med)
		}
	}
	return fmt.Sprintf("%s was not found in your cart.", medName)
}

// CORE DOMAIN STRUCTS

/*
Query represents a user message, including metadata like timestamp and user ID.
It forms part of the conversation context.
*/
type Query struct {
	UserID    int
	Timestamp time.Time
	Content   string
}

// Response represents the chatbot's reply, with a timestamp for logging.
type Response struct {
	Content string
	Time    time.Time
}

/*
Chatbot encapsulates the bot. Contains bot name, access to a medical database,
and stores the ongoing conversation history.
*/
type Chatbot struct {
	Name string
	Base MedDataBase // Composition: Chatbot uses a MedDataBase
	Conv []Query     // Stores the full query log
}

// CHATBOT BEHAVIOR

// Demonstrates behavioral encapsulation: user interaction is handled internally.
func (cb *Chatbot) ReceiveInput(user *User, input string) Query {
	query := Query{
		UserID:    user.ID,
		Timestamp: time.Now(),
		Content:   input,
	}
	cb.Conv = append(cb.Conv, query)
	user.History = append(user.History, input)
	return query
}

// Demonstrates polymorphism and control flow based on user intent.
func (cb *Chatbot) GenerateResponse(query Query, user *User) Response {
	var content string
	lowerInput := strings.ToLower(query.Content)

	// Intent recognition via pattern matching
	switch {
	case lowerInput == "help":
		content = "You can ask me about your symptoms, or you can ask for a list of doctors. For example:\n" +
			"- Type 'list' to see a list of doctors or you can type recommend.\n" +
			"- Ask about symptoms like 'headache' or 'fever' to get advice on how to deal with them.\n" +
			"- You can buy medicines using 'buy <medication>' or remove with 'remove <medication>'.\n" +
			"- Type 'view cart' to see items in your cart or 'checkout' to complete your purchase."
	case strings.Contains(lowerInput, "headache"):
		content = cb.Base.GetInfo("headache")
	case strings.Contains(lowerInput, "fever"):
		content = cb.Base.GetInfo("fever")
	case strings.Contains(lowerInput, "toothache"):
		content = cb.Base.GetInfo("toothache")
	case strings.Contains(lowerInput, "diarrhea"):
		content = cb.Base.GetInfo("diarrhea")
	case strings.Contains(lowerInput, "chill"):
		content = cb.Base.GetInfo("chill")
	case strings.Contains(lowerInput, "runny nose"):
		content = cb.Base.GetInfo("runny nose")
	case strings.Contains(lowerInput, "vomiting"):
		content = cb.Base.GetInfo("vomiting")
	case strings.Contains(lowerInput, "cut"):
		content = cb.Base.GetInfo("cut")
	case strings.Contains(lowerInput, "recommend"):
		content = cb.Base.SuggestDoctor("general")
	case strings.Contains(lowerInput, "list"):
		content = cb.Base.ListDoctors()
	case strings.Contains(lowerInput, "select"):
		content = cb.Base.SelectDoctor(query.Content)
	case strings.Contains(lowerInput, "pharmacy") || strings.HasPrefix(lowerInput, "buy"):
		content = cb.Base.CheckPharmacy(query.Content, user)
	case strings.HasPrefix(lowerInput, "remove"):
		parts := strings.Fields(query.Content)
		if len(parts) < 2 {
			content = "Please specify the medication you want to remove (e.g., 'remove Ibuprofen'.)"
		} else {
			medName := strings.Join(parts[1:], " ")
			content = user.RemoveFromCart(medName)
		}
	case strings.Contains(lowerInput, "view cart"):
		content = user.ViewCart()
	case strings.Contains(lowerInput, "checkout"):
		content = user.Checkout()
	default:
		content = "I don't understand that. Type 'help' for instructions on how to interact with me."
	}

	return Response{
		Content: content,
		Time:    time.Now(),
	}
}

// DATABASE: Encapsulation of Medical Knowledge and Services

type MedDataBase struct{} // Acts as a service layer

// Doctor struct represents a physician. Part of the database, not behavior-heavy.
type Doctor struct {
	Name   string
	Field  string
	Rating string
	City   string
}

// GetInfo returns health advice based on a given symptom keyword.
func (md *MedDataBase) GetInfo(topic string) string {
	data := map[string]string{
		"headache":   "Use a hot or cold compress on your head or neck. Try gentle massage. Drink small amounts of caffeine. Take over-the-counter pain relievers like ibuprofen or aspirin.",
		"fever":      "If you have a fever rest, stay hydrated, and take fever-reducing medicine like acetaminophen or ibuprofen.",
		"toothache":  "Rinse your mouth with warm salt water. Use a cold compress on your cheek. Apply clove oil to the tooth. Take over-the-counter painkillers like ibuprofen.",
		"diarrhea":   "Drink plenty of fluids and oral rehydration solution. Eat bland foods like bananas and rice. Avoid dairy, caffeine, and greasy food. Rest as much as possible.",
		"chill":      "Wear warm clothing and use blankets. Drink hot tea or soup. Rest and check for fever. Use a warm compress if you feel tense or achy.",
		"runny nose": "Drink lots of water. Use a saline spray or rinse. Try a warm compress on your face. Rest and avoid allergens or irritants.",
		"vomiting":   "Sip water or electrolyte drinks slowly. Avoid solid food until vomiting stops. Eat plain food like toast or crackers after. Rest and avoid strong smells.",
		"cut":        "Wash the cut with water and mild soap. Press with a clean cloth to stop bleeding. Apply antiseptic and cover with a bandage. Change the bandage daily.",
	}
	if info, ok := data[topic]; ok {
		return info
	}
	return "Sorry, I have no information on that."
}

// SuggestDoctor recommends a doctor from a hardcoded list based on specialization.
func (md *MedDataBase) SuggestDoctor(field string) string {
	doctors := []Doctor{
		{"Dr. Smith", "General", "4.7", "New York"},
		{"Dr. Life", "Cardiology", "4.5", "San Francisco"},
		{"Dr. Bold", "Pediatrics", "4.9", "Chicago"},
		{"Dr. Rose", "Dermatology", "4.6", "Los Angeles"},
		{"Dr. Green", "Neurology", "4.8", "Boston"},
	}
	for _, d := range doctors {
		if strings.EqualFold(d.Field, field) {
			return "Recommended doctor: " + d.Name + " (" + d.Field + ", Rating: " + d.Rating + ")"
		}
	}
	return "Sorry, no doctor found for that speciality."
}

// ListDoctors returns a string listing all available doctors.
func (md *MedDataBase) ListDoctors() string {
	doctors := []Doctor{
		{"Dr. Smith", "General", "4.7", "New York"},
		{"Dr. Life", "Cardiology", "4.5", "San Francisco"},
		{"Dr. Bold", "Pediatrics", "4.9", "Chicago"},
		{"Dr. Rose", "Dermatology", "4.6", "Los Angeles"},
		{"Dr. Green", "Neurology", "4.8", "Boston"},
	}

	var response strings.Builder
	response.WriteString("Here are the available doctors:\n")
	for i, d := range doctors {
		response.WriteString(
			fmt.Sprintf("%d. %s - %s (Location: %s, Rating: %s)\n",
				i+1, d.Name, d.Field, d.City, d.Rating),
		)
	}
	return response.String()
}

func (md *MedDataBase) SelectDoctor(input string) string {
	doctors := []Doctor{
		{"Dr. Smith", "General", "4.7", "New York"},
		{"Dr. Life", "Cardiology", "4.5", "San Francisco"},
		{"Dr. Bold", "Pediatrics", "4.9", "Chicago"},
		{"Dr. Rose", "Dermatology", "4.6", "Los Angeles"},
		{"Dr. Green", "Neurology", "4.8", "Boston"},
	}

	parts := strings.Fields(input)
	if len(parts) < 2 {
		return "Please provide the number of the doctor you wish to select (e.g., 'select 2')."
	}

	doctorNum, err := strconv.Atoi(parts[1])
	if err != nil || doctorNum < 1 || doctorNum > len(doctors) {
		return "Invalid doctor number. Please select a valid number from the list."
	}

	appointmentDate := time.Now().AddDate(0, 0, 7)
	selectedDoctor := doctors[doctorNum-1]

	return fmt.Sprintf("You selected: %s - %s (Location: %s, Rating: %s)\nYour appointment is scheduled for: %s.",
		selectedDoctor.Name, selectedDoctor.Field, selectedDoctor.City, selectedDoctor.Rating, appointmentDate.Format("Monday, January 2, 2006"))
}

// PHARMACY STRUCTS AND LOGIC

// Pharmacy represents a medication item in the virtual pharmacy.
type Pharmacy struct {
	Med       string
	Price     int
	Available int
	Quantity  int
}

// This also modifies the User's Cart, demonstrating behavior interaction between structs.
func (md *MedDataBase) CheckPharmacy(input string, user *User) string {
	pharmacyItems := []Pharmacy{
		{"Ibuprofen", 10, 20, 0},
		{"Paracetamol", 5, 50, 0},
		{"Cough Syrup", 8, 10, 0},
		{"Aspirin", 7, 30, 0},
	}

	if strings.HasPrefix(strings.ToLower(input), "buy") {
		parts := strings.Fields(input)
		if len(parts) < 2 {
			return "Please specify the medication (e.g., 'buy Ibuprofen 2')."
		}

		medName := parts[1]
		quantity := 1 // Default

		if len(parts) >= 3 {
			q, err := strconv.Atoi(parts[2])
			if err == nil && q > 0 {
				quantity = q
			}
		}

		for i, item := range pharmacyItems {
			if strings.EqualFold(item.Med, medName) {
				if item.Available >= quantity {
					user.AddToCart(item, quantity)
					pharmacyItems[i].Available -= quantity
					return fmt.Sprintf("%d x %s has been added to your cart.", quantity, item.Med)
				}
				return fmt.Sprintf("Sorry, only %d %s available.", item.Available, item.Med)
			}
		}
		return "Medication not found."
	}

	var response strings.Builder
	response.WriteString("Available medications in the pharmacy:\n")
	for _, item := range pharmacyItems {
		response.WriteString(fmt.Sprintf("- %s: $%d (Availability: %d)\n", item.Med, item.Price, item.Available))
	}
	return response.String()
}
