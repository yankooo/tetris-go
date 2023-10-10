package tetris

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math"
	"time"
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
}

func NewGame() *tetrisGame {
	g := &tetrisGame{}
	return g
}

func (g *tetrisGame) Initialize() {
	g.baseSpeed = 0.8
	g.gravitySpeed = 0.8
	g.levelUpTimer = levelLength
	g.initWindows()

	g.board = NewBoard()
	g.board.initResource()
	g.board.AddPiece() // Add initial Piece to game
}

// Initialize the window
func (g *tetrisGame) initWindows() {
	var err error
	windowWidth := 765.0
	windowHeight := 450.0
	cfg := pixelgl.WindowConfig{
		Title:  "俄罗斯方块nb",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}
	g.win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(any(err))
	}
}

// run is the main code for the game. Allows pixelgl to run on main thread
func (g *tetrisGame) Run() {
	last := time.Now()
	for !g.win.Closed() && !g.board.GameOver() {
		// Perform time processing events
		dt := time.Since(last).Seconds()
		last = time.Now()
		g.gravityTimer += dt
		g.levelUpTimer -= dt

		// Time Functions:
		// Gravity
		if g.gravityTimer > g.gravitySpeed {
			g.gravityTimer -= g.gravitySpeed
			didCollide := g.board.applyGravity()
			if !didCollide {
				if g.board.isTouchingFloor() {
					g.gravityTimer -= g.gravitySpeed // Add extra time when touching floor
				}
			} else {
				g.board.AddScore(10)
			}
		}

		if g.leftRightDelay > 0.0 {
			g.leftRightDelay = math.Max(g.leftRightDelay-dt, 0.0)
		}

		// Speed up
		if g.levelUpTimer <= 0 {
			if g.baseSpeed > 0.2 {
				g.baseSpeed = math.Max(g.baseSpeed-speedUpRate, 0.2)
			}
			g.levelUpTimer = levelLength
			g.gravitySpeed = g.baseSpeed
		}

		// Keypress Functions
		if g.win.Pressed(pixelgl.KeyRight) && g.leftRightDelay == 0.0 {
			g.board.movePiece(1)
			if g.moveCounter > 0 {
				g.leftRightDelay = 0.1
			} else {
				g.leftRightDelay = 0.5
			}
			g.moveCounter++
		}
		if g.win.Pressed(pixelgl.KeyLeft) && g.leftRightDelay == 0.0 {
			g.board.movePiece(-1)
			if g.moveCounter > 0 {
				g.leftRightDelay = 0.1
			} else {
				g.leftRightDelay = 0.5
			}
			g.moveCounter++
		}
		if g.win.JustPressed(pixelgl.KeyDown) {
			g.gravitySpeed = 0.08 // TODO: Code could result in bugs if game pause functionality added
			if g.gravityTimer > 0.08 {
				g.gravityTimer = 0.08
			}
		}
		if g.win.JustReleased(pixelgl.KeyDown) {
			g.gravitySpeed = g.baseSpeed // TODO: Code could result in bugs if game pause functionality added
		}
		if g.win.JustPressed(pixelgl.KeyUp) {
			g.board.rotatePiece()
			if g.board.isTouchingFloor() {
				g.gravityTimer = 0 // Make gravity more forgiving when moving pieces
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

		g.win.Clear(colornames.Black)
		g.board.displayBG(g.win)
		g.board.displayText(g.win)
		g.board.displayBoard(g.win)
		g.win.Update()
	}
}
