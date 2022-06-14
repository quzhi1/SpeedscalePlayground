# SpeedscalePlayground
My online lab for speedscale

## Prerquisite
1. Have minikube up and running.
2. Install tilt: `curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash`
3. Go to https://app.speedscale.com/profile, and put your API key in .zshrc:
```bash
export SPEEDSCALE_API_KEY=<your-api-key>
```

## Running services
```bash
tilt up
```

## Speedscale replay
1. Run `kubectl patch deployment hello-world --patch-file patch.yaml`.
2. Go to https://app.speedscale.com/reports?page=1&pageSize=10&sort=1&order=2 and see your result.

## Key findings for Speedscale
### proxy sidecar restarts
When running replays, I see the proxy sidecar restating several times. This is the error log:
```
[event: pod hello-world-854d7f5d7-bbtwp] Successfully pulled image "gcr.io/speedscale/goproxy:v1.0.8" in 625.903666ms
{"L":"ERROR","T":"2022-06-14T01:22:56.177Z","M":"panic recovered","panic":"runtime error: invalid memory address or nil pointer dereference","stack":"goroutine 1 [running]:\nruntime/debug.Stack()\n\t/usr/local/go/src/runtime/debug/stack.go:24 +0x68\ngitlab.com/speedscale/speedscale/lib/log.OnPanic({0x0, 0x0, 0x0})\n\t/src/vendor/gitlab.com/speedscale/speedscale/lib/log/wrapper.go:537 +0x48\npanic({0x1581ee0, 0x2aa2830})\n\t/usr/local/go/src/runtime/panic.go:838 +0x20c\ngitlab.com/speedscale/speedscale/responder/server.(*AMQPv091ReplayServer).TotalTransactions(...)\n\t/src/server/ampqv091.go:421\ngitlab.com/speedscale/speedscale/responder/provider/amqpv091.(*AMQPv091Provider).Stop(0x400050d380, {0x4?, 0x182dff2?})\n\t/src/provider/amqpv091/amqpv091.go:93 +0xc4\ngitlab.com/speedscale/speedscale/responder/server.Run.func2({0x1b99e78, 0x400050d380})\n\t/src/server/server.go:127 +0x44\ngitlab.com/speedscale/speedscale/responder/server.Run({0x1ba5c60?, 0x4000ad5560}, 0x4000773200, {0x1bab2d0?, 0x40008023c0}, {0x1bab190, 0x4000a943c0}, 0x400059cf20, {0x400055af40, 0x3, ...})\n\t/src/server/server.go:240 +0xf30\nmain.main()\n\t/src/main.go:143 +0x11bc\n"}
```

It happens even if there is no port forwarding.

However, I don't see anything is broken due to this error.

### Problem with redis client
When running replays, my hello world server is having a hard time connecting to redis. The error log is
```
[speedscale-goproxy] {"L":"ERROR","T":"2022-06-14T01:22:50.672Z","M":"dial() ERR could not connect","source":"goproxy","version":"v1.0.8","k8sAppPodNamespace":"default","k8sAppPodName":"hello-world-5cb89f8966-sl7bc","k8sAppLabel":"hello-world","addr":"10.108.60.152:6379","error":"dial tcp 10.108.60.152:6379: connect: connection refused"}
```

I didn't port forward redis.

All my replay failures are caused by this error: https://app.speedscale.com/report/468ec0ff-1444-4a8f-a0ab-4b85932f8dd5.
