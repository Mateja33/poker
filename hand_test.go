package poker

import (
	"testing"
)

func TestHandEvaluation(t *testing.T) {
	tests := []struct {
		name     string
		cards    string
		expected HandRank
	}{
		{"Royal Flush", "CT CJ CQ CK CA", RoyalFlush},
		{"Straight Flush (8-Q)", "D8 DQ DJ DT D9", StraightFlush},
		{"Straight Flush (7-J)", "H8 HT HJ H7 H9", StraightFlush},
		{"Four of a Kind (T)", "HT SQ ST DT CT", FourOfAKind},
		{"Four of a Kind (T) v2", "HT SK ST DT CT", FourOfAKind},
		{"Four of a Kind (8)", "H8 SQ S8 D8 C8", FourOfAKind},
		{"Four of a Kind (7)", "H7 SK S7 D7 C7", FourOfAKind},
		{"Full House (2-Q)", "H2 SQ C2 D2 CQ", FullHouse},
		{"Full House (2-J)", "H2 SJ C2 D2 CJ", FullHouse},
		{"Flush (K high)", "HK HQ H2 H4 H5", Flush},
		{"Flush (K high) v2", "D5 D4 D2 DQ DK", Flush},
		{"Straight (3-7)", "H3 S7 H5 D6 H4", Straight},
		{"Straight (7-J)", "C9 CT SJ D7 H8", Straight},
		{"Straight (A-5)", "H4 S5 HA D3 H2", Straight},
		{"Three of a Kind (2)", "H2 SQ S2 D2 CK", ThreeOfAKind},
		{"Three of a Kind (2) v2", "H2 S7 S2 D2 C9", ThreeOfAKind},
		{"Three of a Kind (2) v3", "H2 S8 S2 D2 C9", ThreeOfAKind},
		{"Two Pairs (5-T)", "H5 SQ C5 DT CT", TwoPairs},
		{"Two Pairs (9-T)", "H9 SQ C9 DT CT", TwoPairs},
		{"One Pair (8)", "H3 S8 H5 D8 CA", OnePair},
		{"One Pair (A)", "S4 DA H3 CA HT", OnePair},
		{"High Card (A high)", "H3 S8 H5 DK CA", HighCard},
		{"High Card (K high)", "H3 S8 H5 DK CT", HighCard},
		{"High Card (K high) v2", "H3 S8 H5 DK C2", HighCard},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cards, err := ParseCards(tt.cards)
			if err != nil {
				t.Fatalf("Failed to parse cards: %v", err)
			}

			hand := NewHand(cards)
			hand.Evaluate()

			if hand.Rank != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, hand.Rank)
			}
		})
	}
}

func TestHandComparison(t *testing.T) {
	tests := []struct {
		name      string
		category  string
		community string
		player1   string
		player2   string
		expected  string
		comment   string
	}{
		// High Card
		{"HC_1", "High Card", "D6 S9 H4 S3 C2", "SK CA", "HA SQ", "hand 1 > hand 2", "SK > SQ"},
		{"HC_2", "High Card", "D6 S9 H4 S3 C2", "SK CA", "HA CK", "hand 1 = hand 2", ""},
		{"HC_3", "High Card", "D6 S9 H4 H3 H4", "C7 DQ", "C8 DJ", "hand 1 > hand 2", "DQ > DJ"},

		// One Pair
		{"OP_1", "One Pair", "SK HT C8 C7 D2", "DK C5", "H8 D5", "hand 1 > hand 2", "K > 8"},
		{"OP_2", "One Pair", "SK HT C8 C7 D2", "DK C4", "HK C5", "hand 1 = hand 2", "K = K, same kickers"},
		{"OP_3", "One Pair", "HA DA ST C9 D4", "D5 C6", "H7 C2", "hand 2 > hand 1", "7 > 6"},

		// Two Pairs
		{"TP_1", "Two Pairs", "SA DQ CK D6 H6", "HA C3", "CQ H4", "hand 1 > hand 2", "A > Q"},
		{"TP_2", "Two Pairs", "SA DQ CK D6 H6", "HQ C3", "SQ H4", "hand 1 = hand 2", ""},
		{"TP_3", "Two Pairs", "SA DQ CK D6 H5", "HQ C6", "CA HK", "hand 2 > hand 1", "A > Q"},

		// Three of a Kind
		{"3K_1", "Three of a Kind", "SA D3 H2 C8 SJ", "HJ SJ", "C3 H3", "hand 1 > hand 2", "J > 3"},
		{"3K_2", "Three of a Kind", "SA D3 H3 C8 SJ", "C3 S2", "S3 H2", "hand 1 = hand 2", "3 = 3"},
		{"3K_3", "Three of a Kind", "HA SA DA H3 HT", "S2 S5", "H2 SK", "hand 2 > hand 1", "K > T"},

		// Straight
		{"ST_1", "Straight", "H3 S4 C5 S6 HT", "D7 HA", "H2 SA", "hand 1 > hand 2", "7 > 6"},
		{"ST_2", "Straight", "H3 S4 C5 S6 HT", "D7 HA", "H7 SA", "hand 1 = hand 2", "7 = 7"},
		{"ST_3", "Straight", "H2 H3 S4 C5 HT", "HA S3", "H6 SA", "hand 2 > hand 1", "6 > 5"},

		// Flush
		{"FL_1", "Flush", "D3 D6 DT C5 HQ", "DK DA", "D2 DQ", "hand 1 > hand 2", "A > Q"},
		{"FL_2", "Flush", "D3 D6 DT DJ DK", "C3 HA", "S9 HJ", "hand 1 = hand 2", "community flush"},
		{"FL_3", "Flush", "D3 D6 DT C5 HQ", "D2 D5", "DJ DA", "hand 2 > hand 1", "A > 5"},

		// Full House
		{"FH_1", "Full House", "HQ SQ HT DT C3", "DQ C2", "CT C4", "hand 1 > hand 2", "3xQ > 3xT"},
		{"FH_2", "Full House", "SA HQ SQ HT D8", "HA DQ", "DA CQ", "hand 1 = hand 2", ""},
		{"FH_3", "Full House", "HQ SQ HT DT C3", "ST C2", "CQ C4", "hand 2 > hand 1", "3xQ > 3xT"},

		// Four of a Kind
		{"4K_1", "Four of a Kind", "HT ST CT DT HK", "HA S7", "DJ C5", "hand 1 > hand 2", "A > K"},
		{"4K_2", "Four of a Kind", "S5 D5 C5 H5 HA", "CT HT", "C4 SQ", "hand 1 = hand 2", "community 4K"},
		{"4K_3", "Four of a Kind", "HT ST CT DT S8", "C2 C3", "C5 HK", "hand 2 > hand 1", "K > 8"},

		// Straight Flush
		{"SF_1", "Straight Flush", "H3 H4 H5 H6 HT", "H7 HA", "H2 SA", "hand 1 > hand 2", "7 > 6"},
		{"SF_2", "Straight Flush", "H3 H4 H5 H6 H7", "HA ST", "CQ D6", "hand 1 = hand 2", "community SF"},
		{"SF_3", "Straight Flush", "S7 S8 S9 ST DK", "S6 C2", "SJ D5", "hand 2 > hand 1", "J > T"},

		// Royal Flush
		{"RF_1", "Royal Flush", "DT DJ DQ DK DA", "C2 C3", "H2 H3", "hand 1 = hand 2", "community RF"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			communityCards, err := ParseCards(tt.community)
			if err != nil {
				t.Fatalf("Failed to parse community cards: %v", err)
			}

			player1Cards, err := ParseCards(tt.player1)
			if err != nil {
				t.Fatalf("Failed to parse player 1 cards: %v", err)
			}

			player2Cards, err := ParseCards(tt.player2)
			if err != nil {
				t.Fatalf("Failed to parse player 2 cards: %v", err)
			}

			allCards1 := append(CopyCards(communityCards), player1Cards...)
			allCards2 := append(CopyCards(communityCards), player2Cards...)

			hand1 := GetBestHand(allCards1)
			hand2 := GetBestHand(allCards2)

			result := hand1.Compare(hand2)
			var actual string
			if result > 0 {
				actual = "hand 1 > hand 2"
			} else if result < 0 {
				actual = "hand 2 > hand 1"
			} else {
				actual = "hand 1 = hand 2"
			}

			if actual != tt.expected {
				t.Errorf("\nCategory: %s\nComment: %s\nCommunity: %s\nPlayer 1: %s (%s)\nPlayer 2: %s (%s)\nExpected: %s\nGot: %s",
					tt.category, tt.comment, tt.community,
					tt.player1, hand1.Rank,
					tt.player2, hand2.Rank,
					tt.expected, actual)
			}
		})
	}
}

func BenchmarkHandEvaluation(b *testing.B) {
	cards, _ := ParseCards("CT CJ CQ CK CA")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hand := NewHand(cards)
		hand.Evaluate()
	}
}

func BenchmarkGetBestHand(b *testing.B) {
	cards, _ := ParseCards("CT CJ CQ CK CA H2 H3")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GetBestHand(cards)
	}
}
