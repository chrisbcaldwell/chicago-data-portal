FROM golang:1.23
WORKDIR / app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /taxis

# Ports - still gaining understanding of what to do here
# EXPOSE <port#>

CMD ["/taxis"]