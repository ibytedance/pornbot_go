package BotUti

import (
	"github.com/go-ego/gse"
)

var (	seg    gse.Segmenter)

func Init() {
	// 加载默认词典
	seg.LoadDict()
	seg.LoadDict("word.txt")
	seg.LoadStop()
}


//CutWords 中文分词
func CutWords(string string)  string{
	return seg.CutStr(DeleteSlice2(seg.CutAll(string)), " ")
}

//DeleteSlice2 删除单个汉字
func DeleteSlice2(a []string) []string{
	j := 0
	for _, val := range a {
		if len(val) != 3 {
			a[j] = val
			j++
		}
	}
	return a[:j]
}