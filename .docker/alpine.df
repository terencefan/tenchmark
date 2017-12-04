FROM alpine

COPY . /tenchmark

WORKDIR /tenchmark

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --update go git libc-dev \
    && export GOPATH=`pwd` \
    && export GOBIN=/usr/bin \
    && go get -x \
    && go install \
    && echo $GOPATH \
    && rm -r /tenchmark \
    && apk del --purge go git libc-dev

WORKDIR /
