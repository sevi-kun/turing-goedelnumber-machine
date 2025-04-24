package main

// Tape represents the Turing machine's tape
type Tape struct {
	content      []string
	headPosition int
}

// NewTape creates a new tape with the given content
func NewTape(content []string) *Tape {
	if len(content) == 0 {
		content = []string{"_"}
	}
	return &Tape{
		content:      content,
		headPosition: 0,
	}
}

// Read returns the symbol at the current head position
func (t *Tape) Read() string {
	if t.headPosition < 0 || t.headPosition >= len(t.content) {
		return "_"
	}
	return t.content[t.headPosition]
}

// Write writes a symbol at the current head position
func (t *Tape) Write(symbol string) {
	// Expand tape if needed
	for t.headPosition >= len(t.content) {
		t.content = append(t.content, "_")
	}
	for t.headPosition < 0 {
		t.content = append([]string{"_"}, t.content...)
		t.headPosition = 0
	}
	t.content[t.headPosition] = symbol
}

// MoveLeft moves the head one position to the left
func (t *Tape) MoveLeft() {
	t.headPosition--
	// Expand tape if needed
	if t.headPosition < 0 {
		t.content = append([]string{"_"}, t.content...)
		t.headPosition = 0
	}
}

// MoveRight moves the head one position to the right
func (t *Tape) MoveRight() {
	t.headPosition++
	// Expand tape if needed
	if t.headPosition >= len(t.content) {
		t.content = append(t.content, "_")
	}
}
