FROM golang:latest AS build_base
ARG GIT_COMMIT
ARG GIT_BRANCH
WORKDIR /src
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN CGO_ENABLED=0 go build -ldflags "-X main.gitCommit=$GIT_COMMIT" -o main main.go

FROM rancher/k3s:latest AS k3s_base

FROM alpine:latest
COPY --from=build_base /src/main /main
COPY --from=k3s_base /bin/k3s /bin/k3s
CMD ["/main"]