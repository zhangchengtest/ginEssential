package test

import (
	"fmt"
	"ginEssential/pb"
	"testing"

	// 导入protobuf依赖包
	"github.com/golang/protobuf/proto"
)

func TestAdd(t *testing.T) {
	// 初始化消息
	score_info := &pb.BaseScoreInfoT{}
	score_info.WinCount = 10
	score_info.LoseCount = 1
	score_info.ExceptionCount = 2
	score_info.KillCount = 2
	score_info.DeathCount = 1
	score_info.AssistCount = 3
	score_info.Rating = 120

	// 以字符串的形式打印消息
	fmt.Println(score_info.String())

	// encode, 转换成二进制数据
	data, err := proto.Marshal(score_info)
	if err != nil {
		panic(err)
	}

	// decode, 将二进制数据转换成struct对象
	new_score_info := pb.BaseScoreInfoT{}
	err = proto.Unmarshal(data, &new_score_info)
	if err != nil {
		panic(err)
	}

	// 以字符串的形式打印消息
	fmt.Println(new_score_info.String())
}
