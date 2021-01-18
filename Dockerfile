# syntax=docker/dockerfile:experimental

# Build mkfs as an static
FROM gcc:10.2.0 as mke2fs
RUN wget https://mirrors.edge.kernel.org/pub/linux/kernel/people/tytso/e2fsprogs/v1.45.6/e2fsprogs-1.45.6.tar.gz; tar -xvzf ./e2fsprogs-1.45.6.tar.gz
WORKDIR e2fsprogs-1.45.6/
RUN ./configure --enable-static=yes CFLAGS='-g -O2 -static'
RUN make
WORKDIR misc
RUN make mke2fs.static

# build mkswap as static
FROM gcc:10.2.0 as mkswap
RUN git clone git://git.kernel.org/pub/scm/utils/util-linux/util-linux.git
WORKDIR util-linux/
RUN apt-get update; apt-get install -y gettext bison autopoint
RUN ./autogen.sh; ./configure LDFLAGS="-static"
RUN make LDFLAGS="-all-static" mkswap

# Build rootio
FROM golang:1.15-alpine as rootio
RUN apk add --no-cache git ca-certificates gcc linux-headers musl-dev
COPY . /go/src/github.com/thebsdbox/rootio/
WORKDIR /go/src/github.com/thebsdbox/rootio
ENV GO111MODULE=on
RUN --mount=type=cache,sharing=locked,id=gomod,target=/go/pkg/mod/cache \
    --mount=type=cache,sharing=locked,id=goroot,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o rootio

# Build final image
FROM scratch
COPY --from=mke2fs /e2fsprogs-1.45.6/misc/mke2fs.static /sbin/mke2fs
COPY --from=rootio /go/src/github.com/thebsdbox/rootio/rootio .
ENTRYPOINT ["/rootio"]