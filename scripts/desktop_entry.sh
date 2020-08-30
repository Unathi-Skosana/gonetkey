#!/usr/bin/env bash

if [[ -z "${GOPATH}" ]]
then
  echo -ne "GOPATH not set"
  exit 0
fi

cat > "$HOME/.local/share/applications/Gonetkey.desktop" << EOF
[Desktop Entry]
Version=1.0
Terminal=false
Exec=$GOPATH/bin/gonetkey
Name=Gonetkey
Icon=$GOPATH/src/github.com/unathi-skosana/gonetkey/assets/icon.png
Type=Application
EOF
