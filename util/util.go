package BotUti

import (
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//genIpaddr 随机ip
func genIpaddr() string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
}
//RegexpUtil 匹配两个字符串之间的内容  rep模板 "字符串1(.*?)字符串2"
func RegexpUtil(rep string,content string) string {
	compile := regexp.MustCompile(rep)
	submatch := compile.FindAllStringSubmatch(content, -1)
	for _, text := range submatch {
		return text[1]
	}
	return ""
}

//VideoLen 获取视频 时长单位 s 宽度 高度
func VideoLen(url string) (int,int,int,error)  {
	//获取视频信息
	args := ffmpeg.KwArgs{"rw_timeout":rw_timeout}
	probe, err := ffmpeg.Probe(url,args)
	if err != nil {
		return 0,0,0, err
	}
	var videoIf videoInfo
	err = json.Unmarshal([]byte(probe), &videoIf)
	if err != nil {
		return 0,0,0, err
	}
	float, err := strconv.ParseFloat(videoIf.Format.Duration,64)
	width := videoIf.Streams[0].Width
	height := videoIf.Streams[0].Height
	if err != nil {
		return 0,0,0, err
	}
	return int(float),width,height,err
}


//EscapeMarkDown  markdown转义字符处理 '_'、'*'、'`'、'['
func EscapeMarkDown(markdownStr string) string {
	nowords := []string{"_","*","`","["}
	for _, word := range nowords {
		if strings.Contains(markdownStr,word) {
			markdownStr=strings.ReplaceAll(markdownStr,word,"\\"+word);
		}
	}
	return markdownStr
}

//GetFileSize 返回单位 M
func GetFileSize(path string) (float64,error) {
	stat, err := os.Stat(path)
	if err != nil {
		return 0,err
	}
	return strconv.ParseFloat(fmt.Sprintf("%.2f", float64(stat.Size())/float64(1024)/float64(1024)), 64)
}

//JsParser 执行js文件
func JsParser(filePath string, functionName string, args ...interface{}) (result string) {
	//读入文件
	bytes, _ := ioutil.ReadFile(filePath)
	vm := otto.New()
	_, _ = vm.Run(string(bytes))
	value, _ := vm.Call(functionName, nil, args...)
	return value.String()
}
