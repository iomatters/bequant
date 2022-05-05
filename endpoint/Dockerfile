FROM golang:1.16-alpine

ARG VERSION

RUN apk update && \
	apk add gcc musl-dev

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
COPY *.go ./

RUN 	go mod download && \
	go mod vendor && \
	go mod tidy 

RUN go test -mod=vendor 

RUN go build -mod=vendor -a -ldflags "-w -extldflags "-static" -X main.BuildVersion=$VERSION" -o bin/endpoint . 


FROM node:8.12-alpine
WORKDIR /app

COPY --from=0 /build/bin/endpoint .
COPY config.toml ./

RUN addgroup non-root && \
        adduser -S non-root -G non-root -h /app -s /usr/sbin/nologin && \
        chown -R non-root:non-root /app

USER non-root

EXPOSE 8080

ENTRYPOINT [ "/bin/sh", "-l", "-c" ]
CMD [ "/app/endpoint --config /app/config.toml" ]

