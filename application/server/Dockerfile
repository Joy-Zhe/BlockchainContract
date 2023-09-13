FROM golang:1.18 AS server
ENV GO111MODULE=on
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /root/joyzheng2003
COPY . .
RUN CGO_ENABLED=0 go build -v -o "server" .

FROM scratch
WORKDIR /root/joyzheng2003/
COPY --from=server /root/joyzheng2003/server ./
COPY config.yaml ./
ENTRYPOINT ["./server"]