FROM alpine
COPY bin/amd64/httpserver httpserver
RUN addgroup --gid 5000 newuser && adduser -h /home/newuser -s /bin/sh -k /dev/null newuser --uid 5000 -G newuser  -S newuser
EXPOSE 9090
USER newuser
ENTRYPOINT ["./httpserver"]

#
##多阶段构建 当你没有可执行文件的时候，需要重新去编译，创建二进制文件
#FROM golang:1.16-alpine AS build
#WORKDIR /go/src/project/
#COPY . /go/src/project/
#ENV GOPROXY=https://goproxy.cn,direct
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/httpserver
#
##FROM alpine
#FROM scratch
#COPY --from=build /go/src/project/bin/httpserver httpserver
#ENTRYPOINT ["./httpserver"]