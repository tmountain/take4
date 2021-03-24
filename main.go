package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Piece is encoded as an int.
type Piece int

// Board is a slice of Pieces.
type Board []Piece

// Standard board dimensions 7x6.
const (
	boardWidth     = 7
	boardHeight    = 6
	numSimulations = 1000
)

// Map names to Piece encodings.
const (
	EMPTY Piece = iota
	P1
	P2
)

// Helpful constants.
const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

// BoardDef represents board defaults at init time (mostly for testing).
type BoardDef map[int]Piece

// ValidMove represents moves that can be made.
type ValidMove []bool

// initBoard sets up the initial board state.
// Cells that are not empty are supplied via defs with the key
// representing the 1d index and the value representing the piece encoding.
func initBoard(defs BoardDef) Board {
	board := make(Board, boardWidth*boardHeight)
	for idx := range board {
		if defs[idx] != EMPTY {
			board[idx] = defs[idx]
		} else {
			board[idx] = EMPTY
		}
	}
	return board
}

// coordToIdx converts a cartesian coordinate to a 1d index.
func coordToIdx(x, y int) int {
	return x + (y * boardWidth)
}

// idxToCoord converts a 1d index to a pair of cartesian coordinates.
func idxToCoord(idx int) (int, int) {
	x := idx % boardWidth
	y := idx / boardWidth
	return x, y
}

// collectFour scans the board on the x, y, and x+y axis returning all
// valid combinations of adjacent cells as a slice of strings.
func collectFour(board Board) []string {
	vals := []string{}

	for idx := 0; idx < (boardWidth * boardHeight); idx++ {
		x, y := idxToCoord(idx)

		// Scan x-axis to the right.
		if x+3 < boardWidth {
			vals = append(vals, fmt.Sprintf("%v%v%v%v", board[idx], board[idx+1], board[idx+2], board[idx+3]))
		}

		// Scan y-axis down.
		if y+3 < boardHeight {
			vals = append(vals, fmt.Sprintf("%v%v%v%v", board[idx], board[coordToIdx(x, y+1)], board[coordToIdx(x, y+2)], board[coordToIdx(x, y+3)]))
		}

		// Scan diagonal down and right.
		if x+3 < boardWidth && y+3 < boardHeight {
			vals = append(vals, fmt.Sprintf("%v%v%v%v", board[idx], board[coordToIdx(x+1, y+1)], board[coordToIdx(x+2, y+2)], board[coordToIdx(x+3, y+3)]))
		}

		// Scan diagonal up and right.
		if x+3 < boardWidth && y-3 >= 0 {
			vals = append(vals, fmt.Sprintf("%v%v%v%v", board[idx], board[coordToIdx(x+1, y-1)], board[coordToIdx(x+2, y-2)], board[coordToIdx(x+3, y-3)]))
		}
	}

	return vals
}

// gameOver specifies if the game is over and who won (P1 or P2).
// If Piece is EMPTY, the game is ongoing.
func gameOver(board Board) Piece {
	groups := collectFour(board)
	p1Str := strconv.Itoa(int(P1))
	p2Str := strconv.Itoa(int(P2))

	for _, group := range groups {

		if group == strings.Repeat(p1Str, 4) {
			return P1
		}

		if group == strings.Repeat(p2Str, 4) {
			return P2
		}
	}

	return EMPTY
}

// validMoves returns a slice of valid moves.
func validMoves(board Board) []bool {
	valid := make([]bool, boardWidth)

	for col := 0; col < boardWidth; col++ {
		if board[col] == EMPTY {
			valid[col] = true
		}
	}
	return valid
}

// move places a piece on the board if the position is valid.
func move(board Board, move int, piece Piece) (Board, error) {
	valid := validMoves(board)

	if !valid[move] {
		return board, fmt.Errorf("move: invalid move specified")
	}

	for row := boardHeight - 1; row >= 0; row-- {
		idx := coordToIdx(move, row)

		if board[idx] == EMPTY {
			board[idx] = piece
			return board, nil
		}
	}

	// Should be unreachable.
	return board, fmt.Errorf("move: invalid code path")
}

// nextTurn specifies who goes next.
func nextTurn(currPlayer Piece) Piece {
	if currPlayer == P1 {
		return P2
	}
	return P1
}

// printBoard prints the board.
func printBoard(board Board) {
	bData := [][]string{
		{"1", "2", "3", "4", "5", "6", "7"},
	}

	for row := 0; row < boardHeight; row++ {
		rData := []string{}

		for col := 0; col < boardWidth; col++ {
			idx := coordToIdx(col, row)
			if board[idx] == EMPTY {
				rData = append(rData, "_")
			} else if board[idx] == P1 {
				rData = append(rData, "X")
			} else if board[idx] == P2 {
				rData = append(rData, "O")
			}
		}
		bData = append(bData, rData)
	}

	for _, row := range bData {
		fmt.Println(strings.Join(row, " "))
	}
}

// validInput returns a bool indicating whether the input is valid.
func validInput(input string) bool {
	valid := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	for _, entry := range valid {
		if input == entry {
			return true
		}
	}
	return false
}

func promptInput(player Piece) {
	playerStr := ""

	if player == P1 {
		playerStr = "P1"
	} else {
		playerStr = "P2"
	}

	fmt.Printf("[%v] Enter a move (1-7): ", playerStr)
}

// readInput takes input from stdin, validates it, and returns a Move.
func readInput() int {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		return -1
	}

	input = strings.TrimSpace(input)
	if validInput(input) {
		move, err := strconv.Atoi(input)

		// Should not be possible post-validation.
		if err != nil {
			fmt.Println(err)
			return -1
		}

		// Subtract one to convert to idx value.
		return move - 1
	}

	return -1
}

func endGame(board Board, piece Piece) {
	if piece == P1 {
		fmt.Println("Game over! P1 wins.")
	} else {
		fmt.Println("Game over! P2 wins.")
	}
	printBoard(board)
	os.Exit(0)
}

func filterValid(moves []bool) []int {
	ret := []int{}
	for i, v := range moves {
		if v {
			ret = append(ret, i)
		}
	}
	return ret
}

func getBestNextMove(board Board, currPlayer Piece) int {
	results := make(map[int]int)
	bestMove := -1
	highScore := MinInt

	for i := 0; i < numSimulations; i++ {
		player := currPlayer
		boardCopy := make(Board, len(board))
		copy(boardCopy, board)
		score := boardWidth * boardHeight
		nextMoves := filterValid(validMoves(boardCopy))
		firstMove := -1
		moveCnt := 1

		for len(nextMoves) > 0 {
			moveCnt++
			randMove := nextMoves[rand.Intn(len(nextMoves))]
			boardCopy, err := move(boardCopy, randMove, player)

			if err != nil {
				fmt.Println("Error making move:")
				fmt.Println(err)
			}

			// This move starts the tree.
			if firstMove == -1 {
				firstMove = randMove
			}

			if piece := gameOver(boardCopy); piece != EMPTY {
				// P1 is the human, so a win deducts points
				if piece == P1 {
					score *= -1
				}

				break
			}

			score -= 1
			player = nextTurn(player)
			nextMoves = filterValid(validMoves(boardCopy))
		}

		if firstMove >= 0 {
			results[firstMove] += score
		}
	}

	for move, score := range results {
		if score > highScore {
			bestMove = move
			highScore = score
		}
	}
	return bestMove
}

func main() {
	rand.Seed(time.Now().UnixNano())
	defs := BoardDef{}
	board := initBoard(defs)
	activePlayer := P1
	var moveInput int

	for {
		if piece := gameOver(board); piece != EMPTY {
			endGame(board, piece)
		}

		printBoard(board)

		if activePlayer == P1 {
			promptInput(activePlayer)
			moveInput = readInput()
		} else {
			moveInput = getBestNextMove(board, P2)
			fmt.Printf("CPU chooses %v\n", int(moveInput)+1)
		}

		if moveInput >= 0 {
			newBoard, err := move(board, moveInput, activePlayer)

			if err != nil {
				fmt.Println("Error making move:")
				fmt.Println(err)
			} else {
				activePlayer = nextTurn(activePlayer)
				board = newBoard
			}
		} else {
			fmt.Println("Invalid move.")
		}
	}
}
