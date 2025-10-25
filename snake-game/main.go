package main

import (
	"fmt"
	"math/rand"
	"os"
	"slices"
	"time"

	"github.com/eiannone/keyboard"
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
	board.PlaceFood(5, 8)
	//board.SpawnFood()
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

// Спавн еды вручную
func (b *Board) PlaceFood(x, y int) {
	var food Point

	for {
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
	time.Sleep(1500 * time.Millisecond)
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
	fmt.Printf("newHead: %v\n", newHead)
	return newHead.X < 0 || newHead.X >= b.Width || newHead.Y < 0 || newHead.Y >= b.Height
}

func (b *Board) CollidesWithItself(newHead Point) bool {
	return slices.Contains(b.Snake.body, newHead)
}

func (b *Board) WillEat(newHead Point) bool {
	return b.Food.X == newHead.X && b.Food.Y == newHead.Y
}

func (b *Board) Update() bool {

	newHead := b.Snake.NextHead()

	if b.IsOutOfBounds(newHead) {
		fmt.Println("You can't escape!")
		return false
	}

	if b.CollidesWithItself(newHead) {
		fmt.Println("Ah, the ouroboros — and suddenly me?!")
		return false
	}

	if b.WillEat(newHead) {
		b.Snake.growing = true
		b.SpawnFood()
	}

	return true
}

// Получение rune и key нажатой клавиши
func getKey(ch chan int) {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}

	defer func() {
		_ = keyboard.Close()
	}()

	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}

	for {
		event := <-keysEvents
		if event.Err != nil {
			fmt.Fprintf(os.Stderr, "Error reading key: %v\n", event.Err)
			continue
		}
		fmt.Printf("rune %q, ", event.Rune)

		if event.Key == keyboard.KeyEsc {
			return
		}
		ch <- int(event.Key)
	}
}

func directionForKey(ch chan int, ch2 chan string) {
	for {
		key := <-ch
		if key != 0 {
			switch key {
			case 65517:
				ch2 <- "up"
				fmt.Println("up 65517")
			case 65515:
				ch2 <- "left"
				fmt.Println("left 65515")
			case 65514:
				ch2 <- "right"
				fmt.Println("right 65514")
			case 65516:
				ch2 <- "down"
				fmt.Println("down 65516")
				// default:
				// 	key = 0
			}
		}
	}
}

func (sn *Snake) ChangeDirection(ch2 chan string) {
	sn.direction = <-ch2
}

func main() {
	ch := make(chan int)
	ch2 := make(chan string)
	chanEndGame := make(chan string)
	points := []Point{
		{5, 5}, {4, 5}, {3, 5}, {3, 4}, {4, 4},
	}
	fmt.Printf("Start point: %+v\n", points[0])

	snake := Snake{
		body:    points,
		growing: false,
	}

	board := NewBoard(10, 10, &snake)
	board.ShowBoard()

	Tick()

	go func() {
		for board.Update() {
			snake.ChangeDirection(ch2)
			snake.Move()
			board.ShowBoard()
			Tick()

		}
		chanEndGame <- "end"
	}()

	go getKey(ch)
	go directionForKey(ch, ch2)

	fmt.Println(<-chanEndGame)
}

// 1/2 Eat & grow
// 1/2 Смена направления перемещения
// 1/2 Проверка границ поля
// (Done) Проверка самопересечения
// (Done) Проверка пересечения границ
// Проверка реверса движения (для змейки с телом)
// 3/4 отрисовка поля
// ...
// чтение с клавиатуры
// Цикл на автоматическое перемещение
// автоматическая игра, до полного заполнения карты
