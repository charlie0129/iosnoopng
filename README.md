# iosnoopng

Like iosnoop but with a Web GUI, intended to be run as a daemon.

You can view the total write and read bytes of each process.

![Overview](img/overview.png)

And the read and write bytes of each process to each file.

![Process Details](img/process.png)

## Build

Make sure you have Node.js and Go installed.

```bash
npm i
make
```

## Prerequisites

1. Disable SIP (System Integrity Protection) by booting into recovery mode and running `csrutil disable`.

## Usage (macOS)

> [!WARNING]
> Make sure your Mac **has NOT slept since last boot**, otherwise you will certainly have a system freeze. If you are not sure, reboot your Mac now. Otherwise, save your work and be prepared for a system freeze if it happens.
> 
> This tool uses dtrace internally. Due to an issue the darwin kernel, if your Mac has slept since the last boot, running dtrace may cause a system freeze. This is an issue with the darwin kernel and I have no control over it. See [this](https://forums.developer.apple.com/forums/thread/735939) .
> 

```bash
sudo ./iosnoopng
```

By default, the server will listen on `http://localhost:8092`. You can use launchd or other tools to run this as a daemon.

Command line options:
- `-a`: Listen address (default "127.0.0.1:8092")
- `-l`: Log level (default "info")
- `-o`: Save raw dtrace output to a file (default "")
