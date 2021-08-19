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
	Mp4BoxPath      = "/root/gpac_public/bin/gcc/MP4Box"
	captionTemplate = `标题: %s
收藏: %s
作者: %s `
	keyword  = "幼"
	//telegramId = *****
	telegramId = *********
)


func main() {
	c := cron.New()
	//定时任务
	c.AddFunc(spec, func() {
		cronTaskSendVideo()
	})
	log.Println("定时任务开启")
	c.Start()
	log.Println("bot开启")
	w := &tb.Webhook{
		Listen:         webhookPort,
		MaxConnections: MaxConnections,
		TLS: &tb.WebhookTLS{
			Key:  serverKey,
			Cert: serverCrt,
		},
		Endpoint: &tb.WebhookEndpoint{PublicURL: webhookUrl},
	}
	b, err = tb.NewBot(tb.Settings{
		Token:  token,
		Poller: w,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	b.SetWebhook(w)
	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, "向我发送91视频链接，获取视频")
	})
	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello World!")
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		text := m.Text
		firstName := m.Sender.FirstName
		lastName := m.Sender.LastName
		log.Println("收到消息："+text, "   来自："+firstName)
		if strings.Contains(text, "viewkey") {
			controlSecPage(text, m.Chat.ID)
		} else if strings.Contains(firstName, keyword) ||
			strings.Contains(lastName, keyword)||
			strings.Contains(text,keyword){
			b.Delete(m)
			send, _ := b.Send(m.Chat, "用户名或者内容违规，已删除")
			time.Sleep(6 *time.Second)
			b.Delete(send)
		}

	})
	b.Start()

	select {} //阻塞主线程停止

}

//匹配两个字符串之间的内容  rep模板 "字符串1(.*?)字符串2"
func regexpUtil(rep string,content string) string {
	compile := regexp.MustCompile(rep)
	submatch := compile.FindAllStringSubmatch(content, -1)
	for _, text := range submatch {
		return text[1]
	}
	return ""
}

//controlSecPage 详情页爬取 发送视频
func controlSecPage(url string, chatId int64) {
	videoinfo, _ := BotUti.GetHttpHtmlContent(url, "#useraction > div:nth-child(1) > span:nth-child(2)", "body")
	parser := regexpUtil("strencode2\\((.*?)\\)\\)",videoinfo.HtmlContent )
	parser = JsParser("./md2.js", "strencode2", parser)
	parser= regexpUtil("src='(.*?)' type=", parser)
	if chatId != telegramId {
		if parser == "" {
			b.Send(&tb.Chat{
				ID: chatId,
			}, "请检查视频地址是否正确")
			return
		}
		b.Send(&tb.Chat{
			ID: chatId,
		}, "视频真实地址："+parser)
	}

	log.Println(parser)

	title := videoinfo.Title
	log.Println("开始转换视频到mp4：" + title)

	os.MkdirAll(title, 0755)

	path := title + "/" + title + ".mp4"
	videoinfo.Duration, err = BotUti.ConVtoMp4(parser, path)
	if err != nil {
		return
	}
	//生成缩略图
	ffmpeg.Input(path).Output(title+"/"+title+".jpg", ffmpeg.KwArgs{"vframes": "1"}).OverWriteOutput().Run()

	filesize := getFileSize(path)
	log.Println("视频大小：" + fmt.Sprintf("%f", filesize))
	if filesize <= 50 {
		sendVideo(title+".mp4", videoinfo, chatId)
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
				size, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(file.Size())/float64(1024)/float64(1024)), 64)
				log.Println("分割视频大小：",file.Name(),"-",size,"M")
				sendVideo(file.Name(), videoinfo, chatId)
			}
		}
	}

	//删除临时文件夹
	os.RemoveAll(videoinfo.Title)
}

//markdown转义字符处理
// '_'、'*'、'`'、'['
func escapeMarkDown(markdownStr string) string {
	nowords := []string{"_","*","`","["}
	for _, word := range nowords {
		if strings.Contains(markdownStr,word) {
			markdownStr=strings.ReplaceAll(markdownStr,word,"\\"+word);
		}
	}
	return markdownStr
}

//sendVideo
//filename 文件名（切割视频后文件名）
//videoinfo 视频信息
func sendVideo(filename string, videoinfo entity.VideoInfo, chatId int64) {
	path := videoinfo.Title + "/" + filename
	videoLen, err := BotUti.VideoLen(path)
	if err!=nil {
		panic(err)
	}
	v := &tb.Video{
		File:     tb.FromDisk(path),
		Duration: videoLen,
		Caption:  fmt.Sprintf(captionTemplate, escapeMarkDown(filename), videoinfo.ScCount, videoinfo.Author),
		Thumbnail: &tb.Photo{
			File: tb.FromDisk(videoinfo.Title + "/" + videoinfo.Title + ".jpg"),
		},
		SupportsStreaming: true,
		FileName:          filename,
	}

	_, err  = b.Send(&tb.Chat{
		ID: chatId,
	}, v, tb.ModeMarkdown)

	if err != nil {
		panic(err)
	}
}

//cmd 执行 Mp4Box命令  /root/gpac_public/bin/gcc/MP4Box -splits 20176 aa.mp4 -out  aa%d.mp4
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
	num1, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(stat.Size())/float64(1024)/float64(1024)), 64)
	return num1

}

func cronTaskSendVideo() {
	log.Println("++++++++++++++定时任务开始+++++++++++++++++++")
	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("#wrapper > div.container.container-minheight > div.row > div > div > div > div", func(e *colly.HTMLElement) {
		//fmt.Println(e.ChildAttr("a","href"))
		url := e.ChildAttr("a", "href")
		controlSecPage(url, telegramId)

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
