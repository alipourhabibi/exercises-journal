FROM golang:alpine
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s -w' -o httpfileserver

FROM alpine:latest
COPY --from=0 /app/config/config.yaml /etc/httpfileserver/config.yaml
COPY --from=0 /app/httpfileserver /bin/httpfileserver
CMD ["httpfileserver"]
