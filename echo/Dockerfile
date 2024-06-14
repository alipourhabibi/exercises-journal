FROM golang:alpine
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s -w' -o http-echo

FROM alpine:latest
COPY --from=0 /app/config/config.yaml /etc/http-echo/config.yaml
COPY --from=0 /app/http-echo /bin/http-echo
CMD ["http-echo", "run", "--", "--config", "/etc/http-echo/config.yaml"]
