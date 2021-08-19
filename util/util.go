package BotUti

import (
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io/ioutil"
	"log"
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
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	log.Println(ip)
	return ip
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

//VideoLen 获取视频时长 单位 s
func VideoLen(url string) (int,error)  {
	//获取视频信息
	args := ffmpeg.KwArgs{"rw_timeout":rw_timeout}
	probe, err := ffmpeg.Probe(url,args)
	if err != nil {
		return 0, err
	}
	var videoIf videoInfo
	err = json.Unmarshal([]byte(probe), &videoIf)
	if err != nil {
		return 0, err
	}
	float, err := strconv.ParseFloat(videoIf.Format.Duration,64)
	if err != nil {
		return 0,err
	}
	return int(float),err
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
func GetFileSize(path string) float64 {

	stat, err := os.Stat(path)
	if err != nil {
		return 0
	}
	num1, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(stat.Size())/float64(1024)/float64(1024)), 64)
	return num1

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
