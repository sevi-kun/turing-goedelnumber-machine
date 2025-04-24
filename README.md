# Universal Turing Machine Simulator

A Go implementation of a Universal Turing Machine with an interactive terminal UI built using Bubbletea.

This is a rewrite of the original implementation by [CuddlyBunion341](https://github.com/CuddlyBunion341/universal-touring-machine).

## Overview

This simulator can emulate any Turing Machine by processing a binary encoding of the machine's transition function and initial tape contents. The application features a beautiful, interactive terminal UI that visualizes the execution of the Turing Machine step by step.

## Features

- **Universal Turing Machine Implementation**: Can simulate any Turing Machine given its binary encoding
- **Interactive Terminal UI**: Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) for a delightful user experience
- **Visualization**: See the tape, head position, and state changes in real-time
- **Execution Modes**: Run automatically or step through the execution manually
- **Transition Table View**: View the machine's transition function in a clean, tabular format
- **Colorful Output**: Using [Lipgloss](https://github.com/charmbracelet/lipgloss) for beautiful styling

## Installation

### Prerequisites

- Go 1.16 or higher

### Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/turing-machine.git
   cd turing-machine
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the application:
   ```bash
   go build
   ```

4. Run the application:
   ```bash
   ./turingmachine
   ```

## Usage

### Input Format

The input to the Universal Turing Machine is a binary string consisting of two parts separated by `111`:

```
<program>111<tape>
```

- **Program**: Encodes the transition function as a series of transitions separated by `11`
- **Tape**: The initial content of the tape

Each transition is encoded as five parts separated by `1`:

```
<from_state>1<read_symbol>1<to_state>1<write_symbol>1<move>
```

- **States** are encoded in unary: `0` for q1, `00` for q2, etc.
- **Symbols** are encoded as: `0` for '0', `00` for '1', `000` for '_' (blank)
- **Moves** are encoded as: `0` for left (L), `00` for right (R)

### Example Input

A default example is provided if no input is given. This example implements a simple binary increment machine.

### Navigation

- **Input Screen**: Enter your Turing Machine encoding or leave blank for the default example
- **Mode Selection**: Choose between automatic execution or step-by-step execution
- **Running Screen**: Watch the Turing Machine execute. Press Enter to advance in step-by-step mode
- **Final Screen**: View the final result (accepted or rejected) and the final tape state

## How It Works

The simulator parses the binary input to extract:

1. The transition function of the Turing Machine
2. The initial tape content

It then simulates the machine's execution by:

1. Starting in state q1 with the head at the leftmost position of the input tape
2. For each step, looking up the current (state, symbol) pair in the transition function
3. Performing the specified action (writing a symbol and moving the head)
4. Changing to the new state
5. Repeating until no valid transition is found or maximum steps are reached

The machine accepts if it halts in state q2.

## Project Structure

- **main.go**: Entry point of the application
- **tape.go**: Implementation of the Turing Machine tape
- **turing_machine.go**: Core logic of the Universal Turing Machine
- **ui.go**: Terminal UI implementation using Bubbletea

## Performance

This Go implementation is optimized for performance, making it suitable for simulating complex Turing Machines. The application can handle long sequences and execute thousands of steps efficiently.

## Claude & Zed

This port was created with the help of the amazing text editor [zed](https://zed.dev) and it's built in code assistant Claude 3.7 Sonnet.
