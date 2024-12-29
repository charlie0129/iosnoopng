# iosnoopng

Monitor disk activity on your Mac. Like iosnoop but with a Web GUI, intended to be run as a daemon.

How is it different from the Disk tab of Activity Monitor? Activity Monitor can only show process that are currently running. However, iosnoopng can collect metrics from **transient processes** that are not shown in Activity Monitor (i.e. process that is not long-running, like your C compilers). It can also show the read and write bytes of each process to each file.

## Screenshots

You can view the total write and read bytes of each process.

<img width="998" alt="Total" src="https://github.com/user-attachments/assets/ff53ac2a-3cb0-4fef-9343-19878a12b0f3" />

And the read and write bytes of each process to each file.

<img width="998" alt="Process" src="https://github.com/user-attachments/assets/636eb391-af80-4105-ad0b-7d111269a86a" />

## Build

Make sure you have Node.js and Go installed.

```bash
npm i
make
```

## Prerequisites (macOS)

1. Disable SIP (System Integrity Protection) by booting into recovery mode and running `csrutil disable`.

## Usage (macOS)

> [!WARNING]
> **Make sure your Mac has NOT slept since last boot**, otherwise you will certainly have a system freeze. If you are not sure, reboot your Mac now. Otherwise, save your work and be prepared for a system freeze if it happens.
>
> This tool uses dtrace internally. Due to an issue the darwin kernel, if your Mac has slept since the last boot, running dtrace may cause a system freeze. This is an issue with the darwin kernel and I have no control over it. See [this](https://forums.developer.apple.com/forums/thread/735939) .
>

```bash
sudo ./iosnoopng -d
```

By default, the server will listen on `http://localhost:8092`. You can use launchd or other tools to run this as a daemon.

Command line options:
- `-a`: Listen address (default "127.0.0.1:8092")
- `-d`: Deduplicate launchd processes. Since launchd writes on behalf of other processes if the same file is written by another process it will be counted twice. This option removes the duplicated entries from launchd.
- `-l`: Log level (default "info")
- `-o`: Save raw dtrace output to a file (default "")
- `-r`: Replay previous log file before collecting new data
- `-s`: Only merge entries that are not updated than this number of seconds. (default 3600)
- `-t`: Merge R/W smaller than this number of bytes into a single entry. This is useful for processes that write to many small files to save memory and make the output more readable. (default 33554432)

> [!WARNING]
> Before quitting iosnoopng, **make sure your Mac have NOT slept**, due to the same reason mentioned above. If your Mac has not slept, you can safely quit iosnoopng. However, if your Mac has slept, quitting iosnoopng may cause a system freeze. Either let iosnoopng running or save your work and **force** reboot your Mac.

## Extra Considerations

- Since iosnoopng records paths of files, its memory usage may grow over time. However, it should not be a problem for most users because iosnoopng has small-files-merging enabled by default. So lots of small files will be merged into one entry to save memory.
- The full path is not reported on macOS (only part of it). You will see every file starting with `/` even if the file is not there. This is a limitation of dtrace on macOS. For example, if you see a file `/idindex/IdIndex_inputs_i`, it is not in the root directory. If you need to find its full path, you can use `find . -name IdIndex_inputs_i` to find it. In my case, it is `$HOME/Library/Caches/JetBrains/GoLand2024.2/index/idindex/IdIndex.storage_i`
