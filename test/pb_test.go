package test

import (
	"fmt"
	"ginEssential/pb"
	"ginEssential/util"
	"strings"
	"testing"
	"unicode/utf8"

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

func TestData(t *testing.T) {
	str := `如何说服你的父母不要管自己的婚姻
行程这个东西是要到事情特别的时候才能用到的
算法真的是太难了
这他妈看半天没啥效果
就真的很生气对吧 哎 需要人理解真她妈的太难
特别是家里人
买桌子 要买个合适的桌子
得去宜家看看
客厅的桌子 估计得去宜家看看
哎 连沙发都还没买呢
订阅相关的话题 如果有更新就推送过来 这多有意思啊`
	fmt.Println("start")
	fmt.Println(util.IntToString(len(str)))
	data := strings.ReplaceAll(str, "\n", "")
	data = strings.ReplaceAll(data, " ", "")
	fmt.Println(util.IntToString(utf8.RuneCountInString(str)))
	fmt.Println("stop")
}

func TestData2(t *testing.T) {
	dirPath := "D:\\hanchuancaolu"
	files, err := util.GetAllFiles2(dirPath)
	if err != nil {
		panic(err)
	}

	// 输出所有文件路径和文件名
	//for _, file := range files {
	//	fmt.Println(file)
	//}
	fmt.Println(len(files))

	resutl := util.GetRandomString(files)
	str := strings.ReplaceAll(resutl, dirPath, "http://peer.punengshuo.com")
	str = strings.ReplaceAll(str, "\\", "/")
	fmt.Println(str)
}
