FROM golang:1.20.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

#RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

COPY main.go main.go

CMD ["/docker-gs-ping"]