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
	fmt.Println(config)
	// ...后续代码不变

	for {
		// 生成随机间隔时间（180-300秒之间）
		interval := rand.Intn(config.MaxInterval-config.MinInterval+1) + config.MinInterval
		waitTime := time.Duration(interval) * time.Second

		fmt.Printf("下一次响铃将在 %v 后触发...\n", waitTime)
		<-time.After(waitTime)

		fmt.Println("🔔 响铃开始！")

		fmt.Print("\a")
		wt := time.Duration(config.RingDuration) * time.Second
		<-time.After(wt)
		fmt.Print("\a")
		fmt.Println("🛑 响铃结束")
	}
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
	if config.MinInterval <= 0 || config.MaxInterval <= 0 || config.RingDuration <= 0 {
		return nil, fmt.Errorf("所有时间参数必须大于0")
	}
	if config.MinInterval > config.MaxInterval {
		return nil, fmt.Errorf("minInterval 不能大于 maxInterval")
	}

	return &config, nil
}
