FROM golang:latest AS build_base
ARG GIT_COMMIT
ARG GIT_BRANCH
WORKDIR /src
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN CGO_ENABLED=0 go build -ldflags "-X main.gitCommit=$GIT_COMMIT" -o main main.go

FROM ubuntu:latest
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    wget \
    software-properties-common
RUN VERSION="v1.24.2" && \
wget https://github.com/kubernetes-sigs/cri-tools/releases/download/$VERSION/crictl-$VERSION-linux-amd64.tar.gz && \
tar zxvf crictl-$VERSION-linux-amd64.tar.gz -C /usr/local/bin && \
rm -f crictl-$VERSION-linux-amd64.tar.gz
COPY --from=build_base /src/main /main
CMD ["/main"]