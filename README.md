# Gonetkey

A golang implementation of inetkey (unix only) to access Stellenbosch
University's firewall.

## Getting started

### CLI

```sh
git clone https://github.com/unathi-skosana/gonetkey
cd cmd/gonetkeycli
go build
./gonetkeycli -h
```


### GUI
```sh
git clone https://github.com/unathi-skosana/gonetkey
cd cmd/gonetkey
go build
chmod +x ./build.sh
./build.sh
./gonetkey
```

### DBUS
```sh
git clone https://github.com/unathi-skosana/gonetkey
cd cmd/gonetkeyd
go build
./gonetkeyd -h
```

## Credits
* [getlantern/systray](https://github.com/getlantern/systray)
* [godbus/dbus](https://github.com/godbus/dbus)
* [mkideal/cli](https://github.com/mkideal/cli)
* [fyne.io/fyne](https://github.com/fyne-io/fyne)

## TODO

* [ ]  MacOS build
* [ ]  Proper bundling of assets
* [ ]  Makefile ?

