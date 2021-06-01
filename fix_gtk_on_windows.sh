#!/usr/bin/env bash

declare -a GDK_FILES=(
  "C:/msys64/mingw64/lib/pkgconfig/gdk-2.0.pc"
  "C:/msys64/mingw64/lib/pkgconfig/gdk-3.0.pc"
  "C:/msys64/mingw64/lib/pkgconfig/gdk-win32-2.0.pc"
  "C:/msys64/mingw64/lib/pkgconfig/gdk-win32-3.0.pc"
)

for file in "${GDK_FILES[@]}"; do
    if [ -f $file ]; then
        sed -i -e "s/-Wl,-luuid/-luuid/g" $file
        echo -e "\nLDFLAGS: -Wl" >> $file
    fi
done
