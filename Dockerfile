# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)

RUN go get -v github.com/ohall/amygo

RUN ls $GOPATH/src/github.com/ohall/amygo

# We'll need these to serve
RUN cp -r $GOPATH/src/github.com/ohall/amygo/templates $GOPATH/bin

# Run the outyet command by default when the container starts.
ENTRYPOINT $GOPATH/bin/amygo

# Document that the service listens on port 8000.
EXPOSE 8000

# Fire it up
# CMD amygo