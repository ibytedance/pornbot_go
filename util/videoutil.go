package BotUti

import (
	"encoding/json"
	"github.com/matryer/try"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"strconv"
)

//ConVtoMp4 转换为mp4格式
func ConVtoMp4(videourl string, pathname string) (int,error) {
	var videoLen int
	try.Do(func(attempt int) (retry bool, err error) {
		videoLen, err = VideoLen(videourl)
		if err !=nil{
			log.Println("Run error,重试 - ", err,"-",attempt,"次")
		} else {
			log.Println("Run ok - ", videoLen)
		}
		// 重试5次
		return attempt < 5, err
	})
	
	//大于四分钟 截取前10秒 -rw_timeout 5000000
	args := ffmpeg.KwArgs{"c:v": "libx264", "threads": "2","crf": "21" ,"rw_timeout":"10000000"}
	if videoLen > 240 {
		args = ffmpeg.KwArgs{"c:v": "libx264","threads": "2" ,"crf": "21" ,"rw_timeout":"10000000","ss": "00:00:10"}
	}
	err = try.Do(func(attempt int) (retry bool, err error) {
		err = ffmpeg.Input(videourl).Output(pathname, args).OverWriteOutput().Run()
		if err != nil {
			log.Println("Run error 重试 -", attempt, ":", err)
		} else {
			log.Println("视频转换成功 ")
		}
		// 重试5次
		return attempt < 5, err
	})

	return videoLen,err
}

//VideoLen 获取视频时长 单位 s
func VideoLen(url string) (int,error)  {
	//获取视频信息
	args := ffmpeg.KwArgs{"rw_timeout":"10000000"}
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



type videoInfo struct {
	Streams []Streams `json:"streams"`
	Format Format `json:"format"`
}
type Disposition struct {
	Default int `json:"default"`
	Dub int `json:"dub"`
	Original int `json:"original"`
	Comment int `json:"comment"`
	Lyrics int `json:"lyrics"`
	Karaoke int `json:"karaoke"`
	Forced int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
	VisualImpaired int `json:"visual_impaired"`
	CleanEffects int `json:"clean_effects"`
	AttachedPic int `json:"attached_pic"`
	TimedThumbnails int `json:"timed_thumbnails"`
}
type Tags struct {
	VariantBitrate string `json:"variant_bitrate"`
}
type Streams struct {
	Index int `json:"index"`
	CodecName string `json:"codec_name"`
	CodecLongName string `json:"codec_long_name"`
	Profile string `json:"profile"`
	CodecType string `json:"codec_type"`
	CodecTagString string `json:"codec_tag_string"`
	CodecTag string `json:"codec_tag"`
	Width int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
	CodedWidth int `json:"coded_width,omitempty"`
	CodedHeight int `json:"coded_height,omitempty"`
	ClosedCaptions int `json:"closed_captions,omitempty"`
	HasBFrames int `json:"has_b_frames,omitempty"`
	PixFmt string `json:"pix_fmt,omitempty"`
	Level int `json:"level,omitempty"`
	ChromaLocation string `json:"chroma_location,omitempty"`
	Refs int `json:"refs,omitempty"`
	IsAvc string `json:"is_avc,omitempty"`
	NalLengthSize string `json:"nal_length_size,omitempty"`
	RFrameRate string `json:"r_frame_rate"`
	AvgFrameRate string `json:"avg_frame_rate"`
	TimeBase string `json:"time_base"`
	StartPts int `json:"start_pts"`
	StartTime string `json:"start_time"`
	BitsPerRawSample string `json:"bits_per_raw_sample,omitempty"`
	Disposition Disposition `json:"disposition"`
	Tags Tags `json:"tags"`
	SampleFmt string `json:"sample_fmt,omitempty"`
	SampleRate string `json:"sample_rate,omitempty"`
	Channels int `json:"channels,omitempty"`
	ChannelLayout string `json:"channel_layout,omitempty"`
	BitsPerSample int `json:"bits_per_sample,omitempty"`
}
type Format struct {
	Filename string `json:"filename"`
	NbStreams int `json:"nb_streams"`
	NbPrograms int `json:"nb_programs"`
	FormatName string `json:"format_name"`
	FormatLongName string `json:"format_long_name"`
	StartTime string `json:"start_time"`
	Duration string `json:"duration"`
	Size string `json:"size"`
	BitRate string `json:"bit_rate"`
	ProbeScore int `json:"probe_score"`
}