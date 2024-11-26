package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
)

const (
	ScreenWidth  = 420
	ScreenHeight = 600
	boardSize    = 4
)

// Game represents a game state.
type Game struct {
	board      *Board
	arr        *Arr
	boardImage *ebiten.Image
	scoreImage *ebiten.Image
}

func Start() {
	c := NewArr()
	g := &Game{
		board: NewBoard(boardSize),
		arr:   c,
	}
	//将TPS 设置为1 协助开发 每秒tick数
	ebiten.SetTPS(30)
	ebiten.SetWindowTitle("2048游戏")
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	//2048  启动
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	//更新动画
	tmp := g.arr.getArr()
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			tmp[i][j].Update()
		}
	}
	//监听用户操作
	g.listenKey()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	if g.boardImage == nil {
		g.boardImage = ebiten.NewImage(g.board.Size())
	}

	if g.scoreImage == nil {
		x, _ := g.board.Size()
		g.scoreImage = ebiten.NewImage(x, x/3)
	}
	//绘制背景色
	screen.Fill(backgroundColor)
	g.board.Draw(g.boardImage, g.scoreImage, g.arr)

	//计算底框的位置
	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	bw, bh := g.boardImage.Bounds().Dx(), g.boardImage.Bounds().Dy()
	x := (sw - bw) / 2
	y := (sh-bh)/2 + 60
	op.GeoM.Translate(float64(x), float64(y))

	//计算得分框的位置
	opScore := &ebiten.DrawImageOptions{}
	ScoreY := (sh-bh)/2 - 60
	opScore.GeoM.Translate(float64(x), float64(ScoreY))

	//渲染底框
	screen.DrawImage(g.boardImage, op)
	screen.DrawImage(g.scoreImage, opScore)
}

func (g *Game) listenKey() {
	if g.arr.isMoving == true {
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.arr.vertical(true)
		g.arr.f = g.arr.checkGameOver()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.arr.vertical(false)
		g.arr.f = g.arr.checkGameOver()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.arr.horizontal(true)
		g.arr.f = g.arr.checkGameOver()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.arr.horizontal(false)
		g.arr.f = g.arr.checkGameOver()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if g.arr.f == false {
			g.arr.retry()
		}
	}
}
