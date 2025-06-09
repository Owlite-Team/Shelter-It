FROM golang:alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN go build -o shelter-it ./cmd/app
EXPOSE 8080
CMD ["/app/shelter-it"]
