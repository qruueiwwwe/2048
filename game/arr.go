package game

import (
	"math/rand"
	"sync"
)

type Arr struct {
	mu       sync.RWMutex //避免出现并发问题
	arr      [4][4]*Cell
	score    int64
	f        bool
	isMoving bool
}

func NewArr() *Arr {
	c := Arr{}
	c.init()
	return &c
}

func (c *Arr) init() {
	c.clear()
	c.fillEmptyCell()
	c.fillEmptyCell()
}

func (c *Arr) retry() {
	c.clear()
	c.score = 0
	c.fillEmptyCell()
	c.fillEmptyCell()
}

func (c *Arr) getArr() [4][4]*Cell {
	return c.arr
}

func (c *Arr) clear() {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			//计算每个方块的位置
			tmp := &Cell{}
			tmp.y = i*cellSize + (i+1)*cellMargin
			tmp.x = j*cellSize + (j+1)*cellMargin
			c.arr[i][j] = tmp
		}
	}
	c.f = true
}

// fillEmptyCell 找到当前一个空位置，并塞入一个值
func (c *Arr) fillEmptyCell() {
	c.mu.Lock()
	defer c.mu.Unlock()
	type cell [2]int
	var list [16]cell
	var count = 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if c.arr[i][j].Value == 0 {
				list[count][0] = i
				list[count][1] = j
				count++
			}
		}
	}

	if count > 0 {
		l := list[rand.Intn(count)]
		tmp := c.arr[l[0]][l[1]]
		tmp.Value = c.getValue()
		tmp.isMove = 0
		tmp.isBuild = 0
		tmp.isNew = cellAnimation
	}
}

// getValue 生成数字，90%几率生成2，10%几率生成4 可以自己调整
func (c *Arr) getValue() int {
	r := rand.Intn(10)
	var value = 2
	if r >= 9 {
		value = 4
	}
	return value
}

// checkGameOver 游戏结束还是继续进行。
func (c *Arr) checkGameOver() bool {
	var f bool
	//三个条件：
	//1、当前单元格是否有空的；
	//2、当前单元格和右侧单元格是否相等；
	//3、当前单元格和下方单元格是否相等。
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if c.arr[j][i].Value == 0 {
				f = true
				goto Goon
			}
			if j < 3 && c.arr[i][j].Value == c.arr[i][j+1].Value {
				f = true
				goto Goon
			}
			if i < 3 && c.arr[i][j].Value == c.arr[i+1][j].Value {
				f = true
				goto Goon
			}
		}
	}
Goon:
	if f {
		c.fillEmptyCell()
	}
	return f
}

// resize 每次移动后，计算新的Cells。原理： 将上下左右移动 转化为一维数组的左右移动
func (c *Arr) resize(a []*Cell) []int {
	arr := a
	l := len(arr)
	//去零
	newArr := make([]int, 0)
	for _, k := range arr {
		if k.Value == 0 {
			continue
		}
		newArr = append(newArr, k.Value)
	}

	for {
		tmp := make([]int, 0)
		f := true
		for k, v := range newArr {
			if v == 0 {
				continue
			}

			if k == len(newArr)-1 || v != newArr[k+1] {
				tmp = append(tmp, v)
				continue
			}

			newArr[k] = v + newArr[k+1]
			newArr[k+1] = 0
			tmp = append(tmp, newArr[k])
			f = false
		}
		newArr = tmp
		if f {
			break
		}

	}
	ret := make([]int, l)
	for k, v := range newArr {
		ret[k] = v
	}
	return ret
}

// resizeV1 版本
func resizeV1(a []Cell) []Cell {
	arr := a
	//去零
	newArr := make([]Cell, 0)
	for _, v := range arr {
		v.oy = v.y
		v.ox = v.x
		if v.Value == 0 {
			continue
		}
		newArr = append(newArr, v)
	}

	tmp := make([]Cell, 0)
	for k, v := range newArr {
		if v.Value == 0 {
			continue
		}

		if k == len(newArr)-1 || v.Value != newArr[k+1].Value {
			tmp = append(tmp, v)
			continue
		}

		newArr[k].Value = v.Value + newArr[k+1].Value
		newArr[k].isBuild = cellAnimation
		newArr[k+1].Value = 0
		newArr[k].ox = newArr[k+1].x
		newArr[k].oy = newArr[k+1].y
		tmp = append(tmp, newArr[k])
	}
	newArr = tmp

	for k, _ := range a {
		a[k].Value = 0
	}

	for k, v := range newArr {
		a[k].Value = v.Value
		a[k].isBuild = v.isBuild
		a[k].ox = v.ox
		a[k].oy = v.oy
		if (a[k].ox != a[k].x || a[k].oy != a[k].y) && a[k].Value != 0 {
			a[k].isMove = cellAnimation
			a[k].nx = a[k].x
			a[k].ny = a[k].y
		}
	}

	return a
}

// horizontal 水平移动 将单元格数向左或向右移动，移除零并对相邻相同数进行叠加
func (c *Arr) horizontal(toLeft bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := 0; i < 4; i++ {
		//获取每一行的值
		newArr := make([]Cell, 4)
		for j := 0; j < 4; j++ {
			newArr[j] = *c.arr[i][j]
		}

		if toLeft == false {
			tmp := make([]Cell, 4)
			for k, v := range newArr {
				tmp[4-k-1] = v
			}
			newArr = tmp
			resultArr := resizeV1(newArr)
			for j := 0; j < 4; j++ {
				c.arr[i][j] = &resultArr[4-j-1]
			}
		} else {
			resultArr := resizeV1(newArr)
			for j := 0; j < 4; j++ {
				c.arr[i][j] = &resultArr[j]

			}
		}
	}
	c.setScore()

}

// vertical 垂直移动 将单元格数向下或向上移动，移除零并对相邻相同数进行叠加
func (c *Arr) vertical(toTop bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := 0; i < 4; i++ {
		//获取每一列的值
		newArr := make([]Cell, 4)
		for j := 0; j < 4; j++ {
			newArr[j] = *c.arr[j][i]
		}

		if toTop == false {
			tmp := make([]Cell, 4)
			for k, v := range newArr {
				tmp[4-k-1] = v
			}
			newArr = tmp
			resultArr := resizeV1(newArr)
			for j := 0; j < 4; j++ {
				c.arr[j][i] = &resultArr[4-j-1]
			}
		} else {
			resultArr := resizeV1(newArr)
			for j := 0; j < 4; j++ {
				c.arr[j][i] = &resultArr[j]
			}
		}
	}
	c.setScore()
}

func (c *Arr) setScore() {
	arr := c.arr
	c.score = 0
	for i := 0; i < len(c.arr); i++ {
		for j := 0; j < len(c.arr[i]); j++ {
			if arr[i][j].Value < 4 {
				continue
			}
			if arr[i][j].Value > 128 {
				c.score += int64(arr[i][j].Value) * 2
				continue
			}

			c.score += int64(arr[i][j].Value)
		}
	}
}
