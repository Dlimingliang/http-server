#FROM alpine
#COPY bin/amd64/httpserver httpserver
#EXPOSE 9090
#ENTRYPOINT ["./httpserver"]


#多阶段构建 当你没有可执行文件的时候，需要重新去编译，创建二进制文件
FROM golang:1.16-alpine AS build
WORKDIR /go/src/project/
COPY . /go/src/project/
ENV GOPROXY=https://goproxy.cn,direct
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/httpserver

#FROM alpine
FROM scratch
COPY --from=build /go/src/project/bin/httpserver httpserver
ENTRYPOINT ["./httpserver"]