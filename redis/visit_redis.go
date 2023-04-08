package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func Visit(player string) error {
	conn := pool.Get()
	defer conn.Close()

	//str, err := json.Marshal(player)
	//if err != nil {
	//	return fmt.Errorf("error marshal player: %v", err)
	//}
	fmt.Println("enter here  %v", player)
	_, err := conn.Do("INCR", player)
	if err != nil {
		return fmt.Errorf("error : %v", err)
	}

	return nil
}

func Set(key string, value string, time int) error {
	conn := pool.Get()
	defer conn.Close()

	//str, err := json.Marshal(player)
	//if err != nil {
	//	return fmt.Errorf("error marshal player: %v", err)
	//}
	_, err := conn.Do("SET", key, value, "EX", time)
	if err != nil {
		return fmt.Errorf("error : %v", err)
	}

	return nil
}

func Get(key string) string {
	conn := pool.Get()
	defer conn.Close()

	//str, err := json.Marshal(player)
	//if err != nil {
	//	return fmt.Errorf("error marshal player: %v", err)
	//}
	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		fmt.Println("error : %s", err)
	}

	return value
}
