package BotUti

import (
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/context"
	"log"
	"pornbot/entity"
	"strings"
	"time"
)

var (
	err error
)

//GetHttpHtmlContent 获取详情也内容
// url 地址
//selector 等待selector节点出现
//sel 过滤内容
func GetHttpHtmlContent(url string, selector string, sel interface{}) (entity.VideoInfo, error) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
		chromedp.Flag("mute-audio", false), // 关闭声音
	}
	//初始化参数，先传一个空的数据
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	// create context
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	// 执行一个空task, 用提前创建Chrome实例
	var res string
	err = chromedp.Run(chromeCtx, setheaders(
		"",
		map[string]interface{}{
			"Accept-Language": "zh-cn,zh;q=0.5",
		},
		&res,
	))

	//创建一个上下文，超时时间为40s
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 8*time.Second)
	defer cancel()

	var videoInfo entity.VideoInfo
	err = chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector),
		//标题
		chromedp.TextContent("#videodetails > h4", &videoInfo.Title, chromedp.ByQuery),
		//收藏数
		chromedp.TextContent("#useraction > div:nth-child(1) > span:nth-child(4) > span",
			&videoInfo.ScCount, chromedp.ByQuery),
		//作者
		chromedp.TextContent("#videodetails-content > div:nth-child(2) > span.title-yakov > a:nth-child(1) > span",
			&videoInfo.Author, chromedp.ByQuery),
		//全部文本
		chromedp.OuterHTML(sel, &videoInfo.HtmlContent, chromedp.ByQuery),
	)
	if err != nil {
		log.Println("Run err : %v\n", err)
		return videoInfo, err
	}
	//去除 空格
	videoInfo.Title = strings.Replace(videoInfo.Title, " ", "", -1)
	// 去除换行符
	videoInfo.Title = strings.Replace(videoInfo.Title, "\n", "", -1)
	return videoInfo, nil
}
// setheaders returns a task list that sets the passed headers.
func setheaders(host string, headers map[string]interface{}, res *string) chromedp.Tasks {
	return chromedp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers(headers)),
		chromedp.Navigate(host),
		chromedp.Text(`#result`, res, chromedp.ByID, chromedp.NodeVisible),
	}
}
