go-lvm
=======================================================================

[![GoDoc](https://godoc.org/github.com/nak3/go-lvm?status.svg)](https://godoc.org/github.com/nak3/go-lvm)
[![Build Status](https://travis-ci.org/nak3/go-lvm.svg?branch=master)](https://travis-ci.org/nak3/go-lvm)
[![Go Report Card](https://goreportcard.com/badge/github.com/nak3/go-lvm)](https://goreportcard.com/report/github.com/nak3/go-lvm)

### Overview

go-lvm is a go library to call liblvm API based on python-lvm developed in [LVM2](https://sourceware.org/lvm2/).

### Usage

Please refer to [go-doc](https://godoc.org/github.com/nak3/go-lvm#example-LvObject--Createremove).

### Test run

Let's create a available volume group and create and delete a LV.

#### step-1. set up a free VG
~~~
sudo dd if=/dev/zero of=disk.img bs=1G count=1
export LOOP=`sudo losetup -f`
sudo losetup $LOOP disk.img
sudo vgcreate vg-targetd $LOOP
~~~

#### step-2. Run an example script
~~~
go run cmd/example.go
~~~
