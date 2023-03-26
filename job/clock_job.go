package job

import (
	"fmt"
	"ginEssential/model"
	"ginEssential/util"
	"github.com/Lofanmi/chinese-calendar-golang/calendar"
	strftime "github.com/itchyny/timefmt-go"
	"github.com/zhangchengtest/simple/sqls"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

func DoClock() {

	DB := sqls.DB()

	// 创建用户
	var clocks []model.Clock

	clockvos := make([]model.ClockVO, 0)

	var ss string

	DB.Find(&clocks)

	// 1、年月日
	year := time.Now().Year()
	month := time.Now().Month()
	//或者
	//month := time.Now().In(cstZone).Month().String()
	day := time.Now().Day()

	//查询指定年份指定月份有多少天
	monthDays := util.GetYearMonthToDay(year, int(month))

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
		} else if w.EventType == 4 {

			//按照周
			dd, _ := strconv.Atoi(w.NotifyDate)
			var res int
			var date string

			weekday := time.Now().Weekday()
			weekdayInt := int(weekday)
			if weekday == time.Sunday {
				weekdayInt = 7
			}

			if dd >= weekdayInt {
				//ss = ss + "还差" + strconv.Itoa(dd-day) + "天就要"
				//ss = ss + w.EventDescription + "\n"
				res = dd - weekdayInt
				c := calendar.BySolar(int64(year), int64(month), int64(day+res), 0, 0, 0)

				date = strconv.FormatInt(c.Solar.GetYear(), 10) + "-" + strconv.FormatInt(c.Solar.GetMonth(), 10) + "-" + strconv.FormatInt(c.Solar.GetDay(), 10)
			} else {
				res = 7 - weekdayInt + dd
				c := calendar.BySolar(int64(year), int64(month), int64(day+res), 0, 0, 0)

				date = strconv.FormatInt(c.Solar.GetYear(), 10) + "-" + strconv.FormatInt(c.Solar.GetMonth(), 10) + "-" + strconv.FormatInt(c.Solar.GetDay(), 10)
			}

			vo := model.ClockVO{
				Days:        res,
				EventType:   w.EventType,
				Description: w.EventDescription,
				RealDate:    date,
			}
			clockvos = append(clockvos, vo)

		} else if w.EventType == 2 {

			//按照月
			dd, _ := strconv.Atoi(w.NotifyDate)
			var res int
			var date string

			c := calendar.BySolar(int64(year), int64(month), int64(dd), 0, 0, 0)

			bytes, err := c.ToJSON()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(bytes))

			if dd >= day {
				//ss = ss + "还差" + strconv.Itoa(dd-day) + "天就要"
				//ss = ss + w.EventDescription + "\n"
				res = dd - day
				date = strconv.FormatInt(c.Solar.GetYear(), 10) + "-" + strconv.FormatInt(c.Solar.GetMonth(), 10) + "-" + strconv.FormatInt(c.Solar.GetDay(), 10)
			} else {
				res = monthDays - day + dd
				fmt.Println(monthDays)
				date = strconv.FormatInt(c.Solar.GetYear(), 10) + "-" + strconv.FormatInt(c.Solar.GetMonth()+1, 10) + "-" + strconv.FormatInt(c.Solar.GetDay(), 10)
			}

			vo := model.ClockVO{
				Days:        res,
				EventType:   w.EventType,
				Description: w.EventDescription,
				RealDate:    date,
			}
			clockvos = append(clockvos, vo)

		} else if w.EventType == 3 {
			if w.LunarFlag == 1 {
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
					Days:        int(math.Floor(d.Hours() / 24)),
					EventType:   w.EventType,
					Description: w.EventDescription,
					RealDate:    date,
				}
				clockvos = append(clockvos, vo)
			} else {
				//按照年
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

				//ss = ss + date
				//ss = ss + "还差" + strconv.FormatFloat(d.Hours()/24+1, 'f', 0, 64) + "天"
				//ss = ss + w.EventDescription

				vo := model.ClockVO{
					Days:        int(math.Floor(d.Hours() / 24)),
					EventType:   w.EventType,
					Description: w.EventDescription,
					RealDate:    date,
				}
				clockvos = append(clockvos, vo)
			}

		}

	}

	clockvos = *findCalendar(clockvos)

	sort.Sort(StudentArray(clockvos))

	for _, w := range clockvos {
		if w.Days < 0 {
			continue
		}
		if w.EventType == 1 {
			ss = ss + "今天" + strconv.Itoa(w.EventTime) + "点"
			ss = ss + w.Description + "\n"
		} else {
			ss = ss + strconv.Itoa(w.Days) + "天后"
			ss = ss + w.Description
			ss = ss + ",就在" + w.RealDate + "\n"
		}

	}

	// 创建用户

	var theme model.Topic

	rand.Seed(time.Now().UnixNano())
	aid := rand.Intn(70) + 1
	DB.Where("id = ?", aid).Find(&theme)
	ss = ss + "\n算法主题：" + theme.Name + "\n"

	ss = sport(ss)
	sendArticle(ss)
}

func findCalendar(clockvos []model.ClockVO) *[]model.ClockVO {
	DB := sqls.DB()
	// 创建用户
	var clocks []model.Calendar

	DB.Find(&clocks)

	for _, w := range clocks {
		//按照时钟
		// 3. ByLunar
		// 农历(最后一个参数表示是否闰月)
		yy := w.NotifyDate.Year()
		mm := w.NotifyDate.Month()
		dd := w.NotifyDate.Day()
		c := calendar.BySolar(int64(yy), int64(mm), int64(dd), 0, 0, 0)

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

	return &clockvos

}

func sport(ss string) string {
	// 创建用户
	DB := sqls.DB()

	var theme model.Sport

	rand.Seed(time.Now().UnixNano())
	aid := rand.Intn(7) + 1
	DB.Where("id = ?", aid).Find(&theme)
	res := ss + "\n运动主题：" + theme.Name + "\n"
	return res
}
