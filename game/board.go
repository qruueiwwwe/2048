package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type task func() error

// Board represents the game board.
type Board struct {
	size  int
	tasks []task
}

// NewBoard 初始化
func NewBoard(size int) *Board {
	b := &Board{
		size: size,
	}
	return b
}

// Size 基数*每个方格的大小再加上边框大小
func (b *Board) Size() (int, int) {
	x := b.size*cellSize + (b.size+1)*cellMargin
	y := x
	return x, y
}

// Draw draws the board to the given boardImage.
func (b *Board) Draw(boardImage, score *ebiten.Image, cell *Arr) {
	//添加第二层底色
	boardImage.Clear()
	boardImage.Fill(frameColor)
	b.drawScore(score, cell)
	//绘制栅栏 和 文字
	cell.isMoving = false
	for i := 0; i < len(cell.arr); i++ {
		for j := 0; j < len(cell.arr[i]); j++ {
			v := cell.arr[i][j]
			if v.isMove > 0 {
				cell.isMoving = true
			}
			//绘制方块和文字
			v.Draw(boardImage)
		}
	}
}

// drawScore 绘制得分板
func (b *Board) drawScore(score *ebiten.Image, cell *Arr) {
	//增加当前分数
	if cell.f {
		fo := SmallFont
		strO := fmt.Sprintf("当前分数：%d", cell.score)
		wo := font.MeasureString(fo, strO).Floor()
		ho := (fo.Metrics().Ascent + fo.Metrics().Descent).Floor()
		bw, bh := score.Bounds().Dx(), score.Bounds().Dy()
		x1 := (bw - wo) / 2
		y1 := (bh-ho)/2 + fo.Metrics().Ascent.Floor()
		score.Clear()
		text.Draw(score, strO, fo, x1, y1, wColor)
	} else {
		fo := SmallFont
		strO := fmt.Sprintf("Game Over,Score:%d", cell.score)
		wo := font.MeasureString(fo, strO).Floor()
		ho := (fo.Metrics().Ascent + fo.Metrics().Descent).Floor()
		bw, bh := score.Bounds().Dx(), score.Bounds().Dy()
		x1 := (bw - wo) / 2
		y1 := (bh-ho)/2 + fo.Metrics().Ascent.Floor()
		score.Clear()
		text.Draw(score, strO, fo, x1, y1, rColor)

		str1 := fmt.Sprintf("Press ESC to restart")
		w1 := font.MeasureString(fo, str1).Floor()
		h1 := (fo.Metrics().Ascent + fo.Metrics().Descent).Floor() - ho - 20
		bw1, bh1 := score.Bounds().Dx(), score.Bounds().Dy()
		x2 := (bw1 - w1) / 2
		y2 := (bh1-h1)/2 + fo.Metrics().Ascent.Floor()
		text.Draw(score, str1, fo, x2, y2, wColor)
	}

}
