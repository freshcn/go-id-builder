package main

import (
	"fmt"
	"time"
)

/*
监听服务器中的key列表的变化情况，当有新的列表被添加进入来的时候，产生新的数据
*/

// id值的队列
type idChan chan int64

// 正在被监听的key列表
var watchList = make(map[string]idChan)

// 监听协程
func watch() {

	for {
		makeChan()
		for key, value := range watchList {
			if int64(len(value)) < getPreStep()/3 {
				var max, step int64 = 0, 0
				var err error
				for {
					max, step, err = updateID(key)
					if err != nil {
						fmt.Println(err)
						time.Sleep(5 * time.Millisecond)
					} else {
						break
					}
				}

				if max > 0 {
					for i := int64(max - step + 1); i <= max; i++ {
						value <- i
					}
				}
			}
		}
		time.Sleep(3 * time.Second)
	}

}

// 为数据创建新的监听队列
func makeChan() {
	var list map[string]int64
	var err error
	for {
		list, err = idsList()
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Millisecond)
		} else {
			break
		}
	}

	// 将不存在于当前监听列表中的数据添加到监听列表中
	for key := range list {
		if _, exists := watchList[key]; !exists {
			watchList[key] = make(idChan, getPreStep()*2)
		}
	}

	// 检查存在于当前监听列表中，却已经不再存在数据中的ID
	for key := range watchList {
		if _, exists := list[key]; !exists {
			delete(watchList, key)
		}
	}
}
