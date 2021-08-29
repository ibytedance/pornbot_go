FROM alpine

MAINTAINER jw-star

RUN apk update && \
    apk add --no-cache tzdata bash chromium go ffmpeg && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    rm -rf /var/cache/apk/* /tmp/* /var/tmp/* $HOME/.cache

# go环境变量
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

COPY gpac_public/ gpac_public/

WORKDIR /home


CMD ["sh","start.sh"]
