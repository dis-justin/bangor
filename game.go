package main

import (
	"fmt"
	"strconv"
	"unicode"
)

type Piece struct {
	Idx int
	Val int
}

var tiles [25]*Piece

type Player struct {
	Name   string
	Pieces []Piece
}

func (player *Player) HasPiece(piece *Piece) bool {

	//fmt.Printf("PIECES: %v\n\n", player.Pieces)
	for i := range player.Pieces {
		if player.Pieces[i].Idx == piece.Idx {

			//fmt.Printf("PLAYER %v: HAS PIECE PIECE: %v\n\n", player.Name, piece)
			return true
		}
	}
	return false
}

func (player *Player) Move(to *Piece, from *Piece) {
	// Update the board
	tiles[to.Idx].Val = from.Val
	tiles[from.Idx].Val = 0

	player.Pieces = append(player.Pieces, *to)

	for i := range player.Pieces {
		if player.Pieces[i].Idx == from.Idx {
			fmt.Print("hit")
			player.Pieces = append(player.Pieces[:i], player.Pieces[i:]...)
		}
	}
	fmt.Printf("PIECE: %v\n\n", player.Pieces)
}

func (player *Player) Play() {
	var move string

	fmt.Printf("Play your %d move: %v\n", 1, player.Name)
	fmt.Scanln(&move)

	boardPiece := *indexToPiece(coordToIndex(move))
	if player.HasPiece(&boardPiece) {
		fmt.Println("Options:\n1. Enter the coordinate to move\n2. Enter 'u' to Upgrade\n3. Enter 'c' to Cancel")

		var decision string
		fmt.Scanln(&decision)

		switch decision {
		case "u":
			upgrade(&boardPiece)
			break
		case "c":
			player.Play()
			break
		default:
			var newPiece = *indexToPiece(coordToIndex(decision))
			if newPiece.Idx != -1 {
				player.Move(&newPiece, &boardPiece)
			} else {
				invalidMove(*player)
			}
			break
		}

		displayBoard()
	} else {
		invalidMove(*player)
		player.Play()
	}
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
		player1.Play()
		evaluateWin(player1)

		player2.Play()
		evaluateWin(player2)
	}
}

func upgrade(piece *Piece) {
	tiles[piece.Idx].Val++
}

func evaluateWin(player Player) bool {
	return false
}

func initializeBoard() {
	for i := 0; i < 25; i++ {
		// Starting Tile values
		if i < 5 {
			piece := Piece{Idx: i, Val: 1}
			tiles[i] = &piece
			player1.Pieces = append(player1.Pieces, piece)
		} else if i >= 20 && i <= 24 {
			piece := Piece{Idx: i, Val: 1}
			tiles[i] = &piece
			player2.Pieces = append(player2.Pieces, piece)
		} else if i == 10 {
			piece := Piece{Idx: i, Val: 1}
			tiles[i] = &piece
			player1.Pieces = append(player1.Pieces, piece)
		} else if i == 15 {
			piece := Piece{Idx: i, Val: 1}
			tiles[i] = &piece
			player2.Pieces = append(player2.Pieces, Piece{Idx: i, Val: 1})
		} else {
			piece := Piece{Idx: i, Val: 0}
			tiles[i] = &piece
		}
	}
}

func initializePlayers() {
	player1.Name = "Justin"
	player2.Name = "David"
}

func displayBoard() {
	const red string = "\033[31m"
	const blue string = "\033[34m"
	const yellow string = "\033[33m"
	const reset string = "\033[0m"

	fmt.Printf("Welcome Player %v\nWelcome Player %v", player1.Name, player2.Name)

	fmt.Printf("%-4s\n", player1.Name)
	fmt.Printf("\n %s|A|B|C|D|E|\033[0m\n", yellow)

	var col int = 0
	for idx := range tiles {
		var prefix string = ""
		var suffix string = ""
		if idx%5 == 0 {
			col += 1
			prefix = yellow + strconv.Itoa(col) + reset
		}

		var nl string = ""
		if (idx+1)%5 == 0 && idx != 0 {
			nl = "|" + yellow + strconv.Itoa(col) + reset + "\n"
		}

		var color string = ""
		if player1.HasPiece(indexToPiece(idx)) {
			color = red
		} else if player2.HasPiece(indexToPiece(idx)) {
			color = blue
		}
		fmt.Printf("%s|"+"%s"+"%d"+"\033[0m"+"%s%s", prefix, color, tiles[idx].Val, nl, suffix)
	}

	fmt.Print(" \033[33m|A|B|C|D|E|\033[0m\n")
	fmt.Printf("\n%-4s\n\n", player2.Name)
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
			gridVal = (len(alpha) * rowIntVal) - (len(alpha) - indexOf(alpha, upper))
		}
	}

	return gridVal
}

func indexToPiece(idx int) *Piece {
	for i := range tiles {
		if tiles[i].Idx == idx {
			return tiles[i]
		}
	}
	return nil
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
	player.Play()
}
