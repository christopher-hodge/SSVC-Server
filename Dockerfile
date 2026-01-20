FROM golang:1.22-alpine AS build

WORKDIR /backend
COPY go.mod go.sum ./
ENV GOTOOLCHAIN=auto
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./server

FROM alpine:latest
WORKDIR /backend
COPY --from=build /backend/server .

EXPOSE 8080
CMD ["./server"]