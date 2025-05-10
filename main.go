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
	maxIters     = 1000
	cCoeff       = complex(-0.74543, 0.11301)
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

	CenterX, CenterY float64
	ScaleX, ScaleY   float64

	points      []Point
	needsUpdate bool
}

func NewGame(width, height int) *Game {
	g := &Game{
		width:       width,
		height:      height,
		CenterX:     0,
		CenterY:     0,
		ScaleX:      1,
		ScaleY:      1,
		points:      nil,
		needsUpdate: true,
	}

	return g
}
func (g *Game) Update() error {
	g.handleInput()

	if g.needsUpdate {
		log.Println("Generating fractal...")
		g.generatePoints()
		g.needsUpdate = false
		log.Println("Fractal updated.")
	}

	return nil
}

func (g *Game) handleInput() {
	const scaleStep = 0.3
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.CenterX -= scaleStep * g.ScaleX
		g.needsUpdate = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.CenterX += scaleStep * g.ScaleX
		g.needsUpdate = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.CenterY -= scaleStep * g.ScaleY
		g.needsUpdate = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.CenterY += scaleStep * g.ScaleY
		g.needsUpdate = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.ScaleX *= 1 + scaleStep
		g.ScaleY *= 1 + scaleStep
		g.needsUpdate = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.ScaleX *= 1 - scaleStep
		g.ScaleY *= 1 - scaleStep
		g.needsUpdate = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		g.needsUpdate = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		os.Exit(0)
	}
}

func (g *Game) generatePoints() {
	g.points = make([]Point, 0, g.width*g.height)

	var (
		centerX float64 = float64(g.width) / 2
		centerY float64 = float64(g.height) / 2
		scaleX  float64 = float64(g.width) / 2
		scaleY  float64 = float64(g.height) / 2
	)

	for x := range g.width {
		for y := range g.height {
			nx := remapPoint(float64(x), centerX, scaleX, g.CenterX, g.ScaleX)
			ny := remapPoint(float64(y), centerY, scaleY, g.CenterY, g.ScaleY)
			col := computeColor(nx, ny, maxIters, cCoeff)
			g.points = append(g.points, Point{x, y, col})
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.points {
		screen.Set(p.X, p.Y, p.Color)
	}

	msg := fmt.Sprintf("TPS: %0.2f FPS: %0.2f\n", ebiten.ActualTPS(), ebiten.ActualFPS())
	msg += fmt.Sprintf("Center: (%v,%v), Scale: (%v,%v)\n", g.CenterX, g.CenterY, g.ScaleX, g.ScaleY)
	msg += fmt.Sprintf("maxIters: %v, cCoeff: %v\n", maxIters, cCoeff)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func remapPoint(x, fromCenter, fromScale, toCenter, toScale float64) float64 {
	return (x-fromCenter)*(toScale/fromScale) + toCenter
}

func computeColor(x, y float64, maxIters int, c complex128) color.RGBA {
	z := complex(x, y)

	for i := range maxIters {
		if cmplx.Abs(z) > 2 {
			r := uint8((1*i + 17) % 255)
			g := uint8((2*i + 01) % 255)
			b := uint8((0*i + 01) % 255)
			return color.RGBA{r, g, b, 255}
		}

		z = z*z + c
	}

	return color.RGBA{}
}

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	ebiten.SetWindowTitle("2D Fractal Viever - Demo")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	game := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
