package job

import (
	"encoding/json"
	"flag"
	"fmt"
	"ginEssential/config"
	"ginEssential/push"
	"ginEssential/spider"
	"ginEssential/util"
	"ginEssential/weather"
	"github.com/pmylund/go-bloom"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Log  *config.LogInfo   `mapstructure:"log"`
	Push *[]push.Push      `mapstructure:"push"`
	Info *[]weather.Inform `mapstructure:"noti"`
}

var (
	configName = flag.String("c", "weather.yaml", "config name")
	query      = flag.String("q", "", "query city code")
)

func WeatherJob() {
	queryData()
	var task Task
	config := config.NewConfigByName(*configName)
	readTask(config, &task)
	task.weatherInfo()
	task.remind()
	task.alarm()

}

var f = bloom.New(10000, 0.001)

func (task Task) alarm() {
	hour := time.Now().Hour()
	if hour < 6 || hour > 22 {
		return
	}
	weather.WarningInforms(*task.Info, *task.Push, f)
}

func (task Task) weatherInfo() {
	for _, w := range *task.Info {
		if w.Report {
			ws := weather.GetWeather(w)
			for _, v := range *task.Push {
				v.Push(weather.GetToMsg(ws, w))
			}
		}
	}
}

func (task Task) remind() {
	for _, w := range *task.Info {
		if w.Remind {
			ws := weather.GetWeather(w)
			info := weather.GetRemindInfo(ws)
			if info != nil {
				for e := range *info.Msg {
					for _, v := range *task.Push {
						v.Push((*info.Msg)[e])
					}
				}
			} else {
				log.Info(w.Info, "明天是晴天！")
			}
		} else {
			log.Info(w.Info, "不做提醒！")
		}

	}
}

func queryData() {

	if query != nil && *query != "" {
		m1 := getMap("http://www.weather.com.cn/data/city3jdata/china.html?_=" + strconv.FormatInt(time.Now().Unix(), 10) + "667")
		if m1 == nil {
			return
		}
		c1, v1, q1 := stringCompare(*query, "", *m1)
		if c1 == nil {
			return
		}
		m2 := getMap("http://www.weather.com.cn/data/city3jdata/provshi/" + *c1 + ".html?_=" + strconv.FormatInt(time.Now().Unix(), 10) + "667")
		c2, v2, q2 := stringCompare(*q1, *v1, *m2)
		if c2 == nil {
			return
		}
		m3 := getMap("http://www.weather.com.cn/data/city3jdata/station/" + *c1 + *c2 + ".html?_=" + strconv.FormatInt(time.Now().Unix(), 10) + "667")
		c3, _, _ := stringCompare(*q2, *v2, *m3)
		fmt.Printf("省：%s 市：%s 县区：%s", *c1, *c2, *c3)
	}
}

func stringCompare(all, this string, codeMap map[string]string) (code, value, query *string) {
	if codeMap == nil {
		return nil, nil, nil
	}
	if this != "" {
		p := all[strings.Index(all, this):]
		code, value = codeCompare(codeMap, p)
		if code != nil {
			return code, value, &p
		}
	}
	code, value = codeCompare(codeMap, all)
	return code, value, &all
}

func codeCompare(codeMap map[string]string, query string) (code, value *string) {
	for k, v := range codeMap {
		if strings.Contains(query, v) {
			return &k, &v
		}
	}
	return nil, nil
}

func getMap(url string) *map[string]string {
	codeMap := make(map[string]string)
	body := spider.GetResponseBody(url)
	if err := json.Unmarshal([]byte(body), &codeMap); err != nil {
		fmt.Println(url, "解析数据出错:", body)
		return nil
	}
	return &codeMap
}

func readTask(config *config.WeatherConfig, task *Task) {
	config.GetViperUnmarshal(task)
	for e := range *task.Info {
		(*task.Info)[e].District = util.Add(2, (*task.Info)[e].District)
		(*task.Info)[e].City = util.Add(2, (*task.Info)[e].City)
	}
}
