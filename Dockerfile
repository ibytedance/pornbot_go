FROM alpine

MAINTAINER jw-star

RUN apk update && \
    apk add --no-cache  chromium go ffmpeg && \
    rm -rf /var/cache/apk/*
    
# go环境变量
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH  

COPY gpac_public/ gpac_public/

WORKDIR /home   


CMD ["sh","start.sh"]
