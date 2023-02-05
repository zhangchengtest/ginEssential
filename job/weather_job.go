package job

import (
	"encoding/json"
	"flag"
	"fmt"
	"ginEssential/config"
	"ginEssential/model"
	"ginEssential/push"
	"ginEssential/spider"
	"ginEssential/util"
	"ginEssential/weather"
	"github.com/Lofanmi/chinese-calendar-golang/calendar"
	strftime "github.com/itchyny/timefmt-go"
	"github.com/pmylund/go-bloom"
	log "github.com/sirupsen/logrus"
	"github.com/zhangchengtest/simple/sqls"
	"math"
	"sort"
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
	//task.weatherInfo()
	//task.remind()
	//task.alarm()

	DB := sqls.DB()

	// 创建用户
	var clocks []model.Clock

	clockvos := make([]model.ClockVO, 0)

	var ss string

	DB.Find(&clocks)

	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	// 1、年月日
	year := time.Now().In(cstZone).Year()
	month := time.Now().In(cstZone).Month()
	//或者
	//month := time.Now().In(cstZone).Month().String()
	day := time.Now().In(cstZone).Day()

	for _, w := range clocks {
		if w.EventType == 1 {
			//按当天
			dd, _ := strconv.Atoi(w.NotifyDate)
			//ss = ss + "还差" + strconv.Itoa(dd-day) + "天就要"
			//ss = ss + w.EventDescription + "\n"

			c := calendar.BySolar(int64(year), int64(month), int64(day), 0, 0, 0)

			bytes, err := c.ToJSON()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(bytes))
			date := strconv.FormatInt(c.Solar.GetYear(), 10) + "-" + strconv.FormatInt(c.Solar.GetMonth(), 10) + "-" + strconv.FormatInt(c.Solar.GetDay(), 10)

			vo := model.ClockVO{
				Days:        0,
				EventTime:   dd,
				EventType:   w.EventType,
				Description: w.EventDescription,
				RealDate:    date,
			}
			clockvos = append(clockvos, vo)
		} else if w.EventType == 2 {
			//按照月
			dd, _ := strconv.Atoi(w.NotifyDate)
			if dd >= day {
				//ss = ss + "还差" + strconv.Itoa(dd-day) + "天就要"
				//ss = ss + w.EventDescription + "\n"

				c := calendar.BySolar(int64(year), int64(month), int64(day), 0, 0, 0)

				bytes, err := c.ToJSON()
				if err != nil {
					fmt.Println(err)
				}

				fmt.Println(string(bytes))
				date := strconv.FormatInt(c.Solar.GetYear(), 10) + "-" + strconv.FormatInt(c.Solar.GetMonth(), 10) + "-" + strconv.FormatInt(c.Solar.GetDay(), 10)

				vo := model.ClockVO{
					Days:        dd - day,
					EventType:   w.EventType,
					Description: w.EventDescription,
					RealDate:    date,
				}
				clockvos = append(clockvos, vo)

			}
		} else if w.EventType == 3 {
			//按照年
			// 3. ByLunar
			// 农历(最后一个参数表示是否闰月)
			arr := strings.Split(w.NotifyDate, "-")
			mm, _ := strconv.Atoi(arr[0])
			dd, _ := strconv.Atoi(arr[1])
			c := calendar.ByLunar(int64(year), int64(mm), int64(dd), 0, 0, 0, false)

			bytes, err := c.ToJSON()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(bytes))
			date := strconv.FormatInt(c.Solar.GetYear(), 10) + "-" + strconv.FormatInt(c.Solar.GetMonth(), 10) + "-" + strconv.FormatInt(c.Solar.GetDay(), 10)
			t2, _ := strftime.Parse(date+" 00:00:00", "%Y-%m-%d %H:%M:%S")

			d := t2.Sub(time.Now())

			//ss = ss + date
			//ss = ss + "还差" + strconv.FormatFloat(d.Hours()/24+1, 'f', 0, 64) + "天"
			//ss = ss + w.EventDescription

			vo := model.ClockVO{
				Days:        int(math.Floor(d.Hours()/24 + 1)),
				EventType:   w.EventType,
				Description: w.EventDescription,
				RealDate:    date,
			}
			clockvos = append(clockvos, vo)
		} else if w.EventType == 0 {
			//按照时钟
			// 3. ByLunar
			// 农历(最后一个参数表示是否闰月)
			arr := strings.Split(w.NotifyDate, "-")
			mm, _ := strconv.Atoi(arr[0])
			dd, _ := strconv.Atoi(arr[1])
			c := calendar.BySolar(int64(year), int64(mm), int64(dd), 0, 0, 0)

			bytes, err := c.ToJSON()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(bytes))
			date := strconv.FormatInt(c.Solar.GetYear(), 10) + "-" + strconv.FormatInt(c.Solar.GetMonth(), 10) + "-" + strconv.FormatInt(c.Solar.GetDay(), 10)
			t2, _ := strftime.Parse(date+" 00:00:00", "%Y-%m-%d %H:%M:%S")

			d := t2.Sub(time.Now())

			vo := model.ClockVO{
				Days:        int(math.Floor(d.Hours()/24 + 1)),
				Description: w.EventDescription,
				EventType:   w.EventType,
				RealDate:    date,
			}
			clockvos = append(clockvos, vo)
		}

	}

	sort.Sort(StudentArray(clockvos))

	for _, w := range clockvos {
		if w.EventType == 1 {
			ss = ss + "今天" + strconv.Itoa(w.EventTime) + "点"
			ss = ss + w.Description + "\n"
		} else {
			ss = ss + strconv.Itoa(w.Days) + "天后"
			ss = ss + w.Description
			ss = ss + ",就在" + w.RealDate + "\n"
		}

	}

	sendArticle(ss)
}

type StudentArray []model.ClockVO

func (array StudentArray) Len() int {
	return len(array)
}

func (array StudentArray) Less(i, j int) bool {
	if array[i].Days == array[j].Days {
		return array[i].EventTime < array[j].EventTime
	}
	return array[i].Days < array[j].Days //从小到大， 若为大于号，则从大到小
}

func (array StudentArray) Swap(i, j int) {
	array[i], array[j] = array[j], array[i]
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
