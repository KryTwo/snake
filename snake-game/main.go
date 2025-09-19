package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Point struct {
	X, Y int
}

type Snake struct {
	body      []Point
	direction string
	growing   bool
}

type Board struct {
	Width, Height int
	Snake         *Snake
	Food          Point
}

func (sn *Snake) NewDirection() {
	dirs := [4]string{
		"left",
		"right",
		"up",
		"down",
	}
	sn.direction = dirs[rand.Intn(4)]

}

// Движение змейки (для теста в случайном направлении)
func (sn *Snake) Move() {
	var newHead Point
	head := sn.body[0]

	sn.NewDirection()

	switch sn.direction {
	case "left":
		newHead = Point{X: head.X - 1, Y: head.Y}
	case "right":
		newHead = Point{X: head.X + 1, Y: head.Y}
	case "up":
		newHead = Point{X: head.X, Y: head.Y - 1}
	case "down":
		newHead = Point{X: head.X, Y: head.Y + 1}
	}

	newBody := []Point{newHead}
	newBody = append(newBody, sn.body...)
	if !sn.growing {
		newBody = newBody[:len(newBody)-1]
	}
	sn.growing = false
	sn.body = newBody
}

// рост змейки
func (sn *Snake) Grow() {
	sn.growing = true
}

// Создание поля
func NewBoard(width, height int, snake *Snake) *Board {
	board := &Board{
		Width:  width,
		Height: height,
		Snake:  snake,
	}
	board.SpawnFood()
	return board
}

// Отрисовка поля
func (b *Board) ShowBoard() {
	var point Point
	for y := 0; y <= b.Height; y++ {
		point.Y = y
		for x := 0; x <= b.Width; x++ {
			point.X = x
			if b.Snake.Contains(point) {
				fmt.Print("#")
			} else if b.ContainsFood(point) {
				fmt.Print("$")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

// Спавн еды в случайном месте
func (b *Board) SpawnFood() {
	var x, y int
	var food Point

	for {
		x = rand.Intn(b.Width)
		y = rand.Intn(b.Height)
		food = Point{X: x, Y: y}

		if !b.Snake.Contains(food) {
			break
		}
	}

	b.Food = food
}

// Проверка - есть ли в Point часть змейки
func (sn *Snake) Contains(p Point) bool {
	for _, part := range sn.body {
		if part.X == p.X && part.Y == p.Y {
			return true
		}
	}
	return false
}

// Проверка - есть ли в point еда
func (b *Board) ContainsFood(p Point) bool {
	if b.Food.X == p.X && b.Food.Y == p.Y {
		return true
	}
	return false
}

func main() {
	// rand.New(rand.NewSource(time.Now().UnixNano()))
	point := Point{X: 5, Y: 5}
	fmt.Printf("Точка старта: %+v\n", point)

	snake := Snake{
		body:    []Point{point},
		growing: false,
	}

	board := NewBoard(10, 10, &snake)

	fmt.Printf("еда: x - %v, y - %v\n", board.Food.X, board.Food.Y)
	//snake.Grow()
	for {
		time.Sleep(time.Millisecond * 500)
		snake.Move()
		fmt.Println(snake.body)
		board.ShowBoard()
	}

	//fmt.Println(snake.body)
}

// Цикл на автоматическое перемещение
// Смена направления перемещения
// Проверка границ поля
// Проверка самопересечения
// Проверка реверса движения
// отрисовка поля
