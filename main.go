package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth    = 1024
	screenHeight   = 768
	maxIters       = 10000
	scaleSaveImage = 5

	animation = true
)

var (
	// cCoeff = complex(-0.74543, 0.11301)
	// cCoeff = complex(-0.8, 0.156)
	// cCoeff = complex(0.285, 0.01)
	// cCoeff = complex(-0.008, 0.71)
	cCoeff = complex(-0.008, 0.85)
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

	if animation {
		r1 := 0.0003 * (2*rand.Float64() - 1)
		r2 := 0.0003*(2*rand.Float64()-1) - 0.0001
		cCoeff += complex(r1, r2)
		g.needsUpdate = true
	}

	if g.needsUpdate {
		log.Println("Generating fractal...")
		g.points = generatePoints(g.CenterX, g.CenterY, g.ScaleX, g.ScaleY, g.width, g.height)
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
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		log.Println("Saving image ...")
		go func() {
			if err := saveImage(g.CenterX, g.CenterY, g.ScaleX, g.ScaleY, g.width, g.height); err != nil {
				log.Printf("Failed to save image: %v", err)
			} else {
				log.Printf("Image saved")
			}
		}()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		os.Exit(0)
	}
}

func saveImage(toCenterX, toCenterY, toScaleX, toScaleY float64, width, height int) error {
	width *= scaleSaveImage
	height *= scaleSaveImage

	points := generatePoints(toCenterX, toCenterY, toScaleX, toScaleY, width, height)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for _, p := range points {
		img.Set(p.X, p.Y, p.Color)
	}

	filename := fmt.Sprintf("/tmp/fractal_%s.png", time.Now().Format("20060102-150405"))
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

func generatePoints(toCenterX, toCenterY, toScaleX, toScaleY float64, width, height int) []Point {
	points := make([]Point, 0, width*height)

	var (
		fromCenterX = float64(width) / 2
		fromCenterY = float64(height) / 2
		fromScaleX  = float64(width) / 2
		fromScaleY  = float64(height) / 2
	)

	for x := range width {
		for y := range height {
			nx := remapPoint(float64(x), fromCenterX, fromScaleX, toCenterX, toScaleX)
			ny := remapPoint(float64(y), fromCenterY, fromScaleY, toCenterY, toScaleY)
			col := computeColor(nx, ny, maxIters, cCoeff)
			points = append(points, Point{x, y, col})
		}
	}

	return points
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
			r := uint8((3*i + 17) % 0xFF)
			g := uint8((2*i + 01) % 0xFF)
			b := uint8((1*i + 01) % 0xFF)
			return color.RGBA{r, g, b, 0xFF}
		}

		z = z*z + c
	}

	return color.RGBA{}
}

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	ebiten.SetWindowTitle("2D Fractal Viever - Demo - amyasnikov.com")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	game := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
