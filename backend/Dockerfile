FROM --platform=$BUILDPLATFORM golang:bookworm AS ytdlp_cache
ARG TARGETOS
ARG TARGETARCH
RUN apt update && apt install -y wget tar xz-utils
RUN wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -O /usr/local/bin/yt-dlp
RUN chmod a+rx /usr/local/bin/yt-dlp

RUN wget https://github.com/upx/upx/releases/download/v5.0.0/upx-5.0.0-arm64_linux.tar.xz && \
    tar -xf upx-5.0.0-arm64_linux.tar.xz && \
    cp upx-5.0.0-arm64_linux/upx /usr/local/bin/ && \
    chmod +x /usr/local/bin/upx
ENV DEBIAN_FRONTEND=noninteractive
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o whishper . && \
    upx whishper

RUN chmod a+rx whishper

FROM --platform=$BUILDPLATFORM python:3-bookworm
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY --from=ytdlp_cache /app/whishper ./whishper 
RUN chmod a+rx ./whishper
COPY --from=ytdlp_cache /usr/local/bin/yt-dlp /bin/yt-dlp

RUN mkdir /app/uploads

EXPOSE 8080
ENTRYPOINT [ "/app/whishper", "-addr", ":8080" ]
