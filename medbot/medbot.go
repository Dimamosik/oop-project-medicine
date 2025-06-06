package medbot

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Query struct {
	Content string
}

type Response struct {
	Content string
	Time    time.Time
}

type User struct {
	Name string
}

type Doctor struct {
	ID       string
	Name     string
	Special  string
	Location string
	Ratings  []int
}

type Appointment struct {
	UserName   string
	DoctorName string
	TimeSlot   string
}

type Item struct {
	Name     string
	Quantity int
}

type MedDataBase struct {
	Doctors      []Doctor
	Appointments []Appointment
	Cart         []Item
	Medicines    []string
}

type CommandHandler interface {
	CanHandle(input string) bool
	Handle(query Query, user *User, db *MedDataBase) string
}

type SymptomInfoHandler struct{}

func (s SymptomInfoHandler) CanHandle(input string) bool {
	symptoms := getSymptomData()
	lower := strings.ToLower(input)
	for keyword := range symptoms {
		if strings.Contains(lower, keyword) {
			return true
		}
	}
	return false
}

func (s SymptomInfoHandler) Handle(query Query, user *User, db *MedDataBase) string {
	symptoms := getSymptomData()
	lower := strings.ToLower(query.Content)
	for keyword, info := range symptoms {
		if strings.Contains(lower, keyword) {
			return fmt.Sprintf("Advice for %s: %s", keyword, info)
		}
	}
	return "I'm not sure how to help with that symptom."
}

func getSymptomData() map[string]string {
	return map[string]string{
		"headache":   "Use a cold compress, take ibuprofen, rest in a quiet dark room.",
		"fever":      "Take acetaminophen, stay hydrated, and rest.",
		"toothache":  "Rinse mouth with warm salt water, take ibuprofen, apply cold compress.",
		"diarrhea":   "Drink oral rehydration fluids, eat bananas and rice.",
		"chill":      "Wear warm clothes, drink hot drinks, rest.",
		"runny nose": "Use a saline spray, drink water, rest.",
		"vomiting":   "Sip fluids, rest, avoid food until it stops.",
		"cut":        "Clean with water, apply antiseptic and bandage.",
	}
}

type PharmacyHandler struct{}

func (p PharmacyHandler) CanHandle(input string) bool {
	return strings.ToLower(input) == "pharmacy"
}

func (p PharmacyHandler) Handle(query Query, user *User, db *MedDataBase) string {
	if len(db.Medicines) == 0 {
		return "No medicines available at the moment."
	}
	return "Available medicines: " + strings.Join(db.Medicines, ", ")
}

type RateDoctorHandler struct{}

func (r RateDoctorHandler) CanHandle(input string) bool {
	return strings.HasPrefix(strings.ToLower(input), "rate")
}

func (r RateDoctorHandler) Handle(query Query, user *User, db *MedDataBase) string {
	parts := strings.Fields(query.Content)
	if len(parts) != 3 {
		return "Invalid rating format. Use: rate <DoctorID> <Rating (1-5)>"
	}

	doctorID := parts[1]
	rating, err := strconv.Atoi(parts[2])
	if err != nil || rating < 1 || rating > 5 {
		return "Please enter a valid rating from 1 to 5."
	}

	for i := range db.Doctors {
		if db.Doctors[i].ID == doctorID {
			db.Doctors[i].Ratings = append(db.Doctors[i].Ratings, rating)
			return fmt.Sprintf("Thank you for rating Dr. %s with %d stars.", db.Doctors[i].Name, rating)
		}
	}

	return "Doctor not found."
}

type ListDoctorsHandler struct{}

func (l ListDoctorsHandler) CanHandle(input string) bool {
	cmd := strings.TrimSpace(strings.ToLower(input))
	return cmd == "list" || cmd == "list doctors"
}

func (l ListDoctorsHandler) Handle(query Query, user *User, db *MedDataBase) string {
	if len(db.Doctors) == 0 {
		return "No doctors available."
	}
	var sb strings.Builder
	sb.WriteString("Available doctors:\n")
	for _, doc := range db.Doctors {
		avgRating := "No ratings"
		if len(doc.Ratings) > 0 {
			sum := 0
			for _, r := range doc.Ratings {
				sum += r
			}
			avg := float64(sum) / float64(len(doc.Ratings))
			avgRating = fmt.Sprintf("%.1f", avg)
		}
		sb.WriteString(fmt.Sprintf("- ID: %s, Name: %s, Specialty: %s, Location: %s, Rating: %s\n",
			doc.ID, doc.Name, doc.Special, doc.Location, avgRating))
	}
	return sb.String()
}

type SelectDoctorHandler struct{}

func (s SelectDoctorHandler) CanHandle(input string) bool {
	return strings.HasPrefix(strings.ToLower(input), "select doctor")
}

func (s SelectDoctorHandler) Handle(query Query, user *User, db *MedDataBase) string {
	doctorID := strings.TrimPrefix(strings.ToLower(query.Content), "select doctor ")
	for _, doc := range db.Doctors {
		if strings.EqualFold(doc.ID, doctorID) {
			return fmt.Sprintf("You selected Dr. %s (%s).", doc.Name, doc.Special)
		}
	}
	return "Doctor not found."
}

type BookAppointmentHandler struct{}

func (b BookAppointmentHandler) CanHandle(input string) bool {
	return strings.HasPrefix(strings.ToLower(input), "book")
}

func (b BookAppointmentHandler) Handle(query Query, user *User, db *MedDataBase) string {
	parts := strings.Fields(query.Content)
	if len(parts) != 3 {
		return "Invalid format. Use: book <DoctorID> <TimeSlot>"
	}
	doctorID := parts[1]
	timeSlot := parts[2]
	for _, doc := range db.Doctors {
		if doc.ID == doctorID {
			db.Appointments = append(db.Appointments, Appointment{
				UserName:   user.Name,
				DoctorName: doc.Name,
				TimeSlot:   timeSlot,
			})
			return fmt.Sprintf("Appointment booked with Dr. %s at %s.", doc.Name, timeSlot)
		}
	}
	return "Doctor not found."
}

type AddToCartHandler struct{}

func (a AddToCartHandler) CanHandle(input string) bool {
	return strings.HasPrefix(strings.ToLower(input), "buy")
}

func (a AddToCartHandler) Handle(query Query, user *User, db *MedDataBase) string {
	parts := strings.Fields(query.Content)
	if len(parts) != 3 {
		return "Invalid format. Use: buy <ItemName> <Quantity>"
	}
	itemName := parts[1]
	quantity, err := strconv.Atoi(parts[2])
	if err != nil || quantity <= 0 {
		return "Please enter a valid quantity."
	}
	for i := range db.Cart {
		if db.Cart[i].Name == itemName {
			db.Cart[i].Quantity += quantity
			return fmt.Sprintf("Updated %s quantity to %d.", itemName, db.Cart[i].Quantity)
		}
	}
	db.Cart = append(db.Cart, Item{Name: itemName, Quantity: quantity})
	return fmt.Sprintf("Added %s x%d to your cart.", itemName, quantity)
}

type RemoveFromCartHandler struct{}

func (r RemoveFromCartHandler) CanHandle(input string) bool {
	return strings.HasPrefix(strings.ToLower(input), "remove")
}

func (r RemoveFromCartHandler) Handle(query Query, user *User, db *MedDataBase) string {
	parts := strings.Fields(query.Content)
	if len(parts) <= 2 || len(parts) > 3 {
		return "Invalid format. Use: remove <ItemName> [Quantity]"
	}
	quantity := 1
	itemName := ""

	if len(parts) == 2 {
		itemName = parts[1]
	} else {
		q, err := strconv.Atoi(parts[1])
		if err != nil || q <= 0 {
			return "Please enter a valid quantity to remove."
		}
		quantity = q
		itemName = parts[2]
	}

	for i := range db.Cart {
		if db.Cart[i].Name == itemName {
			if db.Cart[i].Quantity <= quantity {
				db.Cart = append(db.Cart[:i], db.Cart[i+1:]...)
				return fmt.Sprintf("Removed all of %s from your cart.", itemName)
			} else {
				db.Cart[i].Quantity -= quantity
				return fmt.Sprintf("Removed %d %s. Remaining: %d.", quantity, itemName, db.Cart[i].Quantity)
			}
		}
	}
	return "Item not found in your cart."
}

type ViewCartHandler struct{}

func (v ViewCartHandler) CanHandle(input string) bool {
	return strings.ToLower(input) == "view cart"
}

func (v ViewCartHandler) Handle(query Query, user *User, db *MedDataBase) string {
	if len(db.Cart) == 0 {
		return "Your cart is empty."
	}
	var sb strings.Builder
	sb.WriteString("Your cart contains:\n")
	for _, item := range db.Cart {
		sb.WriteString(fmt.Sprintf("- %s x%d\n", item.Name, item.Quantity))
	}
	return sb.String()
}

type ViewAppointmentsHandler struct{}

func (v ViewAppointmentsHandler) CanHandle(input string) bool {
	return strings.ToLower(input) == "view appointments"
}

func (v ViewAppointmentsHandler) Handle(query Query, user *User, db *MedDataBase) string {
	var sb strings.Builder
	for _, a := range db.Appointments {
		if a.UserName == user.Name {
			sb.WriteString(fmt.Sprintf("Appointment with Dr. %s at %s\n", a.DoctorName, a.TimeSlot))
		}
	}
	if sb.Len() == 0 {
		return "You have no appointments."
	}
	return sb.String()
}

type DefaultHandler struct{}

func (d DefaultHandler) CanHandle(input string) bool {
	return true
}

func (d DefaultHandler) Handle(query Query, user *User, db *MedDataBase) string {
	return "I didn't understand that. Try commands like: rate, book, pharmacy, list, or ask about symptoms like 'I have a fever'."
}

type Chatbot struct {
	Base MedDataBase
}

func (cb *Chatbot) GenerateResponse(query Query, user *User) Response {
	handlers := []CommandHandler{
		PharmacyHandler{},
		SymptomInfoHandler{},
		RateDoctorHandler{},
		SelectDoctorHandler{},
		BookAppointmentHandler{},
		AddToCartHandler{},
		RemoveFromCartHandler{},
		ViewCartHandler{},
		ViewAppointmentsHandler{},
		ListDoctorsHandler{},
		DefaultHandler{},
	}

	for _, handler := range handlers {
		if handler.CanHandle(query.Content) {
			return Response{
				Content: handler.Handle(query, user, &cb.Base),
				Time:    time.Now(),
			}
		}
	}

	return Response{Content: "Unexpected error.", Time: time.Now()}
}
