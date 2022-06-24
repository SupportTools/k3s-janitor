FROM golang:latest AS build_base
ARG GIT_COMMIT
ARG GIT_BRANCH
WORKDIR /src
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN CGO_ENABLED=0 go build -ldflags "-X main.gitCommit=$GIT_COMMIT" "-X main.gitBranch=$GIT_BRANCH" -o main main.go

FROM ubuntu:latest
COPY --from=build_base /src/main /main
CMD ["/main"]