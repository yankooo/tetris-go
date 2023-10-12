package tetris

import (
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type tetrisGame struct {
	win   *pixelgl.Window
	board *Board

	gravityTimer   float64
	baseSpeed      float64
	gravitySpeed   float64
	levelUpTimer   float64
	leftRightDelay float64
	moveCounter    int

	isPaused bool
}

func NewGame() *tetrisGame {
	g := &tetrisGame{}
	return g
}

func (g *tetrisGame) Initialize() {
	g.baseSpeed = 0.8
	g.gravitySpeed = g.baseSpeed
	g.levelUpTimer = levelLength // Assuming levelLength is initialized elsewhere
	g.initWindows()

	g.board = NewBoard()
	g.board.initResource()
	g.board.AddPiece()
}

func (g *tetrisGame) initWindows() {
	var err error
	windowWidth := 765.0
	windowHeight := 450.0
	cfg := pixelgl.WindowConfig{
		Title:  "俄罗斯方块",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}
	g.win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
}

func (g *tetrisGame) Run() {
	last := time.Now()
	for !g.win.Closed() && !g.board.GameOver() {
		if g.win.JustPressed(pixelgl.MouseButtonLeft) {
			last = g.togglePause(last)
		}

		if g.isPaused {
			g.displayPausedMessage()
			g.win.Update()
			continue
		}

		dt := time.Since(last).Seconds()
		last = time.Now()
		g.gravityTimer += dt
		g.levelUpTimer -= dt

		if g.gravityTimer > g.gravitySpeed {
			g.gravityTimer -= g.gravitySpeed
			didCollide := g.board.applyGravity()
			if !didCollide {
				if g.board.isTouchingFloor() {
					g.gravityTimer -= g.gravitySpeed
				}
			} else {
				g.board.AddScore(10)
			}
		}
		if g.leftRightDelay > 0.0 {
			g.leftRightDelay = math.Max(g.leftRightDelay-dt, 0.0)
		}

		if g.levelUpTimer <= 0 {
			if g.baseSpeed > 0.2 {
				g.baseSpeed = math.Max(g.baseSpeed-speedUpRate, 0.2) // Assuming speedUpRate is initialized elsewhere
			}
			g.levelUpTimer = levelLength
			g.gravitySpeed = g.baseSpeed
		}

		g.processKeypresses()

		g.win.Clear(colornames.Black)
		g.board.displayBG(g.win)
		g.board.displayText(g.win)
		g.board.displayBoard(g.win)
		g.win.Update()
	}
}

// 为tetrisGame结构体添加一个新方法来处理暂停的逻辑
func (g *tetrisGame) togglePause(last time.Time) time.Time {
	g.isPaused = !g.isPaused
	if !g.isPaused {
		// 重置时间以防止dt累积
		last = time.Now()
	}
	return last
}

// 显示暂停消息的方法
func (g *tetrisGame) displayPausedMessage() {
	g.board.displayPaused(g.win)
}

// Separated keypress handling for clarity
func (g *tetrisGame) processKeypresses() {
	if g.win.Pressed(pixelgl.KeyRight) && g.leftRightDelay == 0.0 {
		g.handleHorizontalMove(1)
	}
	if g.win.Pressed(pixelgl.KeyLeft) && g.leftRightDelay == 0.0 {
		g.handleHorizontalMove(-1)
	}
	if g.win.JustPressed(pixelgl.KeyDown) {
		g.gravitySpeed = 0.08
		if g.gravityTimer > 0.08 {
			g.gravityTimer = 0.08
		}
	}
	if g.win.JustReleased(pixelgl.KeyDown) {
		g.gravitySpeed = g.baseSpeed
	}
	if g.win.JustPressed(pixelgl.KeyUp) {
		g.board.rotatePiece()
		if g.board.isTouchingFloor() {
			g.gravityTimer = 0
		}
	}
	if g.win.JustPressed(pixelgl.KeySpace) {
		g.board.instafall()
		g.board.AddScore(12)
	}
	if !g.win.Pressed(pixelgl.KeyRight) && !g.win.Pressed(pixelgl.KeyLeft) {
		g.moveCounter = 0
		g.leftRightDelay = 0
	}
}

// Handle left/right movement
func (g *tetrisGame) handleHorizontalMove(direction int) {
	g.board.movePiece(direction)
	if g.moveCounter > 0 {
		g.leftRightDelay = 0.1
	} else {
		g.leftRightDelay = 0.5
	}
	g.moveCounter++
}
