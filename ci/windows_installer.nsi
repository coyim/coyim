!include "MUI2.nsh"

!define NAME "CoyIM"

Name "${NAME}"
OutFile "${NAME} Installer.exe"
Unicode True

!define MUI_ICON "ci/win_icon_256x256.ico"
!define MUI_UNICON "ci/win_icon_256x256.ico"

Caption "CoyIM (${VERSION}) Installer"
BrandingText " "

InstallDir "$ProgramFiles\${Name}"

!define MUI_WELCOMEPAGE_TITLE "Welcome to the CoyIM Installer"
!define MUI_WELCOMEPAGE_TEXT "This installer will guide you through the installation of CoyIM.$\r$\n$\r$\n$\r$\n$\r$\n$_CLICK"

!define MUI_LICENSEPAGE_TEXT_BOTTOM "If you accept the terms of the agreement, click I Agree to continue."

!define MUI_FINISHPAGE_NOREBOOTSUPPORT

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "LICENSE"
!define MUI_COMPONENTSPAGE_NODESC
!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

!insertmacro MUI_LANGUAGE "English"

Section "CoyIM"
  SetOutPath "$INSTDIR"

  SectionIn 1 RO

  File /oname=CoyIM.exe win_installer\coyim_windows_amd64.exe
  File win_installer\toast.exe
  File win_installer\*.dll
  File /r win_installer\lib
  File /r win_installer\share
  File /oname=CoyIM.ico ci/win_icon_256x256.ico

  WriteUninstaller "$INSTDIR\Uninstall.exe"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${NAME}"   "DisplayName" "${NAME}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${NAME}"   "UninstallString" "$INSTDIR\Uninstall.exe"
SectionEnd

Section "Start Menu shortcut"
  CreateShortCut "$SMPROGRAMS\${NAME}.lnk" "$INSTDIR\CoyIM.exe" "" "$INSTDIR\CoyIM.ico"
SectionEnd

Section "Desktop shortcut"
  CreateShortCut "$DESKTOP\${NAME}.lnk" "$INSTDIR\CoyIM.exe" "" "$INSTDIR\CoyIM.ico"
SectionEnd

Section "Uninstall"
  Delete "$INSTDIR\CoyIM.exe"
  Delete "$INSTDIR\toast.exe"
  Delete "$INSTDIR\*.dll"
  RMDir /r "$INSTDIR\lib"
  RMDir /r "$INSTDIR\share"

  Delete "$SMPROGRAMS\${NAME}.lnk"
  Delete "$DESKTOP\${NAME}.lnk"
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${NAME}"
  Delete "$INSTDIR\Uninstall.exe"
  RMDir "$INSTDIR"
SectionEnd
