package poker

import "fmt"

type Rank int

const (
	Two Rank = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

// String returns the string representation of a Rank
func (r Rank) String() string {
	ranks := map[Rank]string{
		Two: "2", Three: "3", Four: "4", Five: "5", Six: "6",
		Seven: "7", Eight: "8", Nine: "9", Ten: "T",
		Jack: "J", Queen: "Q", King: "K", Ace: "A",
	}
	if str, ok := ranks[r]; ok {
		return str
	}
	return fmt.Sprintf("Unknown(%d)", r)
}

// ParseRank converts a byte to a Rank
func ParseRank(b byte) (Rank, error) {
	rankMap := map[byte]Rank{
		'2': Two, '3': Three, '4': Four, '5': Five, '6': Six,
		'7': Seven, '8': Eight, '9': Nine, 'T': Ten,
		'J': Jack, 'Q': Queen, 'K': King, 'A': Ace,
	}
	if rank, ok := rankMap[b]; ok {
		return rank, nil
	}
	return 0, fmt.Errorf("invalid rank: %c", b)
}
