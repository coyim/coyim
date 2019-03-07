mkdir %GOPATH%\src\github.com\coyim\coyim\release\
mkdir %GOPATH%\src\github.com\coyim\coyim\release\share
echo ONE
%MSYS_PATH%\usr\bin\bash -lc "ls /mingw%MSYS2_BITS%/**/*.dll"
echo TWO
%MSYS_PATH%\usr\bin\bash -lc "ls /mingw%MSYS2_BITS%/lib/gdk-pixbuf*/**/*.dll"
echo THREE
%MSYS_PATH%\usr\bin\bash -lc "cp -v /mingw%MSYS2_BITS%/**/*.dll /c/gopath/src/github.com/coyim/coyim/release/"
%MSYS_PATH%\usr\bin\bash -lc "cp -r /mingw%MSYS2_BITS%/share/icons /c/gopath/src/github.com/coyim/coyim/release/share/"
%MSYS_PATH%\usr\bin\bash -lc "cp -r /mingw%MSYS2_BITS%/share/glib-2.0 /c/gopath/src/github.com/coyim/coyim/release/share/"
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy sed" > nul
%MSYS_PATH%\usr\bin\bash -lc "cd /c/gopath/src/github.com/coyim/coyim && ci/release"
%MSYS_PATH%\usr\bin\bash -lc "cd /c/gopath/src/github.com/coyim/coyim/release && ls -alF ."
%MSYS_PATH%\usr\bin\bash -lc "cd /c/gopath/src/github.com/coyim/coyim/release && 7z a -tzip coyim.zip *"
%MSYS_PATH%\usr\bin\bash -lc "cd /c/gopath/src/github.com/coyim/coyim/release && rm *.dll"
%MSYS_PATH%\usr\bin\bash -lc "cd /c/gopath/src/github.com/coyim/coyim/release && rm *.exe"
xcopy %GOPATH%\src\github.com\coyim\coyim\release\coyim.zip %APPVEYOR_BUILD_FOLDER%\ /e /i > nul
