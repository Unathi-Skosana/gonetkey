# Gonetkey

A golang implementation of inetkey (unix only) to access Stellenbosch
University's firewall.

## Getting started

### CLI

```sh
go install github.com/unathi-skosana/gonetkey/cmd/gonetkeycli
./gonetkeycli -h
```

### DBUS
```sh
go install github.com/unathi-skosana/gonetkey/cmd/gonetkeyd
./gonetkeyd -h
```

### GUI
```sh
go install github.com/unathi-skosana/gonetkey/cmd/gonetkey
./gonetkey
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
