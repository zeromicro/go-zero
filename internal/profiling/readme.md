# Continuous profiling

Configure continuous profiling through the `Profiling` section of
`service.ServiceConf`, including services that embed it in REST or RPC configs.
A profiling server address is required to enable collection.

```yaml
Name: worker
Profiling:
  ServerAddr: http://localhost:4040
  CpuThreshold: 700
  HeapThreshold: 536870912 # 512 MiB of allocated Go heap objects
  GoroutineThreshold: 10000
  CheckInterval: 10s
  ProfilingDuration: 2m
```

Every `CheckInterval`, an idle profiler starts a collection session if **any**
enabled threshold is reached or exceeded. This can capture memory growth or
goroutine accumulation even when CPU usage stays low. Only one collection
session runs at a time, using the existing profiling duration and upload rate.

| Setting | Measurement | Default | Meaning of zero |
| --- | --- | --- | --- |
| `CpuThreshold` | CPU usage on a scale of 0–1000; 700 means 70% | `700` | Starts profiling unconditionally |
| `HeapThreshold` | Allocated Go heap objects, in bytes | `0` | Disables the heap trigger |
| `GoroutineThreshold` | Number of live goroutines | `0` | Disables the goroutine trigger |

The heap measurement uses Go's `/memory/classes/heap/objects:bytes` runtime
metric. It includes live objects and objects that have not yet been freed by
the garbage collector; it does not measure process RSS or a percentage of a
container memory limit. Use an integer byte count for `HeapThreshold`. Negative
resource thresholds are rejected by the configuration loader.

Existing configurations that omit the new thresholds keep CPU-only triggering.
The new measurements are read only when their thresholds are enabled and the
CPU threshold has not already triggered a session.

`ProfileType` selects the data collected during a session independently of
the trigger. For example, a heap threshold can start a session that captures
CPU, memory, and goroutine profiles together.
