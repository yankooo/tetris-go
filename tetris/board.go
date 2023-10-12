package tetris

import (
	"fmt"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/yankooo/tetris-go/tetris/spritesheet"
	"golang.org/x/image/font/basicfont"
)

// Board is an array containing the entire game board pieces.
type Board struct {
	board        [22][10]Block
	currentPiece Piece
	nextPiece    Piece
	activeShape  Shape // The shape that the player controls
	score        int
	gameOver     bool

	blockGen          func(int) pixel.Picture
	bgImgSprite       pixel.Sprite
	gameBGSprite      pixel.Sprite
	scoreBgSprite     pixel.Sprite
	nextPieceBGSprite pixel.Sprite
}

func NewBoard() *Board {
	b := &Board{}
	b.nextPiece = Piece(rand.Intn(7))
	b.gameOver = false
	return b
}

func (b *Board) AddScore(score int) {
	b.score += score
}

// isTouchingFloor checks if the piece that the user is controlling has a piece
// directly below it. Used to give the user more time when placing block on
// floor
func (b *Board) isTouchingFloor() bool {
	blockType := b.board[b.activeShape[0].row][b.activeShape[0].col]
	b.drawPiece(b.activeShape, Empty)
	isTouching := b.checkCollision(moveShapeDown(b.activeShape))
	b.drawPiece(b.activeShape, blockType)
	return isTouching
}

// rotatePiece rotates the piece that the user is currently moving clockwise by
// 90 degrees. The rotation is made and collision is checked. If the rotation can
// be completed by moving the newly rotated shape, the rotation will also be
// performed. If it is impossible to rotate, does nothing.
func (b *Board) rotatePiece() {
	// The O piece should not be rotated
	if b.currentPiece == OPiece {
		return
	}
	blockType := b.board[b.activeShape[0].row][b.activeShape[0].col]
	// Erase Piece
	b.drawPiece(b.activeShape, Empty)

	// Get the new shape and check for it's collision
	newShape := rotateShape(b.activeShape)
	if b.checkCollision(newShape) {
		if !b.checkCollision(moveShapeRight(newShape)) {
			newShape = moveShapeRight(newShape)
		} else if !b.checkCollision(moveShapeLeft(newShape)) {
			newShape = moveShapeLeft(newShape)
		} else if !b.checkCollision(moveShapeDown(newShape)) {
			newShape = moveShapeDown(newShape)
		} else {
			b.drawPiece(b.activeShape, blockType)
			return
		}
	}
	b.activeShape = newShape
	b.drawPiece(b.activeShape, blockType)
}

// movePiece attemps to move the piece that the user is controlling either
// right or left. +1 signifies a right move while -1 signifies a left move
func (b *Board) movePiece(dir int) {
	blockType := b.board[b.activeShape[0].row][b.activeShape[0].col]

	// Erase old piece
	b.drawPiece(b.activeShape, Empty)

	// Check collision
	didCollide := b.checkCollision(moveShape(0, dir, b.activeShape))
	if !didCollide {
		b.activeShape = moveShape(0, dir, b.activeShape)
	}
	b.drawPiece(b.activeShape, blockType)
}

// drawPiece sets the values of a board, b, to a specific block type, t
// according to shape, s.
func (b *Board) drawPiece(s Shape, t Block) {
	for i := 0; i < 4; i++ {
		b.board[b.activeShape[i].row][b.activeShape[i].col] = t
	}
}

// checkCollision checks if at the 4 points of a shape, s, there is
// nothing but Empty value under it and the position of the shape
// is inside the playing board (10x22 (top two rows invisiable)).
func (b Board) checkCollision(s Shape) bool {
	for i := 0; i < 4; i++ {
		r := s[i].row
		c := s[i].col
		if r < 0 || r > 21 || c < 0 || c > 9 || b.board[r][c] != Empty {
			return true
		}
	}
	return false
}

// applyGravity is the function that moves a piece down. If a collision
// is detected place the piece down and add a new piece. Returns wheather
// a collision was made.
func (b *Board) applyGravity() bool {
	blockType := b.board[b.activeShape[0].row][b.activeShape[0].col]
	// Erase old piece
	b.drawPiece(b.activeShape, Empty)

	// Does the block collide if it moves down?
	didCollide := b.checkCollision(moveShapeDown(b.activeShape))

	if !didCollide {
		b.activeShape = moveShapeDown(b.activeShape)
	}

	b.drawPiece(b.activeShape, blockType)

	if didCollide {
		if isGameOver(b.activeShape) {
			b.gameOver = true
		}
		b.checkRowCompletion(b.activeShape)
		b.AddPiece() // Replace with random piece
		return true
	}
	return false
}

// instafall calls the applyGravity function until a collision is detected.
func (b *Board) instafall() {
	collide := false
	for !collide {
		collide = b.applyGravity()
	}
}

// checkRowCompletion checks if the rows in a given shape are filled (ie should
// be deleted). If full, deletes the rows.
func (b *Board) checkRowCompletion(s Shape) {
	// Ony the rows of the shape can be filled
	rowWasDeleted := true
	// Since when we delete a row it can be shifted down, repeatedly try
	// to delete a row until no more deletes can be made
	var deleteRowCt int
	for rowWasDeleted {
		rowWasDeleted = false
		for i := 0; i < 4; i++ {
			r := s[i].row
			emptyFound := false
			// Look for empty row
			for c := 0; c < 10; c++ {
				if b.board[r][c] == Empty {
					emptyFound = true
					continue
				}
			}
			// If no empty cell was found in row delete row
			if !emptyFound {
				b.deleteRow(r)
				rowWasDeleted = true
				b.score += 200
				deleteRowCt++
			}
		}
	}
	// Bonus score for combos over one
	if deleteRowCt > 1 {
		b.score += (deleteRowCt - 1) * 200
	}
}

// deleteRow remoes a row by shifting everything above it down by one.
func (b *Board) deleteRow(row int) {
	for r := row; r < 21; r++ {
		for c := 0; c < 10; c++ {
			b.board[r][c] = b.board[r+1][c]
		}
	}
}

// setPiece sets a value in the game board to a specific block type.
func (b *Board) setPiece(r, c int, val Block) {
	b.board[r][c] = val
}

// fillShape sets
func (b *Board) fillShape(s Shape, val Block) {
	for i := 0; i < 4; i++ {
		b.setPiece(s[i].row, s[i].col, val)
	}
}

// addPiece creates a piece at the top of the screen at a random position
// and sets it to the piece that the player is controlling
// (ie b.activeShape).
func (b *Board) AddPiece() {
	var offset int
	if b.nextPiece == IPiece {
		offset = rand.Intn(7)
	} else if b.nextPiece == OPiece {
		offset = rand.Intn(9)
	} else {
		offset = rand.Intn(8)
	}
	baseShape := getShapeFromPiece(b.nextPiece)
	baseShape = moveShape(20, offset, baseShape)
	b.fillShape(baseShape, piece2Block(b.nextPiece))
	b.currentPiece = b.nextPiece
	b.activeShape = baseShape
	b.nextPiece = Piece(rand.Intn(7))
}

// displayBoard displays a particular game board with all of its pieces
// onto a given window, win
func (b *Board) displayBoard(win *pixelgl.Window) {
	boardBlockSize := 20.0 //win.Bounds().Max.X / 10
	pic := b.blockGen(0)
	imgSize := pic.Bounds().Max.X
	scaleFactor := float64(boardBlockSize) / float64(imgSize)

	for col := 0; col < BoardCols; col++ {
		for row := 0; row < BoardRows-2; row++ {
			val := b.board[row][col]
			if val == Empty {
				continue
			}

			x := float64(col)*boardBlockSize + boardBlockSize/2
			y := float64(row)*boardBlockSize + boardBlockSize/2
			pic := b.blockGen(block2spriteIdx(val))
			sprite := pixel.NewSprite(pic, pic.Bounds())
			sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, scaleFactor).Moved(pixel.V(x+282, y+25)))
		}
	}

	// Display Shadow
	pieceType := b.board[b.activeShape[0].row][b.activeShape[0].col]
	ghostShape := b.activeShape
	b.drawPiece(b.activeShape, Empty)
	for {
		if b.checkCollision(moveShapeDown(ghostShape)) {
			break
		}
		ghostShape = moveShapeDown(ghostShape)
	}
	b.drawPiece(b.activeShape, pieceType)

	gpic := b.blockGen(block2spriteIdx(Gray))
	sprite := pixel.NewSprite(gpic, gpic.Bounds())
	for i := 0; i < 4; i++ {
		if b.board[ghostShape[i].row][ghostShape[i].col] == Empty {
			x := float64(ghostShape[i].col)*boardBlockSize + boardBlockSize/2
			y := float64(ghostShape[i].row)*boardBlockSize + boardBlockSize/2
			sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, scaleFactor/2).Moved(pixel.V(x+282, y+25)))
		}
	}
}

// block2spriteIdx associates a blocks color (b Block) with its index in the sprite sheet.
func block2spriteIdx(b Block) int {
	return int(b) - 1
}

// piece2Block associates a pieces shape (Piece) with it's color/image (Block).
func piece2Block(p Piece) Block {
	switch p {
	case LPiece:
		return Goluboy
	case IPiece:
		return Siniy
	case OPiece:
		return Pink
	case TPiece:
		return Purple
	case SPiece:
		return Red
	case ZPiece:
		return Yellow
	case JPiece:
		return Green
	}
	panic(any("piece2Block: Invalid piece passed in"))
	return GraySpecial // Return strange value value
}

func (b *Board) displayPaused(win *pixelgl.Window) {
	// Text Generator
	scoreTextLocX := 315.0
	scoreTextLocY := 215.0
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	scoreTxt := text.New(pixel.V(scoreTextLocX, scoreTextLocY), basicAtlas)
	fmt.Fprintf(scoreTxt, "Game Pause")
	scoreTxt.Draw(win, pixel.IM.Scaled(scoreTxt.Orig, 2))
}

func (b *Board) displayText(win *pixelgl.Window) {
	// 玩法说明
	b.displayIntroduction(win)

	// 分数
	scoreTextLocX := 100.0
	scoreTextLocY := 400.0
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	scoreTxt := text.New(pixel.V(scoreTextLocX, scoreTextLocY), basicAtlas)
	fmt.Fprintf(scoreTxt, "Score")
	scoreTxt.Draw(win, pixel.IM.Scaled(scoreTxt.Orig, 2))

	score := text.New(pixel.V(scoreTextLocX, scoreTextLocY-30), basicAtlas)
	fmt.Fprintf(score, "%d", b.score)
	score.Draw(win, pixel.IM.Scaled(scoreTxt.Orig, 1.5))

	nextPieceTextLocX := 100.0
	nextPieceTextLocY := 300.0
	nextPieceTxt := text.New(pixel.V(nextPieceTextLocX, nextPieceTextLocY), basicAtlas)

	fmt.Fprintf(nextPieceTxt, "Next Piece")
	nextPieceTxt.Draw(win, pixel.IM.Scaled(scoreTxt.Orig, 2))
}

func (b *Board) displayIntroduction(win *pixelgl.Window) {
	scoreTextLocX := 460.0
	scoreTextLocY := 350.0
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	scoreTxt := text.New(pixel.V(scoreTextLocX, scoreTextLocY), basicAtlas)
	fmt.Fprintf(scoreTxt, `
	Down : Drop faster

	Left : Move one block left
	
	Right: Move one block right
	
	Space: Drop all the way down
	
	Click: Pause	
	`)
	scoreTxt.Draw(win, pixel.IM.Scaled(scoreTxt.Orig, 1.28))
}

func (b *Board) displayBG(win *pixelgl.Window) {
	b.bgImgSprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	b.gameBGSprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	b.scoreBgSprite.Draw(win, pixel.IM.Moved(pixel.V(200, 360)))
	b.nextPieceBGSprite.Draw(win, pixel.IM.Moved(pixel.V(150, 120)))

	baseShape := getShapeFromPiece(b.nextPiece)
	pic := b.blockGen(block2spriteIdx(piece2Block(b.nextPiece)))
	sprite := pixel.NewSprite(pic, pic.Bounds())
	boardBlockSize := 20.0
	scaleFactor := float64(boardBlockSize) / pic.Bounds().Max.Y
	shapeWidth := getShapeWidth(baseShape) + 1
	shapeHeight := 2

	for i := 0; i < 4; i++ {
		r := baseShape[i].row
		c := baseShape[i].col
		x := float64(c)*boardBlockSize + boardBlockSize/2
		y := float64(r)*boardBlockSize + boardBlockSize/2
		sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, scaleFactor).Moved(pixel.V(x+150-(float64(shapeWidth)*10), y+120-(float64(shapeHeight)*10))))
	}
}

func (b *Board) initResource() {
	var err error
	// Load Various Resources:
	b.blockGen, err = spritesheet.InitBlock("blocks.png", 2, 8)
	if err != nil {
		panic(any(err))
	}

	bgPic, err := spritesheet.LoadPicture("bg_whitecanvas.png")
	if err != nil {
		panic(any(err))
	}
	b.bgImgSprite = *pixel.NewSprite(bgPic, bgPic.Bounds())

	// tetrisGame Background
	blackPic := spritesheet.GetPlayBGPic()
	b.gameBGSprite = *pixel.NewSprite(blackPic, blackPic.Bounds())

	// Score BG
	scoreBgPic := spritesheet.GetScoreBGPic()
	b.scoreBgSprite = *pixel.NewSprite(scoreBgPic, scoreBgPic.Bounds())

	// Next Piece BG
	nextPiecePic := spritesheet.GetNextPieceBGPic()
	b.nextPieceBGSprite = *pixel.NewSprite(nextPiecePic, nextPiecePic.Bounds())
}

func (b *Board) GameOver() bool {
	return b.gameOver
}
