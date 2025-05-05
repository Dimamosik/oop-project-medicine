package medbot

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID      int
	Name    string
	Cart    []Pharmacy
	History []string
}

type Query struct {
	UserID    int
	Timestamp time.Time
	Content   string
}

type Response struct {
	Content string
	Time    time.Time
}

type Chatbot struct {
	Name string
	Base MedDataBase
	Conv []Query
}

type MedDataBase struct{}

type Doctor struct {
	Name   string
	Field  string
	Rating string
	City   string
}

type Appointment struct {
	DoctorName  string
	DoctorField string
	Date        time.Time
}

type Pharmacy struct {
	Med       string
	Price     int
	Available int
}

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

func (cb *Chatbot) GenerateResponse(query Query, user *User) Response {
	var content string
	lowerInput := strings.ToLower(query.Content)

	switch {
	case lowerInput == "help":
		content = "You can ask me about your symptoms, or you can ask for a list of doctors. For example:\n" +
			"- Type 'list' to see a list of doctors.\n" +
			"- Ask about symptoms like 'headache' or 'fever' to get advice on how to deal with them.\n" +
			"- You can also inquire about medications available in our pharmacy or buy them directly.\n" +
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
	case strings.Contains(lowerInput, "doctor"):
		content = cb.Base.SuggestDoctor("general")
	case strings.Contains(lowerInput, "list"):
		content = cb.Base.ListDoctors()
	case strings.Contains(lowerInput, "select"):
		content = cb.Base.SelectDoctor(query.Content)
	case strings.Contains(lowerInput, "pharmacy") || strings.HasPrefix(lowerInput, "buy"):
		content = cb.Base.CheckPharmacy(query.Content, user)
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

func (db *MedDataBase) GetInfo(topic string) string {
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

func (md *MedDataBase) CheckPharmacy(input string, user *User) string {
	pharmacyItems := []Pharmacy{
		{"Ibuprofen", 10, 20},
		{"Paracetamol", 5, 50},
		{"Cough Syrup", 8, 10},
		{"Aspirin", 7, 30},
	}

	if strings.Contains(input, "buy") {
		parts := strings.Fields(input)
		if len(parts) < 2 {
			return "Please specify the medication you want to buy (e.g., 'buy Ibuprofen')."
		}

		medName := strings.Join(parts[1:], " ")
		for i, item := range pharmacyItems {
			if strings.EqualFold(item.Med, medName) {
				if item.Available > 0 {
					user.AddToCart(item)
					pharmacyItems[i].Available--
					return fmt.Sprintf("%s has been added to your cart.", item.Med)
				} else {
					return fmt.Sprintf("Sorry, %s is out of stock.", item.Med)
				}
			}
		}
		return "Medication not found."
	}

	var response strings.Builder
	response.WriteString("Available medications in the pharmacy:\n")
	for _, item := range pharmacyItems {
		response.WriteString(
			fmt.Sprintf("- %s: $%d (Availability: %d)\n", item.Med, item.Price, item.Available),
		)
	}

	return response.String()
}

func (user *User) AddToCart(med Pharmacy) {
	user.Cart = append(user.Cart, med)
}

func (user *User) ViewCart() string {
	var response strings.Builder
	if len(user.Cart) == 0 {
		return "Your cart is empty."
	}

	response.WriteString("Your cart contains the following items:\n")
	for _, item := range user.Cart {
		response.WriteString(fmt.Sprintf("- %s: $%d\n", item.Med, item.Price))
	}

	return response.String()
}

func (user *User) Checkout() string {
	var total int
	for _, item := range user.Cart {
		total += item.Price
	}

	user.Cart = []Pharmacy{}

	return fmt.Sprintf("Your total is $%d. Thank you for your purchase!", total)
}
