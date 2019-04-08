FROM golang:latest 
WORKDIR $GOPATH/src/github.com/FourthState/plasma-mvp-sidechain
COPY . .
RUN curl -L -s https://github.com/golang/dep/releases/download/v0.5.1/dep-linux-amd64 -o $GOPATH/bin/dep
RUN chmod +x $GOPATH/bin/dep
RUN dep ensure -vendor-only
WORKDIR $GOPATH/src/github.com/FourthState/plasma-mvp-sidechain/server/plasmad
RUN go install
WORKDIR $GOPATH/src/github.com/FourthState/plasma-mvp-sidechain/client/plasmacli
RUN go install
WORKDIR $GOPATH/src/github.com/FourthState/plasma-mvp-sidechain
RUN plasmad init
CMD tail -f /dev/null
