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
	user.AddToCart(med)

	cart := user.ViewCart()
	if !strings.Contains(cart, "Ibuprofen") {
		t.Errorf("Expected 'Ibuprofen' in cart, got: %s", cart)
	}
}

// TestRemoveFromCart test the removal functionality
func TestRemoveFromCart(t *testing.T) {
	user := &User{ID: 1, Name: "Test User"}
	user.AddToCart(Pharmacy{Med: "Ibuprofen", Price: 10})

	msg := user.RemoveFromCart("Ibubprofen")
	if !strings.Contains(msg, "removed from your cart") {
		t.Errorf("Expected removal confirmation, got: %s", msg)
	}

	if len(user.Cart) != 0 {
		t.Errorf("Expected cart to be empty after removal")
	}
}

// TestCheckout tests the checkout functionality
func TestCheckout(t *testing.T) {
	user := &User{ID: 1, Name: "Test User"}
	user.AddToCart(Pharmacy{Med: "Paracetamol", Price: 5})

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
	result := db.SelectDoctor("select 2")

	if !strings.Contains(result, "Dr. Life") {
		t.Errorf("Expected to select Dr. Life, got: %s", result)
	}
}
