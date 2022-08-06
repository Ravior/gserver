package gastar

import (
	"fmt"
	"testing"
)

func TestAStar_Find(t *testing.T) {
	// 地图映射数组
	var maps = [5][6]int32{
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
	}

	var canPassFunc = func(x, y int32) bool {
		return maps[y][x] == 0
	}

	var astarObj = NewAStar(5, 6, canPassFunc)

	var findWay = func(startPosX, startPosY, endPosX, endPosY int32) []*Node {
		// 初始化
		var (
			i int32 = 0
			j int32 = 0
		)

		// 设置起点，终点,开始寻路
		start := NewNode(startPosX, startPosY)
		end := NewNode(endPosX, endPosY)

		path := astarObj.Find(start, end)

		// 打印路径图形
		for i = 0; i < 5; i++ {
			for j = 0; j < 6; j++ {
				found := false
				for index := 0; index < len(path); index++ {
					if path[index].X == j && path[index].Y == i {
						found = true
						break
					}
				}

				if found {
					print("1 ")
				} else {
					print("* ")
				}
			}
			print("\n")
		}

		return path
	}

	// 设置起点，终点,开始寻路
	path := findWay(2, 2, 4, 3)
	// 打印一共花费步数
	fmt.Printf("path step is %v\n\n", len(path))
	if len(path) != 3 {
		t.Fail()
	}

	// 设置起点，终点,开始寻路
	path2 := findWay(2, 2, 4, 3)
	// 打印一共花费步数
	fmt.Printf("path step is %v\n\n", len(path2))
	if len(path2) != 3 {
		t.Fail()
	}

	// 设置起点，终点,开始寻路
	path3 := findWay(0, 0, 3, 3)
	// 打印一共花费步数
	fmt.Printf("path step is %v\n\n", len(path3))
	if len(path3) != 6 {
		t.Fail()
	}
}
