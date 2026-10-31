FROM golang:1.27-alpine

RUN --mount=type=cache,target=/var/cache/apk \
  apk --update add tzdata git && \
  cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
  apk del tzdata

WORKDIR /app

RUN apk upgrade && apk add git

RUN go install github.com/air-verse/air@latest

COPY . .

CMD ["air", "-c", ".air.toml"]
