package main

import (
	"bytes"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/robertkrimen/otto"
	"github.com/robfig/cron"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	tb "gopkg.in/tucnak/telebot.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"pornbot/entity"
	"pornbot/util"
	_ "pornbot/util"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	err         error
	//定时任务的cron表达式
	spec        = "0 0 5 * * ?"
	token       = "telegram Token"
	webhookUrl  = "webhook url"
	webhookPort = ":443"
	serverCrt   = "server.crt"
	serverKey   = "server.key"
	//bot并发数
	MaxConnections = 40
	b           *tb.Bot
)
const (
	Mp4BoxPath = "/root/gpac_public/bin/gcc/MP4Box"
	//视频描述模板
	captionTemplate = `标题: %s
收藏: %s
作者: %s `
	//定时任务发送的群组Id
	telegramId = -222222
)


//另一种机器人
func main() {
	c := cron.New() 
	//定时任务
	c.AddFunc(spec, func() {
		cronTaskSendVideo()
	})
	log.Println("定时任务开启")
	c.Start()
	log.Println("bot开启")
	b, err = tb.NewBot(tb.Settings{
		Token: token,
		Poller: &tb.Webhook{
			Listen:  webhookPort,
			MaxConnections: MaxConnections,
			TLS: &tb.WebhookTLS{
				Key:  serverKey,
				Cert: serverCrt,
			},
			Endpoint: &tb.WebhookEndpoint{PublicURL: webhookUrl},
		},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, "向我发送91视频链接，获取视频")
	})
	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello World!")
	})
	b.Handle(tb.OnText, func(m *tb.Message) {
		text := m.Text
		log.Println("收到消息：" + text)
		if strings.Contains(text,"viewkey") {
			controlSecPage(text,m.Chat.ID)
		}else {
			b.Send(m.Sender, "请发送视频链接给我")
		}

	})
	b.Start()

	select {}  //阻塞主线程停止

}

//controlSecPage 详情页爬取 发送视频
func controlSecPage(url string,chatId int64) {
	videoinfo, _ := BotUti.GetHttpHtmlContent(url, "#useraction > div:nth-child(1) > span:nth-child(2)", "body")
	//fmt.Println(content)
	reg, _ := regexp.Compile(`strencode2\((.*?)\)\)`)
	findString := reg.FindString(videoinfo.HtmlContent)
	findString = strings.Replace(findString, `strencode2("`, "", -1)
	findString = strings.Replace(findString, `"))`, "", -1)
	parser := JsParser("./md2.js", "strencode2", findString)
	parser = strings.Replace(parser, `<source src='`, "", -1)
	parser = strings.Replace(parser, `' type='application/x-mpegURL'>`, "", -1)
	log.Println(parser)
	title := videoinfo.Title
	log.Println("开始转换视频到mp4：" + title)

	os.MkdirAll(title, 0755)

	path := title + "/" + title + ".mp4"
	videoinfo.Duration = BotUti.ConVtoMp4(parser, path)
	//生成缩略图
	ffmpeg.Input(path).Output(title+"/"+title+".jpg", ffmpeg.KwArgs{"vframes": "1"}).OverWriteOutput().Run()

	filesize := getFileSize(path)
	log.Println("视频大小：" + fmt.Sprintf("%f", filesize))
	if filesize <= 50 {
		sendVideo(title+".mp4", videoinfo,chatId)
	} else {
		//切割视频
		cmd(path)
		files, err := ioutil.ReadDir(title)
		if err != nil {
			panic(err)
		}
		// 获取文件，并输出它们的名字
		for _, file := range files {
			if strings.Contains(file.Name(), title+"_") {
				println(file.Name())
				log.Println("触发发送")
				sendVideo(file.Name(), videoinfo,chatId)
			}
		}
	}

}

//sendVideo
//filename 文件名（切割视频后文件名）
//videoinfo 视频信息
func sendVideo(filename string, videoinfo entity.VideoInfo,chatId int64) {
	log.Println("发送视频:" + filename)
	path := videoinfo.Title + "/" + filename

	v := &tb.Video{
		File:     tb.FromDisk(path),
		Duration: BotUti.VideoLen(path),
		Caption: fmt.Sprintf(captionTemplate, videoinfo.Title, videoinfo.ScCount, videoinfo.Author),

		Thumbnail: &tb.Photo{
			File: tb.FromDisk(videoinfo.Title + "/" + videoinfo.Title + ".jpg"),
		},

		SupportsStreaming: true,
		FileName:          filename,
	}

	b.Send(&tb.Chat{
		ID: chatId,
	}, v, tb.ModeMarkdown)

	//删除临时文件夹
	os.RemoveAll(videoinfo.Title)
}

//cmd 执行 Mp4Box命令
func cmd(pathname string) {
	//51200
	cmd := exec.Command(Mp4BoxPath, "-splits", "51200", pathname)
	var stdoutProcessStatus bytes.Buffer
	cmd.Stdout = io.MultiWriter(ioutil.Discard, &stdoutProcessStatus)
	done := make(chan struct{})
	go func() {
		tick := time.NewTicker(time.Second)
		defer tick.Stop()
		for {
			select {
			case <-done:
				return
			case <-tick.C:
				log.Printf("downloaded: %d", stdoutProcessStatus.Len())
			}
		}
	}()
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to call Run(): %v", err)
	}
	close(done)
}

//getFileSize  返回单位 M
func getFileSize(path string) float64 {

	stat, err := os.Stat(path)
	if err != nil {
		return 0
	}
	log.Println(stat.Size())
	num1, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(stat.Size())/float64(1024)/float64(1024)), 64)
	return num1

}

func cronTaskSendVideo() {
	log.Println("++++++++++++++定时任务开始+++++++++++++++++++")
		c := colly.NewCollector()
		// Find and visit all links
		c.OnHTML("#wrapper > div.container.container-minheight > div.row > div > div > div > div", func(e *colly.HTMLElement) {
			//fmt.Println(e.ChildAttr("a","href"))
			log.Println(e.ChildText("a .video-title"))
			url := e.ChildAttr("a", "href")
			log.Println(url)
			controlSecPage(url,telegramId)

		})

		c.OnRequest(func(r *colly.Request) {
			r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
			fmt.Println("Visiting", r.URL)
		})

		c.Visit("http://91porn.com/index.php")


	log.Println("++++++++++++++定时任务结束+++++++++++++++++++")
}

func JsParser(filePath string, functionName string, args ...interface{}) (result string) {
	//读入文件
	bytes, _ := ioutil.ReadFile(filePath)
	vm := otto.New()
	_, _ = vm.Run(string(bytes))
	value, _ := vm.Call(functionName, nil, args...)
	return value.String()
}
