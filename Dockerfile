FROM golang:1.17 AS builder

ENV GOPATH  /go_workspace
ENV APP_DIR $GOPATH/src/
RUN mkdir -p $APP_DIR
WORKDIR $APP_DIR


COPY . .
RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 go install -a std

RUN go build -o app

FROM alpine:3.15.6
RUN apk --update add ca-certificates
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
ENV APP_HOME /go/src/app
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY --from=builder /go_workspace/src/app .
COPY conf conf
# ARG FILESRC
# COPY ${FILESRC}/conf conf
# RUN go mod download

# RUN go build -o /marketing

EXPOSE 8081
EXPOSE 10443

# USER nonroot:nonroot
# RUN groupadd nonroot
# RUN useradd -u 8877 nonroot
# RUN usermod -a -G nonroot nonroot
# USER nonroot



ENTRYPOINT ["/go/src/app/app"]

