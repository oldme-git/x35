package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Config struct {
	MinInterval  int `json:"minInterval"`  // 最小间隔（秒）
	MaxInterval  int `json:"maxInterval"`  // 最大间隔（秒）
	RingDuration int `json:"ringDuration"` // 响铃持续时间（秒）
	OutDuration  int `json:"outDuration"`  // 程序运行持续时间（秒）
}

func main() {
	configPath := "config.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	config, err := loadConfig(configPath)
	if err != nil {
		fmt.Printf("配置加载错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("\a")
	log("学习开始")
	fmt.Printf("当前配置：最小间隔: %ds 最大间隔: %ds 响铃时长: %ds 持续时长: %ds\n",
		config.MinInterval, config.MaxInterval, config.RingDuration, config.OutDuration)

	timeoutCh := time.After(time.Duration(config.OutDuration) * time.Second)

	for {
		// 生成随机间隔时间
		interval := rand.Intn(config.MaxInterval-config.MinInterval+1) + config.MinInterval
		waitTime := time.Duration(interval) * time.Second

		select {
		case <-timeoutCh:
			fmt.Print("\a")
			time.Sleep(1000 * time.Millisecond)
			fmt.Print("\a")
			log("学习结束")
			os.Exit(1)
		case <-time.After(waitTime):
			log("冥想开始")
			fmt.Print("\a")

			wt := time.Duration(config.RingDuration) * time.Second
			<-time.After(wt)

			fmt.Print("\a")
			log("冥想结束")
		}
	}
}

func log(t string) {
	n := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(n + " " + t)
}

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("无法打开配置文件: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("配置文件解析失败: %v", err)
	}

	// 验证配置有效性
	if config.MinInterval <= 0 || config.MaxInterval <= 0 || config.RingDuration <= 0 || config.OutDuration <= 0 {
		return nil, fmt.Errorf("所有时间参数必须大于0")
	}
	if config.MinInterval > config.MaxInterval {
		return nil, fmt.Errorf("minInterval 不能大于 maxInterval")
	}

	return &config, nil
}
