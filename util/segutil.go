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

	cut := seg.CutAll(string)
	return seg.CutStr(cut, " ")
}