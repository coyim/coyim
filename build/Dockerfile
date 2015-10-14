FROM ubuntu:vivid

RUN apt-get update && apt-get upgrade -y -o Dpkg::Options::="--force-confold"

RUN DEBIAN_FRONTEND=noninteractive apt-get -y install \
  faketime build-essential curl gccgo git

RUN DEBIAN_FRONTEND=noninteractive apt-get -y install \
  libgtk2.0-dev libgtk-3-dev #gtk2.0 gtk+3.0

RUN update-alternatives --set go /usr/bin/go-5

ADD  setup-reproducible.sh /root/setup-reproducible
ADD  build-golang.sh /root/build-golang

# Download Go 1.5.1
# SHA1: 46eecd290d8803887dec718c691cc243f2175fe0
# https://golang.org/dl/
RUN mkdir -p /pkgs && cd /pkgs &&\
    curl https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz -O &&\
    echo "46eecd290d8803887dec718c691cc243f2175fe0 /pkgs/go1.5.1.linux-amd64.tar.gz" | sha1sum -c -

RUN /root/build-golang

# Clean up APT when done.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ADD  build.sh /root/build
ADD  test.sh /root/test

VOLUME /src

# This should be provided by the build script
ENV REFERENCE_DATETIME ""
ENV GOPATH /go
ENV GO_PKG "github.com/twstrike/coyim"

CMD bash -c /root/test && /root/build 
