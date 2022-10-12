package model

import (
	"gorm.io/gorm"
	"time"
)

// ClassTrans 课程的爬虫json形式
type ClassTrans struct {
	JSXM     string `json:"JSXM"`     //上课老师 注意去重
	KKXND    string `json:"KKXND"`    //上课年份 2022-2023
	KKXQM    string `json:"KKXQM"`    //上课学期  1表示上学期
	ZZZ      int    `json:"ZZZ"`      //结束周    表示第几周结课
	SKXQ     string `json:"SKXQ"`     //星期几上课  1-7 表示星期一到星期天
	QSZ      int    `json:"QSZ"`      //搞不懂什么意思
	CXJC     string `json:"CXJC"`     //上几节课  “4”    表示上四节课
	SKZC     string `json:"SKZC"`     //一个1010字符串 第几位表示第几周有课
	KCMC     string `json:"KCMC"`     //课程名称 需要注意去重
	JXDD     string `json:"JXDD"`     //教学地点 注意去重
	SKJC     string `json:"SKJC"`     //上课开始节数，第几节开始上课
	KCH      string `json:"KCH"`      //课程编号  A031231
	SKZ      string `json:"SKZ"`      //上课周   第一周-第十周
	ColorNum string `json:"colorNum"` //颜色数值
}

// TimeTable 定义一周课表
// 每个用户都会与一个时间表关联，每个时间表会保存至数据库
//一个时间表将与多个课程进行关联
type TimeTable struct {
	gorm.Model
	StdNumber string
	Courses   []*Course
}

// Course 课程，对爬虫获取的课程的解析结果，
// 每次到时间点之前都要从缓存中遍历数据，挑选出合适的课程进行报送
type Course struct {
	gorm.Model
	WeekDay   int       //星期几上课 1~7表示周一到周日
	ClassName string    //课程名
	StartTime time.Time //课程开始时间
	EndTime   time.Time //课程结束时间
	place     string    //上课地点
	teacher   string    //任课老师
}
