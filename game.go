package main

import (
	"fmt"
	"strconv"
	"unicode"
)

type Player struct {
	Name   string
	Pieces []Piece
	Pos    int
}

var player1 Player
var player2 Player

type Piece struct {
	Idx int
	Val int
}

func (piece *Piece) Upgrade() {
	tiles[piece.Idx].Val++
}

func (piece *Piece) IsOnBackline(player *Player) bool {
	if player == &player1 {
		return piece.Idx >= 0 && piece.Idx <= 4
	} else if player == &player2 {
		return piece.Idx >= 20 && piece.Idx <= 25
	}
	return false
}

var tiles [25]*Piece

func (player *Player) HasPiece(piece *Piece) bool {
	for i := range player.Pieces {
		if player.Pieces[i].Idx == piece.Idx {
			return true
		}
	}
	return false
}

func (player *Player) Move(to *Piece, from *Piece) {

	// Battle time...
	otherPlayer := OtherPlayer(player)
	if otherPlayer.HasPiece(to) {
		player.Battle(otherPlayer, to, from)
	} else {
		// Update the board
		tiles[to.Idx].Val += from.Val
		tiles[from.Idx].Val = 0

		to.Val = tiles[to.Idx].Val
		player.AddPiece(to)
		player.Pieces = player.RemovePiece(from)
	}

}

func (player *Player) RemovePiece(piece *Piece) []Piece {
	for i := range player.Pieces {
		if player.Pieces[i].Idx == *&piece.Idx {
			return player.Pieces[:i+copy(player.Pieces[i:], player.Pieces[i+1:])]
		}
	}
	return player.Pieces
}

func (player *Player) AddPiece(piece *Piece) {
	player.Pieces = append(player.Pieces, *piece)
	tiles[piece.Idx] = piece
}

func (attacker *Player) Battle(defender *Player, to *Piece, from *Piece) {
	tTo := to.Val
	tFrom := from.Val

	if to.Val-tFrom < 0 {
		to.Val = 0
	} else {
		to.Val -= tFrom
	}

	if from.Val-tTo < 0 {
		from.Val = 0
	} else {
		from.Val -= tTo
	}

	tiles[to.Idx].Val = to.Val
	tiles[from.Idx].Val = from.Val

	if to.Val <= 0 {
		defender.Pieces = defender.RemovePiece(to)
	}
	if from.Val <= 0 {
		attacker.Pieces = attacker.RemovePiece(from)
	}
}

func (player *Player) Combine(to *Piece, from *Piece) {
	to.Val += from.Val
	from.Val = 0

	tiles[to.Idx] = to
	tiles[from.Idx] = from
}

func (player *Player) Play() {
	var move string

	fmt.Printf("Play your %d move: %v\n", 1, player.Name)
	fmt.Scanln(&move)

	var idx int = coordToIndex(move)
	var boardPiece Piece
	if idx > -1 {
		boardPiece = *indexToPiece(idx)
	}

	if player.HasPiece(&boardPiece) {
		fmt.Println("Options:\n1. Enter the coordinate to move\n2. Enter 'u' to Upgrade\n3. Enter 'c' to Cancel")

		var decision string
		fmt.Scanln(&decision)

		switch decision {
		case "u":
			boardPiece.Upgrade()
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
	} else if boardPiece.IsOnBackline(player) {
		boardPiece.Val = 1
		player.AddPiece(&boardPiece)
	} else {
		invalidMove(*player)
	}

	displayBoard()
}

func (player *Player) Welcome() {
	fmt.Printf("Welcome Player %v!\n", player.Name)
}

func OtherPlayer(player *Player) *Player {
	if player == &player1 {
		return &player2
	} else if player == &player2 {
		return &player1
	} else {
		return nil // Who could it be Shaggy?
	}
}

func main() {
	initializeBoard()
	initializePlayers()

	player1.Welcome()
	player2.Welcome()

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

func evaluateWin(player Player) bool {
	return false
}

func initializeBoard() {
	for i := 0; i < 25; i++ {
		// Starting Tile values
		if i < 4 {
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

	var gridVal int = -1

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
