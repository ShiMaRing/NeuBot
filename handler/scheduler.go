package handler

import (
	"NeuBot/internal/service"
	"NeuBot/model"
	"NeuBot/pkg/spider"
	"fmt"
	"log"
	"strings"
	"time"
)

/*
1节 8:30-9:20
2节 9:30-10:20
3节 10:40-11:30
4节 11:40-12:30
5节14:00-14:50
6节 15:00-15:50
7节 16:10-17:00
8节 17:10-18:00
9节 18:30-19:20
10节 19:30-20:20
11节 20:30-21:20
12节 21:30-22:20
*/
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

// NewSchedulerHandler schedulerHandler构造函数
func NewSchedulerHandler() (*MessageHandler, error) {
	userService, err := service.NewUserService()
	if err != nil {
		return nil, err
	}
	return &MessageHandler{srv: userService}, nil
}

func (h *schedulerHandler) submission() {
	//每次寻找对应时间的消息，需要获取当前时间
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
	}
	for i := range users {
		go func(i int) {
			user := users[i]
			//校验用户权限
			if user.Perm&model.CoursePerm == 0 {
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
			}
		}(i)
	}
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