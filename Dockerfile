FROM golang:1.18 AS BUILDER

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN go build -o azuretts .

FROM debian:11.4-slim

ENV DEBIAN_FRONTEND noninteractive

RUN sed -i s@/[a-z]*.debian.org/@/mirrors.tuna.tsinghua.edu.cn/@g /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -y ca-certificates tzdata && \
    ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY templates /app/templates
COPY static /app/static
COPY ssml.xml /app/ssml.xml
COPY --from=BUILDER /app/azuretts /app/azuretts

RUN chmod +x /app/azuretts

ENTRYPOINT ["/app/azuretts"]