FROM golang:1.16

# Depending on the golang version GO111MODULE can be removed as env variable
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on

# Set the Current Working Directory inside the container
WORKDIR /app/go-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/go-app .

RUN chmod a+x ./out/go-app

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["./out/go-app"]