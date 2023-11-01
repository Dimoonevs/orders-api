package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Dimoonevs/orders-api.git/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	RedisClient *redis.Client
}

func orderIDKye(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("Marshal: %v", err)
	}

	kye := orderIDKye(order.OrderID)

	txn := r.RedisClient.TxPipeline()
	res := txn.SetNX(ctx, kye, string(data), 0)

	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("SetNX: %v", err)
	}

	if err := txn.SAdd(ctx, "orders", kye).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("SAdd: %v", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		txn.Discard()
		return fmt.Errorf("Exec: %v", err)
	}
	return nil
}

var ErrNotExist = errors.New("Order does not exist")

func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) {
	kye := orderIDKye(id)

	value, err := r.RedisClient.Get(ctx, kye).Result()
	if errors.Is(err, redis.Nil) {
		return model.Order{}, ErrNotExist
	} else if err != nil {
		return model.Order{}, fmt.Errorf("Get: %v", err)
	}

	var order model.Order
	err = json.Unmarshal([]byte(value), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("Unmarshal: %v", err)
	}
	return order, nil
}

func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	kye := orderIDKye(id)
	txn := r.RedisClient.TxPipeline()
	err := txn.Del(ctx, kye).Err()
	if errors.Is(err, redis.Nil) {
		txn.Discard()
		return ErrNotExist
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("Del: %v", err)
	}

	if err := txn.SRem(ctx, "orders", kye).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("SRem: %v", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		txn.Discard()
		return fmt.Errorf("Exec: %v", err)
	}
	return nil
}

func (r *RedisRepo) Update(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)

	if err != nil {
		return fmt.Errorf("Marshal: %v", err)
	}
	kye := orderIDKye(order.OrderID)
	err = r.RedisClient.SetXX(ctx, kye, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("SetXX: %v", err)
	}
	return nil
}

type FindAllPage struct {
	Size   uint64
	Offset uint64
}

type FindAllResult struct {
	Orders []model.Order
	Cursor uint64
}

func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (FindAllResult, error) {
	res := r.RedisClient.SScan(ctx, "orders", page.Offset, "*", int64(page.Size))

	keys, cursor, err := res.Result()
	if err != nil {
		return FindAllResult{}, fmt.Errorf("SScan: %v", err)
	}
	if len(keys) == 0 {
		return FindAllResult{
			Orders: []model.Order{},
		}, nil
	}
	xs, err := r.RedisClient.MGet(ctx, keys...).Result()
	if err != nil {
		return FindAllResult{}, fmt.Errorf("MGet: %v", err)
	}
	//fmt.Println(xs)
	orders := make([]model.Order, len(xs))

	for i, x := range xs {
		x := x.(string)
		var order model.Order
		err := json.Unmarshal([]byte(x), &order)
		if err != nil {
			return FindAllResult{}, fmt.Errorf("Unmarshal: %v", err)
		}
		orders[i] = order
	}
	return FindAllResult{
		Orders: orders,
		Cursor: cursor,
	}, nil
}
