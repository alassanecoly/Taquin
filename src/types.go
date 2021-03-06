package main

import (
	"math"
)

type Puzzle struct {
	board     [][]int
	path      []int
	dimension []int // dimension[0] represents the number of lines. dimension[1.txt], the number of columns
	lastMove  int
	distance  int
}

func (puzzle Puzzle) getBlankSpacePosition() []int {
	for i := 0; i < puzzle.dimension[0]; i++ {
		for j := 0; j < puzzle.dimension[1]; j++ {
			if puzzle.board[i][j] == 0 {
				return []int{i, j}
			}
		}
	}
	return nil
}

func (puzzle Puzzle) swap(i1 int, j1 int, i2 int, j2 int) {
	puzzle.board[i1][j1], puzzle.board[i2][j2] = puzzle.board[i2][j2], puzzle.board[i1][j1]
}

func (puzzle Puzzle) getMove(piece int) string {
	var blankSpacePosition = puzzle.getBlankSpacePosition()
	var line = blankSpacePosition[0]
	var column = blankSpacePosition[1]
	switch {
	case line > 0 && piece == puzzle.board[line-1][column]:
		return DOWN
	case line < puzzle.dimension[0]-1 && piece == puzzle.board[line+1][column]:
		return UP
	case column > 0 && piece == puzzle.board[line][column-1]:
		return RIGHT
	case column < puzzle.dimension[1]-1 && piece == puzzle.board[line][column+1]:
		return LEFT
	}
	return ""
}

func (puzzle *Puzzle) isGoalState() bool {

	for i := 0; i < puzzle.dimension[0]; i++ {
		for j := 0; j < puzzle.dimension[1]; j++ {
			piece := puzzle.board[i][j]
			if piece != 0 {
				originalLine := int(math.Floor(float64((piece - 1) / puzzle.dimension[0])))
				originalColumn := (piece - 1) % puzzle.dimension[1]
				if i != originalLine || j != originalColumn {
					return false
				}
			}
		}
	}
	return true
}

func (puzzle Puzzle) getCopy() *Puzzle {
	newPuzzle := new(Puzzle)
	n, m := len(puzzle.board), len(puzzle.board[0])

	newBoardDuplicate := make([][]int, n)
	data := make([]int, n*m)
	for i := range puzzle.board {
		start := i * m
		end := start + m
		newBoardDuplicate[i] = data[start:end:end]
		copy(newBoardDuplicate[i], puzzle.board[i])
	}
	newDimensionDuplicate := make([]int, n)
	copy(newDimensionDuplicate, puzzle.dimension)
	newPathDuplicate := make([]int, len(puzzle.path))
	copy(newPathDuplicate, puzzle.path)

	newPuzzle.board = newBoardDuplicate
	newPuzzle.path = newPathDuplicate
	newPuzzle.dimension = newDimensionDuplicate
	return newPuzzle
}

func (puzzle Puzzle) getAllowedMoves() []int {
	var allowedMoves = make([]int, 0)

	for i := 0; i < puzzle.dimension[0]; i++ {
		for j := 0; j < puzzle.dimension[1]; j++ {
			piece := puzzle.board[i][j]
			if puzzle.getMove(piece) != "" {
				allowedMoves = append(allowedMoves, piece)
			}
		}
	}
	return allowedMoves
}

func (puzzle Puzzle) visit() []*Puzzle {
	children := make([]*Puzzle, 0)
	allowedMoves := puzzle.getAllowedMoves()

	for i := 0; i < len(allowedMoves); i++ {
		var move = allowedMoves[i]
		if move != puzzle.lastMove {
			var newInstance = puzzle.getCopy()
			newInstance.move(move)
			newInstance.lastMove = move
			newInstance.path = append(newInstance.path, move)
			children = append(children, newInstance)
		}
	}
	return children
}

func (puzzle Puzzle) move(piece int) string {
	var move = puzzle.getMove(piece)
	if move != "" {
		blankSpacePosition := puzzle.getBlankSpacePosition()
		line, column := blankSpacePosition[0], blankSpacePosition[1]
		switch move {
		case LEFT:
			puzzle.swap(line, column, line, column+1)
		case RIGHT:
			puzzle.swap(line, column, line, column-1)
		case UP:
			puzzle.swap(line, column, line+1, column)
		case DOWN:
			puzzle.swap(line, column, line-1, column)
		}

		return move
	}
	return ""
}

func (puzzle Puzzle) getManhattanDistance() int {
	var distance = 0
	for i := 0; i < puzzle.dimension[0]; i++ {
		for j := 0; j < puzzle.dimension[1]; j++ {
			piece := puzzle.board[i][j]
			if piece != 0 {
				originalLine := int(math.Floor(float64((piece - 1) / puzzle.dimension[0])))
				originalColumn := int((piece - 1) % puzzle.dimension[1])
				distance += int(math.Abs(float64(i-originalLine)) + math.Abs(float64(j-originalColumn)))
			}
		}
	}
	return distance
}

func (puzzle Puzzle) countMisplaced() int {
	count := 0
	for i := 0; i < puzzle.dimension[0]; i++ {
		for j := 0; j < puzzle.dimension[1]; j++ {
			piece := puzzle.board[i][j]
			if piece != 0 {
				originalLine := int(math.Floor(float64((piece - 1) / puzzle.dimension[0])))
				originalColumn := int((piece - 1) % puzzle.dimension[1])
				if i != originalLine || j != originalColumn {
					count++
				}
			}
		}
	}
	return count
}

func (puzzle Puzzle) getCost() int {
	if selectedAlgorithm == MANHATTAN {
		return puzzle.getManhattanDistance()
	} else if selectedAlgorithm == MISPLACED {
		return puzzle.countMisplaced()
	} else {
		return puzzle.countMisplaced() + puzzle.getManhattanDistance()
	}
}

func (puzzle Puzzle) isSolvable() bool {
	flattenPuzzle := make([]int, 0)

	for i := 0; i < puzzle.dimension[0]; i++ {
		for j := 0; j < puzzle.dimension[1]; j++ {
			flattenPuzzle = append(flattenPuzzle, puzzle.board[i][j])
		}
	}

	ic := puzzle.getInversionsCount(flattenPuzzle)

	if puzzle.dimension[0]%2 == 1 || puzzle.dimension[1]%2 == 1 {
		return ic%2 == 0
	} else if ic%2 == 0 {
		return false
	} else {
		return true
	}
}

func (puzzle Puzzle) getInversionsCount(flattenPuzzle []int) int {
	ic := 0
	for i := 0; i < puzzle.dimension[0]*puzzle.dimension[1]-1; i++ {
		for j := i + 1; j < puzzle.dimension[0]*puzzle.dimension[1]; j++ {
			if flattenPuzzle[i] > flattenPuzzle[j] {
				ic++
			}
		}
	}
	return ic
}
