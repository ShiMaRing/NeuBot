package spider

import (
	"NeuBot/model"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	GetLearnWeek         = "https://portal.neu.edu.cn/tp_up/up/classScheduleP/getLearnweekbyDate"
	GetClassesByUserInfo = "https://portal.neu.edu.cn/tp_up/up/widgets/getClassbyUserInfo"
	GetClassesByTime     = "https://portal.neu.edu.cn/tp_up/up/widgets/getClassbyTime"
)

type TmpClass struct {
	JSXM  string `json:"JSXM"`
	KKXND string `json:"KKXND"`
	KKXQM string `json:"KKXQM"`
	ZZZ   int    `json:"ZZZ"`
	SKXQ  string `json:"SKXQ"`
	QSZ   int    `json:"QSZ"`
	CXJC  string `json:"CXJC"`
	SKZC  string `json:"SKZC"`
	KCMC  string `json:"KCMC"`
	JXDD  string `json:"JXDD"`
	SKJC  string `json:"SKJC"`
	KCH   string `json:"KCH"`
}

type LearnWeek struct {
	SchoolYear string `json:"schoolYear"`
	LearnWeek  string `json:"learnWeek"`
	Semester   string `json:"semester"`
}

// GetCourse 该方法需要使用传入的user信息获取绑定的table-model
func GetCourse(user *model.User) (model.TimeTable, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	//传入的信息可能会过期，需要提前判断
	cookie := &http.Cookie{}
	cookie.Name = "tp_up"
	success, token, err := AuthWithAccount(user.StdNumber, user.Password)
	if !success {
		return nil, model.PasswordIncorrectError
	}
	user.Token = token
	cookie.Value = token
	req := buildPostReq(GetLearnWeek, nil, cookie)
	data, err := fetch(req)
	if err != nil {
		return nil, err
	}
	learnWeek := &LearnWeek{}
	err = json.Unmarshal(data, learnWeek)
	if err != nil {
		return nil, err
	}
	weekData, err := json.Marshal(learnWeek)
	req = buildPostReq(GetClassesByUserInfo, weekData, cookie)
	classesData, err := fetch(req)
	if err != nil {
		return nil, err
	}
	tmpClasses := make([]TmpClass, 0)
	err = json.Unmarshal(classesData, &tmpClasses)
	if err != nil {
		return nil, err
	}
	//获取到classData后，构建最后的数据表格
	type PostData struct {
		LearnWeek
		ClassList []TmpClass `json:"classList"`
	}
	//携带该请求
	formData := PostData{
		LearnWeek: *learnWeek,
		ClassList: tmpClasses,
	}
	tmp, err := json.Marshal(formData)
	if err != nil {
		return nil, err
	}
	req = buildPostReq(GetClassesByTime, tmp, cookie)
	result, err := fetch(req)
	if err != nil {
		return nil, err
	}
	return process(result, user)
}

//处理json数据，获取课程表的数据结构
func process(data []byte, user *model.User) (model.TimeTable, error) {
	originCourses := make([]*model.ClassTrans, 0)
	err := json.Unmarshal(data, &originCourses)
	if err != nil {
		return nil, err
	}
	courses := make([]*model.Course, 0)
	//此时获取了所有的课程信息，需要进行处理
	for i := range originCourses {
		originCourse := originCourses[i]
		//获取原始课程信息
		course := convert2Course(originCourse)
		courses = append(courses, course)
	}
	return courses, nil
}

//信息转化
func convert2Course(tmp *model.ClassTrans) *model.Course {
	weekDay, _ := strconv.Atoi(tmp.SKXQ)
	className := removeDup(tmp.KCMC)
	place := removeDup(tmp.JXDD)
	teacher := removeDup(tmp.JSXM)
	start, _ := strconv.Atoi(tmp.SKJC)
	length, _ := strconv.Atoi(tmp.CXJC)
	return &model.Course{
		WeekDay:     weekDay,
		ClassName:   className,
		Start:       start,
		ClassLength: length,
		Place:       place,
		Teacher:     teacher,
	}
}

//去重
func removeDup(dup string) string {
	m := make(map[string]struct{})
	dup = strings.TrimSpace(dup)
	tmps := strings.Split(dup, "/")
	for i := range tmps {
		m[tmps[i]] = struct{}{}
	}
	res := make([]string, 0)

	for k := range m {
		res = append(res, k)
	}
	return strings.Join(res, "/")
}

func fetch(req *http.Request) ([]byte, error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func buildPostReq(url string, data []byte, cookie *http.Cookie) *http.Request {
	if data == nil {
		//说明携带空请求
		data, _ = json.Marshal(struct{}{})
	}

	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.AddCookie(cookie)
	return request
}
