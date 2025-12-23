package poker

import (
	"fmt"
	"sort"
)

type HandRank int

const (
	HighCard HandRank = iota
	OnePair
	TwoPairs
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

// String returns the string representation of a HandRank
func (hr HandRank) String() string {
	names := []string{
		"High Card", "One Pair", "Two Pairs", "Three of a Kind",
		"Straight", "Flush", "Full House", "Four of a Kind",
		"Straight Flush", "Royal Flush",
	}
	if int(hr) < len(names) {
		return names[hr]
	}
	return fmt.Sprintf("Unknown(%d)", hr)
}

type Hand struct {
	Cards  []*Card
	Rank   HandRank
	Values []Rank // Kicker values for comparison
}

// NewHand creates a new Hand from a slice of cards
func NewHand(cards []*Card) *Hand {
	return &Hand{
		Cards: CopyCards(cards),
	}
}

// Evaluate analyzes the hand and sets its rank and values
func (h *Hand) Evaluate() {
	h.sortByRank()

	rankCounts := h.getRankCounts()
	suitCounts := h.getSuitCounts()

	isFlush := h.checkFlush(suitCounts)
	isStraight, straightHigh := h.checkStraight()

	// Royal Flush
	if isFlush && isStraight && straightHigh == Ace {
		h.Rank = RoyalFlush
		h.Values = []Rank{Ace}
		return
	}

	// Straight Flush
	if isFlush && isStraight {
		h.Rank = StraightFlush
		h.Values = []Rank{straightHigh}
		return
	}

	// Four of a Kind
	if rank, ok := h.findNOfAKind(rankCounts, 4); ok {
		h.Rank = FourOfAKind
		kickers := h.getKickers([]Rank{rank}, 1)
		h.Values = append([]Rank{rank}, kickers...)
		return
	}

	// Full House
	threeRank, hasThree := h.findNOfAKind(rankCounts, 3)
	pairRank, hasPair := h.findNOfAKind(rankCounts, 2)
	if hasThree && hasPair {
		h.Rank = FullHouse
		h.Values = []Rank{threeRank, pairRank}
		return
	}

	// Flush
	if isFlush {
		h.Rank = Flush
		h.Values = h.getSortedRanks()
		return
	}

	// Straight
	if isStraight {
		h.Rank = Straight
		h.Values = []Rank{straightHigh}
		return
	}

	// Three of a Kind
	if rank, ok := h.findNOfAKind(rankCounts, 3); ok {
		h.Rank = ThreeOfAKind
		kickers := h.getKickers([]Rank{rank}, 2)
		h.Values = append([]Rank{rank}, kickers...)
		return
	}

	// Two Pairs
	pairs := h.findAllPairs(rankCounts)
	if len(pairs) == 2 {
		sort.Slice(pairs, func(i, j int) bool { return pairs[i] > pairs[j] })
		h.Rank = TwoPairs
		kickers := h.getKickers(pairs, 1)
		h.Values = append(pairs, kickers...)
		return
	}

	// One Pair
	if rank, ok := h.findNOfAKind(rankCounts, 2); ok {
		h.Rank = OnePair
		kickers := h.getKickers([]Rank{rank}, 3)
		h.Values = append([]Rank{rank}, kickers...)
		return
	}

	// High Card
	h.Rank = HighCard
	h.Values = h.getSortedRanks()
}

// Compare compares this hand with another hand
// Returns: 1 if this hand wins, -1 if other hand wins, 0 if tie
func (h *Hand) Compare(other *Hand) int {
	if h.Rank > other.Rank {
		return 1
	}
	if h.Rank < other.Rank {
		return -1
	}

	minLen := len(h.Values)
	if len(other.Values) < minLen {
		minLen = len(other.Values)
	}

	for i := 0; i < minLen; i++ {
		if h.Values[i] > other.Values[i] {
			return 1
		}
		if h.Values[i] < other.Values[i] {
			return -1
		}
	}

	if len(h.Values) > len(other.Values) {
		return 1
	}
	if len(h.Values) < len(other.Values) {
		return -1
	}

	return 0
}

// sortByRank sorts the cards by rank in descending order
func (h *Hand) sortByRank() {
	sort.Slice(h.Cards, func(i, j int) bool {
		return h.Cards[i].Rank > h.Cards[j].Rank
	})
}

// getRankCounts returns a map of rank frequencies
func (h *Hand) getRankCounts() map[Rank]int {
	counts := make(map[Rank]int)
	for _, card := range h.Cards {
		counts[card.Rank]++
	}
	return counts
}

// getSuitCounts returns a map of suit frequencies
func (h *Hand) getSuitCounts() map[Suit]int {
	counts := make(map[Suit]int)
	for _, card := range h.Cards {
		counts[card.Suit]++
	}
	return counts
}

// checkFlush checks if all cards are of the same suit
func (h *Hand) checkFlush(suitCounts map[Suit]int) bool {
	for _, count := range suitCounts {
		if count == 5 {
			return true
		}
	}
	return false
}

// checkStraight checks if the cards form a straight
func (h *Hand) checkStraight() (bool, Rank) {
	if len(h.Cards) != 5 {
		return false, 0
	}

	// Ace-low straight (A-2-3-4-5)
	if h.Cards[0].Rank == Ace && h.Cards[1].Rank == Five &&
		h.Cards[2].Rank == Four && h.Cards[3].Rank == Three &&
		h.Cards[4].Rank == Two {
		return true, Five
	}

	// Normal straight
	for i := 0; i < 4; i++ {
		if h.Cards[i].Rank-h.Cards[i+1].Rank != 1 {
			return false, 0
		}
	}
	return true, h.Cards[0].Rank
}

// findNOfAKind finds the highest rank that appears exactly n times
func (h *Hand) findNOfAKind(rankCounts map[Rank]int, n int) (Rank, bool) {
	var best Rank
	found := false
	for rank, count := range rankCounts {
		if count == n {
			if !found || rank > best {
				best = rank
				found = true
			}
		}
	}
	return best, found
}

// findAllPairs finds all ranks that appear exactly twice
func (h *Hand) findAllPairs(rankCounts map[Rank]int) []Rank {
	var pairs []Rank
	for rank, count := range rankCounts {
		if count == 2 {
			pairs = append(pairs, rank)
		}
	}
	return pairs
}

// getKickers returns the n highest kickers excluding the given ranks
func (h *Hand) getKickers(exclude []Rank, n int) []Rank {
	excludeMap := make(map[Rank]bool)
	for _, r := range exclude {
		excludeMap[r] = true
	}

	var kickers []Rank
	seen := make(map[Rank]bool)

	for _, card := range h.Cards {
		if !excludeMap[card.Rank] && !seen[card.Rank] {
			kickers = append(kickers, card.Rank)
			seen[card.Rank] = true
		}
	}

	sort.Slice(kickers, func(i, j int) bool { return kickers[i] > kickers[j] })

	if len(kickers) > n {
		return kickers[:n]
	}
	return kickers
}

// getSortedRanks returns all ranks in descending order
func (h *Hand) getSortedRanks() []Rank {
	var ranks []Rank
	for _, card := range h.Cards {
		ranks = append(ranks, card.Rank)
	}
	sort.Slice(ranks, func(i, j int) bool { return ranks[i] > ranks[j] })
	return ranks
}

// GetBestHand finds the best 5-card hand from the given cards
func GetBestHand(cards []*Card) *Hand {
	if len(cards) <= 5 {
		hand := NewHand(cards)
		hand.Evaluate()
		return hand
	}

	combos := generateCombinations(cards, 5)
	if len(combos) == 0 {
		return &Hand{Rank: HighCard}
	}

	bestHand := NewHand(combos[0])
	bestHand.Evaluate()

	for i := 1; i < len(combos); i++ {
		hand := NewHand(combos[i])
		hand.Evaluate()
		if hand.Compare(bestHand) > 0 {
			bestHand = hand
		}
	}

	return bestHand
}

// generateCombinations generates all k-sized combinations from cards
func generateCombinations(cards []*Card, k int) [][]*Card {
	var result [][]*Card
	var current []*Card
	generateCombosHelper(cards, k, 0, current, &result)
	return result
}

func generateCombosHelper(cards []*Card, k, start int, current []*Card, result *[][]*Card) {
	if len(current) == k {
		combo := make([]*Card, k)
		copy(combo, current)
		*result = append(*result, combo)
		return
	}

	for i := start; i < len(cards); i++ {
		newCurrent := make([]*Card, len(current), len(current)+1)
		copy(newCurrent, current)
		newCurrent = append(newCurrent, cards[i])
		generateCombosHelper(cards, k, i+1, newCurrent, result)
	}
}
