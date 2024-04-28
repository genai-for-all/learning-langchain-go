FROM golang:1.22.1-alpine as buildernext
WORKDIR /app
COPY main.go .
COPY go.mod .
COPY go.sum .
RUN go build

FROM scratch
WORKDIR /app
COPY *.md .
COPY --from=buildernext /app/chronicles-loader .
CMD ["./chronicles-loader"]
