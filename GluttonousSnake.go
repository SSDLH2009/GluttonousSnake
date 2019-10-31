package main

import (
	"container/list"
	"fmt"
	"github.com/jan-bar/golibs"
	"math/rand"
	"time"
)

const (
	LengthX  = 20 /* 矩阵行数 */
	LengthY  = 40 /* 矩阵列数 */
	NoneFlag = 0  /* 空白标识 */
	FoodFlag = 1  /* 食物标识 */
	BodyFlag = 2  /* 蛇身标识 */
)

var (
	MatrixData [LengthX][LengthY]byte
	api        *golibs.Win32Api
)

/**
初始化
*/
func init() {
	api = golibs.NewWin32Api()
	api.Clear()                      //清屏
	rand.Seed(time.Now().UnixNano()) //初始化随机数种子
	api.SetWindowText("go版本贪食蛇")
	api.CenterWindowOnScreen(550, 400) //居中显示 设置窗口大小
	api.ShowHideCursor(false)

	var i int

	for i = 1; i <= LengthY; i++ {
		api.GotoXY(i, 0)
		fmt.Print("-") //画上边横线
		api.GotoXY(i, LengthX)
		fmt.Print("-") //画下边横线
	}
	for i = 1; i < LengthX; i++ {
		api.GotoXY(0, i)
		fmt.Print("|") //画左边竖线
		api.GotoXY(LengthY+1, i)
		fmt.Print("|") //画右边竖线
	}
}

/**
主程序
*/
func main() {
	var (
		headPos       = golibs.Coord{LengthX / 2, LengthY / 2}
		dirNew  int32 = golibs.KeyUp
		dirOld int32 = golibs.KeyDown
		direction = make(chan int32)
		tLoop     = time.NewTimer(time.Second)
		Snake     = list.New()
		score     = 0
	)
	DrawSnake(headPos, BodyFlag)
	Snake.PushFront(headPos)
	RandFood()
	PrintScore(score)

	go func() {
		for {
			direction <- golibs.WaitKeyBoard()
		}
	}()

	for {
		headPos = Snake.Front().Value.(golibs.Coord)
		select {
		case dir, ok := <-direction:
			tLoop.Stop()
			if ok {
				dirNew = dir
			}
		case <-tLoop.C:
		}
		tLoop.Reset(time.Second)

		switch dirNew {
		case golibs.KeyUp: /* 新方向朝上 */
			if dirOld == golibs.KeyDown {
				dirNew = golibs.KeyDown
				headPos.X++ /* 旧方向如果朝下,则继续向下走 */
			} else { /* 旧方向为: 左右上 */
				dirOld = golibs.KeyUp
				headPos.X--
			}
		case golibs.KeyDown:
			if dirOld == golibs.KeyUp {
				dirNew = golibs.KeyUp
				headPos.X-- /* 旧方向如果朝上,则继续向上走 */
			} else { /* 旧方向为: 左右下 */
				dirOld = golibs.KeyDown
				headPos.X++
			}
		case golibs.KeyLeft:
			if dirOld == golibs.KeyRight {
				dirNew = golibs.KeyRight
				headPos.Y++ /* 旧方向如果朝右,则继续向右走 */
			} else { /* 旧方向为: 上下左 */
				dirOld = golibs.KeyLeft
				headPos.Y--
			}
		case golibs.KeyRight:
			if dirOld == golibs.KeyLeft {
				dirNew = golibs.KeyLeft
				headPos.Y-- /* 旧方向如果朝左,则继续向左走 */
			} else { /* 旧方向为: 上下右 */
				dirOld = golibs.KeyRight
				headPos.Y++
			}
		}

		if headPos.X < 0 || headPos.Y < 0 || headPos.X >= LengthX || headPos.Y >= LengthY || MatrixData[headPos.X][headPos.Y] == BodyFlag {
			api.GotoXY(LengthY+3, 7)
			fmt.Print("你输了,按回车退出!")
			break /* 越界了,这个位置是蛇身,游戏结束 */
		}
		if MatrixData[headPos.X][headPos.Y] == NoneFlag {
			var end = Snake.Back()                        /* 如果头部为空白,则尾部也要去掉一个 */
			DrawSnake(end.Value.(golibs.Coord), NoneFlag) /* 将蛇尾置位为空白 */
			Snake.Remove(end)                             /* 下一步是空白,尾部删除一个蛇尾 */
		} else { /* 吃到一个食物,则需要再初始化一个食物 */
			score++ /* 吃一个食物加1分 */
			PrintScore(score)
			RandFood()
		}
		Snake.PushFront(headPos) /* 头部加入链表 */
		DrawSnake(headPos, BodyFlag)
	}

	fmt.Scanln() /* 避免一闪而逝 */
}

/**
打印分数
*/
func PrintScore(score int) {
	api.GotoXY(LengthY+3, 5)
	fmt.Print("分数:", score)
}

/**
* 在指定位置画整个界面数据
* 包括空白,蛇,食物
 */
func DrawSnake(pos golibs.Coord, dType byte) {
	MatrixData[pos.X][pos.Y] = dType
	api.GotoXY(pos.Y+1, pos.X+1)
	switch dType {
	case NoneFlag:
		fmt.Print(" ")
	case BodyFlag:
		fmt.Print("#")
	case FoodFlag:
		fmt.Print("+")
	}
}

/**
* 在空白位置
* 随机画一个食物
**/
func RandFood() {
	var i, j int
	for {
		i = rand.Intn(LengthX)
		j = rand.Intn(LengthY)
		if NoneFlag == MatrixData[i][j] {
			DrawSnake(golibs.Coord{i, j}, FoodFlag)
			break
		}
	}
}
