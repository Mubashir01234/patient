FROM golang:1.21.0-bullseye
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go install github.com/cosmtrek/air@latest
COPY . /app
EXPOSE 8001
CMD ["air", "-c", ".air.toml"]
