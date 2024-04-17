package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Board represents a 9x9 Sudoku board.
type Board [9][9]int

type BoardRequest struct {
	Board Board `json:"board"`
}

var field Board
var webMode bool

func init() {
	flag.BoolVar(&webMode, "web", false, "set to true to run as web server")
}

func main() {
	flag.Parse()

	if webMode {
		runWebServer()
	} else {
		ArgsFromTerminal()
	}
}

func ArgsFromTerminal() {

	processArgs()

	preserveOriginal := field

	diagonalField := invertDiagonally(field)

	horizonField := invertHorizontally(field)

	verticalField := invertVertically(field)
	// drawBoard(field)

	if !solve(&field, 0, 0) {
		fmt.Println("Error")
		return
	}

	solve(&diagonalField, 0, 0)

	solve(&horizonField, 0, 0)

	solve(&verticalField, 0, 0)

	// if !solve(&verticalField, 0, 0) {
	// 	onlyOne = true
	// }

	// Initialize solvedBoard with the solved state
	var solvedBoard Board = field

	boardToCheck := preserveOriginal

	// drawBoard(field)

	hasAnotherSolution, secondSolution := isThereAnotherSolution(&boardToCheck, &solvedBoard, 0, 0)

	if hasAnotherSolution {
		// fmt.Println("Another solution exists:")
		drawBoard(*secondSolution)
		fmt.Println("Error")
		_ = secondSolution
		return
	} else {
		_ = secondSolution
		// fmt.Println("No other solutions found.")
	}

	horizonField = invertHorizontally(horizonField)

	diagonalField = invertDiagonally(diagonalField)

	verticalField = invertVertically(verticalField)

	if !equalBoards(field, horizonField) || !equalBoards(field, diagonalField) || !equalBoards(field, verticalField) {
		fmt.Println("Error")
		return
	}

	drawBoard(field)
}

func runWebServer() {
	router := mux.NewRouter()
	router.HandleFunc("/solve", solveHTTPHandler).Methods("POST")

	corsObj := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowCredentials(),
	)

	http.Handle("/", corsObj(router))
	log.Println("Starting web server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func solveHTTPHandler(w http.ResponseWriter, r *http.Request) {
	var inputBoard Board
	var input BoardRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	inputBoard = input.Board

	if solve(&inputBoard, 0, 0) {
		w.Header().Set("Content-Type", "application/json")
		response := struct {
			SolvedBoard Board `json:"board"`
		}{
			SolvedBoard: inputBoard,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "No solution found", http.StatusUnprocessableEntity)
	}
}

func processArgs() {

	// there are 2 possibilities, either we have 1 arg passed or 9.
	if len(os.Args) != 2 && len(os.Args) != 10 {
		fmt.Println("Error")
		fmt.Println("Usage: go run main.go [--web] or go run main.go <board>")
		os.Exit(0)
	}

	// for edge case handling where whole board is the second argument
	if len(os.Args) == 2 {
		args := strings.Split(os.Args[1], " ")
		for i, arg := range args {
			args[i] = strings.ReplaceAll(arg, "\"", "")
		}
		insertToMatrix(args)
	}

	if len(os.Args) == 10 {
		insertToMatrix(os.Args[1:])
	}
}

func insertToMatrix(args []string) {
	for _, arg := range args {
		if len(arg) != 9 {
			fmt.Println("Error")
			os.Exit(0)
		}
	}

	for i, arg := range args {
		replaceDotBy0(&arg)
		for j, ch := range arg {
			num, err := strconv.Atoi(string(ch))
			if err != nil {
				fmt.Printf("Error: Invalid character '%c' in input.\n", ch)
				os.Exit(2)
			}
			field[i][j] = num
		}
	}

	// check if any of the digits from 1-9 are in the same row repeated twice
	for i := 0; i < 9; i++ {
		digitsFound := make(map[int]bool)
		for j := 0; j < 9; j++ {
			if field[i][j] != 0 {
				if digitsFound[field[i][j]] {
					fmt.Println("Error")
					os.Exit(0)
				}
				digitsFound[field[i][j]] = true
			}
		}
	}
}

func replaceDotBy0(s *string) {
	*s = strings.ReplaceAll(*s, ".", "0")
}

func drawBoard(board Board) {
	for _, row := range board {
		for _, num := range row {
			fmt.Printf("%d ", num)
		}
		fmt.Println()
	}
}

func solve(board *Board, x, y int) bool {
	if y >= 9 {
		return true
	}

	nextX, nextY := next(x, y)
	if board[y][x] != 0 {
		return solve(board, nextX, nextY)
	}

	for v := 1; v <= 9; v++ {
		if canPut(*board, x, y, v) {
			board[y][x] = v
			if solve(board, nextX, nextY) {
				return true
			}
			board[y][x] = 0
		}
	}
	return false
}

func canPut(board Board, x, y, value int) bool {
	return !alreadyInRow(board, y, value) && !alreadyInColumn(board, x, value) && !alreadyInSquare(board, x, y, value)
}

func alreadyInRow(board Board, y, value int) bool {
	for _, val := range board[y] {
		if val == value {
			return true
		}
	}
	return false
}

func alreadyInColumn(board Board, x, value int) bool {
	for _, row := range board {
		if row[x] == value {
			return true
		}
	}
	return false
}

func alreadyInSquare(board Board, x, y, value int) bool {
	startRow, startCol := y/3*3, x/3*3
	for i := startRow; i < startRow+3; i++ {
		for j := startCol; j < startCol+3; j++ {
			if board[i][j] == value {
				return true
			}
		}
	}
	return false
}

func next(x, y int) (int, int) {
	nextX := (x + 1) % 9
	nextY := y
	if nextX == 0 {
		nextY++
	}
	return nextX, nextY
}

// this one inverts the board diagonally the board
func invertDiagonally(board Board) Board {
	invField := Board{}
	for i, row := range board {
		for j := range row {
			invField[8-i][8-j] = board[i][j]
		}
	}
	return invField
}

// this one inverts the board horizontally
func invertHorizontally(board Board) Board {
	invField := Board{}
	for i, row := range board {
		for j := range row {
			invField[i][8-j] = board[i][j]
		}
	}
	return invField
}

// this one inverts the board vertically
func invertVertically(board Board) Board {
	invField := Board{}
	for i, row := range board {
		for j := range row {
			invField[8-i][j] = board[i][j]
		}
	}
	return invField
}

// check if 2 boards are equal
func equalBoards(board1, board2 [9][9]int) bool {
	for i := range [9]int{} {
		for j := range [9]int{} {
			if board1[i][j] != board2[i][j] {
				return false
			}
		}
	}
	return true
}

func isThereAnotherSolution(board, solvedBoard *Board, x, y int) (bool, *Board) {
	if y == 9 {
		return true, board
	}
	if board[y][x] != 0 {
		nx, ny := next(x, y)
		return isThereAnotherSolution(board, solvedBoard, nx, ny)
	} else {
		for i := range [9]int{} {
			var v = i + 1
			if canPut(*board, x, y, v) {
				// Skip values that are the same as the original solved board
				if solvedBoard[y][x] == v {
					continue
				}
				board[y][x] = v
				nx, ny := next(x, y)
				if hasAnother, secondSolution := isThereAnotherSolution(board, solvedBoard, nx, ny); hasAnother {
					return true, secondSolution
				}
				board[y][x] = 0
			}
		}
		return false, nil
	}
}
