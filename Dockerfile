FROM golang:1.18 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN GOOS=linux GOARCH=amd64 go build ./cmds/kcl-go

FROM kcllang/kcl

WORKDIR /app

COPY --from=builder /app/kcl-go .

ENV PATH="/app/bin:${PATH}"
ENV LANG=en_US.utf8

CMD ["bash"]
