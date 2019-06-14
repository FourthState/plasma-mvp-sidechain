# Simple usage with a mounted data directory:
# > docker build -t plasma .
# 
# It is important to link the right volume to the container. The volume contains configuration files used to launch the server dameon
#
# plasmad
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.plasmad:/root/.plasmad
FROM golang:1.12-alpine3.9 AS builder

RUN apk add git make npm curl gcc libc-dev && \
    mkdir -p /root/plasma-mvp-sidechain

# install dependencies
WORKDIR /root/plasma-mvp-sidechain
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# install plasmad and plasmacli
RUN go install -mod=readonly ./cmd/plasmad ./cmd/plasmacli

### Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over the plasmad and plasmacli binaries from the build-env
COPY --from=builder /go/bin/plasmad /usr/bin/plasmad
COPY --from=builder /go/bin/plasmacli /usr/bin/plasmacli

# As an executable, the dameon will simply start
CMD ["plasmad", "start"]
