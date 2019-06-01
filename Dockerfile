FROM golang:1.12.5-alpine3.9 AS builder

WORKDIR /go/src/github.com/arkticman/go-armored-warship

# Install Mage so we can build correctly.
RUN apk add git && \
    go get -u -d github.com/magefile/mage && \
    cd $GOPATH/src/github.com/magefile/mage && \
    go run bootstrap.go
    

# Copy in our workspace
COPY . .
RUN mage build:production 

FROM scratch 
WORKDIR /games
COPY --from=builder /go/src/github.com/arkticman/go-armored-warship/publish/battleship-linux-amd64 .
CMD ["/games/battleship-linux-amd64"] 