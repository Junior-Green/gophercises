package blackjack

import (
	"testing"

	"github.com/Junior-Green/gophercises/deck"
)

func TestSoft17(t *testing.T) {

	t.Run("ACE + SIX", func(t *testing.T) {
		hand := []deck.Card{
			{Suit: deck.HEART, Type: deck.ACE},
			{Suit: deck.HEART, Type: deck.SIX},
		}

		h := Hand{Hand: hand}
		if h.Value() != 17 {
			t.Fatalf("expected 17, got %v", h.Value())
		}
	})

	t.Run("ACE + ACE + FIVE", func(t *testing.T) {
		hand := []deck.Card{
			{Suit: deck.HEART, Type: deck.ACE},
			{Suit: deck.DIAMOND, Type: deck.ACE},
			{Suit: deck.HEART, Type: deck.FIVE},
		}

		h := Hand{Hand: hand}
		if h.Value() != 17 {
			t.Fatalf("expected 17, got %v", h.Value())
		}
	})

	t.Run("ACE + TWO + FOUR", func(t *testing.T) {
		hand := []deck.Card{
			{Suit: deck.HEART, Type: deck.ACE},
			{Suit: deck.DIAMOND, Type: deck.TWO},
			{Suit: deck.HEART, Type: deck.FOUR},
		}

		h := Hand{Hand: hand}
		if h.Value() != 17 {
			t.Fatalf("expected 17, got %v", h.Value())
		}
	})
}

func TestBlackJack(t *testing.T) {
	hand := []deck.Card{
		{Suit: deck.HEART, Type: deck.ACE},
		{Suit: deck.SPADE, Type: deck.QUEEN},
		{Suit: deck.HEART, Type: deck.KING},
	}

	h := Hand{Hand: hand}
	if h.Value() != 21 {
		t.Fatalf("expected 21, got %v", h.Value())
	}
}

func TestNumberCards(t *testing.T) {
	t.Run("THREE + SIX + FIVE", func(t *testing.T) {
		hand := []deck.Card{
			{Suit: deck.HEART, Type: deck.THREE},
			{Suit: deck.HEART, Type: deck.SIX},
			{Suit: deck.HEART, Type: deck.FIVE},
		}

		h := Hand{Hand: hand}
		if h.Value() != 14 {
			t.Fatalf("expected 14, got %v", h.Value())
		}
	})

	t.Run("NINE + NINE + TWO", func(t *testing.T) {
		hand := []deck.Card{
			{Suit: deck.HEART, Type: deck.NINE},
			{Suit: deck.DIAMOND, Type: deck.NINE},
			{Suit: deck.DIAMOND, Type: deck.TWO},
		}

		h := Hand{Hand: hand}
		if h.Value() != 20 {
			t.Fatalf("expected 20, got %v", h.Value())
		}
	})

	t.Run("FIVE + EIGHT + THREE", func(t *testing.T) {
		hand := []deck.Card{
			{Suit: deck.HEART, Type: deck.FIVE},
			{Suit: deck.HEART, Type: deck.EIGHT},
			{Suit: deck.HEART, Type: deck.THREE},
		}

		h := Hand{Hand: hand}
		if h.Value() != 16 {
			t.Fatalf("expected 16, got %v", h.Value())
		}
	})
}

func TestTripleAce(t *testing.T) {
	hand := []deck.Card{
		{Suit: deck.HEART, Type: deck.ACE},
		{Suit: deck.DIAMOND, Type: deck.ACE},
		{Suit: deck.CLUB, Type: deck.ACE},
	}

	h := Hand{Hand: hand}
	if h.Value() != 13 {
		t.Fatalf("expected 13, got %v", h.Value())
	}
}

func TestFaceCards(t *testing.T) {
	hand := []deck.Card{
		{Suit: deck.HEART, Type: deck.KING},
		{Suit: deck.DIAMOND, Type: deck.QUEEN},
		{Suit: deck.CLUB, Type: deck.JACK},
	}

	h := Hand{Hand: hand}
	if h.Value() != 30 {
		t.Fatalf("expected 30, got %v", h.Value())
	}
}

func TestQuadAces(t *testing.T) {

	t.Run("4x ACE", func(t *testing.T) {
		hand := []deck.Card{
			{Suit: deck.HEART, Type: deck.ACE},
			{Suit: deck.DIAMOND, Type: deck.ACE},
			{Suit: deck.CLUB, Type: deck.ACE},
			{Suit: deck.SPADE, Type: deck.ACE},
		}

		h := Hand{Hand: hand}
		if h.Value() != 14 {
			t.Fatalf("expected 14, got %v", h.Value())
		}
	})

	t.Run("4x ACE + 2x TWO", func(t *testing.T) {
		hand := []deck.Card{
			{Suit: deck.HEART, Type: deck.ACE},
			{Suit: deck.DIAMOND, Type: deck.ACE},
			{Suit: deck.CLUB, Type: deck.ACE},
			{Suit: deck.SPADE, Type: deck.ACE},
			{Suit: deck.CLUB, Type: deck.TWO},
			{Suit: deck.SPADE, Type: deck.TWO},
		}

		h := Hand{Hand: hand}
		if h.Value() != 18 {
			t.Fatalf("expected 8, got %v", h.Value())
		}
	})
}
