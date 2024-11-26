// Package deck provides a programmatic reconstruction of a deck of cards
// with utilites to manipulate and interact with it.
package deck

import (
	"fmt"
	"math/rand"
	"slices"
)

// Type Suit represents a card's suit. (e.g Heart)
//
//go:generate stringer -type=Suit
type Suit int8

// Type represents a card's type (e.g Ace)
//
//go:generate stringer -type=Type
type Type int8

// Type CardComparator wraps func(a, b Card) int.
// Used as a sort comparator to pass as an option when generating
// a new deck of cards. Returns a negative number when a < b, a positive
// number when a > b, and 0 otherwise. The default comparator DefaultComparator is already provided
// that sorts the deck conventionally.
type CardComparator func(a, b Card) int

// type OptionFunc acts a wrapper for functional
// options used for configuration in NewDeck
type OptionFunc func(*DeckOptions)

// Default comparator used in NewDeck
var DefaultComparator CardComparator = func(a, b Card) int {
	if a.Suit == b.Suit {
		return int(a.Type - b.Type)
	}
	return int(a.Suit - b.Suit)
}

// Card represents an instance of a card.
// Suit represents Hearts, Spades, Diamonds, Clubs, Jokers.
// Values represent numerical values of the card starting from
// Ace to King
type Card struct {
	Suit Suit
	Type Type
}

func (c Card) String() string {
	if c.Suit == JOKER {
		return "Joker"
	}
	return fmt.Sprintf("%s of %s", c.Suit, c.Type)
}

// Type DeckOptions is used to configure deck creation behaviour
// in NewDeck.
//
// shuffle: option to shuffle the deck. Note that shuffling is done last
// so any modifications done with numJoker and combineWith options will be shuffled as well.
// Enable with WithShuffle function
//
// numJokers: option used to add n amount of Jokers to the deck. NewDeck does not add any Jokers
// by default. Enable with WithJokers function
//
// filter: option that takes type func(Card) bool that returns true if the Card is to be kept and
// false otherwise. Enable with WithFilter function
//
// comparator: option that takes tyep CardComparator that is used to sorts the deck using the underlying
// function and is invoked before right before shuffling is done if its enabled to true. Enable with
// WithSort function
//
// combineWith: option used to append Card slice to the existing 52 card standard deck. Note that
// that combineWith is the first deck option exersized and will be effected by shuffle, filter, and
// comparator deck options. Enable with WithCombineDeck function
type DeckOptions struct {
	Shuffle     bool
	NumJokers   int
	Filter      func(Card) bool
	Comparator  CardComparator
	CombineWith []Card
}

// Option to shuffle the deck. Note that shuffling is done last
// so any modifications done with numJoker and combineWith options will be shuffled as well.
func WithShuffle() OptionFunc {
	return func(o *DeckOptions) {
		o.Shuffle = true
	}
}

// Option used to add n amount of Jokers to the deck. NewDeck does not add any Jokers
// by default.
func WithJokers(numJokers int) OptionFunc {
	return func(o *DeckOptions) {
		o.NumJokers = numJokers
	}
}

// Option that takes type func(Card) bool that returns true if the Card is to be kept and
// false otherwise.
func WithFilter(filterFunc func(Card) bool) OptionFunc {
	return func(o *DeckOptions) {
		o.Filter = filterFunc
	}
}

// Option that takes tyep CardComparator that is used to sorts the deck using the underlying
// function and is invoked before right before shuffling is done if its enabled to true.
func WithSort(cmp CardComparator) OptionFunc {
	return func(o *DeckOptions) {
		o.Comparator = cmp
	}
}

// Option used to append Card slice to the existing 52 card standard deck. Note that
// that combineWith is the first deck option exersized and will be effected by shuffle, filter, and
// comparator deck options.
func WithCombineDeck(deck []Card) OptionFunc {
	return func(o *DeckOptions) {
		o.CombineWith = deck
	}
}

const (
	SPADE Suit = iota
	DIAMOND
	CLUB
	HEART
	JOKER
)

const (
	NONE Type = iota
	ACE
	TWO
	THREE
	FOUR
	FIVE
	SIX
	SEVEN
	EIGHT
	NINE
	TEN
	JACK
	QUEEN
	KING
)

// NewDeck generates and returns a slice of cards that serves as a deck.
// By default NewDeck generates a standard 52 card deck sorted in conventional order
// from aces to kings grouped by suits. Functional options provided in the package are used to
// customize deck generation. See DeckOptions for more information.
func NewDeck(options ...OptionFunc) []Card {
	deckOptions := DeckOptions{
		Shuffle:     false,
		NumJokers:   0,
		Filter:      nil,
		Comparator:  DefaultComparator,
		CombineWith: nil,
	}
	for _, option := range options {
		option(&deckOptions)
	}

	deck := createStandardDeck()

	//combineWith
	if deckOptions.CombineWith != nil {
		deck = AddCards(deck, deckOptions.CombineWith...)
	}
	//numJokers
	for i := 0; i < deckOptions.NumJokers; i++ {
		deck = AddCards(deck, NewCard(JOKER, NONE))
	}
	//filter
	if deckOptions.Filter != nil {
		deck = slices.DeleteFunc(deck, deckOptions.Filter)
	}
	//comparators
	slices.SortStableFunc(deck, deckOptions.Comparator)
	//shuffle
	ShuffleDeck(deck, 3)

	return deck
}

// AddCards appends cards to the end of the deck.
func AddCards(deck []Card, cards ...Card) []Card {
	deck = append(deck, cards...)
	return deck
}

// NewCard creates new card defined by its Suit and Type.
func NewCard(suit Suit, cardType Type) Card {
	return Card{suit, cardType}
}

func createStandardDeck() []Card {
	deck := make([]Card, 0, 52)
	for i := SPADE; i <= HEART; i++ {
		for j := ACE; j <= KING; j++ {
			AddCards(deck, Card{i, j})
		}
	}
	return deck
}

// ShuffleDeck pseudo randomizes the deck n times.
func ShuffleDeck(deck []Card, n int) {
	for i := 0; i < n; i++ {
		rand.Shuffle(len(deck), func(i, j int) {
			deck[i], deck[j] = deck[j], deck[i]
		})
	}
}
