FROM golang:1.14 AS builder

# enable Go modules support
#ENV GO111MODULE=on

WORKDIR /build

# manage dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy src code from the host and compile it
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app main.go

FROM alpine:3.9
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/app /bin
CMD ["/bin/app"]