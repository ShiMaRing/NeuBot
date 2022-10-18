package handler

import (
	"NeuBot/internal/service"
	"NeuBot/model"
	"NeuBot/pkg/spider"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"strings"
	"sync"
	"time"
)

var classTimeMap map[int]string
var hourToSerial map[int]int //小时转序号

func init() {
	classTimeMap = map[int]string{
		1:  "8:30-9:20",
		2:  "9:30-10:20",
		3:  "10:40-11:30",
		4:  "11:40-12:30",
		5:  "14:00-14:50",
		6:  "15:00-15:50",
		7:  "16:10-17:00",
		8:  "17:10-18:00",
		9:  "18:30-19:20",
		10: "19:30-20:20",
		11: "20:30-21:20",
		12: "21:30-22:20",
	}
	hourToSerial = map[int]int{
		8:  1,
		9:  2,
		10: 3,
		11: 4,
		13: 5,
		14: 6,
		16: 7,
		17: 8,
		18: 9,
		19: 10,
		20: 11,
		21: 12,
	}
}

type schedulerHandler struct {
	srv *service.UserService
}

// newSchedulerHandler schedulerHandler构造函数
func newSchedulerHandler() (*schedulerHandler, error) {
	userService, err := service.NewUserService()
	if err != nil {
		return nil, err
	}
	return &schedulerHandler{srv: userService}, nil
}

func (h *schedulerHandler) submission() {
	//每次寻找对应时间的消息，需要获取当前时间
	log.Println("开始执行: ", time.Now())
	now := time.Now()
	weekday := now.Weekday() //注意星期天为0，而数据中为7
	if weekday == time.Sunday {
		weekday = 7
	}
	//需要正确找出对应的课程
	hour := now.Hour() //[0_23]
	serial := hourToSerial[hour]
	users, err := h.srv.GetAllUser()
	//此时用户应当已经与课程完成关联
	if err != nil {
		log.Println(err) //查找失败
		return
	}
	group := sync.WaitGroup{}
	for i := range users {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			user := users[i]
			//校验用户权限
			if user.Perm&model.CoursePerm == 0 || user.State == model.LOGOUT {
				//说明没有权限
				return
			}
			user.Mu.Lock()
			courses := user.TimeTable
			user.Mu.Unlock()
			for idx := range courses {
				course := courses[idx]
				//已经发送过或者不是本星期，直接结束
				if course.IsSubmission == true || course.WeekDay != int(weekday) {
					continue
				}
				if course.Start != serial {
					continue
				}
				//允许发送
				msg := buildMsg(course)
				course.IsSubmission = true //表示已经发送完成
				replyMsg(user.QQ, msg, false)
				fmt.Println("发送消息", user.QQ, msg)
			}
		}(i)
	}
	group.Wait()
}

//刷新课程信息
func (h *schedulerHandler) refreshCourse() error {
	users, err := h.srv.GetAllUser()
	if err != nil {
		return err
	}
	for i := range users {
		go func(i int) {
			user := users[i]
			courses, err := spider.GetCourse(user)
			if err != nil { //继续执行下一个用户的课程获取任务
				log.Println(err)
				return
			}
			user.Mu.Lock()
			user.TimeTable = courses
			user.Mu.Unlock()
			h.srv.UpdateUser(user) //更新用户信息，需要进行持久化操作,对于缓冲来说没什么必要
		}(i)
	}
	return nil
}

//构造发送的消息
func buildMsg(course *model.Course) string {
	builder := strings.Builder{}
	builder.WriteString("----上课提醒----\n")
	builder.WriteString(fmt.Sprintf("课程名： %s\n", course.ClassName))
	builder.WriteString(fmt.Sprintf("教师： %s\n", course.Teacher))
	builder.WriteString(fmt.Sprintf("地点： %s\n", course.Place))
	startTime := classTimeMap[course.Start]
	endTime := classTimeMap[course.Start+course.ClassLength-1]
	timeMessage := fmt.Sprintf("%s-%s", strings.Split(startTime, "-")[0], strings.Split(endTime, "-")[1])
	builder.WriteString(fmt.Sprintf("时间： %s", timeMessage))
	return builder.String()
}

// StartSchedule 开启调度任务
func StartSchedule() error {
	var err error
	handler, err := newSchedulerHandler()
	handler.srv.GetAllUser()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	c := cron.New(cron.WithSeconds())
	_, err = c.AddFunc("* 20 8,9,18,19,20,21 * * ? ", func() { handler.submission() })
	if err != nil {
		return err
	}
	_, err = c.AddFunc("* 30 10,11 * * ?", func() { handler.submission() })
	if err != nil {
		return err
	}
	_, err = c.AddFunc("* 50 13,14 * * ?", func() { handler.submission() })
	if err != nil {
		return err
	}
	_, err = c.AddFunc("10 0 16,17 * * ?", func() { handler.submission() })
	if err != nil {
		return err
	}
	_, err = c.AddFunc("@midnight", func() { handler.srv.CleanUp() })
	if err != nil {
		return err
	}
	_, err = c.AddFunc("* * 1 ? * 0 ", func() {
		handler.srv.CleanAllCourse() //先清除
		handler.refreshCourse()
	})
	if err != nil {
		return err
	}
	c.Start()
	return nil
}
