package main

import (
	"fmt"
	"strconv"
	"strings"
)

// TuringMachine represents the universal Turing machine
type TuringMachine struct {
	tape         *Tape
	transitions  map[string][]string
	currentState string
	isFinished   bool
}

// NewTuringMachine creates a new Turing machine with the given sequence
func NewTuringMachine(sequence string) *TuringMachine {
	parts := strings.SplitN(sequence, "111", 2)

	var programPart, tapePart string
	programPart = parts[0]
	if len(parts) > 1 {
		tapePart = parts[1]
	}

	tapeContent := strings.Split(tapePart, "")
	tape := NewTape(tapeContent)

	transitions := make(map[string][]string)
	transitionParts := strings.Split(programPart, "11")

	for _, transition := range transitionParts {
		if transition == "" {
			continue
		}

		parts := strings.SplitN(transition, "1", 5)
		if len(parts) != 5 {
			fmt.Println("Invalid transition format:", transition)
			continue
		}

		fromState := parseState(parts[0])
		readSymbol := parseSymbol(parts[1])
		toState := parseState(parts[2])
		writeSymbol := parseSymbol(parts[3])
		move := parseMove(parts[4])

		key := fromState + "," + readSymbol
		transitions[key] = []string{toState, writeSymbol, move}
	}

	return &TuringMachine{
		tape:         tape,
		transitions:  transitions,
		currentState: "q1",
		isFinished:   false,
	}
}

// Step executes one step of the Turing machine
func (tm *TuringMachine) Step() {
	if tm.isFinished {
		return
	}

	key := tm.currentState + "," + tm.tape.Read()
	transition, exists := tm.transitions[key]

	if !exists {
		tm.isFinished = true
		return
	}

	toState := transition[0]
	writeSymbol := transition[1]
	move := transition[2]

	tm.tape.Write(writeSymbol)

	if move == "L" {
		tm.tape.MoveLeft()
	} else if move == "R" {
		tm.tape.MoveRight()
	}

	tm.currentState = toState
}

// IsAccepting returns true if the machine is in an accepting state
func (tm *TuringMachine) IsAccepting() bool {
	return tm.currentState == "q2"
}

// GetTransitionsTable returns the transitions formatted as a table
func (tm *TuringMachine) GetTransitionsTable() [][]string {
	table := make([][]string, 0, len(tm.transitions)+1)

	// Header row
	table = append(table, []string{"#", "From State", "Read", "To State", "Write", "Move"})

	// Data rows
	i := 1
	for key, transition := range tm.transitions {
		parts := strings.Split(key, ",")
		row := []string{
			strconv.Itoa(i),
			parts[0],
			parts[1],
			transition[0],
			transition[1],
			transition[2],
		}
		table = append(table, row)
		i++
	}

	return table
}

// Parse helpers
func parseState(unaryCode string) string {
	return "q" + strconv.Itoa(len(unaryCode))
}

func parseSymbol(unaryCode string) string {
	switch len(unaryCode) {
	case 1:
		return "0"
	case 2:
		return "1"
	case 3:
		return "_"
	default:
		return strconv.Itoa(len(unaryCode))
	}
}

func parseMove(unaryCode string) string {
	if unaryCode == "0" {
		return "L"
	} else if unaryCode == "00" {
		return "R"
	}
	return "X" // Invalid move
}
