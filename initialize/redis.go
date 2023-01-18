package initialize

import (
	"context"
	"fmt"
	"metalflow/pkg/global"
	"time"

	"github.com/go-redis/redis"
)

// Redis 初始化redis数据库
func Redis() {
	init := false
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(global.Conf.System.ConnectTimeout)*time.Second)
	defer cancel()
	go func() {
		for { //nolint:gosimple
			select {
			case <-ctx.Done():
				if !init {
					panic(fmt.Sprintf("初始化redis异常: 连接超时(%ds)", global.Conf.System.ConnectTimeout))
				}
				// 此处需return避免协程空跑
				return
			}
		}
	}()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.Conf.Redis.Host, global.Conf.Redis.Port),
		DB:       global.Conf.Redis.Database,
		Password: global.Conf.Redis.Password,
	})
	err := client.Ping().Err()
	if err != nil {
		panic(fmt.Sprintf("初始化redis异常: %v", err))
	}
	init = true
	global.Redis = client
	global.Log.Info("初始化redis完成")
}
