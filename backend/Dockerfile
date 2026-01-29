FROM --platform=$BUILDPLATFORM golang:bookworm AS ytdlp_cache
ARG TARGETOS
ARG TARGETARCH
ARG BUILDPLATFORM
RUN apt update && apt install -y wget tar xz-utils

RUN BUILDARCH=$(echo "$BUILDPLATFORM" | cut -d'/' -f2); \
    wget https://github.com/upx/upx/releases/download/v5.0.0/upx-5.0.0-${BUILDARCH}_linux.tar.xz && \
    tar -xf upx-5.0.0-${BUILDARCH}_linux.tar.xz && \
    cp upx-5.0.0-${BUILDARCH}_linux/upx /usr/local/bin/ && \
    chmod +x /usr/local/bin/upx
ENV DEBIAN_FRONTEND=noninteractive
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o whishper . && \
    upx whishper

RUN chmod a+rx whishper

FROM --platform=$BUILDPLATFORM python:3-bookworm
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY --from=ytdlp_cache /app/whishper ./whishper 
RUN chmod a+rx ./whishper
RUN pip install yt-dlp && ln -s /usr/local/bin/yt-dlp /bin/yt-dlp
RUN python -m pip install -U --pre "yt-dlp[default]"

RUN mkdir /app/uploads

EXPOSE 8080
ENTRYPOINT [ "/app/whishper", "-addr", ":8080" ]
