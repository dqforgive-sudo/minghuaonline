package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/redmask-hb/GoSimplePrint/goPrint"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type replyAddRes struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
	Status bool        `json:"status"`
}
type courseRes struct {
	Code   int           `json:"code"`
	Msg    string        `json:"msg"`
	Result courseResList `json:"result"`
	Status bool          `json:"status"`
}

type courseResList struct {
	Announcement interface{} `json:"announcement"` //公告的信息
	Banner       interface{} `json:"banner"`       //banner 横幅广告的信息
	Enlist       interface{} `json:"enlist"`       //全校各学院开放的课程
	Finish       interface{} `json:"finish"`       //已经结束的课程
	List         []courseObj `json:"list"`         //当前正在进行的课程
}

type courseObj struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	ClassTeacher string  `json:"classTeacher"`
	Progress     float32 `json:"progress"`
}

type chapterRes struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result chapterList `json:"result"`
	Status bool        `json:"status"`
}

type chapterList struct {
	List []chapterListObj `json:"list"`
}

type chapterListObj struct {
	Id       int              `json:"id"`
	Name     string           `json:"name"`
	NodeList []chapterNodeObj `json:"nodeList"`
}

type chapterNodeObj struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	VideoDuration string `json:"videoDuration"`
	VideoState    int    `json:"videoState"`
	TabVideo      bool   `json:"tabVideo"`
}

type studyRes struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result studyResObj `json:"result"`
	Status bool        `json:"status"`
}

type studyResObj struct {
	Data studyObj `json:data`
}

type studyObj struct {
	StudyId int `json:"studyId"`
}
type Online struct {
	Code   int         `json:"_code"`
	Msg    string      `json:"msg"`
	Result studyResObj `json:"result"`
	Status bool        `json:"status"`
}

var CourseResObj = courseRes{}
var reply bool
var c = studyObj{}

//学校列表结构体
type schoolList struct {
	Code   int           `json:"code"`
	Msg    string        `json:"msg"`
	Result schoolResList `json:"result"`
	Status bool          `json:"status"`
}

type schoolResList struct {
	List []schoolObj `json:"list"`
}

type schoolObj struct {
	Badge string `json:"badge"`
	Host  string `json:"host2"`
	Id    string `json:"id"`
	Ident string `json:"ident"`
	Name  string `json:"name"`
}

//登陆响应结构体
type LoginRes struct {
	Code   int       `json:"code"`
	Msg    string    `json:"msg"`
	Result LoginInfo `json:"result"`
	Status bool      `json:"status"`
}

type LoginInfo struct {
	Data LoginObj `json:"data"`
}

type LoginObj struct {
	Id          int    `json:"id"`
	Token       string `json:"token"`
	Name        string `json:"name"`
	ClassId     int    `json:"classId"`
	CollegeId   int    `json:"collegeId"`
	Point       int    `json:"point"`
	Rank        int    `json:"rank"`
	ClassName   string `json:"className"`
	CollegeName string `json:"collegeName"`
}

type SignObj struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
	Status bool        `json:"status"`
}

//获取课程列表
func GetCourseList() {
	url := Host + "/api/course/list.json"
	var bufReader bytes.Buffer
	mpWriter := multipart.NewWriter(&bufReader)
	mpWriter.WriteField("token", Token)
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, &bufReader)
	request.Header.Add("Content-Type", "multipart/form-data; boundary="+mpWriter.Boundary())
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &CourseResObj)
	if err != nil {
		log.Fatal(err)
	}
	if CourseResObj.Status {
		//log.Printf("当前正在学习课程获取成功\n\n")
		fmt.Printf("\n%s %s %s\n", "课程序号", "完成度", "课程名称")
		var courseIndex = 0
		for _, item := range CourseResObj.Result.List {
			fmt.Printf("%03d %.2f%% %s\n", item.Id, item.Progress*100, item.Name)
			courseIndex++
		}
		var selectIndex int
		fmt.Printf("\n请输入课程的序号: ")
		fmt.Scanf("%d", &selectIndex)
		getCourseChapter(selectIndex)
	} else {
		log.Fatalf("课程获取失败 %s\n", CourseResObj.Msg)
	}
}

//获取课程的详情，但如果提交课程，并不需要此请求
func getCourseDetail(courseId int) {
	url := Host + "/api/course/detail.json"
	var bufReader bytes.Buffer
	mpWriter := multipart.NewWriter(&bufReader)
	mpWriter.WriteField("courseId", strconv.Itoa(courseId))
	mpWriter.WriteField("token", Token)
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, &bufReader)
	request.Header.Add("Content-Type", "multipart/form-data; boundary="+mpWriter.Boundary())
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
}

//获取课程章节
func getCourseChapter(courseId int) {
	url := Host + "/api/course/chapter.json"
	var bufReader bytes.Buffer
	mpWriter := multipart.NewWriter(&bufReader)
	mpWriter.WriteField("courseId", strconv.Itoa(courseId))
	mpWriter.WriteField("token", Token)
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, &bufReader)
	request.Header.Add("Content-Type", "multipart/form-data; boundary="+mpWriter.Boundary())
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	m := chapterRes{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
	}
	if m.Status {
		var replyYes string
		fmt.Printf("是否开启学习评论 格式为 学习 [课程名称] 记录打卡 ? [y/N] ")
		fmt.Scan(&replyYes)
		switch replyYes {
		case "y":
			fmt.Printf("已经开启学习后自动评论 !\n")
			reply = true
		case "n":
			fmt.Printf("关闭学习后自动评论 !\n")
			reply = false
		default:
			fmt.Printf("不开启学习后自动评论 !\n")
			reply = false
		}
		for _, course := range m.Result.List {
			chapterName := course.Name
			for _, chapter := range course.NodeList {
				if chapter.TabVideo {
					if chapter.VideoState == 2 {
						log.Printf("%s[%s] 已经完成学习，自动略过", chapterName, chapter.Name)
					} else {
						doChapter(chapter)
					}
				} else {
					log.Printf("%s[%s] 不是视频教学 !", chapterName, chapter.Name)
				}
			}
		}
	} else {
		log.Printf("课程章节获取出错 %s 请确保输入了正确的课程序号", m.Msg)
	}
	GetCourseList()
}

//判断是否需要学习
func doChapter(chapter chapterNodeObj) {
	firststudyChapter(chapter)
	go lg()
	cc := chapter.VideoDuration
	if s, err := strconv.Atoi(cc); err == nil {
		fmt.Printf("视频时长：%v\n", s)
		bar := goPrint.NewBar(s)
		bar.SetNotice("进度：")
		bar.SetGraph(">")
		for i := 1; i <= s; i++ {
			bar.PrintBar(i)
			time.Sleep(time.Second)
		}
		bar.PrintEnd("Finish!")
		//time.Sleep(time.Duration(s) * time.Second)
		endstudyChapter(chapter)
	}

}

//完成章节学习
func firststudyChapter(chapter chapterNodeObj) {
	url := Host + "/api/node/study.json"
	var bufReader bytes.Buffer
	mpWriter := multipart.NewWriter(&bufReader)
	mpWriter.WriteField("platform", "Android")
	mpWriter.WriteField("version", "1.3.2")
	mpWriter.WriteField("nodeId", strconv.Itoa(chapter.Id))
	mpWriter.WriteField("terminal", "Android")
	mpWriter.WriteField("studyTime", chapter.VideoDuration)
	mpWriter.WriteField("token", Token)
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, &bufReader)
	request.Header.Add("Content-Type", "multipart/form-data; boundary="+mpWriter.Boundary())
	request.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 9; Mi Note 3 Build/PKQ1.181007.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/78.0.3904.96 Mobile Safari/537.36;sino-web.com")
	request.Header.Add("Connection", "Keep-Alive")
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	m := studyRes{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
	}
	if m.Status {
		if reply {
			replyMsg := AddReply(chapter)
			log.Printf("%s %s %s\n", chapter.Name, m.Msg, replyMsg)
		} else {
			log.Printf("%s %s studyid: %d\n", chapter.Name, m.Msg, m.Result.Data.StudyId)
			c.StudyId = m.Result.Data.StudyId
		}
	} else {
		log.Printf("%s 学习失败 ! %s\n", chapter.Name, m.Msg)
	}
}
func endstudyChapter(chapter chapterNodeObj) {
	url := Host + "/api/node/study.json"
	var bufReader bytes.Buffer
	mpWriter := multipart.NewWriter(&bufReader)
	mpWriter.WriteField("platform", "Android")
	mpWriter.WriteField("version", "1.3.2")
	mpWriter.WriteField("nodeId", strconv.Itoa(chapter.Id))
	mpWriter.WriteField("terminal", "Android")
	mpWriter.WriteField("studyTime", chapter.VideoDuration)
	mpWriter.WriteField("studyId", strconv.Itoa(c.StudyId))
	mpWriter.WriteField("token", Token)
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, &bufReader)
	request.Header.Add("Content-Type", "multipart/form-data; boundary="+mpWriter.Boundary())
	request.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 9; Mi Note 3 Build/PKQ1.181007.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/78.0.3904.96 Mobile Safari/537.36;sino-web.com")
	request.Header.Add("Connection", "Keep-Alive")
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	m := studyRes{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
	}
	if m.Status {
		if reply {
			replyMsg := AddReply(chapter)
			log.Printf("%s %s %s\n", chapter.Name, m.Msg, replyMsg)
		} else {
			log.Printf("%s %s studyid: %d\n", chapter.Name, m.Msg, m.Result.Data.StudyId)
			c.StudyId = m.Result.Data.StudyId
		}
	} else {
		log.Printf("%s 学习失败 ! %s\n", chapter.Name, m.Msg)
	}
}
func online() {
	url := Host + "/api/online.json"
	var bufReader bytes.Buffer
	mpWriter := multipart.NewWriter(&bufReader)
	mpWriter.WriteField("platform", "Android")
	mpWriter.WriteField("version", "1.3.2")
	mpWriter.WriteField("schoolId", SchoolId)
	mpWriter.WriteField("token", Token)
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, &bufReader)
	request.Header.Add("Content-Type", "multipart/form-data; boundary="+mpWriter.Boundary())
	request.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 9; Mi Note 3 Build/PKQ1.181007.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/78.0.3904.96 Mobile Safari/537.36;sino-web.com")
	request.Header.Add("Connection", "Keep-Alive")
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	m := Online{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
	}
	if m.Status {
		log.Printf("%s\n", m.Msg)
	}
}
func lg() {
	for range time.Tick(120 * time.Second) {
		online()
	}
}

var Host string     //域名
var Token string    //token
var SchoolId string //学校id

// 获取学校
func GetSchool() {
	url := "http://mooc.yinghuaonline.com/api/login/school.json"
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	m := schoolList{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
	}
	if m.Status {
		log.Printf("获取学校信息列表成功\n\n")
		fmt.Printf("%s %s\n", "序号", "学校名称")
		var schoolIndex = 0
		for _, item := range m.Result.List {
			fmt.Printf("%03d %s\n", schoolIndex, item.Name)
			schoolIndex++
		}
		var selectIndex int
		fmt.Printf("\n请输入学校的序号: ")
		//reader := bufio.NewReader(os.Stdin)
		fmt.Scanf("%d\n", &selectIndex)
		for {
			if selectIndex < len(m.Result.List) {
				break
			}
			fmt.Printf("请输入学校的序号: ")
			fmt.Scanf("%d\n", &selectIndex)
		}
		fmt.Printf("选择了 %s\n", m.Result.List[(selectIndex)].Name)
		Host = m.Result.List[selectIndex].Host
		fmt.Println(Host)
		SchoolId = m.Result.List[selectIndex].Id
	} else {
		log.Fatalf("学校列表获取失败 %s\n", m.Msg)
	}
}

// 登录
func Login() {
	var username string
	var password string
	//reader := bufio.NewReader(os.Stdin)
	fmt.Printf("请输入学号: ")
	fmt.Scanf("%s\n", &username)
	//username, _ = reader.ReadString('\n')
	fmt.Printf("请输入密码: ")
	fmt.Scanf("%s\n", &password)
	//password, _ = reader.ReadString('\n')
	url := Host + "/api/login.json"
	var bufReader bytes.Buffer
	mpWriter := multipart.NewWriter(&bufReader)
	mpWriter.WriteField("username", username)
	mpWriter.WriteField("password", password)
	mpWriter.WriteField("schoolId", SchoolId)
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, &bufReader)
	request.Header.Add("Content-Type", "multipart/form-data; boundary="+mpWriter.Boundary())
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	if err != nil {
		log.Fatal(err)
	}
	m := LoginRes{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
	}
	if m.Status {
		Token = m.Result.Data.Token
		log.Printf("%s 登录成功 ! 来自 %s-%s 的 %s %s\n", username, m.Result.Data.CollegeName, m.Result.Data.ClassName, m.Result.Data.Name, signIn())
	} else {
		log.Fatalf("%s 登录失败! %s", username, m.Msg)
	}
	GetCourseList()
}

// 签到
func signIn() string {
	url := Host + "/api/user/sign_in.json"
	var bufReader bytes.Buffer
	mpWriter := multipart.NewWriter(&bufReader)
	mpWriter.WriteField("token", Token)
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, &bufReader)
	request.Header.Add("Content-Type", "multipart/form-data; boundary="+mpWriter.Boundary())
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	m := SignObj{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
	}
	return m.Msg
}
func GetAnswer() {
	url := Host + "/api/answer.json"
	var bufReader bytes.Buffer
	mpWriter := multipart.NewWriter(&bufReader)
	mpWriter.WriteField("schoolId", "6")
	mpWriter.WriteField("questionId", "1")
	mpWriter.WriteField("token", Token)
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, &bufReader)
	request.Header.Add("Content-Type", "multipart/form-data; boundary="+mpWriter.Boundary())
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
}

//添加课程的评论
func AddReply(chapter chapterNodeObj) string {
	url := Host + "/api/node/add_reply.json"
	var bufReader bytes.Buffer
	an := regexp.MustCompile("[1-9]\\d*")
	chapterName := an.ReplaceAllString(chapter.Name, "")
	content := "学习 " + chapterName + " 记录打卡 !"
	mpWriter := multipart.NewWriter(&bufReader)
	mpWriter.WriteField("nodeId", strconv.Itoa(chapter.Id))
	mpWriter.WriteField("content", content)
	mpWriter.WriteField("token", Token)
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, &bufReader)
	request.Header.Add("Content-Type", "multipart/form-data; boundary="+mpWriter.Boundary())
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	m := replyAddRes{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
	}
	return m.Msg
}
