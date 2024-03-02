package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

const (
	screenWidth   = 320
	screenHeight  = 240
	circleRadius  = 16
	paddleWidth   = 8
	paddleHeight  = 48
	movementSpeed = 4
	speedIncrease = 0.00001 // Adjust this value to control the speed increase rate
)

type Paddle struct {
	position ebiten.GeoM
}

func (p *Paddle) Update(upKey, downKey ebiten.Key) {
	if ebiten.IsKeyPressed(upKey) && p.position.Element(1, 2)-paddleHeight/2 > 0 {
		p.position.Translate(0, -movementSpeed)
	}
	if ebiten.IsKeyPressed(downKey) && p.position.Element(1, 2)+paddleHeight/2 < screenHeight {
		p.position.Translate(0, movementSpeed)
	}
}

func (p *Paddle) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.position.Element(0, 2), p.position.Element(1, 2))

	// Draw the paddle
	ebitenutil.DrawRect(screen, p.position.Element(0, 2)-paddleWidth/2, p.position.Element(1, 2)-paddleHeight/2, float64(paddleWidth), float64(paddleHeight), color.RGBA{0x80, 0xa0, 0xc0, 0xff})
}

type Ball struct {
	position ebiten.GeoM
	speedX   float64
	speedY   float64
}

func (b *Ball) Update(player, enemy Paddle, score *int, speedMultiplier *float64) {
	b.position.Translate(b.speedX, b.speedY)

	// Bounce off the top and bottom walls
	if b.position.Element(1, 2)-circleRadius/2 < 0 || b.position.Element(1, 2)+circleRadius/2 > screenHeight {
		b.speedY = -b.speedY
	}

	// Bounce off the paddles
	if b.position.Element(0, 2)-circleRadius/2 < player.position.Element(0, 2)+paddleWidth/2 &&
		b.position.Element(0, 2)+circleRadius/2 > player.position.Element(0, 2)-paddleWidth/2 &&
		b.position.Element(1, 2)-circleRadius/2 < player.position.Element(1, 2)+paddleHeight/2 &&
		b.position.Element(1, 2)+circleRadius/2 > player.position.Element(1, 2)-paddleHeight/2 {
		b.speedX = -b.speedX
	}

	if b.position.Element(0, 2)-circleRadius/2 < enemy.position.Element(0, 2)+paddleWidth/2 &&
		b.position.Element(0, 2)+circleRadius/2 > enemy.position.Element(0, 2)-paddleWidth/2 &&
		b.position.Element(1, 2)-circleRadius/2 < enemy.position.Element(1, 2)+paddleHeight/2 &&
		b.position.Element(1, 2)+circleRadius/2 > enemy.position.Element(1, 2)-paddleHeight/2 {
		b.speedX = -b.speedX
	}

	// Score point if the ball passes the enemy's paddle
	if b.position.Element(0, 2)+circleRadius/2 > screenWidth {
		*score++
		b.Reset()
	}

	// Reset if the ball goes out of bounds on the player's side
	if b.position.Element(0, 2)-circleRadius/2 < 0 {
		b.Reset()
	}

	// Increase speed over time
	b.speedX *= *speedMultiplier
	b.speedY *= *speedMultiplier
}

func (b *Ball) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.position.Element(0, 2), b.position.Element(1, 2))

	// Draw the ball as a circle
	ebitenutil.DrawRect(screen, b.position.Element(0, 2)-circleRadius/2, b.position.Element(1, 2)-circleRadius/2, float64(circleRadius), float64(circleRadius), color.RGBA{0xff, 0x00, 0x00, 0xff})
}

func (b *Ball) Reset() {
	b.position.SetElement(0, 2, screenWidth/2)
	b.position.SetElement(1, 2, screenHeight/2)
	b.speedX = rand.Float64()*2 - 1 // Random horizontal speed
	b.speedY = rand.Float64()*2 - 1 // Random vertical speed
}

type Game struct {
	count          int
	player         Paddle
	enemy          Paddle
	ball           Ball
	score          int
	scoreX         float64
	scoreY         float64
	scoreOp        ebiten.DrawImageOptions
	speedMultiplier float64
}

var game *Game

func (g *Game) Update() error {
	g.count++

	// Update player's paddle
	g.player.Update(ebiten.KeyW, ebiten.KeyS)

	// Update enemy's paddle
	g.enemy.Update(ebiten.KeyI, ebiten.KeyK)

	// Update ball
	g.ball.Update(g.player, g.enemy, &g.score, &g.speedMultiplier)

	// Increase speed multiplier over time
	g.speedMultiplier += speedIncrease

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)
	g.enemy.Draw(screen)
	g.ball.Draw(screen)

	// Draw score text
	scoreText := fmt.Sprintf("Score: %d", g.score)
	ebitenutil.DebugPrintAt(screen, scoreText, int(g.scoreX), int(g.scoreY))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Pong (Ebitengine Demo)")

	// Initialize the player's paddle with an initial value
	player := Paddle{
		position: ebiten.GeoM{},
	}
	player.position.Translate(screenWidth/6, screenHeight/2)

	// Initialize the enemy's paddle with an initial value
	enemy := Paddle{
		position: ebiten.GeoM{},
	}
	enemy.position.Translate(screenWidth-screenWidth/6, screenHeight/2)

	// Initialize the ball with an initial value
	ball := Ball{
		position: ebiten.GeoM{},
		speedX:   rand.Float64()*2 - 1,
		speedY:   rand.Float64()*2 - 1,
	}
	ball.position.Translate(screenWidth/2, screenHeight/2)

	game = &Game{
		player:          player,
		enemy:           enemy,
		ball:            ball,
		scoreX:          screenWidth / 2,
		scoreY:          20,
		speedMultiplier: 1.0,
	}

	rand.Seed(time.Now().UnixNano()) // Seed for random values

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
