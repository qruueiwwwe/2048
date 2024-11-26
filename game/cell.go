package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"image/color"
	"strconv"
)

const (
	cellSize      = 80 //单个方块的大小
	cellMargin    = 4  //单个方块的边框
	cellAnimation = 6
	maxScale      = 1.2
)

// 绘制单个方块
var (
	cellImage = ebiten.NewImage(cellSize, cellSize)
)

// init 默认为白色
func init() {
	cellImage.Fill(color.White)
}

// Cell 相关结构和操作
type Cell struct {
	Value, OldValue        int
	x, y                   int
	ox, oy, nx, ny         int //旧的位置
	isMove, isBuild, isNew int
}

func (c *Cell) Update() {
	if c.isNew > 0 {
		c.isNew--
	}

	if c.isMove > 0 {
		c.isMove--
	}

	if c.isBuild > 0 {
		c.isBuild--
	}

}

// Draw 绘制所属分格的底色和动画
func (c *Cell) Draw(boardImage *ebiten.Image) {

	//创建一个新的图片

	op := &ebiten.DrawImageOptions{}
	//加点特技
	c.animation(op)

	op.GeoM.Translate(float64(c.x), float64(c.y))
	op.ColorScale.ScaleWithColor(tileBackgroundColor(c.Value))
	boardImage.DrawImage(cellImage, op)

	if c.Value == 0 {
		return
	}
	str := strconv.Itoa(c.Value)
	f := getFont(str)
	w := font.MeasureString(f, str).Floor()
	h := (f.Metrics().Ascent + f.Metrics().Descent).Floor()
	x := c.x + (cellSize-w)/2
	y := c.y + (cellSize-h)/2 + f.Metrics().Ascent.Floor()
	text.Draw(boardImage, str, f, x, y, tileColor(c.Value))

}

func (c *Cell) animation(do *ebiten.DrawImageOptions) {
	//新增
	if c.isNew > 0 {
		rate := 1 - float64(c.isNew)/float64(cellAnimation)
		scale := meanF(0.0, 1.0, rate)
		do.GeoM.Translate(float64(-cellSize/2), float64(-cellSize/2))
		do.GeoM.Scale(scale, scale)
		do.GeoM.Translate(float64(cellSize/2), float64(cellSize/2))
	}

	//构建
	if c.isBuild > 0 {
		rate := float64(c.isBuild) / float64(cellAnimation*2/3)
		if cellAnimation*2/3 <= c.isBuild {
			rate = 1 - float64(c.isBuild-2*cellAnimation/3)/float64(cellAnimation/3)
		}

		scale := meanF(1.0, maxScale, rate)
		do.GeoM.Translate(float64(-cellSize/2), float64(-cellSize/2))
		do.GeoM.Scale(scale, scale)
		do.GeoM.Translate(float64(cellSize/2), float64(cellSize/2))
	}

	if c.isMove > 0 {
		if c.isMove == 1 {
			c.x = c.nx
			c.y = c.ny
		} else {
			rate := 1 - float64(c.isMove)/cellAnimation
			c.x = mean(c.ox, c.nx, rate)
			c.ox = c.x
			c.y = mean(c.oy, c.ny, rate)
			c.oy = c.y
		}
	}
}

func mean(a, b int, rate float64) int {
	return int(float64(a)*(1-rate) + float64(b)*rate)
}

func meanF(a, b float64, rate float64) float64 {
	return a*(1-rate) + b*rate
}
