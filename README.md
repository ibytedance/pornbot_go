# pornbot_go


### 软件安装

环境 debian10

#### 安装 ffmpeg
```
apt install -y ffmpeg
```

#### 安装 go


下载  https://golang.org/dl/
```
wget https://golang.org/dl/go1.16.7.linux-amd64.tar.gz
```

解压
```
tar -C /usr/local -zxvf  go1.16.7.linux-amd64.tar.gz
```

环境变量

```
vi /etc/profile
# 在最后一行添加
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin
# 保存退出后source一下（vim 的使用方法可以自己搜索一下）
source /etc/profile
```
go版本号查看

```
go version
```
#### 安装 chrome

```
wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
```

```
apt install  -y ./google-chrome-stable_current_amd64.deb
```

最后提示下面信息没关系

```
N: Download is performed unsandboxed as root as file '/root/pornbot/google-chrome-stable_current_amd64.deb' couldn't be accessed by user '_apt'. - pkgAcquire::Run (13: Permission denied)
```

#### 获取MP4Box包

debian10编译好的
```
https://github.com/jw-star/myFigurebed/releases/download/1.00/gpac.tar.gz
```


### 配置信息
注意事项参考
https://core.telegram.org/bots/api#setwebhook
```

1.只要设置了传出 webhook，您就无法使用getUpdates接收更新。
2.要使用自签名证书，您需要使用证书参数上传您的公钥证书。请上传为 InputFile，发送字符串将不起作用。3. Webhooks当前支持的端口：443, 80, 88, 8443。

```

证书申请

```

curl  https://get.acme.sh | sh -s email=xxxxx@xxx.xxx


export CF_Key="xxxxxx"

export CF_Email="xxx@xxx.xxx"

acme.sh   --issue   --dns dns_cf   -d xxxx.xxxx.com 
//设置证书位置到项目下
acme.sh  --installcert  -d  xxxx.xxxx.com     \
        --key-file   /root/porn/server.key \
        --fullchain-file /root/porn/server.crt

```


```
var (
	err         error
	//定时任务的cron表达式
	spec        = "0 0 5 * * ?"
  // telegram Token
	token       = "telegram Token"
  // webhook url
	webhookUrl  = "webhook url"
  //webhook 端口  443, 80, 88, 8443
	webhookPort = ":443"
	serverCrt   = "server.crt"
	serverKey   = "server.key"
	//bot并发数
	MaxConnections = 40
	b           *tb.Bot
)
const (
  //Mp4Box路径
	Mp4BoxPath = "/root/gpac_public/bin/gcc/MP4Box"
	//视频描述模板
	captionTemplate = `标题: %s
收藏: %s
作者: %s `
	//定时任务发送的群组Id
	telegramId = -222222
)

````


### 运行

```
go run main.go
```
后台运行 来自 https://github.com/icattlecoder/godaemon
```
go build main.go
./main -d=true
```


### 测试

发送 /hello 到机器人

得到回复  `Hello World!`


### 鸣谢

https://github.com/acmesh-official/acme.sh/wiki/%E8%AF%B4%E6%98%8E

https://github.com/tucnak/telebot

https://github.com/gocolly/colly

https://github.com/chromedp/chromedp



