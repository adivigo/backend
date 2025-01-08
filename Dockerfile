FROM golang

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o backend main.go

EXPOSE 8888

ENTRYPOINT ["/app/backend"]