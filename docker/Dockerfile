FROM docker.io/golang:1.24.1

WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ .

WORKDIR /app/cmd/summit
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X github.com/RhykerWells/Summit/common.VERSION=$(git describe --tags --abbrev=0)"

EXPOSE 80

CMD ["/app/cmd/summit/summit"]