Unicode true
!include "MUI2.nsh"
!include "x64.nsh"
!ifndef VERSION
  !error "VERSION is required"
!endif
!ifndef PAYLOAD
  !error "PAYLOAD is required"
!endif
!ifndef OUTPUT
  !error "OUTPUT is required"
!endif

Name "takedown"
OutFile "${OUTPUT}"
InstallDir "$LOCALAPPDATA\Programs\takedown"
RequestExecutionLevel user
SetCompressor /SOLID lzma
!define MUI_ICON "..\..\assets\icon.ico"
!define MUI_UNICON "..\..\assets\icon.ico"
VIProductVersion "${VERSION}.0"
VIAddVersionKey "ProductName" "takedown"
VIAddVersionKey "FileDescription" "takedown per-user installer"
VIAddVersionKey "FileVersion" "${VERSION}"

# NSIS invokes this on the generated uninstaller BEFORE embedding it in setup.
!ifdef SIGN_SCRIPT
  !uninstfinalize 'pwsh.exe -NoLogo -NoProfile -NonInteractive -File "${SIGN_SCRIPT}" -Path "%1"' = 0
!endif

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_LANGUAGE "English"

Function .onInit
  ${IfNot} ${RunningX64}
    MessageBox MB_ICONSTOP "This package requires 64-bit Windows."
    Abort
  ${EndIf}
  SetShellVarContext current
FunctionEnd

Section "Install"
  SetOutPath "$INSTDIR"
  File "${PAYLOAD}\*"
  WriteUninstaller "$INSTDIR\uninstall.exe"
  CreateDirectory "$SMPROGRAMS\takedown"
  CreateShortcut "$SMPROGRAMS\takedown\takedown.lnk" "$INSTDIR\takedown.exe" "gui"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\takedown" "DisplayName" "takedown"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\takedown" "DisplayVersion" "${VERSION}"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\takedown" "InstallLocation" "$INSTDIR"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\takedown" "UninstallString" '$\"$INSTDIR\uninstall.exe$\"'
  WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\takedown" "NoModify" 1
  WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\takedown" "NoRepair" 1
SectionEnd

Section "Uninstall"
  SetShellVarContext current
  # At most two user-requested retries after the first failure; silent runs fail
  # immediately. Never schedule a reboot or discard the recovery entry on a lock.
  StrCpy $0 2
  uninstall_try_payload:
  ClearErrors
  Delete "$INSTDIR\takedown.exe"
  IfErrors uninstall_payload_locked uninstall_remove_docs
  uninstall_payload_locked:
  IfSilent uninstall_failed
  IntCmp $0 0 uninstall_failed uninstall_failed uninstall_retry_prompt
  uninstall_retry_prompt:
  IntOp $0 $0 - 1
  MessageBox MB_RETRYCANCEL|MB_ICONEXCLAMATION "Close takedown and retry. Cancel preserves the installer and uninstall registration so removal can be retried later." IDRETRY uninstall_try_payload
  Goto uninstall_failed
  uninstall_remove_docs:
  ClearErrors
  Delete "$INSTDIR\README.md"
  Delete "$INSTDIR\PRIVACY.md"
  Delete "$INSTDIR\PRIVACY.html"
  Delete "$INSTDIR\BUILD-INFO.txt"
  Delete "$INSTDIR\LICENSE"
  Delete "$INSTDIR\NOTICE"
  Delete "$INSTDIR\THIRD_PARTY_NOTICES.md"
  Delete "$INSTDIR\LICENSE.md"
  Delete "$INSTDIR\LICENSE.txt"
  Delete "$INSTDIR\COPYING"
  IfErrors uninstall_failed
  ClearErrors
  Delete "$SMPROGRAMS\takedown\takedown.lnk"
  IfErrors uninstall_failed
  ClearErrors
  Delete "$INSTDIR\uninstall.exe"
  IfErrors uninstall_failed
  RMDir "$INSTDIR"
  RMDir "$SMPROGRAMS\takedown"
  DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\takedown"
  SetErrorLevel 0
  Goto uninstall_done
  uninstall_failed:
  DetailPrint "Removal incomplete. Close takedown and rerun the registered uninstaller. Recovery metadata was retained."
  SetErrorLevel 10
  Quit
  uninstall_done:
  # Deliberately retain archives, app state, and Windows Credential Manager tokens.
SectionEnd
