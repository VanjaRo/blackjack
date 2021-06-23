package main

import (
	"fmt"

	"github.com/VanjaRo/blackjack_go/bjgame"
)

func main() {
	game := bjgame.NewGame()
	winnings := game.Play(bjgame.HumanPlayer())
	fmt.Print(winnings)
}
