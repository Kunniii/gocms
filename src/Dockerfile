FROM golang:1.21
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod vendor
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/main ./
ENTRYPOINT ["/out/main"]
