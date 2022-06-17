FROM golang:latest AS build_base
WORKDIR /src
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN export GIT_COMMIT=$(git rev-list -1 HEAD) && \
CGO_ENABLED=0 go build -ldflags "-X main.GitCommit=$GIT_COMMIT" -o main main.go

FROM scratch
COPY --from=build_base /src/main /main
CMD ["/main"]