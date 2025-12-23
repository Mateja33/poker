package poker

import (
	"fmt"
	"strings"
)

type Card struct {
	Suit Suit
	Rank Rank
}

// NewCard creates a new Card
func NewCard(suit Suit, rank Rank) *Card {
	return &Card{Suit: suit, Rank: rank}
}

// String returns the string representation of a Card
func (c *Card) String() string {
	return fmt.Sprintf("%s%s", c.Suit.String(), c.Rank.String())
}

// ParseCard parses a card string (e.g., "H7" for 7 of Hearts)
func ParseCard(s string) (*Card, error) {
	if len(s) != 2 {
		return nil, fmt.Errorf("invalid card string: %s", s)
	}

	suit, err := ParseSuit(s[0])
	if err != nil {
		return nil, err
	}

	rank, err := ParseRank(s[1])
	if err != nil {
		return nil, err
	}

	return &Card{Suit: suit, Rank: rank}, nil
}

// ParseCards parses multiple cards from a space-separated string
func ParseCards(s string) ([]*Card, error) {
	var cards []*Card
	for _, cardStr := range strings.Fields(s) {
		card, err := ParseCard(cardStr)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

// CopyCards creates a deep copy of a slice of cards
func CopyCards(cards []*Card) []*Card {
	result := make([]*Card, len(cards))
	for i, card := range cards {
		result[i] = &Card{Suit: card.Suit, Rank: card.Rank}
	}
	return result
}
