if "%METHOD%"=="ci" SET MSYS_PATH=c:\msys64
if "%METHOD%"=="cross" SET MSYS_PATH=%APPVEYOR_BUILD_FOLDER%\msys%MSYS2_BITS%
mkdir %GOPATH%\src\github.com\coyim\
xcopy %APPVEYOR_BUILD_FOLDER%\* %GOPATH%\src\github.com\coyim\coyim /e /i /s /EXCLUDE:%MSYS_PATH% > nul
xcopy %APPVEYOR_BUILD_FOLDER%\.git\* %GOPATH%\src\github.com\coyim\coyim\.git /e /i /s /r /h > nul
dir %GOPATH%\src\github.com\coyim\coyim
SET PATH=%MSYS_PATH%\usr\bin;%PATH%
SET PATH=%MSYS_PATH%\mingw%MSYS2_BITS%\bin;%PATH%
if "%METHOD%"=="cross" appveyor DownloadFile http://kent.dl.sourceforge.net/project/msys2/Base/%MSYS2_ARCH%/msys2-base-%MSYS2_ARCH%-%MSYS2_BASEVER%.tar.xz
if "%METHOD%"=="cross" 7z x msys2-base-%MSYS2_ARCH%-%MSYS2_BASEVER%.tar.xz > nul
if "%METHOD%"=="cross" 7z x msys2-base-%MSYS2_ARCH%-%MSYS2_BASEVER%.tar > nul
%MSYS_PATH%\usr\bin\bash -lc "echo update-core starting..."
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Syuu" > nul
%MSYS_PATH%\usr\bin\bash -lc "echo install-deps starting..."
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy autoconf" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy automake" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy make" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy mingw-w64-%MSYS2_ARCH%-libiconv" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy mingw-w64-%MSYS2_ARCH%-gcc" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy mingw-w64-%MSYS2_ARCH%-gdb" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy mingw-w64-%MSYS2_ARCH%-make" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy zlib-devel" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy mingw-w64-%MSYS2_ARCH%-pango" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy mingw-w64-%MSYS2_ARCH%-gtk3" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy mingw-w64-%MSYS2_ARCH%-pkg-config" > nul
%MSYS_PATH%\usr\bin\bash -lc "yes|pacman --noconfirm -Sc" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman -Ql mingw-w64-%MSYS2_ARCH%-gdk-pixbuf2"
%MSYS_PATH%\usr\bin\bash -lc "echo ZERO"
%MSYS_PATH%\usr\bin\bash -lc "ls /mingw%MSYS2_BITS%/**/*.dll"
%MSYS_PATH%\usr\bin\bash -lc "echo ZERO.bin"
%MSYS_PATH%\usr\bin\bash -lc "ls /mingw%MSYS2_BITS%/bin"
%MSYS_PATH%\usr\bin\bash -lc "echo ZERO.lib"
%MSYS_PATH%\usr\bin\bash -lc "ls /mingw%MSYS2_BITS%/lib"
%MSYS_PATH%\usr\bin\bash -lc "echo ZERO.1"
%MSYS_PATH%\usr\bin\bash -lc "find / -name '*.dll'"
if "%METHOD%"=="cross" %MSYS_PATH%\autorebase.bat > nul
