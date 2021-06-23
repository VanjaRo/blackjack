package bjgame

import (
	"fmt"

	"github.com/VanjaRo/deck_of_cards/deck"
)

type Player interface {
	Bet() int
	Result(player [][]deck.Card, dealer []deck.Card)
	Play(player []deck.Card, dealer deck.Card) Move
}

func HumanPlayer() humanPlayer {
	return humanPlayer{}
}

type humanPlayer struct {
}

type dealerPlayer struct{}

func (dp dealerPlayer) Bet() int {
	return 1
}

func (dp dealerPlayer) Play(player []deck.Card, dealer deck.Card) Move {
	dScore := Score(player...)
	if dScore <= 16 || Soft(17, player...) {
		return MoveHit
	}
	return MoveStand

}

func (dp dealerPlayer) Result(player [][]deck.Card, dealer []deck.Card) {}

func (hp humanPlayer) Bet() int {
	return 1
}

func (hp humanPlayer) Play(player []deck.Card, dealer deck.Card) Move {
	for {
		fmt.Println("PLayer: ", player)
		fmt.Println("Dealer: ", dealer)
		fmt.Println("What will you do? (h)it/ (s)tand")
		var input string
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		default:
			fmt.Print("Invalid input: ", input)
		}
	}

}

func (hp humanPlayer) Result(player [][]deck.Card, dealer []deck.Card) {
	fmt.Println("\n===FINAL HANDS===")
	fmt.Println("PLayer: ", player)
	fmt.Println("Dealer: ", dealer)
}
