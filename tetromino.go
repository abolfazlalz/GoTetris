package main

import (
	"image/color"
	"math/rand"
)

type TetrominoType int

const (
	I TetrominoType = iota
	O
	T
	S
	Z
	J
	L
)

type Tetromino struct {
	Shape [][]int // 4x4 matrix
	Color color.Color
	X     int // position on board (column)
	Y     int // position on board (row)
	Type  TetrominoType
}

var tetrominoShapes = map[TetrominoType][][]int{
	I: {
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	O: {
		{1, 1, 0, 0},
		{1, 1, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	T: {
		{0, 1, 0, 0},
		{1, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	S: {
		{0, 1, 1, 0},
		{1, 1, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	Z: {
		{1, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	J: {
		{1, 0, 0, 0},
		{1, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	L: {
		{0, 0, 1, 0},
		{1, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
}

var tetrominoColors = map[TetrominoType]color.Color{
	I: color.RGBA{0, 255, 255, 255}, // Cyan
	O: color.RGBA{255, 255, 0, 255}, // Yellow
	T: color.RGBA{128, 0, 128, 255}, // Purple
	S: color.RGBA{0, 255, 0, 255},   // Green
	Z: color.RGBA{255, 0, 0, 255},   // Red
	J: color.RGBA{0, 0, 255, 255},   // Blue
	L: color.RGBA{255, 165, 0, 255}, // Orange
}

func NewRandomTetromino(boardWidth int) *Tetromino {
	t := TetrominoType(rand.Intn(7))

	return &Tetromino{
		Type:  t,
		Shape: tetrominoShapes[t],
		Color: tetrominoColors[t],
		X:     boardWidth/2 - 2, // وسط بورد
		Y:     0,
	}
}
