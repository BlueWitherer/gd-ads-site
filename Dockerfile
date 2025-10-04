# Build stage for Go backend
FROM        --platform=$TARGETOS/$TARGETARCH golang:1.24-alpine
WORKDIR     /build
COPY        bridge/ ./
RUN         go build -o bridge .

# Main stage with Node.js
FROM        --platform=$TARGETOS/$TARGETARCH node:24-bookworm-slim

RUN         apt update \
            && apt -y install ffmpeg iproute2 git sqlite3 libsqlite3-dev python3 python3-dev ca-certificates dnsutils tzdata zip tar curl build-essential libtool iputils-ping libnss3 tini \
            && useradd -m -d /home/container container

RUN         npm install --global npm@latest typescript ts-node @types/node

# Copy the built Go binary from the builder stage
COPY        --from=go-builder /build/bridge /usr/local/bin/bridge

USER        container
ENV         USER=container HOME=/home/container
WORKDIR     /home/container

STOPSIGNAL SIGINT

COPY        --chown=container:container ./../entrypoint.sh /entrypoint.sh
RUN         chmod +x /entrypoint.sh
ENTRYPOINT    ["/usr/bin/tini", "-g", "--"]
CMD         ["/entrypoint.sh"]