FROM golang:latest AS build_base
ARG GIT_COMMIT
WORKDIR /src
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN CGO_ENABLED=0 go build -ldflags "-X main.GitCommit=$GIT_COMMIT" -o main main.go

FROM ubuntu:latest
COPY --from=build_base /src/main /main
CMD ["/main"]