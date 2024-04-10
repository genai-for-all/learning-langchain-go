FROM golang:1.22.1-alpine as builder
WORKDIR /app
COPY main.go .
COPY go.mod .
COPY go.sum .
RUN go build

FROM scratch
WORKDIR /app
COPY public ./public
COPY --from=builder /app/genai-app .
CMD ["./genai-app"]
