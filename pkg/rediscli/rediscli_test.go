package rediscli

import (
	"reflect"
	"testing"

	"github.com/go-redis/redis"
)

func TestRedisInit(t *testing.T) {
	type args struct {
		addrs    string
		password string
	}
	tests := []struct {
		name string
		args args
		want *redis.Client
	}{
		// TODO: Add test cases.
		{"test", args{"192.168.0.80:6379", ""}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RedisInit(tt.args.addrs, tt.args.password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet(t *testing.T) {
	rediscli := RedisInit("192.168.0.80:6379", "")
	key := "tmp_key"
	value := "tmp_value"
	rediscli.Set(key, value, 0)
	t.Logf("Hset set. \n")
}

func TestGet(t *testing.T) {
	rediscli := RedisInit("192.168.0.80:6379", "")
	key := "tmp_key"
	val, _ := rediscli.Get(key).Result()
	t.Logf("val : %s \n", val)
}

func TestHSet(t *testing.T) {
	rediscli := RedisInit("192.168.0.80:6379", "")
	key := "pl-7"
	field := "status"
	rediscli.HSet(key, field, "tmpvalue")
	t.Logf("Hset set. \n")
}

func TestHGet(t *testing.T) {
	rediscli := RedisInit("192.168.0.80:6379", "")
	key := "pl-7"
	field := "status"
	val := rediscli.HGet(key, field)
	t.Logf("val : %s \n", val)
}

func TestDelete(t *testing.T) {
	rediscli := RedisInit("192.168.0.80:6379", "")
	key := "tmpkeydel"
	field := "status"

	rediscli.HSet(key, field, "123")
	val, _ := rediscli.HGet(key, field).Result()
	t.Logf("val : %s \n", val)
	rediscli.Del(key)
	val, _ = rediscli.HGet(key, field).Result()
	t.Logf("val : %s \n", val)
}

func TestHGetAll(t *testing.T) {
	rediscli := RedisInit("192.168.0.80:6379", "")
	key := "pl-7"
	val, _ := rediscli.HGetAll(key).Result()
	for k, v := range val {
		t.Logf("%s : %s \n", k, v)
	}
}
