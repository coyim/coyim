FROM debian:buster

RUN apt-get update && apt-get upgrade -y -o Dpkg::Options::="--force-confold" --no-install-recommends

RUN DEBIAN_FRONTEND=noninteractive apt-get -y install --no-install-recommends \
  faketime build-essential curl golang-go gccgo git

# libgtk2.0-dev libgtk-3-dev gtk2.0 gtk+3.0
RUN DEBIAN_FRONTEND=noninteractive apt-get -y install --no-install-recommends \
  libgtk-3-dev ruby

# Clean up APT when done.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ADD setup-reproducible.sh /root/setup-reproducible

ADD build.sh /root/build

VOLUME /src
