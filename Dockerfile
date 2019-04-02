# Simple usage with a mounted data directory:
# > docker build -t plasma .
# 
# It is important to link the right volume to the container. The volume contains configuration files used to launch the server dameon
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.plasmad:/root/.plasmad -v ~/.plasmacli:/root/.plasmacli plasmad plasmad init
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.plasmad:/root/.plasmad -v ~/.plasmacli:/root/.plasmacli plasmad plasmad start
FROM golang:1.11-alpine3.8 AS builder

RUN apk add git make npm curl gcc libc-dev && \
    mkdir -p $GOPATH/src/github.com/FourthState/plasma-mvp-sidechain

WORKDIR $GOPATH/src/github.com/FourthState/plasma-mvp-sidechain
COPY Gopkg.toml Gopkg.lock ./

# install dep, vendor dependencies, and copy over the repository
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && dep ensure -vendor-only
COPY . .

# install plasmad
RUN cd server/plasmad && go install && cd ../../client/plasmacli && go install

### Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over server and client binaries from the build-env
COPY --from=builder /go/bin/plasmad /usr/bin/plasmad
COPY --from=builder /go/bin/plasmacli /usr/bin/plasmacli

# Run plasmad by default, omit entrypoint to ease using container with plasmacli
CMD ["plasmad"]
