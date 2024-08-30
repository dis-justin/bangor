package main

import (
	"fmt"
	"strconv"
	"unicode"
)

var tiles [25]int

type Player struct {
	name   string
	pieces [5]int
}

var player1 Player
var player2 Player

func main() {
	initializeBoard()
	initializePlayers()
	displayBoard()
	runGameLoop()
}

func runGameLoop() {
	var winner bool = false
	for winner == false {
		playRound(player1)
		evaluateWin(player1)

		playRound(player2)
		evaluateWin(player2)
	}
}

func playRound(player Player) {
	var move string

	fmt.Printf("Play your %d move: %v\n", 1, player.name)
	fmt.Scanln(&move)

	var boardCoord int = coordToIndex(move)
	if hasPiece(player, boardCoord) {
		fmt.Println("Options:\n1. Enter the coordinate to move\n2. Enter 'u' to Upgrade\n3. Enter 'c' to Cancel")

		var decision string
		fmt.Scanln(&decision)

		switch decision {
		case "u":
			upgrade(boardCoord)
			break
		case "c":
			playRound(player)
			break
		default:
			var newCoord = coordToIndex(decision)
			if newCoord != -1 {
				moveTile(newCoord, boardCoord, player)
			} else {
				invalidMove(player)
			}
			break
		}

		displayBoard()
	} else {
		invalidMove(player)
		playRound(player)
	}
}

func moveTile(to int, from int, player Player) {
	tiles[to] = tiles[from]
	tiles[from] = 0

	//player.pieces[indexOfTile(player.pieces, from)] = to
}

func upgrade(coord int) {
	tiles[coord]++
}

func evaluateWin(player Player) bool {
	return false
}

func initializeBoard() {
	for i := 0; i < 25; i++ {
		// Starting Tile values
		if i < 5 {
			tiles[i] = 1
			player1.pieces[i] = i
		} else if i >= 20 && i <= 24 {
			tiles[i] = 1
			player2.pieces[i%5] = i
		} else {
			tiles[i] = 0
		}
	}
}

func initializePlayers() {
	player1.name = "Justin"
	player2.name = "David"
}

func displayBoard() {
	fmt.Printf("Welcome Player %v\n", player1.name)
	fmt.Printf("Welcome Player %v\n", player2.name)

	fmt.Printf("   %s\n", player1.name)

	fmt.Print("\n |A|B|C|D|E|\n")

	var col int = 0
	for idx := range tiles {
		var prefix string = ""
		var suffix string = ""
		if idx%5 == 0 {
			col += 1
			prefix = strconv.Itoa(col)
		}

		var nl string = ""
		if (idx+1)%5 == 0 && idx != 0 {
			nl = "|" + strconv.Itoa(col) + "\n"
		}
		fmt.Printf("%s|%d%s%s", prefix, tiles[idx], nl, suffix)
	}

	fmt.Print(" |A|B|C|D|E|\n")
	fmt.Printf("\n\n   %s\n\n", player2.name)
}

func coordToIndex(move string) int {
	var alpha = []byte{65, 66, 67, 68, 69}
	col := move[0]
	row := move[1]

	var gridVal int = 0

	upper := byte(unicode.ToUpper(rune(col)))
	rowIntVal, _ := strconv.Atoi(string(row))

	for idx := range alpha {
		if upper == alpha[idx] {
			gridVal = rowIntVal*(len(alpha)-indexOf(alpha, upper)) - len(alpha)
		}
	}

	return gridVal
}

func hasPiece(player Player, gridVal int) bool {
	for i := range player.pieces {
		if player.pieces[i] == gridVal {
			return true
		}
	}
	return false
}

func indexOf(slice []byte, value byte) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}

func indexOfTile(slice []int, value int) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}

func invalidMove(player Player) {
	fmt.Print("Invalid move, try again...\n\n")
	playRound(player)
}
