package main

import (
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 开始启动监听key变化监听
	go watch()

	bind, b := config.Get("", "bind")
	if !b {
		bind = ":3002"
	}

	http.HandleFunc("/", requestID)
	http.HandleFunc("/*	", requestID)

	// 开始http处理
	err := http.ListenAndServe(bind, nil)

	if err != nil {
		panic(err)
	}
}

// 处理用户的网络请求
func requestID(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	name := r.Form.Get("name")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("name can`t be empty"))
		return
	}

	num := 1

	if tmpNum := r.Form.Get("num"); tmpNum != "" {
		tmpNum2, err := strconv.Atoi(tmpNum)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("num not a int"))
			return
		}
		num = tmpNum2
	}

	if tmpMaxRequestNum, b := config.Get("", "max_request_num"); b {
		if tmpMaxRequestNumInt, err := strconv.Atoi(tmpMaxRequestNum); err == nil {
			if num > tmpMaxRequestNumInt {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("num more than max request id num"))
				return
			}
		} else {
			errMessage := err.Error()
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errMessage))
			return
		}
	}

	if value, exists := watchList[name]; exists {
		w.WriteHeader(http.StatusOK)
		if num > 1 {
			arr := make([]string, 0)
			for i := 0; i < num; i++ {
				//  从通信队列中获取数据时添加5秒的超时时间
				select {
				case rs := <-value:
					arr = append(arr, strconv.FormatInt(rs, 10))
				case <-time.After(5 * time.Second):
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("server create id timeout"))
				}

			}
			w.Write([]byte(strings.Join(arr, ",")))
		} else {
			//  从通信队列中获取数据时添加5秒的超时时间
			select {
			case rs := <-value:
				w.Write([]byte(strconv.FormatInt(rs, 10)))
			case <-time.After(5 * time.Second):
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("server create id timeout"))
			}

		}
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("name can`t found"))
}
