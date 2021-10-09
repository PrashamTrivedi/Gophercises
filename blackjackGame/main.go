package main

import (
	"fmt"
	deck "prashamhtrivedi/blackjackLib"
	"strings"
)

type Hand []deck.Card

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i := range h {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}

func (h Hand) DealerString() string {
	return fmt.Sprintf("%s, ***HIDDEN***", h[0].String())
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func (h Hand) MinScore() int {
	score := 0
	for _, card := range h {
		score += min(int(card.Rank), 10)
	}
	return score
}

func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}
	for _, card := range h {
		if card.Rank == deck.Ace {
			//Ace is 1 right now, We are adding 10 more points to make it worth 11
			return minScore + 10
		}
	}
	return minScore
}

func Shuffle(gs GameState) GameState {
	ret := clone(gs)
	ret.Deck = deck.New(deck.Deck(3), deck.Shuffle)
	return ret
}

func Deal(gs GameState) GameState {
	ret := clone(gs)
	ret.Player = make(Hand, 0, 5)
	ret.Dealer = make(Hand, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, ret.Deck = draw(ret.Deck)
		ret.Player = append(ret.Player, card)

		card, ret.Deck = draw(ret.Deck)
		ret.Dealer = append(ret.Dealer, card)

	}
	ret.State = int(StatePlayerTurn)
	return ret
}

func Hit(gs GameState) GameState {
	ret := clone(gs)
	var card deck.Card
	hand := ret.CurrentPlayer()
	card, ret.Deck = draw(ret.Deck)
	*hand = append(*hand, card)
	if hand.Score() > 21 {
		return Stand(ret)
	}
	return ret
}

func Stand(gs GameState) GameState {
	ret := clone(gs)
	ret.State++
	return ret
}

func EndGame(gs GameState) GameState {
	ret := clone(gs)
	pScore, dScore := ret.Player.Score(), ret.Dealer.Score()
	fmt.Println("===FINAL HANDS===")
	fmt.Printf("Dealer: %s \nScore: %d\n", ret.Dealer, dScore)
	fmt.Printf("Player: %s \nScore: %d\n", ret.Player, pScore)

	switch {
	case pScore > 21:
		fmt.Println("You busted!!")
	case dScore > 21:
		fmt.Println("Dealer busted!!!")
	case pScore > dScore:
		fmt.Println("Yay!! You won")
	case dScore > pScore:
		fmt.Println("Dealer Won")
	case dScore == pScore:
		fmt.Println("Draw")
	}

	fmt.Println()
	ret.Player = nil
	ret.Dealer = nil
	return ret
}

func main() {

	var gs GameState

	gs = Shuffle(gs)

	gs = Deal(gs)

	var input string

	for gs.State == int(StatePlayerTurn) {

		fmt.Println("Dealer:", gs.Dealer.DealerString())
		fmt.Println("Player:", gs.Player)
		fmt.Println("What will you do? (h)it, (s)tand")
		fmt.Scanf("%s", &input)
		switch input {
		case "h":
			gs = Hit(gs)
		case "s":
			gs = Stand(gs)
			fmt.Println(gs.State)
		default:
			fmt.Println("Not a valid option.")
		}
	}

	for gs.State == int(StateDealerTurn) {
		if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
			gs = Hit(gs)
		} else {
			gs = Stand(gs)
		}
	}

	gs = EndGame(gs)

}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

type State int8

const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

type GameState struct {
	Deck   []deck.Card
	State  int
	Player Hand
	Dealer Hand
}

func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case int(StatePlayerTurn):
		return &gs.Player
	case int(StateDealerTurn):
		return &gs.Dealer
	default:
		panic("It isn't currently any player's turn")
	}
}

func clone(g GameState) GameState {
	gs := GameState{
		Deck:   make([]deck.Card, len(g.Deck)),
		Player: make(Hand, len(g.Player)),
		Dealer: make(Hand, len(g.Dealer)),
		State:  g.State,
	}
	copy(gs.Dealer, g.Dealer)
	copy(gs.Deck, g.Deck)
	copy(gs.Player, g.Player)
	return gs
}
