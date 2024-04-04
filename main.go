// Fronton 2D game with one ball (accelerated), similar to famous Pong.
// Divertimento in go with ebiten. Thanks to GO team and Hajimehoshi.
//
// Marc Riart, 202404.

package main

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	windowTitle  = "Fronton"
	screenWidth  = 400
	screenHeight = 600
	racketWidth  = screenWidth / 4
	racketHeight = 5
	winnerScore  = 8
)

type Game struct {
	State int // Defines the game state: 0 not initiated, 1 started, 2 over
	Score struct {
		Player int // Each hit in the racket, increases 1
		CPU    int // Each loss, increases 1
	}
	Ball struct {
		X      int
		Y      int
		SpeedX int
		SpeedY int
	}
	Racket struct {
		X     int
		Speed int // Constant, but can be dynamically change in the future
	}
}

func (g *Game) Initialize() {
	g.State = 0
	g.Score.Player = 0
	g.Score.CPU = 0
	g.Ball.X, g.Ball.Y = rand.IntN(screenWidth), 5 // Initial ball position
	g.Ball.SpeedX, g.Ball.SpeedY = 3, 3            // Initial ball speed
	g.Racket.X = screenWidth/2 - racketWidth/2
	g.Racket.Speed = 5
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	// Logic for start the game over. Space to start
	if g.State == 0 {
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			os.Exit(0)
		}

		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.State = 1
		}

		return nil
	}

	// Logic for exit. Escape to finish
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	// Logic for re-start once game over. Space to re-start
	if g.State == 2 {
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			os.Exit(0)
		}

		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.Initialize()
		}

		return nil
	}

	// Logic for the match, g.State = 1

	// Move the racket
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.Racket.X -= g.Racket.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.Racket.X += g.Racket.Speed
	}

	// Move the ball
	g.Ball.X += g.Ball.SpeedX
	g.Ball.Y += g.Ball.SpeedY

	// Ball above to the field line (not hitting racket or missing yet)
	if (g.Ball.Y + 5) < (screenHeight - racketHeight - 5) {
		// Ball hits the walls
		if (g.Ball.X-5 <= 0) || (g.Ball.X+5) >= (screenWidth-1) {
			g.Ball.SpeedX = -g.Ball.SpeedX
		}

		// Ball hits the ceil
		if (g.Ball.Y - 5) <= 0 {
			g.Ball.SpeedY = -g.Ball.SpeedY
		}

		return nil
	}

	// Ball below the field line, hits or misses the racket
	// Hits the racket, else misses
	if g.Ball.X >= g.Racket.X && g.Ball.X <= (g.Racket.X+racketWidth) {
		g.Ball.SpeedX = accelerate(g.Ball.SpeedX)
		g.Ball.SpeedY = accelerateNRevers(g.Ball.SpeedY)
		g.Score.Player++
	} else {
		g.Ball.X, g.Ball.Y = rand.IntN(screenWidth), 5
		g.Score.CPU++
	}

	if g.Score.Player == winnerScore || g.Score.CPU == winnerScore {
		g.State = 2
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Debug: fmt.Println(g.Ball.X, g.Ball.Y, g.Ball.SpeedX, g.Ball.SpeedY)
	// Draw initial screen, pre-start
	if g.State == 0 {
		vector.DrawFilledCircle(screen, screenWidth/2, screenHeight/2, 5, color.RGBA{R: 0, G: 255, B: 0, A: 0}, true)
		vector.DrawFilledRect(screen, float32(g.Racket.X), screenHeight-racketHeight-5, racketWidth, racketHeight, color.White, false)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Press SPACE to star, ESC at any time to finish.\nMove racket with horizontal arrows.\nFirst to score %d wins. Enjoy!", winnerScore))
		return
	}

	// Draw game over
	if g.State == 2 {
		if g.Score.Player == winnerScore {
			ebitenutil.DebugPrint(screen, "Game over. You won!")
		} else {
			ebitenutil.DebugPrint(screen, "Game over. I won!")
		}
		ebitenutil.DebugPrint(screen, "\n\nPress SPACE to star again, ESC to exit.")
		return
	}

	// Draw the match

	// Draw ball
	vector.DrawFilledCircle(screen, float32(g.Ball.X), float32(g.Ball.Y), 5, color.RGBA{R: 0, G: 255, B: 0, A: 0}, true)

	// Draw racket
	vector.DrawFilledRect(screen, float32(g.Racket.X), screenHeight-racketHeight-5, racketWidth, racketHeight, color.White, false)

	// Print score
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d/%d", g.Score.Player, g.Score.CPU))
}

func main() {
	g := Game{}
	g.Initialize()

	ebiten.SetWindowTitle(windowTitle)
	ebiten.SetWindowSize(screenWidth, screenHeight)

	err := ebiten.RunGame(&g)
	if err != nil {
		panic(err)
	}
}

func accelerateNRevers(x int) int {
	if x > 0 {
		return -(x + 1)
	}
	if x < 0 {
		return -(x - 1)
	}
	return 3 // Case x=0, re-start
}

func accelerate(x int) int {
	if x > 0 {
		return x + 1
	}
	if x < 0 {
		return x - 1
	}
	return 3 // Case x=0, re-start
}
