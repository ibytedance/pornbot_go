# pornbot_go

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




