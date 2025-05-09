package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
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
	X, Y float64
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
		dx0:    10,
		dy0:    10,
		points: nil,
	}

	return g
}
func (g *Game) Update() error {
	if g.points == nil {
		var x, y float64
		for range 1000000 {
			x, y = calc(x, y)
			g.points = append(g.points, Point{x, y})
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.x0 -= 0.2 * g.dx0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.x0 += 0.2 * g.dx0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.y0 -= 0.2 * g.dy0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.y0 += 0.2 * g.dy0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.dx0 += 0.2 * g.dx0
		g.dy0 += 0.2 * g.dy0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.dx0 -= 0.2 * g.dx0
		g.dy0 -= 0.2 * g.dy0
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
	msg := fmt.Sprintf("TPS: %0.2f FPS: %0.2f\n", ebiten.ActualTPS(), ebiten.ActualFPS())
	msg += fmt.Sprintf("p=(%v,%v) d=(%v,%v)\n", g.x0, g.y0, g.dx0, g.dy0)
	ebitenutil.DebugPrint(screen, msg)

	x1, y1 := float64(g.width)/2, float64(g.height)/2
	dx1, dy1 := float64(g.width)/2, float64(g.height/2)

	m := make(map[[2]int]int, 0)

	for _, point := range g.points {
		x := int(norm(point.X, g.x0, g.dx0, x1, dx1))
		y := int(norm(point.Y, g.y0, g.dy0, y1, dy1))
		m[[2]int{x, y}]++
	}

	var maxCount int
	for _, v := range m {
		maxCount = max(v, maxCount)
	}

	for p, count := range m {
		x := p[0]
		y := p[1]

		iterRatio := math.Pow(float64(count)/float64(maxCount), 0.3)
		r := uint8(80 * iterRatio)
		g := uint8(255 * iterRatio)
		b := uint8(30 * iterRatio)
		screen.Set(x, y, color.RGBA{r, g, b, 0})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("2d Demo")

	game := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func calc(x, y float64) (float64, float64) {
	var a, b, c, d, e, f float64

	r := rand.Float64()
	switch {
	case r < 0.01:
		a, b, c, d, e, f = 0, 0, 0, 0.16, 0, 0
	case r < 0.01+0.85:
		a, b, c, d, e, f = 0.85, 0.04, -0.04, 0.85, 0, 1.6
	case r < 0.01+0.85+0.07:
		a, b, c, d, e, f = 0.20, -0.26, 0.23, 0.22, 0, 1.6
	case r < 0.01+0.85+0.07+0.07:
		a, b, c, d, e, f = -1.15, 0.28, 0.26, 0.24, 0, 0.44
	}

	x1 := a*x + b*y + e
	y1 := c*x + d*y + f

	return x1, y1
}

func norm(x, x0, dx0, x1, dx1 float64) float64 {
	x -= x0
	x *= dx1 / dx0
	x += x1
	return x
}
