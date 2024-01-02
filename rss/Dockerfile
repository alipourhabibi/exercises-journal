FROM golang:alpine
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s -w' -o rss

FROM alpine:latest
COPY --from=0 /app/config/config.yaml /etc/rss/config.yaml
COPY --from=0 /app/config/rss.yaml /etc/rss/rss.yaml
COPY --from=0 /app/rss /bin/rss
CMD ["rss", "run", "--", "--config", "/etc/rss/config.yaml", "--rssfile", "/etc/rss/rss.yaml"]
