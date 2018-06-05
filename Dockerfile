FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/theswope/bank

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install github.com/theswope/bank

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/bank