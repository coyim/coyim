FROM ubuntu:zesty

RUN apt-get update && apt-get upgrade -y -o Dpkg::Options::="--force-confold"

RUN DEBIAN_FRONTEND=noninteractive apt-get -y install \
  build-essential curl gccgo git pkg-config libxml-parser-perl flex bison cmake llvm clang

RUN update-alternatives --install "/usr/bin/go" "go" "/usr/bin/go-6" 0
RUN update-alternatives --set go /usr/bin/go-6

# Clean up APT when done.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN mkdir -p /usr/repackaged /root/installers

ENV CC=clang CXX=clang++ PKG_CONFIG_PATH="/usr/repackaged/lib/pkgconfig:/usr/repackaged/share/pkgconfig" STANDARD_CFLAGS="-I/usr/repackaged/include -fno-omit-frame-pointer -fPIE" PIC_CFLAGS="-I/usr/repackaged/include -fno-omit-frame-pointer -fPIC" MSAN_CFLAGS="-I/usr/repackaged/include -fno-omit-frame-pointer -fPIE -fsanitize=memory -fsanitize-memory-track-origins -fsanitize-recover=all" LDFLAGS="-L/usr/repackaged/lib" PATH="/usr/repackaged/bin:${PATH}" LD_LIBRARY_PATH="/usr/repackaged/lib" MSAN_OPTIONS=halt_on_error=0,exitcode=0 MSAN_PIC_CFLAGS="-I/usr/repackaged/include -fno-omit-frame-pointer -fPIE -fsanitize=memory -fsanitize-memory-track-origins -fsanitize-recover=all" MSAN_LDFLAGS="-L/usr/repackaged/lib -fno-omit-frame-pointer -fPIE -fsanitize=memory -fsanitize-memory-track-origins -fsanitize-recover=all"

ADD standard-download.sh standard-verified-download.sh standard-build.sh download-and-build.sh download-verified-and-build.sh clean-binaries.sh standard-rebuild.sh /root/installers/

RUN /root/installers/download-verified-and-build.sh zlib    https://zlib.net/zlib-1.2.11.tar.gz c3e5e9fdd5004dcb542feda5ee4f0ff0744628baf8ed2dd5d66f8ca1197cb1a1 &&\
 /root/installers/download-and-build.sh          libffi  ftp://sourceware.org/pub/libffi/libffi-3.2.1.tar.gz &&\
 /root/installers/download-and-build.sh          gettext https://ftp.gnu.org/pub/gnu/gettext/gettext-0.19.8.tar.xz
RUN /root/installers/download-and-build.sh       ncurses-6.0 https://ftp.gnu.org/pub/gnu/ncurses/ncurses-6.0.tar.gz "--with-termlib --with-shared --without-tests"
RUN /root/installers/download-and-build.sh          util-linux https://www.kernel.org/pub/linux/utils/util-linux/v2.29/util-linux-2.29.1.tar.xz "" "-I/usr/repackaged/include/ncurses" &&\
 /root/installers/download-and-build.sh          pcre ftp://ftp.csx.cam.ac.uk/pub/software/programming/pcre/pcre-8.39.tar.bz2 "--enable-utf --enable-unicode-properties" &&\
 /root/installers/download-and-build.sh          libpng ftp://ftp-osl.osuosl.org/pub/libpng/src/libpng16/libpng-1.6.32.tar.xz &&\
 /root/installers/download-and-build.sh          freetype http://download.savannah.gnu.org/releases/freetype/freetype-2.6.3.tar.bz2 "--without-harfbuzz" "" "" "true"
RUN /root/installers/download-verified-and-build.sh harfbuzz https://www.freedesktop.org/software/harfbuzz/release/harfbuzz-1.4.2.tar.bz2 8f234dcfab000fdec24d43674fffa2fdbdbd654eb176afbde30e8826339cb7b3 &&\
 /root/installers/standard-build.sh              freetype "" "" "" "true" &&\
 /root/installers/download-and-build.sh          expat https://downloads.sourceforge.net/project/expat/expat/2.2.4/expat-2.2.4.tar.bz2 &&\
 /root/installers/download-and-build.sh          fontconfig https://www.freedesktop.org/software/fontconfig/release/fontconfig-2.11.94.tar.bz2
RUN /root/installers/download-verified-and-build.sh glib http://ftp.gnome.org/pub/gnome/sources/glib/2.53/glib-2.53.6.tar.xz e01296a9119c09d2dccb37ad09f5eaa1e0c5570a473d8fed04fc759ace6fb6cc &&\
 /root/installers/download-verified-and-build.sh atk-2 http://ftp.gnome.org/pub/gnome/sources/atk/2.22/atk-2.22.0.tar.xz d349f5ca4974c9c76a4963e5b254720523b0c78672cbc0e1a3475dbd9b3d44b6 &&\
 /root/installers/download-and-build.sh          dbus https://dbus.freedesktop.org/releases/dbus/dbus-1.10.10.tar.gz &&\
 /root/installers/download-and-build.sh          intltool https://launchpad.net/intltool/trunk/0.51.0/+download/intltool-0.51.0.tar.gz
RUN /root/installers/download-and-build.sh          xextproto https://xorg.freedesktop.org/archive/individual/proto/xextproto-7.3.0.tar.bz2 &&\
 /root/installers/download-and-build.sh          xtrans https://xorg.freedesktop.org/archive/individual/lib/xtrans-1.3.5.tar.bz2 &&\
 /root/installers/download-and-build.sh          xcb-proto https://xorg.freedesktop.org/archive/individual/xcb/xcb-proto-1.12.tar.bz2 &&\
 /root/installers/download-and-build.sh          libpthread-stubs https://xorg.freedesktop.org/archive/individual/xcb/libpthread-stubs-0.4.tar.bz2
RUN /root/installers/download-and-build.sh          kbproto https://xorg.freedesktop.org/archive/individual/proto/kbproto-1.0.7.tar.bz2 &&\
 /root/installers/download-and-build.sh          inputproto https://xorg.freedesktop.org/archive/individual/proto/inputproto-2.3.2.tar.bz2 &&\
 /root/installers/download-and-build.sh          recordproto https://xorg.freedesktop.org/archive/individual/proto/recordproto-1.14.2.tar.bz2 &&\
 /root/installers/download-and-build.sh          xproto https://xorg.freedesktop.org/archive/individual/proto/xproto-7.0.31.tar.bz2
RUN /root/installers/download-and-build.sh          libXau https://xorg.freedesktop.org/archive/individual/lib/libXau-1.0.8.tar.bz2 &&\
 /root/installers/download-and-build.sh          libxcb https://xorg.freedesktop.org/archive/individual/xcb/libxcb-1.12.tar.bz2 &&\
 /root/installers/download-and-build.sh          libX11 https://xorg.freedesktop.org/archive/individual/lib/libX11-1.6.4.tar.bz2 &&\
 /root/installers/download-and-build.sh          libXext https://xorg.freedesktop.org/archive/individual/lib/libXext-1.3.3.tar.bz2
RUN /root/installers/download-and-build.sh          fixesproto https://xorg.freedesktop.org/archive/individual/proto/fixesproto-5.0.tar.bz2 &&\
 /root/installers/download-and-build.sh          libXfixes https://xorg.freedesktop.org/archive/individual/lib/libXfixes-5.0.3.tar.bz2 &&\
 /root/installers/download-and-build.sh          libXi https://www.x.org/archive/individual/lib/libXi-1.7.9.tar.bz2 &&\
 /root/installers/download-and-build.sh          libXtst https://xorg.freedesktop.org/archive/individual/lib/libXtst-1.2.3.tar.bz2
RUN /root/installers/download-and-build.sh          glproto https://xorg.freedesktop.org/archive/individual/proto/glproto-1.4.17.tar.bz2 &&\
 /root/installers/download-verified-and-build.sh at-spi2-core http://ftp.gnome.org/pub/gnome/sources/at-spi2-core/2.22/at-spi2-core-2.22.0.tar.xz 415ea3af21318308798e098be8b3a17b2f0cf2fe16cecde5ad840cf4e0f2c80a "--x-includes=/usr/repackaged/include --x-libraries=/usr/repackaged/lib" &&\
 /root/installers/download-verified-and-build.sh at-spi2-atk http://ftp.gnome.org/pub/gnome/sources/at-spi2-atk/2.22/at-spi2-atk-2.22.0.tar.xz e8bdedbeb873eb229eb08c88e11d07713ec25ae175251648ad1a9da6c21113c1 &&\
 /root/installers/download-and-build.sh          libpciaccess https://xorg.freedesktop.org/archive/individual/lib/libpciaccess-0.13.4.tar.bz2
RUN /root/installers/download-and-build.sh          libdrm https://dri.freedesktop.org/libdrm/libdrm-2.4.76.tar.bz2 &&\
 /root/installers/download-and-build.sh          dri2proto https://xorg.freedesktop.org/archive/individual/proto/dri2proto-2.8.tar.bz2 &&\
 /root/installers/download-and-build.sh          dri3proto https://xorg.freedesktop.org/archive/individual/proto/dri3proto-1.0.tar.bz2 &&\
 /root/installers/download-and-build.sh          presentproto https://xorg.freedesktop.org/archive/individual/proto/presentproto-1.1.tar.bz2
RUN /root/installers/download-and-build.sh          libxshmfence https://xorg.freedesktop.org/archive/individual/lib/libxshmfence-1.2.tar.bz2 &&\
 /root/installers/download-and-build.sh          damageproto https://xorg.freedesktop.org/archive/individual/proto/damageproto-1.2.1.tar.bz2 &&\
 /root/installers/download-and-build.sh          libXdamage https://xorg.freedesktop.org/archive/individual/lib/libXdamage-1.1.4.tar.bz2 &&\
 /root/installers/download-and-build.sh          mesa ftp://ftp.freedesktop.org/pub/mesa/mesa-17.0.3.tar.xz
RUN /root/installers/download-and-build.sh          libepoxy https://github.com/anholt/libepoxy/releases/download/v1.3.1/libepoxy-1.3.1.tar.bz2 &&\
 /root/installers/download-and-build.sh          Python https://www.python.org/ftp/python/2.7.13/Python-2.7.13.tar.xz "--with-system-expat" "-I/usr/repackaged/include/ncurses" "true" &&\
 /root/installers/download-verified-and-build.sh gobject-introspection https://ftp.gnome.org/pub/gnome/sources/gobject-introspection/1.52/gobject-introspection-1.52.0.tar.xz 9fc6d1ebce5ad98942cb21e2fe8dd67b722dcc01981840632a1b233f7d0e2c1e &&\
 /root/installers/standard-build.sh              at-spi2-core "--x-includes=/usr/repackaged/include --x-libraries=/usr/repackaged/lib"
RUN /root/installers/download-and-build.sh          pixman https://www.cairographics.org/releases/pixman-0.34.0.tar.gz &&\
 /root/installers/download-and-build.sh          cairo https://www.cairographics.org/releases/cairo-1.14.8.tar.xz &&\
 /root/installers/standard-build.sh              harfbuzz &&\
 /root/installers/download-and-build.sh          libxml2 ftp://xmlsoft.org/libxml2/libxml2-2.9.4.tar.gz "" "" "true"
RUN /root/installers/download-verified-and-build.sh libcroco http://ftp.gnome.org/pub/GNOME/sources/libcroco/0.6/libcroco-0.6.11.tar.xz 132b528a948586b0dfa05d7e9e059901bca5a3be675b6071a90a90b81ae5a056 &&\
 /root/installers/download-and-build.sh          tiff ftp://download.osgeo.org/libtiff/tiff-4.0.7.tar.gz &&\
 /root/installers/download-and-build.sh          jpeg http://www.ijg.org/files/jpegsrc.v8c.tar.gz &&\
 /root/installers/download-verified-and-build.sh pango http://ftp.gnome.org/pub/gnome/sources/pango/1.40/pango-1.40.7.tar.xz 517645c00c4554e82c0631e836659504d3fd3699c564c633fccfdfd37574e278
RUN /root/installers/download-verified-and-build.sh gdk-pixbuf https://ftp.gnome.org/pub/gnome/sources/gdk-pixbuf/2.36/gdk-pixbuf-2.36.5.tar.xz 7ace06170291a1f21771552768bace072ecdea9bd4a02f7658939b9a314c40fc &&\
 /root/installers/download-verified-and-build.sh librsvg https://download.gnome.org/sources/librsvg/2.40/librsvg-2.40.16.tar.xz d48bcf6b03fa98f07df10332fb49d8c010786ddca6ab34cbba217684f533ff2e &&\
 /root/installers/standard-build.sh              atk-2 &&\
 /root/installers/download-verified-and-build.sh gtk http://ftp.gnome.org/pub/gnome/sources/gtk+/3.22/gtk+-3.22.17.tar.xz a6c1fb8f229c626a3d9c0e1ce6ea138de7f64a5a6bc799d45fa286fe461c3437

RUN ln -s /usr/bin/llvm-symbolizer-4.0 /usr/bin/llvm-symbolizer

RUN /root/installers/standard-rebuild.sh  zlib
RUN /root/installers/standard-rebuild.sh  libffi
RUN /root/installers/standard-rebuild.sh  gettext
RUN /root/installers/standard-rebuild.sh  ncurses-6.0
RUN /root/installers/standard-rebuild.sh  util-linux "-I/usr/repackaged/include/ncurses"
RUN /root/installers/standard-rebuild.sh  pcre
RUN /root/installers/standard-rebuild.sh  libpng
RUN /root/installers/standard-rebuild.sh  harfbuzz
RUN /root/installers/standard-rebuild.sh  freetype
RUN /root/installers/standard-rebuild.sh  expat
RUN /root/installers/standard-rebuild.sh  fontconfig
RUN /root/installers/standard-rebuild.sh  glib
RUN /root/installers/standard-rebuild.sh  atk-2
RUN /root/installers/standard-rebuild.sh  dbus
RUN /root/installers/standard-rebuild.sh  intltool
RUN /root/installers/standard-rebuild.sh  xextproto
RUN /root/installers/standard-rebuild.sh  xtrans
RUN /root/installers/standard-rebuild.sh  xcb-proto
RUN /root/installers/standard-rebuild.sh  libpthread-stubs
RUN /root/installers/standard-rebuild.sh  kbproto
RUN /root/installers/standard-rebuild.sh  inputproto
RUN /root/installers/standard-rebuild.sh  recordproto
RUN /root/installers/standard-rebuild.sh  xproto
RUN /root/installers/standard-rebuild.sh  libXau
RUN /root/installers/standard-rebuild.sh  libxcb
RUN /root/installers/standard-rebuild.sh  libX11
RUN /root/installers/standard-rebuild.sh  libXext
RUN /root/installers/standard-rebuild.sh  fixesproto
RUN /root/installers/standard-rebuild.sh  libXfixes
RUN /root/installers/standard-rebuild.sh  libXi
RUN /root/installers/standard-rebuild.sh  libXtst
RUN /root/installers/standard-rebuild.sh  glproto
RUN /root/installers/standard-rebuild.sh  at-spi2-core
RUN /root/installers/standard-rebuild.sh  at-spi2-atk
RUN /root/installers/standard-rebuild.sh  libpciaccess
RUN /root/installers/standard-rebuild.sh  libdrm
RUN /root/installers/standard-rebuild.sh  dri2proto
RUN /root/installers/standard-rebuild.sh  dri3proto
RUN /root/installers/standard-rebuild.sh  presentproto
RUN /root/installers/standard-rebuild.sh  libxshmfence
RUN /root/installers/standard-rebuild.sh  damageproto
RUN /root/installers/standard-rebuild.sh  libXdamage
RUN /root/installers/standard-rebuild.sh  mesa
RUN /root/installers/standard-rebuild.sh  libepoxy
RUN /root/installers/standard-rebuild.sh  Python "-I/usr/repackaged/include/ncurses" "true"
RUN /root/installers/standard-rebuild.sh  gobject-introspection
RUN /root/installers/standard-rebuild.sh  pixman
RUN /root/installers/standard-rebuild.sh  cairo
RUN /root/installers/standard-rebuild.sh  libxml2 "" "true"
RUN /root/installers/standard-rebuild.sh  libcroco
RUN /root/installers/standard-rebuild.sh  tiff
RUN /root/installers/standard-rebuild.sh  jpeg
RUN /root/installers/standard-rebuild.sh  pango
RUN /root/installers/standard-rebuild.sh  gdk-pixbuf
RUN /root/installers/standard-rebuild.sh  librsvg
RUN /root/installers/standard-rebuild.sh  gtk

ADD download-golang-1.9.sh /root/download-golang
ADD build-golang.sh /root/build-golang
RUN /root/download-golang && /root/build-golang

ADD build-coyim.sh /root/build-coyim

VOLUME /src
