package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// 工具函数集

// 返回配置中的每次最大申请ID值
func getPreStep() int64 {
	var preStep int64 = 1000
	tmpPreStep, b := config.Get("", "per_step")
	if b {
		intPreStep, err := strconv.ParseInt(tmpPreStep, 10, 64)
		if err != nil {
			return preStep
		}
		preStep = intPreStep
	}
	return preStep
}

// 返回当前的运行目录
func runPath() string {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}

	return filepath.Dir(path)
}
