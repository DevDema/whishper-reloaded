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
# ffmpeg muxes YouTube's separate video and audio streams (stream copy).
# Deno is the JavaScript runtime yt-dlp needs to read current YouTube formats.
RUN apt-get update && apt-get install -y --no-install-recommends \
        ffmpeg ca-certificates curl unzip \
    && case "${TARGETARCH}" in \
         amd64) DENO_TARGET=x86_64-unknown-linux-gnu ;; \
         arm64) DENO_TARGET=aarch64-unknown-linux-gnu ;; \
         *) echo "unsupported architecture: ${TARGETARCH}" && exit 1 ;; \
       esac \
    && curl -fsSL -o /tmp/deno.zip "https://github.com/denoland/deno/releases/download/v2.9.7/deno-${DENO_TARGET}.zip" \
    && unzip /tmp/deno.zip -d /usr/local/bin \
    && chmod +x /usr/local/bin/deno \
    && rm /tmp/deno.zip \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir /app/uploads

EXPOSE 8080
ENTRYPOINT [ "/app/whishper", "-addr", ":8080" ]
