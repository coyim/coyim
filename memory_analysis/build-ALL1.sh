#!/bin/bash

set -xe

/root/installers/build-zlib.sh

/root/installers/build-libffi.sh

/root/installers/build-gettext.sh

/root/installers/build-ncurses.sh

/root/installers/build-libmount.sh

/root/installers/build-libpcre.sh

/root/installers/build-libpng.sh

/root/installers/build-libfreetype.sh

/root/installers/build-libharfbuzz.sh

/root/installers/rebuild1-libfreetype.sh

/root/installers/build-expat.sh

/root/installers/build-libfontconfig.sh

/root/installers/build-gtk-glib.sh

/root/installers/build-libatk.sh

/root/installers/build-dbus.sh

/root/installers/build-intltool.sh

/root/installers/build-x11.sh

/root/installers/build-atspi2.sh

/root/installers/build-libatk-bridge.sh

/root/installers/build-libpciaccess.sh

/root/installers/build-libdrm.sh

/root/installers/build-dri2proto.sh

/root/installers/build-dri3proto.sh

/root/installers/build-presentproto.sh

/root/installers/build-xshmfence.sh

/root/installers/build-damageproto.sh

/root/installers/build-xdamage.sh

/root/installers/build-mesa.sh

/root/installers/build-libepoxy.sh
