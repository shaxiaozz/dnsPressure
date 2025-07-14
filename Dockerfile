FROM golang:1.24.2-alpine as builder
WORKDIR /data/dnsPressure-code
RUN apk add --no-cache upx ca-certificates tzdata
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o dnsPressure

FROM shaxiaozz/centos:7.9-google-chrome as runner
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /data/dnsPressure-code/dnsPressure /dnsPressure
RUN chmod +x /dnsPressure
CMD ["/dnsPressure","-c","10","www.baidu.com"]