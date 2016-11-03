FROM ubuntu:vivid

ADD apt_preferences /etc/apt/preferences

RUN apt-get update && apt-get upgrade -y -o Dpkg::Options::="--force-confold"

RUN DEBIAN_FRONTEND=noninteractive apt-get -y install \
  faketime build-essential curl gccgo git

# libgtk2.0-dev libgtk-3-dev gtk2.0 gtk+3.0
RUN DEBIAN_FRONTEND=noninteractive apt-get -y install \
  libgtk-3-dev ruby

RUN update-alternatives --set go /usr/bin/go-5

# Clean up APT when done.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ADD setup-reproducible.sh /root/setup-reproducible
ADD download-golang-1.7.sh /root/download-golang
ADD build-golang.sh /root/build-golang

# Build golang as part of the image
RUN /root/download-golang && /root/build-golang

ADD build.sh /root/build

VOLUME /src
