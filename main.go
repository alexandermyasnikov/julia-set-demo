package main

import (
	"fmt"
	"image/color"
	"log"
	"math/cmplx"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 640
	maxIt        = 128
)

type Point struct {
	X, Y  int
	Color color.RGBA
}

type ImageData struct {
	Points []int
}

type Game struct {
	width  int
	height int

	x0, y0, dx0, dy0 float64

	points []Point
}

func NewGame(width, height int) *Game {
	g := &Game{
		width:  width,
		height: height,
		x0:     0,
		y0:     0,
		dx0:    1,
		dy0:    1,
		points: nil,
	}

	return g
}
func (g *Game) Update() error {
	if g.points == nil {
		log.Println("update")
		for x := range g.width {
			for y := range g.height {
				xx := norm(float64(x), float64(g.width)/2, float64(g.width)/2, g.x0, g.dx0)
				yy := norm(float64(y), float64(g.width)/2, float64(g.width)/2, g.y0, g.dy0)

				col := calc(xx, yy)
				g.points = append(g.points, Point{x, y, col})
			}
		}
		log.Println("update ok")
	}

	scaleSize := 0.5
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.points = nil
		g.x0 -= scaleSize * g.dx0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.points = nil
		g.x0 += scaleSize * g.dx0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.points = nil
		g.y0 -= scaleSize * g.dy0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.points = nil
		g.y0 += scaleSize * g.dy0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.points = nil
		g.dx0 += scaleSize * g.dx0
		g.dy0 += scaleSize * g.dy0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.points = nil
		g.dx0 -= scaleSize * g.dx0
		g.dy0 -= scaleSize * g.dy0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		g.points = nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		os.Exit(1)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, point := range g.points {
		screen.Set(point.X, point.Y, point.Color)
	}

	msg := fmt.Sprintf("TPS: %0.2f FPS: %0.2f\n", ebiten.ActualTPS(), ebiten.ActualFPS())
	msg += fmt.Sprintf("p=(%v,%v) d=(%v,%v)\n", g.x0, g.y0, g.dx0, g.dy0)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("2d Demo")

	game := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func norm(x, x0, dx0, x1, dx1 float64) float64 {
	x -= x0
	x *= dx1 / dx0
	x += x1
	return x
}

func calc(x, y float64) color.RGBA {
	c := complex(-0.74543, 0.11301)
	z := complex(x, y)
	maxIters := 1000

	for i := range maxIters {
		if cmplx.Abs(z) > 2 {
			r := uint8((1*i + 17) % 255)
			g := uint8((2*i + 01) % 255)
			b := uint8((0*i + 01) % 255)
			return color.RGBA{r, g, b, 0}
		}
		z = z*z + c
	}

	return color.RGBA{}
}
