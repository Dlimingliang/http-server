#FROM alpine
#COPY bin/amd64/httpserver httpserver
#EXPOSE 9090
#ENTRYPOINT ["./httpserver"]


#多阶段构建 当你没有可执行文件的时候，需要重新去编译，创建二进制文件
FROM golang:1.16-alpine AS build
RUN apk add --no-cache git
RUN go get github.com/Dlimingliang/http-server

FROM alpine
#FROM scrath
COPY --from=build /go/bin/http-server httpserver
ENTRYPOINT ["./httpserver"]