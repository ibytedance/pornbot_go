package main


import (
	"bytes"
	"fmt"
	"github.com/go-ego/gse"
	"github.com/gocolly/colly"
	"github.com/matryer/try"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	tb "gopkg.in/tucnak/telebot.v2"
	"github.com/robfig/cron"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"pornbot/entity"
	"pornbot/util"
	_ "pornbot/util"
	"strconv"
	"strings"
	"time"
)
var (
	seg    gse.Segmenter
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
	Mp4BoxPath      = "/gpac_public/bin/gcc/MP4Box"
	captionTemplate = `标题: %s
收藏: %s
作者: %s `
	keyword  = "幼"
	//telegramId = *****
	telegramId = *********
)

func init() {
	BotUti.Init()
}

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
		Token:  token,
		Poller: &tb.Webhook{
			Listen:         webhookPort,
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
		b.Send(m.Sender, `向我发送91视频链接，获取视频,有问题请留言 @bzhzq`)
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


//controlSecPage 详情页爬取 发送视频
func controlSecPage(url string, chatId int64) {
	var videoinfo entity.VideoInfo
	try.Do(func(attempt int) (retry bool, err error) {
		flag:=false
		if attempt>3 {
			flag=true
		}
		videoinfo, err = BotUti.GetHttpHtmlContent(url, "#useraction > div:nth-child(1) > span:nth-child(2)", "body",flag)
		if err !=nil{
			log.Println("Run error,重试 - ", err,"-",attempt,"次")
		} else {
			log.Println("Run ok - ", "详情页爬取成功")
		}
		// 重试5次
		return attempt < 5, err
	})

	parser := BotUti.RegexpUtil("strencode2\\((.*?)\\)\\)",videoinfo.HtmlContent )
	parser = BotUti.JsParser("./md2.js", "strencode2", parser)
	parser= BotUti.RegexpUtil("src='(.*?)' type=", parser)
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

	filesize := 0.0
	filesize, err = BotUti.GetFileSize(path)
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



//sendVideo
//filename 文件名（切割视频后文件名）
//videoinfo 视频信息
func sendVideo(filename string, videoinfo entity.VideoInfo, chatId int64) {
	path := videoinfo.Title + "/" + filename
	videoLen,w,h, err := BotUti.VideoLen(path)
	log.Println("视频宽度：",w)
	if err!=nil {
		panic(err)
	}
	newFileName := strings.ReplaceAll(filename, ".mp4", "")
	title := BotUti.EscapeMarkDown(newFileName)
	//中文分词
	words := BotUti.CutWords(newFileName)
	log.Println("分词结果: ", words)
	v := &tb.Video{
		File:     tb.FromDisk(path),
		Duration: videoLen,
		Caption:  fmt.Sprintf(captionTemplate, title, videoinfo.ScCount,videoinfo.Author, BotUti.EscapeMarkDown(words)),
		Thumbnail: &tb.Photo{
			File: tb.FromDisk(videoinfo.Title + "/" + videoinfo.Title + ".jpg"),
		},
		SupportsStreaming: true,
		FileName:          filename,
		Width: w,
		Height: h,
	}

	_, err  = b.Send(&tb.Chat{
		ID: chatId,
	}, v, tb.ModeMarkdown)

	if err != nil {
		panic(err)
	}
	log.Println("发送视频成功")
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
		log.Println("Visiting", r.URL)
	})

	c.Visit("http://91porn.com/index.php")

	log.Println("++++++++++++++定时任务结束+++++++++++++++++++")
}

