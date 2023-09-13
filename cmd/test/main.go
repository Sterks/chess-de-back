package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/notnil/chess"
)

func main() {

	str := "1. Nf3 d5 2. g3 Bg4 3. b3 Nd7 4. Bb2 Nf3"
	pgnReader := strings.NewReader(str)
	pgn, err := chess.PGN(pgnReader)
	if err != nil {
		log.Println(err)
	}
	game := chess.NewGame(pgn)
	fmt.Println(game.Position().Board().Draw())
	// fmt.Println(game.Position().Board().Draw())
	// fmt.Println(steps)
	// // print outcome and game PGN
	// fmt.Println(game.Position().Board().Draw())
	// fmt.Printf("Game completed. %s by %s.\n", game.Outcome(), game.Method())
	// fmt.Println(game.String())
}
