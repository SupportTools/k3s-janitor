FROM golang:latest AS build_base
ARG GIT_COMMIT
ARG GIT_BRANCH
WORKDIR /src
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN CGO_ENABLED=0 go build -ldflags "-X main.gitCommit=$GIT_COMMIT" -o main main.go

FROM rancher/k3s:latest
COPY --from=build_base /src/main /main
RUN echo export PATH='$PATH:/var/lib/rancher/k3s/data/current/bin' >> /etc/profile
CMD ["/main"]