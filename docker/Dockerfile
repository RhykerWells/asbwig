FROM docker.io/golang:1.23.3

WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ .

WORKDIR /app/cmd/ASBWIG
RUN CGO_ENABLED=0 GOOS=linux go build

CMD ["/app/cmd/ASBWIG/ASBWIG"]