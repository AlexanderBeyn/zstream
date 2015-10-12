# Introduction

zstream is a utility to convert video data from Zmodo cameras (tested on model ZM-SH75D001-WA) to a format usable by ffmpeg. The primary goal of this utility is to integrate the above camera into ZoneMinder.

# Usage
```
Usage of zstream:
  -c string
    	Camera address (address:port)
  -l string
    	Listen address (default ":8888")
```

* `-c string` _(Required)_
  * Address of the camera, including the address and the port separated by a colon. The port is normally 8000.
* `-l string` _(Optional)_
  * Address and TCP port on which to listen. Defaults to port 8888 on all addresses.

# ZoneMinder
Setup a new source with type `Ffmpeg` and source path `tcp://localhost:8888` (assuming the default port and zstream running on the same host as ZoneMinder).
