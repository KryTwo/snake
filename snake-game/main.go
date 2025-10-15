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
func (sn *Snake) MoveManually(dir string) {
	var newHead Point
	head := sn.body[0]

	sn.direction = dir

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

	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			p := Point{X: x, Y: y}
			if b.Snake.Contains(p) {
				//змейка
				if p == b.Snake.body[0] {
					fmt.Print("^")
				} else {
					fmt.Print("#")
				}

			} else if b.ContainsFood(p) {
				fmt.Print("$") //еда
			} else {
				fmt.Print(".") //пустое поле
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

// time.sleep
func Tick() {
	time.Sleep(900 * time.Millisecond)
}

// func (b *Board) isValidMove() bool {
// 	return true
// }

func (sn *Snake) NextHead() Point {
	//sn.NewDirection()
	head := sn.body[0]

	switch sn.direction {
	case "left":
		return Point{X: head.X - 1, Y: head.Y}
	case "right":
		return Point{X: head.X + 1, Y: head.Y}
	case "up":
		return Point{X: head.X, Y: head.Y - 1}
	default:
		return Point{X: head.X, Y: head.Y + 1}
	}
}

func (b *Board) IsOutOfBounds(newHead Point) bool {
	fmt.Printf("новая голова %v\n", newHead)
	return newHead.X < 0 || newHead.X >= b.Width || newHead.Y < 0 || newHead.Y >= b.Height
}

func (b *Board) CollidesWithItself(newHead Point) bool {
	return false
}

func (b *Board) Update() bool {

	newHead := b.Snake.NextHead()

	if b.IsOutOfBounds(newHead) {
		fmt.Println("Игра окончена, конец карты!")
		return false
	}
	return true

}

func main() {
	// rand.New(rand.NewSource(time.Now().UnixNano()))
	point := Point{X: 2, Y: 2}
	fmt.Printf("Точка старта: %+v\n", point)

	snake := Snake{
		body:    []Point{point},
		growing: false,
	}

	board := NewBoard(10, 10, &snake)

	//fmt.Printf("еда: x - %v, y - %v\n", board.Food.X, board.Food.Y)
	//snake.Grow()

	//time.Sleep(time.Millisecond * 500)
	Tick()
	//snake.NewDirection()

	snake.direction = "left"

	for board.Update() {
		snake.MoveManually(snake.direction)
		board.ShowBoard()
		Tick()
	}

	//fmt.Println(snake.body)
}

// Цикл на автоматическое перемещение
// Смена направления перемещения
// 1/2 Проверка границ поля
// Проверка самопересечения
// Проверка реверса движения (для змейки с телом)
// 3/4 отрисовка поля
// ...
// чтение с клавиатуры
