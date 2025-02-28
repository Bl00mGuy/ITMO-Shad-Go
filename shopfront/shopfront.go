//go:build !solution

package shopfront

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type ShopCounters struct {
	redisClient *redis.Client
}

func New(redisClient *redis.Client) Counters {
	return ShopCounters{redisClient: redisClient}
}

func (c ShopCounters) RecordView(ctx context.Context, itemID ItemID, userID UserID) error {
	itemKey := c.generateItemKey(itemID)
	viewCountKey := c.generateViewCountKey(itemID)
	userIdentifier := int64ToString(int64(userID))
	pipe := c.redisClient.TxPipeline()

	pipe.Incr(ctx, viewCountKey)
	pipe.SAdd(ctx, itemKey, userIdentifier)

	_, err := pipe.Exec(ctx)
	return err
}

func (c ShopCounters) GetItems(ctx context.Context, itemIDs []ItemID, userID UserID) ([]Item, error) {
	itemKeys := c.generateItemKeys(itemIDs)
	viewCountKeys := c.generateViewCountKeys(itemIDs)
	userIdentifier := int64ToString(int64(userID))
	pipe := c.redisClient.Pipeline()
	viewCountCommands := make([]*redis.StringCmd, len(viewCountKeys))
	viewedCommands := make([]*redis.BoolCmd, len(itemKeys))

	for i, itemKey := range itemKeys {
		viewCountKey := viewCountKeys[i]
		viewCountCommands[i] = pipe.Get(ctx, viewCountKey)
		viewedCommands[i] = pipe.SIsMember(ctx, itemKey, userIdentifier)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	return c.buildItemsFromRedisResults(viewCountCommands, viewedCommands)
}

func (c ShopCounters) buildItemsFromRedisResults(viewCountCommands []*redis.StringCmd, viewedCommands []*redis.BoolCmd) ([]Item, error) {
	items := make([]Item, len(viewCountCommands))

	for i := range items {
		viewCount, err := viewCountCommands[i].Int()
		if err != nil && err != redis.Nil {
			return nil, err
		}
		isViewed := viewedCommands[i].Val()
		items[i] = Item{
			ViewCount: viewCount,
			Viewed:    isViewed,
		}
	}
	return items, nil
}

func (c ShopCounters) generateItemKey(itemID ItemID) string {
	return "item_" + int64ToString(int64(itemID))
}

func (c ShopCounters) generateItemKeys(itemIDs []ItemID) []string {
	keys := make([]string, len(itemIDs))
	for i, itemID := range itemIDs {
		keys[i] = c.generateItemKey(itemID)
	}
	return keys
}

func (c ShopCounters) generateViewCountKey(itemID ItemID) string {
	return "item_" + int64ToString(int64(itemID)) + "_count"
}

func (c ShopCounters) generateViewCountKeys(itemIDs []ItemID) []string {
	keys := make([]string, len(itemIDs))
	for i, itemID := range itemIDs {
		keys[i] = c.generateViewCountKey(itemID)
	}
	return keys
}

func int64ToString(number int64) string {
	return strconv.FormatInt(number, 10)
}
