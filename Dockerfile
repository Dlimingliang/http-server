FROM alpine
COPY bin/amd64/httpserver httpserver
EXPOSE 9090
ENTRYPOINT ["./httpserver"]