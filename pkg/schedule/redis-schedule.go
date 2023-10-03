package schedule

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisScheduler struct {
	redis            redis.UniversalClient
	redisScheduleKey string
	member           string
	quit             chan uint8
}

func (r *RedisScheduler) StartAsync(cron string, onExecuting func()) error {

	cronParse, err := CronParse(cron, time.UTC, false)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// 将任务添加到 Redis 有序集合
	_, err = r.redis.ZAdd(ctx, r.redisScheduleKey, redis.Z{
		Score:  float64(cronParse.Next(time.Now().UTC()).Unix()),
		Member: r.member,
	}).Result()

	if err != nil {
		fmt.Println("Failed to add task to Redis:", err)
		return err
	}

	// 循环检查定时任务
	go func() {

		for {
			// 获取需要执行的任务
			results, err := r.redis.ZRangeByScore(ctx, r.redisScheduleKey, &redis.ZRangeBy{
				Min:    "0",
				Max:    fmt.Sprintf("%d", time.Now().UTC().Unix()),
				Offset: 0,
				Count:  1,
			}).Result()
			if err != nil {
				fmt.Println("Failed to retrieve tasks from Redis:", err)
				return
			}

			// 如果有任务需要执行，执行任务并删除任务
			if len(results) > 0 {
				task := results[0]
				fmt.Printf("Executing task: %s\n", task)

				// 在这里执行任务的具体逻辑
				onExecuting()

				// 将任务重新添加到 Redis 有序集合
				_, err = r.redis.ZAdd(ctx, r.redisScheduleKey, redis.Z{
					Score:  float64(cronParse.Next(time.Now().UTC()).Unix()),
					Member: r.member,
				}).Result()

				if err != nil {
					fmt.Println("Failed to restart task from Redis:", err)
					return
				}
			}

			// 等待一段时间再进行下一次检查
			time.Sleep(1 * time.Second)
			// 当需要终止goroutine时
			select {
			case <-r.quit:
				fmt.Println("go func loop is stopped....")
				return
			default:
				continue
			}
		}
	}()
	return nil
}

func (r *RedisScheduler) Stop() error {
	defer close(r.quit)
	_, err := r.redis.ZRem(context.Background(), r.redisScheduleKey, r.member).Result()
	if err != nil {
		return err
	}
	return nil
}

func NewRedisScheduler(addr []string, user, pass, masterName string, db int) IScheduler {

	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:      addr,
		Username:   user,
		Password:   pass,
		DB:         db,
		MasterName: masterName,
	})

	return &RedisScheduler{
		redis:            rdb,
		redisScheduleKey: "redis-schedule-key",
		member:           "timer",
		quit:             make(chan uint8),
	}
}
