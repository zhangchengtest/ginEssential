package redis

import (
	"fmt"
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
