package medbot

import (
	"strings"
	"testing"
)

func setupBot() (*Chatbot, *User) {
	db := MedDataBase{
		Doctors: []Doctor{
			{ID: "D1", Name: "Dr. Smith", Special: "General", Location: "New York"},
			{ID: "D2", Name: "Dr. Life", Special: "Cardiology", Location: "San Francisco"},
			{ID: "D3", Name: "Dr. Bold", Special: "Pediatrics", Location: "Chicago"},
			{ID: "D4", Name: "Dr. Rose", Special: "Dermatology", Location: "Los Angeles"},
			{ID: "D5", Name: "Dr. Green", Special: "Neurology", Location: "Boston"},
		},
	}
	bot := &Chatbot{Base: db}
	user := &User{Name: "TestUser"}
	return bot, user
}

func TestRateDoctorHandler(t *testing.T) {
	bot, user := setupBot()
	resp := bot.GenerateResponse(Query{Content: "rate D1 5"}, user)
	if !strings.Contains(resp.Content, "Thank you for rating Dr. Smith") {
		t.Errorf("unexpected response: %s", resp.Content)
	}
}

func TestSelectDoctorHandler(t *testing.T) {
	bot, user := setupBot()
	resp := bot.GenerateResponse(Query{Content: "select doctor D2"}, user)
	if !strings.Contains(resp.Content, "Dr. Life") {
		t.Errorf("unexpected response: %s", resp.Content)
	}
}

func TestBookAppointmentHandler(t *testing.T) {
	bot, user := setupBot()
	resp := bot.GenerateResponse(Query{Content: "book D3 2PM"}, user)
	if !strings.Contains(resp.Content, "Appointment booked") {
		t.Errorf("unexpected response: %s", resp.Content)
	}
}

func TestAddToCartHandler(t *testing.T) {
	bot, user := setupBot()
	resp := bot.GenerateResponse(Query{Content: "add aspirin 2"}, user)
	if !strings.Contains(resp.Content, "Added aspirin x2") {
		t.Errorf("unexpected response: %s", resp.Content)
	}
}

func TestRemoveFromCartHandler(t *testing.T) {

	bot, user := setupBot()

	_ = bot.GenerateResponse(Query{Content: "buy aspirin 2"}, user)

	resp1 := bot.GenerateResponse(Query{Content: "remove 1 aspirin"}, user)
	if !strings.Contains(resp1.Content, "Removed 1 aspirin. Remaining: 1.") {
		t.Errorf("unexpected response after first removal: %s", resp1.Content)
	}

	resp2 := bot.GenerateResponse(Query{Content: "remove 1 aspirin"}, user)
	if !strings.Contains(resp2.Content, "Removed all of aspirin from your cart.") {
		t.Errorf("unexpected response after final removal: %s", resp2.Content)
	}
}

func TestViewCartHandler(t *testing.T) {
	bot, user := setupBot()
	_ = bot.GenerateResponse(Query{Content: "add ibuprofen 3"}, user)
	resp := bot.GenerateResponse(Query{Content: "view cart"}, user)
	if !strings.Contains(resp.Content, "ibuprofen x3") {
		t.Errorf("unexpected response: %s", resp.Content)
	}
}

func TestViewAppointmentsHandler(t *testing.T) {
	bot, user := setupBot()
	_ = bot.GenerateResponse(Query{Content: "book D1 10AM"}, user)
	resp := bot.GenerateResponse(Query{Content: "view appointments"}, user)
	if !strings.Contains(resp.Content, "Appointment with Dr. Smith") {
		t.Errorf("unexpected response: %s", resp.Content)
	}
}

func TestGetInfoHandler(t *testing.T) {
	bot, user := setupBot()
	resp := bot.GenerateResponse(Query{Content: "info headache"}, user)
	if !strings.Contains(resp.Content, "Use a hot or cold compress") {
		t.Errorf("unexpected response: %s", resp.Content)
	}
}

func TestListDoctorsHandler(t *testing.T) {
	bot, user := setupBot()
	resp := bot.GenerateResponse(Query{Content: "list doctors"}, user)
	if !strings.Contains(resp.Content, "Dr. Smith") || !strings.Contains(resp.Content, "General") {
		t.Errorf("unexpected response: %s", resp.Content)
	}
}

func TestUnknownCommandHandler(t *testing.T) {
	bot, user := setupBot()
	resp := bot.GenerateResponse(Query{Content: "xyz unknown"}, user)
	if !strings.Contains(resp.Content, "I didn't understand that") {
		t.Errorf("unexpected response: %s", resp.Content)
	}
}
