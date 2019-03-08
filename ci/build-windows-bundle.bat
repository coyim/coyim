mkdir %GOPATH%\src\github.com\coyim\coyim\release\
mkdir %GOPATH%\src\github.com\coyim\coyim\release\share
%MSYS_PATH%\usr\bin\bash -lc "echo 'FINDING ALL'"
%MSYS_PATH%\usr\bin\bash -lc "find /mingw%MSYS2_BITS% -name '*.dll'"
%MSYS_PATH%\usr\bin\bash -lc "echo 'Copying bin/*.dll'"
%MSYS_PATH%\usr\bin\bash -lc "cp -v /mingw%MSYS2_BITS%/bin/*.dll /c/gopath/src/github.com/coyim/coyim/release/"
%MSYS_PATH%\usr\bin\bash -lc "echo 'Copying lib/*.dll'"
%MSYS_PATH%\usr\bin\bash -lc "cp -v /mingw%MSYS2_BITS%/lib/*.dll /c/gopath/src/github.com/coyim/coyim/release/"
%MSYS_PATH%\usr\bin\bash -lc "echo 'Copying lib/**/*.dll'"
%MSYS_PATH%\usr\bin\bash -lc "cp -v /mingw%MSYS2_BITS%/lib/**/*.dll /c/gopath/src/github.com/coyim/coyim/release/"
%MSYS_PATH%\usr\bin\bash -lc "echo 'Copying lib/gdk-pixbuf-2.0/*/loaders/*.dll'"
%MSYS_PATH%\usr\bin\bash -lc "cp -v /mingw%MSYS2_BITS%/lib/gdk-pixbuf-2.0/*/loaders/*.dll /c/gopath/src/github.com/coyim/coyim/release/"
%MSYS_PATH%\usr\bin\bash -lc "cp -r /mingw%MSYS2_BITS%/share/icons /c/gopath/src/github.com/coyim/coyim/release/share/"
%MSYS_PATH%\usr\bin\bash -lc "cp -r /mingw%MSYS2_BITS%/share/glib-2.0 /c/gopath/src/github.com/coyim/coyim/release/share/"
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy sed" > nul
%MSYS_PATH%\usr\bin\bash -lc "cd /c/gopath/src/github.com/coyim/coyim && ci/release"
%MSYS_PATH%\usr\bin\bash -lc "cd /c/gopath/src/github.com/coyim/coyim/release && ls -alF ."
%MSYS_PATH%\usr\bin\bash -lc "cd /c/gopath/src/github.com/coyim/coyim/release && 7z a -tzip coyim.zip *"
%MSYS_PATH%\usr\bin\bash -lc "cd /c/gopath/src/github.com/coyim/coyim/release && rm *.dll"
%MSYS_PATH%\usr\bin\bash -lc "cd /c/gopath/src/github.com/coyim/coyim/release && rm *.exe"
xcopy %GOPATH%\src\github.com\coyim\coyim\release\coyim.zip %APPVEYOR_BUILD_FOLDER%\ /e /i > nul
