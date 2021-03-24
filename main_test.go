package main

import (
	"fmt"
	"testing"
)

func TestCoordToIdx(t *testing.T) {
	x, y := idxToCoord(0)
	if x != 0 || y != 0 {
		t.Fatalf("expected x=0, y=0 but got x=%v, y=%v", x, y)
	}

	x, y = idxToCoord(6)
	if x != 6 || y != 0 {
		t.Fatalf("expected x=6, y=0, but got x=%v, y=%v", x, y)
	}

	x, y = idxToCoord(7)
	if x != 0 || y != 1 {
		t.Fatalf("expected x=0, y=1, but got x=%v, y=%v", x, y)
	}
}

func TestIdxToCoord(t *testing.T) {
	idx := coordToIdx(0, 0)
	if idx != 0 {
		t.Fatalf("expected idx=0, but got idx=%v", idx)
	}

	idx = coordToIdx(6, 0)
	if idx != 6 {
		t.Fatalf("expected idx=6, but got idx=%v", idx)
	}

	idx = coordToIdx(0, 1)
	if idx != 7 {
		t.Fatalf("expected idx=7, but got idx=%v", idx)
	}
}

func TestInitBoard(t *testing.T) {
	defs := BoardDef{1: P1, 2: P2}
	board := initBoard(defs)

	if board[1] != P1 {
		t.Fatalf("expected board[1]=P1, but got board[1]=%v", board[1])
	}

	if board[2] != P2 {
		t.Fatalf("expected board[2]=P2, but got board[2]=%v", board[2])
	}
}

func TestCollectFour(t *testing.T) {
	defs := BoardDef{1: P1, 2: P2, 5: P1, 6: P1, 7: P2, 8: P2}
	board := initBoard(defs)
	groups := collectFour(board)
	wanted := []string{
		"0120",
		"0200",
		"0200",
		"1200",
		"1200",
		"1000",
		"2001",
		"2000",
		"2000",
		"0011",
	}

	for idx, val := range wanted {
		if groups[idx] != val {
			t.Fatalf("expected groups[%v]=%v, but got groups[%v]=%v", idx, val, idx, groups[idx])
		}
	}

	// Test up-right diagonal.
	if groups[50] != "0001" {
		t.Fatalf("expected groups[50]=0001, but got groups[50]=%v", groups[50])
	}
}

func TestValidMoves(t *testing.T) {
	defs := BoardDef{1: P1, 2: P2, 5: P1, 6: P1}
	board := initBoard(defs)
	valid := validMoves(board)

	if !valid[0] {
		t.Fatalf("expected 0 valid, but invalid")
	}

	if valid[1] {
		t.Fatalf("expected 1 invalid, but valid")
	}

	if valid[2] {
		t.Fatalf("expected 2 invalid, but valid")
	}

	if !valid[3] {
		t.Fatalf("expected 3 valid, but invalid")
	}

	if !valid[4] {
		t.Fatalf("expected 4 valid, but invalid")
	}

	if valid[5] {
		t.Fatalf("expected 5 invalid, but valid")
	}

	if valid[6] {
		t.Fatalf("expected 6 invalid, but valid")
	}
}

func TestGameOver(t *testing.T) {
	defs := BoardDef{}
	board := initBoard(defs)
	board, _ = move(board, 0, P1)

	if piece := gameOver(board); piece != EMPTY {
		t.Fatalf("expected EMPTY, but got %v", piece)
	}
	board, _ = move(board, 0, P1)
	board, _ = move(board, 0, P1)
	board, _ = move(board, 0, P1)

	if piece := gameOver(board); piece != P1 {
		printBoard(board)
		groups := collectFour(board)

		for _, group := range groups {
			fmt.Println(group)
		}

		t.Fatalf("expected P1, but got %v", piece)
	}
}

func TestMove(t *testing.T) {
	defs := BoardDef{}
	board := initBoard(defs)
	// Drop a piece in column 1.
	board, err := move(board, 0, P1)

	if err != nil {
		t.Fatalf("expected valid, but got error")
	}

	if board[coordToIdx(0, 5)] != P1 {
		t.Fatalf("expected P1, but got %v", board[coordToIdx(0, 5)])
	}

	// Drop a piece in column 1.
	board, err = move(board, 0, P1)

	if err != nil {
		t.Fatalf("expected valid, but got error")
	}

	if board[coordToIdx(0, 4)] != P1 {
		t.Fatalf("expected P1, but got %v", board[coordToIdx(0, 4)])
	}

	// Fill column 1.
	board, _ = move(board, 0, P1)
	board, _ = move(board, 0, P1)
	board, _ = move(board, 0, P1)
	board, _ = move(board, 0, P1)
	// Ensure invalid.
	valid := validMoves(board)

	if valid[0] {
		t.Fatalf("expected 0 invalid, but got valid")
	}

	// Make sure a column won't overflow.
	_, err = move(board, 0, P1)
	if err == nil {
		t.Fatalf("expecting invalid making move, but got valid")
	}

	// Make sure both piece types work.
	board, err = move(board, 1, P2)
	if err != nil {
		t.Fatalf("expected valid, but got error")
	}

	if board[coordToIdx(1, 5)] != P2 {
		t.Fatalf("expected P2, but got %v", board[coordToIdx(1, 5)])
	}
	// Test outer boundary.
	board, err = move(board, 6, P2)
	if err != nil {
		t.Fatalf("expected valid, but got error")
	}

	if board[coordToIdx(6, 5)] != P2 {
		t.Fatalf("expected P2, but got %v", board[coordToIdx(6, 5)])
	}
}
