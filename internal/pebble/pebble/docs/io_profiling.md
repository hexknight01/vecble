# I/O Profiling

Linux provide extensive kernel profiling capabilities, including the
ability to trace operations at the block I/O layer. These tools are
incredibly powerful, though sometimes overwhelming in their
flexibility. This document captures some common recipes for profiling
Linux I/O.

* [Perf](#perf)
* [Blktrace](#blktrace)

## Perf

The Linux `perf` command can instrument CPU performance counters, and
the extensive set of kernel trace points. A great place to get started
understanding `perf` are Brendan Gregg's [perf
examples](http://www.brendangregg.com/perf.html).

The two modes of operation are "live" reporting via `perf top`, and
record and report via `perf record` and `perf
{report,script}`. 

Recording the stack traces for `block:block_rq_insert` event allows
determination of what Pebble level code is generating block requests.

### Installation

Ubuntu AWS installation:

```
sudo apt-get install linux-tools-common linux-tools-4.4.0-1049-aws linux-cloud-tools-4.4.0-1049-aws
```

### Recording

`perf record` (and `perf top`) requires read and write access to
`/sys/kernel/debug/tracing`. Running as root as an easiest way to get
the right permissions.

```
# Trace all block device (disk I/O) requests with stack traces, until Ctrl-C.
sudo perf record -e block:block_rq_insert -ag

# Trace all block device (disk I/O) issues and completions with stack traces, until Ctrl-C.
sudo perf record -e block:block_rq_issue -e block:block_rq_complete -ag
```

The `-a` flag records events on all CPUs (almost always desirable).

The `-g` flag records call graphs (a.k.a stack traces). Capturing the
stack trace makes the recording somewhat more expensive, but it
enables determining the originator of the event. Note the stack traces
include both the kernel and application code, allowing pinpointing the
source of I/O as due to flush, compaction, WAL writes, etc.

The `-e` flag controls which events are instrumented. The list of
`perf` events is enormous. See `sudo perf list`.

The `-o` flag controls where output is recorded. The default is
`perf.data`.

In order to record events for a specific duration, you can append `--
sleep <duration>` to the command line.

```
# Trace all block device (disk I/O) requests with stack traces for 10s.
sudo perf record -e block:block_rq_insert -ag -- sleep 10
```

### Reporting

The recorded perf data (`perf.data`) can be explored using `perf
report` and `perf script`.

```
# Show perf.data in an ncurses browser.
sudo perf report

# Show perf.data as a text report.
sudo perf report --stdio
```

As an example, `perf report --stdio` from perf data gathered using
`perf record -e block:block_rq_insert -ag` will show something like:

```
    96.76%     0.00%  pebble          pebble             [.] runtime.goexit
                    |
                    ---runtime.goexit
                       |
                       |--85.58%-- github.com/cockroachdb/pebble/v2/internal/record.NewLogWriter.func2
                       |          runtime/pprof.Do
                       |          github.com/cockroachdb/pebble/v2/internal/record.(*LogWriter).flushLoop-fm
                       |          github.com/cockroachdb/pebble/v2/internal/record.(*LogWriter).flushLoop
                       |          github.com/cockroachdb/pebble/v2/internal/record.(*LogWriter).flushPending
                       |          github.com/cockroachdb/pebble/v2/vfs.(*syncingFile).Sync
                       |          github.com/cockroachdb/pebble/v2/vfs.(*syncingFile).syncFdatasync-fm
                       |          github.com/cockroachdb/pebble/v2/vfs.(*syncingFile).syncFdatasync
                       |          syscall.Syscall
                       |          entry_SYSCALL_64_fastpath
                       |          sys_fdatasync
                       |          do_fsync
                       |          vfs_fsync_range
                       |          ext4_sync_file
                       |          filemap_write_and_wait_range
                       |          __filemap_fdatawrite_range
                       |          do_writepages
                       |          ext4_writepages
                       |          blk_finish_plug
                       |          blk_flush_plug_list
                       |          blk_mq_flush_plug_list
                       |          blk_mq_insert_requests
```

This is showing that `96.76%` of block device requests on the entire
system were generated by the `pebble` process, and `85.58%` of the
block device requests on the entire system were generated from WAL
syncing within this `pebble` process.

The `perf script` command provides access to the raw request
data. While there are various pre-recorded scripts that can be
executed, it is primarily useful for seeing call stacks along with the
"trace" data. For block requests, the trace data shows the device, the
operation type, the offset, and the size.

```
# List all events from perf.data with recommended header and fields.
sudo perf script --header -F comm,pid,tid,cpu,time,event,ip,sym,dso,trace
...
pebble  6019/6019  [008] 16492.555957: block:block_rq_insert: 259,0 WS 0 () 3970952 + 256 [pebble]
            7fff813d791a blk_mq_insert_requests
            7fff813d8878 blk_mq_flush_plug_list
            7fff813ccc96 blk_flush_plug_list
            7fff813cd20c blk_finish_plug
            7fff812a143d ext4_writepages
            7fff8119ea1e do_writepages
            7fff81191746 __filemap_fdatawrite_range
            7fff8119188a filemap_write_and_wait_range
            7fff81297c41 ext4_sync_file
            7fff81244ecb vfs_fsync_range
            7fff81244f8d do_fsync
            7fff81245243 sys_fdatasync
            7fff8181ae6d entry_SYSCALL_64_fastpath
                  3145e0 syscall.Syscall
                  6eddf3 github.com/cockroachdb/pebble/v2/vfs.(*syncingFile).syncFdatasync
                  6f069a github.com/cockroachdb/pebble/v2/vfs.(*syncingFile).syncFdatasync-fm
                  6ed8d2 github.com/cockroachdb/pebble/v2/vfs.(*syncingFile).Sync
                  72542f github.com/cockroachdb/pebble/v2/internal/record.(*LogWriter).flushPending
                  724f5c github.com/cockroachdb/pebble/v2/internal/record.(*LogWriter).flushLoop
                  72855e github.com/cockroachdb/pebble/v2/internal/record.(*LogWriter).flushLoop-fm
                  7231d8 runtime/pprof.Do
                  727b09 github.com/cockroachdb/pebble/v2/internal/record.NewLogWriter.func2
                  2c0281 runtime.goexit
```

Let's break down the trace data:

```
259,0 WS 0 () 3970952 + 256
 |     |         |       |
 |     |         |       + size (sectors)
 |     |         |
 |     |         + offset (sectors)
 |     |
 |     +- flags: R(ead), W(rite), B(arrier), S(ync), D(iscard), N(one)
 |
 +- device: <major>, <minor>
```

The above is indicating that a synchronous write of `256` sectors was
performed starting at sector `3970952`. The sector size is device
dependent and can be determined with `blockdev --report <device>`,
though it is almost always `512` bytes. In this case, the sector size
is `512` bytes indicating that this is a write of 128 KB.

## Blktrace

The `blktrace` tool records similar info to `perf`, but is targeted to
the block layer instead of being general purpose. The `blktrace`
command records data, while the `blkparse` command parses and displays
data. The `btrace` command is a shortcut for piping the output from
`blktrace` directly into `blkparse.

### Installation

Ubuntu AWS installation:

```
sudo apt-get install blktrace
```

## Usage

```
# Pipe the output of blktrace directly into blkparse.
sudo blktrace -d /dev/nvme1n1 -o - | blkparse -i -

# Equivalently.
sudo btrace /dev/nvme1n1
```

The information captured by `blktrace` is similar to what `perf` captures:

```
sudo btrace /dev/nvme1n1
...
259,0    4      186     0.016411295 11538  Q  WS 129341760 + 296 [pebble]
259,0    4      187     0.016412100 11538  Q  WS 129342016 + 40 [pebble]
259,0    4      188     0.016412200 11538  G  WS 129341760 + 256 [pebble]
259,0    4      189     0.016412714 11538  G  WS 129342016 + 40 [pebble]
259,0    4      190     0.016413148 11538  U   N [pebble] 2
259,0    4      191     0.016413255 11538  I  WS 129341760 + 256 [pebble]
259,0    4      192     0.016413321 11538  I  WS 129342016 + 40 [pebble]
259,0    4      193     0.016414271 11538  D  WS 129341760 + 256 [pebble]
259,0    4      194     0.016414860 11538  D  WS 129342016 + 40 [pebble]
259,0   12      217     0.016687595     0  C  WS 129341760 + 256 [0]
259,0   12      218     0.016700021     0  C  WS 129342016 + 40 [0]
```

The standard format is:

```
<device> <cpu> <seqnum> <timestamp> <pid> <action> <RWBS> <start-sector> + <size> [<command>]
```

See `man blkparse` for an explanation of the actions.

The `blktrace` output can be used to highlight problematic I/O
patterns. For example, it can be used to determine there are an
excessive number of small sequential read I/Os indicating that dynamic
readahead is not working correctly.
