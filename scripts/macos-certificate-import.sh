#!/usr/bin/env bash
# Sourced by the signer and its synthetic native import regression test.
macos_certificate_rewrap() {
  # Use PKCS#12 PBE/MAC algorithms understood by SecKeychainItemImport (the API
  # which honors nonextractable key attributes), rather than OpenSSL defaults.
  openssl rand -hex 32 > "$3" || return $?
  openssl pkcs12 -export -in "$1" -out "$2" -passout "file:$3" \
    -keypbe PBE-SHA1-3DES -certpbe PBE-SHA1-3DES -macalg sha1
}

macos_certificate_import() {
  local password path
  IFS= read -r password < "$3" || return 1
  [[ $password =~ ^[a-f0-9]{64}$ ]] || return 1
  for path in "$1" "$2"; do
    case $path in *'"'*|*\\*|*$'\n'*|*$'\r'*) return 1 ;; esac
  done
  # The strong transfer password is never an argv value or exported variable.
  # Let SecKeychainItemImport detect PKCS#12. The explicit -f pkcs12 branch in
  # newer security versions calls SecPKCS12Import without forwarding -x.
  printf 'import "%s" -k "%s" -x -P "%s" -T /usr/bin/codesign\n' \
    "$1" "$2" "$password" | security -i
}

macos_certificate_import_diagnostic() {
  # Only exact, known Security constants leave the private command log. Their
  # corresponding OSStatus is mapped, never scraped from arbitrary private text.
  local line status=unknown cause=unclassified
  while IFS= read -r line || [[ -n $line ]]; do
    line=${line#security> }
    case $line in
      'security: SecKeychainItemImport: The user name or passphrase you entered is not correct.') status=-25293; cause=authentication-failed ;;
      'security: SecKeychainItemImport: MAC verification failed during PKCS12 import (wrong password?)') status=-25264; cause=pkcs12-mac-verification-failed ;;
      'security: SecKeychainItemImport: Unable to decode the provided data.') status=-26275; cause=decode-failed ;;
      'security: SecKeychainItemImport: Unknown format in import.') status=-25257; cause=unknown-format ;;
      'security: SecKeychainItemImport: User interaction is not allowed.') status=-25308; cause=interaction-not-allowed ;;
      'security: SecKeychainItemImport: One or more parameters passed to a function were not valid.') status=-50; cause=invalid-parameters ;;
    esac
  done < "$1"
  printf 'sign-macos: keychain-import security_status=%s cause=%s (allowlisted classification; raw output withheld)\n' "$status" "$cause" >&2
}
