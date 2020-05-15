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
Exec=$GOPATH/src/github.com/unathi-skosana/gonetkey/cmd/gonetkey/gonetkey
Name=Gonetkey
Icon=$GOPATH/src/github.com/unathi-skosana/gonetkey/build/gonetkey.png
Type=Application
EOF

# TODO: MacOS
