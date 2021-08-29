# pornbot_go


###  特点

破解91视频的播放限制、理论上可以无限下载

切除长视频(4分钟)播放开始的静态帧(10秒)

由于电报Bot单次发送最大50M文件，切割发送视频(MP4Box大法好!!!)

为标题添加中文分词，解决电报对中文搜索的问题

重试机制，网络超时重试

bot采用 webhook 方式，并发（同时可以做负载均衡）

向机器人([@porn_91Bot](https://t.me/porn_91Bot))发送链接，可以 `获取视频真实地址` 并 `下载视频`



### 软件安装



#### 安装 docker
```
curl -fsSL get.docker.com -o get-docker.sh && sh get-docker.sh --mirror Aliyun&&systemctl enable docker&&systemctl start docker

```

#### 拉取项目到本地
```
 git clone https://github.com/jw-star/pornbot_go.git .
```

#### 启动docker

```
docker run -d -i -p 443:443 -v /pornbot:/home --restart=always --name pornbot jwstar/pornbot_tg
```
#### 常用docker命令

```
删除容器
docker rm -f pornbot
删除镜像
docker rmi jwstar/pornbot_tg
进入容器
docker exec -it pornbot /bin/bash
```

### 配置信息
注意事项参考
https://core.telegram.org/bots/api#setwebhook
```

1.只要设置了传出 webhook，您就无法使用getUpdates接收更新。
2.要使用自签名证书，您需要使用证书参数上传您的公钥证书。请上传为 InputFile，发送字符串将不起作用。3. Webhooks当前支持的端口：443, 80, 88, 8443。

```

证书申请
根据 cloudflare dns申请 ，其他方法参考 https://github.com/acmesh-official/acme.sh/wiki/%E8%AF%B4%E6%98%8E
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
  //Mp4Box路径（docker镜像已包含）
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

杀死后台
```
ps -ef|grep main
```

```
kill pid
```


### 测试

发送 /hello 到机器人

得到回复  `Hello World!`


### 鸣谢

https://github.com/acmesh-official/acme.sh/wiki/%E8%AF%B4%E6%98%8E

https://github.com/tucnak/telebot

https://github.com/gocolly/colly

https://github.com/chromedp/chromedp

https://github.com/go-ego/gse

