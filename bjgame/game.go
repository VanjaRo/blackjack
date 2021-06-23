package bjgame

import (
	"fmt"
	"time"

	"github.com/VanjaRo/deck_of_cards/deck"
)

type state uint8

const (
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

type Game struct {
	// unexported fields
	deck     []deck.Card
	state    state
	player   []deck.Card
	dealer   []deck.Card
	dealerAI Player
	humanAI  Player
	balance  int
}

func NewDeal(g *Game) {
	g.player = make([]deck.Card, 0, 9)
	g.dealer = make([]deck.Card, 0, 9)
	g.player = append(g.player, draw(&g.deck, 2)...)
	g.dealer = append(g.dealer, draw(&g.deck, 2)...)
	g.state = statePlayerTurn
}

func NewGame() Game {
	return Game{
		state:    statePlayerTurn,
		dealerAI: dealerPlayer{},
		humanAI:  humanPlayer{},
		balance:  0,
	}
}

func (g *Game) Play(p Player) int {

	g.deck = deck.New(deck.Decks(3), deck.Shuffle(time.Now().Unix()))
	for i := 0; i < 2; i++ {
		NewDeal(g)
		for g.state == statePlayerTurn {
			playerHand := make([]deck.Card, len(g.player))
			copy(playerHand, g.player)
			move := p.Play(playerHand, g.dealer[0])
			move(g)
		}

		for g.state == stateDealerTurn {
			dealerHand := make([]deck.Card, len(g.dealer))
			copy(dealerHand, g.dealer)
			move := g.dealerAI.Play(dealerHand, g.dealer[0])
			move(g)
		}

		EndHand(g)
	}
	return g.balance
}

func (g *Game) CurrentPlayer() *[]deck.Card {
	switch g.state {
	case statePlayerTurn:
		return &g.player
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("it isn't currently any player's turn")
	}
}

type Move func(*Game)

func MoveHit(g *Game) {
	hand := g.CurrentPlayer()
	*hand = append(*hand, draw(&g.deck, 1)...)
	if Score(*hand...) > 21 {
		MoveStand(g)
	}
}

func MoveStand(g *Game) {
	g.state++
}

func draw(cards *[]deck.Card, number int) []deck.Card {
	var card_sl []deck.Card
	card_sl = append(card_sl, (*cards)[:number]...)
	*cards = (*cards)[number:]
	return card_sl
}

func Score(hand ...deck.Card) int {
	score := minScore(hand...)
	if score > 11 {
		return score
	}
	for _, c := range hand {
		if c.Rank == deck.Ace {
			return score + 10
		}
	}
	return score
}

func Soft(n int, hand ...deck.Card) bool {
	hScore := Score(hand...)
	hMinScore := minScore(hand...)
	return hScore == n && hScore != hMinScore
}

func minScore(hand ...deck.Card) int {
	score := 0
	for _, c := range hand {
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

func EndHand(g *Game) {
	pScore, dScore := Score(g.player...), Score(g.dealer...)
	if pScore > 21 {
		pScore = 0
	}
	if dScore > 21 {
		dScore = 0
	}
	switch {
	case pScore > dScore || dScore == 0:
		fmt.Print("You win!")
		g.balance++
	case pScore < dScore || pScore == 0:
		fmt.Print("Dealer wins")
		g.balance--
	case pScore == dScore:
		fmt.Print("Draw")
	}
	g.humanAI.Result([][]deck.Card{g.player}, g.dealer)
	g.player = nil
	g.dealer = nil
}
