package poker

import "fmt"

type Suit int

const (
	Hearts Suit = iota
	Spades
	Clubs
	Diamonds
)

// String returns the string representation of a Suit
func (s Suit) String() string {
	suits := map[Suit]string{
		Hearts:   "H",
		Spades:   "S",
		Clubs:    "C",
		Diamonds: "D",
	}
	if str, ok := suits[s]; ok {
		return str
	}
	return fmt.Sprintf("Unknown(%d)", s)
}

// ParseSuit converts a byte to a Suit
func ParseSuit(b byte) (Suit, error) {
	suitMap := map[byte]Suit{
		'H': Hearts,
		'S': Spades,
		'C': Clubs,
		'D': Diamonds,
	}
	if suit, ok := suitMap[b]; ok {
		return suit, nil
	}
	return 0, fmt.Errorf("invalid suit: %c", b)
}
