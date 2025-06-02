package medbot

import (
	"strings"
	"testing"
	"time"
)

// TestAddToCartAndViewCart tests adding items to the cart and viewing it
func TestAddToCartAndViewCart(t *testing.T) {
	user := &User{ID: 1, Name: "Test User"}

	med := Pharmacy{Med: "Ibuprofen", Price: 10, Available: 5}
	user.AddToCart(med, 1)

	cart := user.ViewCart()
	if !strings.Contains(cart, "Ibuprofen") {
		t.Errorf("Expected 'Ibuprofen' in cart, got: %s", cart)
	}
}

// TestRemoveFromCart test the removal functionality
func TestRemoveSpecificQuantityFromCart(t *testing.T) {
	user := &User{ID: 1, Name: "Test User"}
	user.AddToCart(Pharmacy{Med: "Ibuprofen", Price: 10}, 3)

	msg := user.RemoveFromCart("Ibuprofen", 2)
	if !strings.Contains(msg, "2 unit(s) of Ibuprofen removed") {
		t.Errorf("Unexpected message: %s", msg)
	}
	if user.Cart[0].Quantity != 1 {
		t.Errorf("Expected quantity to be 1, got: %d", user.Cart[0].Quantity)
	}
}

// TestCheckout tests the checkout functionality
func TestCheckout(t *testing.T) {
	user := &User{ID: 1, Name: "Test User"}
	user.AddToCart(Pharmacy{Med: "Paracetamol", Price: 5}, 3)

	checkout := user.Checkout()
	if !strings.Contains(checkout, "Your total is $5") {
		t.Errorf("Unexpected checkout message: %s", checkout)
	}
	if len(user.Cart) != 0 {
		t.Errorf("Cart should be empty after checkout")
	}
}

// TestGenerateResponseHelp tests the chatbot 'help' response
func TestGenerateResponseHelp(t *testing.T) {
	cb := &Chatbot{
		Name: "TestBot",
		Base: MedDataBase{},
	}
	user := &User{ID: 1, Name: "Tester"}
	query := Query{UserID: user.ID, Timestamp: time.Now(), Content: "help"}

	response := cb.GenerateResponse(query, user)

	if !strings.Contains(response.Content, "ask me about your symptoms") {
		t.Errorf("Expected help instructions in response, got: %s", response.Content)
	}
}

// TestGetInfo tests symptom advice retrieval
func TestGetInfo(t *testing.T) {
	db := MedDataBase{}
	info := db.GetInfo("headache")
	if !strings.Contains(info, "compress") {
		t.Errorf("Unexpected headache info: %s", info)
	}
}

// TestCheckPharmacyBuy tests buying medication from the pharmacy
func TestCheckPharmacyBuy(t *testing.T) {
	db := MedDataBase{}
	user := &User{ID: 1, Name: "Tester"}

	result := db.CheckPharmacy("buy Ibuprofen", user)

	if !strings.Contains(result, "added to your cart") {
		t.Errorf("Expected item to be added to cart, got: %s", result)
	}
	if len(user.Cart) == 0 || user.Cart[0].Med != "Ibuprofen" {
		t.Errorf("Item not correctly added to cart")
	}
}

// TestSelectDoctor tests doctor selection logic
func TestSelectDoctor(t *testing.T) {
	db := MedDataBase{}
	user := &User{ID: 1, Name: "Tester"}
	result := db.SelectDoctor("select 2", user)
	if !strings.Contains(result, "Dr. Life") {
		t.Errorf("Expected to select Dr. Life, got: %s", result)
	}
	if !strings.Contains(result, "Dr. Life") {
		t.Errorf("Expected to select Dr. Life, got: %s", result)
	}
}

// TestBookAppointment checks the booking output
func TestBookAppointment(t *testing.T) {
	db := &MedDataBase{}
	user := &User{ID: 1, Name: "Tester"}
	out := db.BookAppointment(user)
	if !strings.Contains(out, "Available date:") || !strings.Contains(out, "confirm <slot number>") {
		t.Errorf("BookAppointment output missing expected content: %s", out)
	}
}

// TestConfirmAppointment checks various confirmations
func TestConfirmAppointment(t *testing.T) {
	cb := &Chatbot{}
	user := &User{ID: 1, Name: "Tester"}
	cb.Base.BookAppointment(user)
	if !strings.Contains(cb.ConfirmAppointment("confirm 2", user), "confirmed") {
		t.Errorf("Expected confirmation for valid slot")
	}
	if !strings.Contains(cb.ConfirmAppointment("confirm 10", user), "Invalid slot") {
		t.Errorf("Expected error for invalid slot")
	}
	if !strings.Contains(cb.ConfirmAppointment("confirm", user), "specify the slot number") {
		t.Errorf("Expected prompt for missing slot number")
	}
}
func TestRateDoctor(t *testing.T) {
	db := &MedDataBase{}
	user := &User{ID: 1, Name: "Tester"}

	response := db.RateDoctor(user, "rate 5")
	expected := "Please select doctor before rate(e.g., 'select 2')."
	if response != expected {
		t.Errorf("Expected: %s, got: %s", expected, response)
	}

	user.SelectDoctor = &Doctor{Name: "Dr. Test", RatingList: []int{}}

	response = db.RateDoctor(user, "rate 4")
	if !strings.Contains(response, "You have rated 4 for Dr. Test") {
		t.Errorf("Unexpected response: %s", response)
	}

	response = db.RateDoctor(user, "rate")
	if !strings.Contains(response, "Please enter a rating") {
		t.Errorf("Unexpected response for missing rating: %s", response)
	}

	response = db.RateDoctor(user, "rate great")
	if !strings.Contains(response, "Invalid rating") {
		t.Errorf("Unexpected response for non-numeric: %s", response)
	}

	response = db.RateDoctor(user, "rate 6")
	if !strings.Contains(response, "Invalid rating") {
		t.Errorf("Unexpected response for out-of-range: %s", response)
	}

	db.RateDoctor(user, "rate 2")
	if user.SelectDoctor.Rating != "3.0" {
		t.Errorf("Expected average 3.0, got %s", user.SelectDoctor.Rating)
	}
}
