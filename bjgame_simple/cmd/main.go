package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/VanjaRo/deck_of_cards/deck"
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
	return h[0].String() + ", ***HIDDEN***"
}

func (h Hand) MinScore() int {
	score := 0
	for _, c := range h {
		score += min(int(c.Rank), 10)
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (h Hand) Score() int {
	score := h.MinScore()
	if score > 11 {
		return score
	}
	for _, c := range h {
		if c.Rank == deck.Ace {
			return score + 10
		}
	}
	return score
}

func NewShuffledDeck(gs GameState) GameState {
	ret := clone(gs)
	ret.Deck = deck.New(deck.Decks(3), deck.Shuffle(time.Now().Unix()))
	return ret
}

func InitDeal(gs GameState) GameState {
	ret := clone(gs)
	ret.Player = make(Hand, 0, 9)
	ret.Dealer = make(Hand, 0, 9)
	ret.Player = append(ret.Player, draw(&ret.Deck, 2)...)
	ret.Dealer = append(ret.Dealer, draw(&ret.Deck, 2)...)
	return ret
}

func Hit(gs GameState) GameState {
	ret := clone(gs)
	hand := ret.CurrentPlayer()
	*hand = append(*hand, draw(&gs.Deck, 1)...)
	if hand.Score() > 21 {
		return Stand(gs)
	}
	return ret
}

func Stand(gs GameState) GameState {
	ret := clone(gs)
	ret.State++
	return ret
}

func EndHand(gs GameState) GameState {
	ret := clone(gs)
	pScore, dScore := gs.Player.Score(), gs.Dealer.Score()
	if pScore > 21 {
		pScore = 0
	}
	if dScore > 21 {
		dScore = 0
	}
	fmt.Println("===FINAL HANDS===")
	fmt.Println("PLayer: ", gs.Player, "\nScore", pScore)
	fmt.Println("Dealer: ", gs.Dealer, "\nScore", dScore)
	switch {
	case pScore > dScore || dScore == 0:
		fmt.Print("You win!")
	case pScore < dScore || pScore == 0:
		fmt.Print("Dealer wins")
	case pScore == dScore:
		fmt.Print("Draw")
	}
	ret.Player = nil
	ret.Dealer = nil
	return ret
}

func main() {
	var gs GameState
	gs = NewShuffledDeck(gs)
	gs = InitDeal(gs)

	var input string
	for gs.State == StatePlayerTurn {
		fmt.Println("PLayer: ", gs.Player)
		fmt.Println("Dealer: ", gs.Dealer.DealerString())
		fmt.Println("What will you do? (h)it/ (s)tand")
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			gs = Hit(gs)
		case "s":
			gs = Stand(gs)
		default:
			fmt.Print("Invalid input: ", input)
		}

	}

	for gs.State == StateDealerTurn {
		dScore := gs.Dealer.Score()
		dMinScore := gs.Dealer.MinScore()
		if dScore <= 16 || dScore == 17 && dMinScore != dScore {
			gs = Hit(gs)
		} else {
			gs = Stand(gs)
		}
	}

	gs = EndHand(gs)
}

func draw(cards *[]deck.Card, number int) []deck.Card {
	var card_sl []deck.Card
	card_sl = append(card_sl, (*cards)[:number]...)
	*cards = (*cards)[number:]
	return card_sl
}

type State uint8

const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

type GameState struct {
	Deck   []deck.Card
	State  State
	Player Hand
	Dealer Hand
}

func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Player
	case StateDealerTurn:
		return &gs.Dealer
	default:
		panic("it isn't currently any player's turn")
	}
}

func clone(gs GameState) GameState {
	ret := GameState{
		Deck:   make([]deck.Card, len(gs.Deck)),
		State:  gs.State,
		Player: make(Hand, len(gs.Player)),
		Dealer: make(Hand, len(gs.Dealer)),
	}
	copy(ret.Dealer, gs.Dealer)
	copy(ret.Deck, gs.Deck)
	copy(ret.Player, gs.Player)
	return ret
}
