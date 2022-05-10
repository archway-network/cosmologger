# FROM golang:alpine3.15 AS development
FROM golang:alpine AS development
ARG arch=x86_64

# ENV CGO_ENABLED=0
WORKDIR /go/src/app/
COPY . /go/src/app/

RUN set -eux; \
    apk add --no-cache \
    git \
    openssh \
    ca-certificates \
    build-base \
    && mkdir -p /build/ 

# We need these configs to handle the private repos
# If you do not have the my_keys directory it will be ignored 
# and will raise an error if a private repo found
# get your keys copied in the `my_keys` directory: 
#   mkdir -p ./my_keys && mkdir -p ./my_keys/.ssh
#   cp ~/.ssh/id_rsa ./my_keys/.ssh && cp ~/.gitconfig ./my_keys
# COPY . /go/src/app/
# COPY main.go my_keys/.ssh/id_rsa* /root/.ssh/
# COPY main.go my_keys/.gitconfig* /root/
# RUN chmod 700 /root/* \
#     && echo -e "\nHost github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config \
#     && cd ~ \
#     && git config --global url.ssh://git@github.com/.insteadOf https://github.com/


# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.0.0-beta10/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.0.0-beta10/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 5b7abfdd307568f5339e2bea1523a6aa767cf57d6a8c72bc813476d790918e44 \
    && sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 2f44efa9c6c1cda138bd1f46d8d53c5ebfe1f4a53cf3457b01db86472c4917ac \
    # Copy the library you want to the final location that will be found by the linker flag `-lwasmvm_muslc`
    && cp /lib/libwasmvm_muslc.${arch}.a /lib/libwasmvm_muslc.a

# Archwayd binary
RUN git clone https://github.com/archway-network/archway /root/archway\
    && cd /root/archway \
    && LEDGER_ENABLED=false BUILD_TAGS=muslc LINK_STATICALLY=true make build \
    && echo "Ensuring binary is statically linked ..." \
    && (file /root/archway/build/archwayd | grep "statically linked") \
    && cp /root/archway/build/archwayd /build

# Let's keep it in a separate layer
RUN go get github.com/go-delve/delve/cmd/dlv \
    && go build -mod=readonly -buildvcs=false -o /build/app . \
    && cp conf.json /build/

ENV PATH=$PATH:/build

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

ENTRYPOINT ["./app"]