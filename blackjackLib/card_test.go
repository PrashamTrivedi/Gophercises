package deck

import (
	"fmt"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Suit: Heart, Rank: King})
	fmt.Println(Card{Suit: Joker, Rank: King})
	fmt.Println(Card{Suit: Spade, Rank: King})
	fmt.Println(Card{Suit: Club, Rank: Seven})
	fmt.Println(Card{Suit: Heart, Rank: Ace})

	//Output:
	//King of Hearts
	//Joker
	//King of Spades
	//Seven of Clubs
	//Ace of Hearts

}

func TestNew(t *testing.T) {
	cards := New()
	if len(cards) != 13*4 {
		t.Error("Wrong number of card in new deck")

	}

}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)

	exp := Card{Rank: Ace, Suit: Spade}
	if cards[0] != exp {
		t.Errorf("Sorting didn't work, expected: %s, got: %s", exp.String(), cards[0].String())
	}
}

func TestDefaultAsCustomSort(t *testing.T) {
	cards := New(Sort(Less))
	exp := Card{Rank: Ace, Suit: Spade}
	if cards[0] != exp {
		t.Errorf("Sorting didn't work, expected: %s, got: %s", exp.String(), cards[0].String())
	}
}
func TestCustomSort(t *testing.T) {
	less := func(cards []Card) func(i, j int) bool {
		return func(i, j int) bool {
			return cards[i].Suit < cards[j].Suit
		}
	}

	cards := New(Sort(less))

	expectedSuite := Diamond

	if cards[14].Suit != expectedSuite {
		t.Errorf("Sorting didn't work, expected %s, got %s", expectedSuite.String(), cards[14].Suit.String())
	}

}

func TestJoker(t *testing.T) {
	expectedNumberOfJokers := 5
	cards := New(Jokers(expectedNumberOfJokers))
	count := 0
	for _, c := range cards {
		if c.Suit == Joker {
			count++
		}
	}

	if count != expectedNumberOfJokers {
		t.Errorf("Expected %d Jokers, but got %d Jokers", expectedNumberOfJokers, count)
	}

}

func TestFilter(t *testing.T) {
	filter := func(card Card) bool {
		return card.Rank == Two || card.Rank == Three
	}
	cards := New(Filter(filter))
	for _, card := range cards {
		if card.Rank == Two || card.Rank == Three {
			t.Error("Expected all twos and threes to filter out")
		}
	}
}

func TestDeck(t *testing.T) {
	cards := New(Deck(3))
	if len(cards) != 13*4*3 {
		t.Errorf("Expected %d cards, got %d cards", 13*4*3, len(cards))
	}
}
