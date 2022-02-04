FROM golang:alpine3.15 AS development

ENV CGO_ENABLED=0
COPY . /go/src/app/
WORKDIR /go/src/app/
ENV GOPATH=/go/

RUN apk add --no-cache \
    # git \
    # curl \
    # gcc \
    # zip \
    && mkdir /build/ \
    # && cp -r docs /build \
    && go get github.com/go-delve/delve/cmd/dlv


# Let's keep it in a separate layer
RUN go build -mod=readonly -o /build/app .
# ENTRYPOINT [ "dlv", "debug", "--headless", "--log", "--listen=:2345", "--api-version=2"]

ENTRYPOINT ["tail", "-f", "/dev/null"]

#----------------------------#

FROM development AS test

WORKDIR /go/src/app/

ENV EXEC_PATH=/go/src/app/

ENTRYPOINT ["go", "test", "-v", "./..."]

#----------------------------#

FROM alpine:latest AS production

WORKDIR /app/
COPY --from=development /build .
RUN apk --no-cache add \
    curl 

# COPY ui/node_modules/react/umd ui/node_modules/react/umd
# COPY ui/node_modules/react-dom/umd ui/node_modules/react-dom/umd
# COPY ui/index.html \
#     ui/favicon.ico \
#     ui/
# COPY ui/dist ui/dist

ENTRYPOINT ["./app"]

#  go get -v  golang.org/x/tools/cmd/godoc
# godoc -http=:8080 -goroot /go/