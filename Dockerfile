FROM golang:1.15-alpine AS builder
ENV GO111MODULE=on
RUN apk add git
WORKDIR /biometrics
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN cd src/main/ && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o server .

FROM alpine:3.8
RUN apk add --no-cache --update ca-certificates tzdata curl
RUN cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime
RUN echo "Asia/Ho_Chi_Minh" >  /etc/timezone
COPY --from=builder /biometrics/src/main/server  /server
COPY ./config.yml  /config.yml

WORKDIR /
EXPOSE 18502
CMD ["/server"]