package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	width  = 20
	height = 10
)

type Point struct {
	X, Y int
}

type Snake struct {
	Body []Point
}

type Game struct {
	Snake      Snake
	Food       Point
	Direction  Point
	IsGameOver bool
}

func main() {
	err := termbox.Init()
	if err != nil {
		fmt.Println("Error initializing termbox:", err)
		return
	}
	defer termbox.Close()

	game := NewGame()
	game.Start()

mainLoop:
	for !game.IsGameOver {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				game.SetDirection(0, -1)
			case termbox.KeyArrowDown:
				game.SetDirection(0, 1)
			case termbox.KeyArrowLeft:
				game.SetDirection(-1, 0)
			case termbox.KeyArrowRight:
				game.SetDirection(1, 0)
			case termbox.KeyEsc:
				break mainLoop
			}
		}

		game.Update()
		game.Render()
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("Game Over. Your score:", len(game.Snake.Body)-1)
}

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	game := &Game{}
	game.Snake = Snake{Body: []Point{{width / 2, height / 2}}}
	game.Food = game.generateFood()
	game.Direction = Point{1, 0} // 初始方向向右
	game.IsGameOver = false
	return game
}

func (game *Game) generateFood() Point {
	for {
		food := Point{rand.Intn(width), rand.Intn(height)}
		if !game.isSnakeBody(food) {
			return food
		}
	}
}

func (game *Game) isSnakeBody(point Point) bool {
	for _, body := range game.Snake.Body {
		if body == point {
			return true
		}
	}
	return false
}

func (game *Game) SetDirection(dx, dy int) {
	game.Direction = Point{dx, dy}
}

func (game *Game) Start() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	game.Render()
}

func (game *Game) Update() {
	if game.IsGameOver {
		return
	}

	// 移动蛇头
	newHead := Point{
		X: game.Snake.Body[0].X + game.Direction.X,
		Y: game.Snake.Body[0].Y + game.Direction.Y,
	}

	// 检查是否碰到墙或自身
	if newHead.X < 0 || newHead.X >= width || newHead.Y < 0 || newHead.Y >= height || game.isSnakeBody(newHead) {
		game.IsGameOver = true
		return
	}

	// 判断是否吃到食物
	if newHead == game.Food {
		game.Snake.Body = append([]Point{newHead}, game.Snake.Body...)
		game.Food = game.generateFood()
	} else {
		// 移动蛇尾
		game.Snake.Body = append([]Point{newHead}, game.Snake.Body[:len(game.Snake.Body)-1]...)
	}
}

func (game *Game) Render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// 渲染蛇身
	for _, body := range game.Snake.Body {
		termbox.SetCell(body.X, body.Y, 'O', termbox.ColorDefault, termbox.ColorGreen)
	}

	// 渲染食物
	termbox.SetCell(game.Food.X, game.Food.Y, 'F', termbox.ColorDefault, termbox.ColorRed)

	termbox.Flush()
}
