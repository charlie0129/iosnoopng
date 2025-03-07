# iosnoopng

Monitor disk activity on your Mac. Like iosnoop but with a Web GUI, intended to be run as a daemon.

How is it different from the Disk tab of Activity Monitor? Activity Monitor can only show process that are currently running. However, iosnoopng can collect metrics from **transient processes** that are not shown in Activity Monitor (i.e. process that is not long-running, like your C compilers). It can also show the read and write bytes of each process to each file. It also has a Prometheus metrics endpoint so you can scrape the data using Prometheus and visualize it in Grafana.

## Screenshots

You can view the total write and read bytes of each process.

<img width="998" alt="Total" src="https://github.com/user-attachments/assets/ff53ac2a-3cb0-4fef-9343-19878a12b0f3" />

And the read and write bytes of each process to each file.

<img width="998" alt="Process" src="https://github.com/user-attachments/assets/636eb391-af80-4105-ad0b-7d111269a86a" />

Or use Grafana to visualize results:

<img width="1627" alt="image" src="https://github.com/user-attachments/assets/40d91a59-7d3b-414c-9450-bbfa55e88c26" />

## Build

Make sure you have Node.js and Go installed.

```bash
npm i
make
```

## Prerequisites (macOS)

You must disable SIP (System Integrity Protection) by booting into recovery mode and running `csrutil disable`. This is because iosnoopng uses dtrace, which requires SIP to be disabled.

## Usage (macOS)

> [!WARNING]
> **Make sure your Mac has NOT slept since last boot**, otherwise you will almost certainly have a system freeze. If you are not sure, reboot your Mac now or save your work and be prepared for a system freeze if it happens.
>
> This tool uses dtrace internally. Due to an issue the darwin kernel, if your Mac has slept since the last boot, running dtrace may cause a system freeze. This is an issue with the darwin kernel and I have no control over it. See [this](https://forums.developer.apple.com/forums/thread/735939) .
>

```bash
sudo ./iosnoopng -d
```

By default, the server will listen on `http://localhost:8092`. You can use launchd or other tools like PM2 to run this as a daemon.

```
Usage of ./iosnoopng:
  -a string
    	Listen address (default "127.0.0.1:8092")
  -d	Deduplicate launchd processes. Since launchd writes on behalf of other processes if the same file is written by another process it will be counted twice. This option removes the duplicated entries from launchd.
  -f string
    	Load previously saved process stat file as startpoint
  -l string
    	Log level (default "info")
  -n	Do not run dtrace, only start the HTTP server
  -o string
    	Save raw dtrace output to a file
  -r	Replay previous log file before collecting new data
  -s int
    	Only merge entries that are not updated more than this number of seconds. (default 3600)
  -t int
    	Merge R/W smaller than this number of bytes into a single entry. This is useful for processes that write to many small files to save memory and make the output more readable. (default 33554432)
```

> [!WARNING]
> Before quitting iosnoopng, **make sure your Mac have NOT slept**, due to the same reason mentioned above. If your Mac has not slept, you can safely quit iosnoopng (for example, by ^C). However, if your Mac has slept, quitting iosnoopng may cause a system freeze. Either let iosnoopng running or save your work and **force** reboot your Mac (yes, normal reboot won't work, you have to press and hold the power button to force reboot).

## Extra Considerations

- Since iosnoopng records paths of files, its memory usage may grow over time. However, it should not be a problem for most users because iosnoopng has small-files-merging enabled by default. So lots of small files will be merged into one entry to save memory.
- The process list will continue to grow because it keeps the RW data of every process ever launched. If you have a lot of processes that *have distinct names* and you run iosnoopng for a long time (like over a month), the memory usage may grow. Currently, you can restart iosnoopng to clear the process list. Or file an issue if you need an HTTP API to clear the process list without restarting iosnoopng and I will add it.
- The full path is not reported on macOS (only part of it). You will see every file starting with `/` even if the file is not there. This is a limitation of dtrace on macOS. For example, if you see a file `/idindex/IdIndex_inputs_i`, it is not in the root directory. If you need to find its full path, you can use `find . -name IdIndex_inputs_i` to find it. In my case, it is `$HOME/Library/Caches/JetBrains/GoLand2024.2/index/idindex/IdIndex.storage_i`

## Scraping Metrics with Prometheus and Visualizing with Grafana

You can scrape the metrics with Prometheus. Here is an example configuration:

```yaml
scrape_configs:
  - job_name: "iosnoopng"
    static_configs:
      - targets: ["127.0.0.1:8092"]
```
You will see metrics names like `iosnoopng_read_bytes_total` and `iosnoopng_written_bytes_total`.

Once you have Grafana and Prometheus set up, you can import the example dashboard from `dashboards/process-read-and-writes.json`.
