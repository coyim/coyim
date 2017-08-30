FROM ubuntu:zesty

RUN apt-get update && apt-get upgrade -y -o Dpkg::Options::="--force-confold"

RUN DEBIAN_FRONTEND=noninteractive apt-get -y install \
  build-essential curl gccgo git clang llvm pkg-config libxml-parser-perl flex bison

RUN update-alternatives --install "/usr/bin/go" "go" "/usr/bin/go-6" 0
RUN update-alternatives --set go /usr/bin/go-6

# Clean up APT when done.
#RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN mkdir -p /usr/repackaged
RUN mkdir -p /root/installers

ENV CC=clang
ENV CXX=clang++
ENV PKG_CONFIG_PATH=/usr/repackaged/lib/pkgconfig:/usr/repackaged/share/pkgconfig

ADD download-ALL1.sh download-atspi2.sh download-damageproto.sh download-dbus.sh download-dri2proto.sh download-dri3proto.sh download-expat.sh download-gettext.sh download-golang-1.9.sh download-gtk-glib.sh download-intltool.sh download-libatk-bridge.sh download-libatk.sh download-libdrm.sh download-libepoxy.sh download-libffi.sh download-libfontconfig.sh download-libfreetype.sh download-libharfbuzz.sh download-libmount.sh download-libpciaccess.sh download-libpcre.sh download-libpng.sh download-mesa.sh download-ncurses.sh download-presentproto.sh download-x11.sh download-xdamage.sh download-xshmfence.sh download-zlib.sh /root/installers/
ADD build-ALL1.sh build-atspi2.sh build-coyim.sh build-damageproto.sh build-dbus.sh build-dri2proto.sh build-dri3proto.sh build-expat.sh build-gettext.sh build-golang.sh build-gtk-glib.sh build-intltool.sh build-libatk-bridge.sh build-libatk.sh build-libdrm.sh build-libepoxy.sh build-libffi.sh build-libfontconfig.sh build-libfreetype.sh build-libharfbuzz.sh build-libmount.sh build-libpciaccess.sh build-libpcre.sh build-libpng.sh build-mesa.sh build-ncurses.sh build-presentproto.sh build-x11.sh build-xdamage.sh build-xshmfence.sh build-zlib.sh /root/installers/
ADD rebuild1-libfreetype.sh /root/installers/

RUN /root/installers/download-ALL1.sh
RUN /root/installers/build-ALL1.sh



ADD download-python.sh /root/installers/
RUN /root/installers/download-python.sh

ADD build-python.sh /root/installers/
RUN /root/installers/build-python.sh

ADD download-gobject-introspection.sh /root/installers/
RUN /root/installers/download-gobject-introspection.sh

ADD build-gobject-introspection.sh /root/installers/
RUN /root/installers/build-gobject-introspection.sh

RUN /root/installers/build-libatk.sh

ADD download-gtk-pixman.sh /root/installers/
RUN /root/installers/download-gtk-pixman.sh

ADD build-gtk-pixman.sh /root/installers/
RUN /root/installers/build-gtk-pixman.sh

ADD download-gtk-cairo.sh /root/installers/
RUN /root/installers/download-gtk-cairo.sh

ADD build-gtk-cairo.sh /root/installers/
RUN /root/installers/build-gtk-cairo.sh

ADD download-gtk-pango.sh /root/installers/
RUN /root/installers/download-gtk-pango.sh

ADD rebuild1-libharfbuzz.sh /root/installers/
RUN /root/installers/rebuild1-libharfbuzz.sh

ADD download-librsvg.sh /root/installers/
RUN /root/installers/download-librsvg.sh

ADD download-libxml.sh /root/installers/
RUN /root/installers/download-libxml.sh

ADD build-libxml.sh /root/installers/
RUN /root/installers/build-libxml.sh

ADD download-libcroco.sh /root/installers/
RUN /root/installers/download-libcroco.sh

ADD build-libcroco.sh /root/installers/
RUN /root/installers/build-libcroco.sh

ADD download-libtiff.sh /root/installers/
RUN /root/installers/download-libtiff.sh

ADD build-libtiff.sh /root/installers/
RUN /root/installers/build-libtiff.sh

ADD download-libjpeg.sh /root/installers/
RUN /root/installers/download-libjpeg.sh

ADD build-libjpeg.sh /root/installers/
RUN /root/installers/build-libjpeg.sh

ADD build-gtk-pango.sh /root/installers/
RUN /root/installers/build-gtk-pango.sh

ADD download-libgdk-pixbuf.sh /root/installers/
RUN /root/installers/download-libgdk-pixbuf.sh

ADD build-libgdk-pixbuf.sh /root/installers/
RUN /root/installers/build-libgdk-pixbuf.sh

ADD build-librsvg.sh /root/installers/
RUN /root/installers/build-librsvg.sh

ADD download-gtk-gtk.sh /root/installers/
RUN /root/installers/download-gtk-gtk.sh

ADD build-gtk-gtk.sh /root/installers/
RUN /root/installers/build-gtk-gtk.sh




# RUN ln -s /usr/bin/llvm-symbolizer-4.0 /usr/bin/llvm-symbolizer

# ENV MSAN_OPTIONS=halt_on_error=0,exitcode=0

# ADD rebuild-gtk-glib.sh /root/rebuild-gtk-glib
# RUN /root/rebuild-gtk-glib

# ADD rebuild-gtk-pixman.sh /root/rebuild-gtk-pixman
# RUN /root/rebuild-gtk-pixman

# ADD rebuild-gtk-cairo.sh /root/rebuild-gtk-cairo
# RUN /root/rebuild-gtk-cairo

# ADD rebuild-gtk-pango.sh /root/rebuild-gtk-pango
# RUN /root/rebuild-gtk-pango

# ADD rebuild-gtk-gtk.sh /root/rebuild-gtk-gtk
# RUN /root/rebuild-gtk-gtk



# ADD download-golang-1.9.sh /root/download-golang
# ADD build-golang.sh /root/build-golang
# RUN /root/download-golang && /root/build-golang


# ADD build-coyim.sh /root/build-coyim

# VOLUME /src
