package main

import (
	"context"
	"os"

	"github.com/go-redis/redis/v9"
)

var redisClient *redis.Client

func initializeRedis() {
	opt, err := redis.ParseURL(os.Getenv("REDIS_CONNECTION_URL"))

	if err != nil {
		panic(err)
	}

	redisClient = redis.NewClient(opt)
}

func addStoryLineToDb(roomId string, storyLine string) error {
	ctx := context.Background()

	_, err := redisClient.RPush(ctx, roomId, storyLine).Result()

	if err != nil {
		return err
	}

	return nil
}

func getStoryLinesFromDb(roomId string) ([]string, error) {
	ctx := context.Background()

	var res []string

	storyLength, err := redisClient.LLen(ctx, roomId).Result()

	if err != nil {
		return res, err
	}

	storyLines, err := redisClient.LRange(ctx, roomId, 0, storyLength).Result()

	if err != nil {
		return res, err
	}

	res = storyLines

	return res, nil
}
