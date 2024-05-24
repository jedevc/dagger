# Repro

This reproduces the issue found in https://github.com/dagger/dagger/pull/7272#issuecomment-2127636616

To run this, at the top-level of the dagger repo, run `./hack/dev` (ensure
you're on the `fix-weird-stall` branch)

Then:

    $ cd ./examples/sdk/go/repro
    $ _EXPERIMENTAL_DAGGER_CLI_BIN=../../../../bin/dagger ../../../../bin/dagger run --silent go run main.go
    linux/amd64
    trying to close

The last command should hang forever.

This situation occurs "naturally" during `sdk_python.go` - the module sets the
session information, etc, as well as `TRACEPARENT`. Then, the
`test_download.py` test *unsets* the environment variable to force starting a
new session binary - the test appears to pass, but stalls on shutdown.

I managed to get a stack dump of the session spun up:

`../../../../bin/dagger session --label dagger.io/sdk.name:go --label dagger.io/sdk.version:0.9.3` (ignore the version number, it's `replace`d in go.mod)

```
goroutine 1 [semacquire]:
sync.runtime_Semacquire(0x1d58050?)
	/usr/lib/go/src/runtime/sema.go:62 +0x25
sync.(*WaitGroup).Wait(0x0?)
	/usr/lib/go/src/sync/waitgroup.go:116 +0x48
golang.org/x/sync/errgroup.(*Group).Wait(0xc0005b4bc0)
	/home/jedevc/go/pkg/mod/golang.org/x/sync@v0.7.0/errgroup/errgroup.go:56 +0x25
github.com/dagger/dagger/engine/client.(*Client).Close(0xc000218a80)
	/home/jedevc/Documents/Projects/dagger/engine/client/client.go:470 +0x3e5
main.withEngine({0x1444e40, 0xc000493a10}, {{0x0, 0x0}, {0x0, 0x0}, {0xc00049a000, 0x24}, {0xc0004241b0, 0x29}, ...}, ...)
	/home/jedevc/Documents/Projects/dagger/cmd/dagger/engine.go:44 +0x2a2
main.EngineSession(0xc00040a908?, {0x1286629?, 0x4?, 0x128662d?})
	/home/jedevc/Documents/Projects/dagger/cmd/dagger/session.go:93 +0x3b6
github.com/spf13/cobra.(*Command).execute(0xc00040a908, {0xc000127880, 0x4, 0x4})
	/home/jedevc/go/pkg/mod/github.com/spf13/cobra@v1.8.0/command.go:983 +0xaca
github.com/spf13/cobra.(*Command).ExecuteC(0x1d7c1e0)
	/home/jedevc/go/pkg/mod/github.com/spf13/cobra@v1.8.0/command.go:1115 +0x3ff
github.com/spf13/cobra.(*Command).Execute(...)
	/home/jedevc/go/pkg/mod/github.com/spf13/cobra@v1.8.0/command.go:1039
github.com/spf13/cobra.(*Command).ExecuteContext(...)
	/home/jedevc/go/pkg/mod/github.com/spf13/cobra@v1.8.0/command.go:1032
main.main.func1({0x1444e40, 0xc0003cdbf0})
	/home/jedevc/Documents/Projects/dagger/cmd/dagger/main.go:297 +0x458
github.com/dagger/dagger/dagql/idtui.(*frontendPlain).Run(0xc000001b00, {0x1444d60, 0x1df1b00}, {0x30?, 0x1f?, 0xc000060008?}, 0x13072a0)
	/home/jedevc/Documents/Projects/dagger/dagql/idtui/frontend_plain.go:153 +0x135
main.main()
	/home/jedevc/Documents/Projects/dagger/cmd/dagger/main.go:270 +0x289

goroutine 6 [select]:
github.com/dagger/dagger/dagql/idtui.(*frontendPlain).Run.func1()
	/home/jedevc/Documents/Projects/dagger/dagql/idtui/frontend_plain.go:136 +0xca
created by github.com/dagger/dagger/dagql/idtui.(*frontendPlain).Run in goroutine 1
	/home/jedevc/Documents/Projects/dagger/dagql/idtui/frontend_plain.go:133 +0x11e

goroutine 14 [select]:
github.com/dagger/dagger/telemetry/inflight.(*batchSpanProcessor).processQueue(0xc0003f8420)
	/home/jedevc/Documents/Projects/dagger/telemetry/inflight/batch_processor.go:329 +0x111
github.com/dagger/dagger/telemetry/inflight.NewBatchSpanProcessor.func1()
	/home/jedevc/Documents/Projects/dagger/telemetry/inflight/batch_processor.go:126 +0x57
created by github.com/dagger/dagger/telemetry/inflight.NewBatchSpanProcessor in goroutine 1
	/home/jedevc/Documents/Projects/dagger/telemetry/inflight/batch_processor.go:124 +0x339

goroutine 15 [select]:
github.com/dagger/dagger/telemetry/sdklog.(*batchLogProcessor).processQueue(0xc00040dd60)
	/home/jedevc/Documents/Projects/dagger/telemetry/sdklog/batch_processor.go:232 +0x11d
github.com/dagger/dagger/telemetry/sdklog.NewBatchLogProcessor.func1()
	/home/jedevc/Documents/Projects/dagger/telemetry/sdklog/batch_processor.go:113 +0x54
created by github.com/dagger/dagger/telemetry/sdklog.NewBatchLogProcessor in goroutine 1
	/home/jedevc/Documents/Projects/dagger/telemetry/sdklog/batch_processor.go:111 +0x2e5

goroutine 16 [select]:
github.com/dagger/dagger/analytics.(*CloudTracker).start(0xc000405c80)
	/home/jedevc/Documents/Projects/dagger/analytics/analytics.go:168 +0xac
created by github.com/dagger/dagger/analytics.New in goroutine 1
	/home/jedevc/Documents/Projects/dagger/analytics/analytics.go:124 +0x119

goroutine 52 [syscall]:
os/signal.signal_recv()
	/usr/lib/go/src/runtime/sigqueue.go:152 +0x29
os/signal.loop()
	/usr/lib/go/src/os/signal/signal_unix.go:23 +0x13
created by os/signal.Notify.func1.1 in goroutine 1
	/usr/lib/go/src/os/signal/signal.go:151 +0x1f

goroutine 41 [chan receive]:
main.EngineSession.func2()
	/home/jedevc/Documents/Projects/dagger/cmd/dagger/session.go:81 +0x2b
created by main.EngineSession in goroutine 1
	/home/jedevc/Documents/Projects/dagger/cmd/dagger/session.go:80 +0x185

goroutine 29 [select]:
google.golang.org/grpc/internal/transport.(*recvBufferReader).readClient(0xc0002f8050, {0xc00021ea90, 0x5, 0x5})
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/transport/transport.go:193 +0x95
google.golang.org/grpc/internal/transport.(*recvBufferReader).Read(0xc0002f8050, {0xc00021ea90?, 0xc00029c018?, 0xc000179730?})
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/transport/transport.go:173 +0x12d
google.golang.org/grpc/internal/transport.(*transportReader).Read(0xc00021ea20, {0xc00021ea90?, 0xc0001797a8?, 0xa51d65?})
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/transport/transport.go:525 +0x2c
io.ReadAtLeast({0x1439b80, 0xc00021ea20}, {0xc00021ea90, 0x5, 0x5}, 0x5)
	/usr/lib/go/src/io/io.go:335 +0x90
io.ReadFull(...)
	/usr/lib/go/src/io/io.go:354
google.golang.org/grpc/internal/transport.(*Stream).Read(0xc000290240, {0xc00021ea90, 0x5, 0x5})
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/transport/transport.go:509 +0x96
google.golang.org/grpc.(*parser).recvMsg(0xc00021ea80, 0x1000000)
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/rpc_util.go:614 +0x46
google.golang.org/grpc.recvAndDecompress(0xc00021ea80, 0xc000290240, {0x0, 0x0}, 0x1000000, 0x0, {0x0, 0x0})
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/rpc_util.go:753 +0x85
google.golang.org/grpc.recv(0x0?, {0x7e2f70f6dbb0, 0x1df1b00}, 0x0?, {0x0?, 0x0?}, {0x116fcc0, 0xc00070eb80}, 0x0?, 0x0, ...)
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/rpc_util.go:833 +0x7d
google.golang.org/grpc.(*csAttempt).recvMsg(0xc0002b49c0, {0x116fcc0, 0xc00070eb80}, 0xc000713f80?)
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/stream.go:1085 +0x289
google.golang.org/grpc.(*clientStream).RecvMsg.func1(0x40?)
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/stream.go:928 +0x1f
google.golang.org/grpc.(*clientStream).withRetry(0xc0002130e0, 0xc000179bd0, 0xc000179c18)
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/stream.go:761 +0x3ae
google.golang.org/grpc.(*clientStream).RecvMsg(0xc0002130e0, {0x116fcc0?, 0xc00070eb80?})
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/stream.go:927 +0x110
github.com/dagger/dagger/telemetry.(*tracesSourceSubscribeClient).Recv(0xc0004d62c0)
	/home/jedevc/Documents/Projects/dagger/telemetry/servers_grpc.pb.go:69 +0x46
github.com/dagger/dagger/engine/client.(*Client).exportTraces.func1()
	/home/jedevc/Documents/Projects/dagger/engine/client/client.go:497 +0xfc
golang.org/x/sync/errgroup.(*Group).Go.func1()
	/home/jedevc/go/pkg/mod/golang.org/x/sync@v0.7.0/errgroup/errgroup.go:78 +0x56
created by golang.org/x/sync/errgroup.(*Group).Go in goroutine 1
	/home/jedevc/go/pkg/mod/golang.org/x/sync@v0.7.0/errgroup/errgroup.go:75 +0x96

goroutine 46 [IO wait]:
internal/poll.runtime_pollWait(0x7e2f794ab2e8, 0x72)
	/usr/lib/go/src/runtime/netpoll.go:345 +0x85
internal/poll.(*pollDesc).wait(0xc000692360?, 0xc000432000?, 0x1)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:84 +0x27
internal/poll.(*pollDesc).waitRead(...)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:89
internal/poll.(*FD).Read(0xc000692360, {0xc000432000, 0x8000, 0x8000})
	/usr/lib/go/src/internal/poll/fd_unix.go:164 +0x27a
os.(*File).read(...)
	/usr/lib/go/src/os/file_posix.go:29
os.(*File).Read(0xc000694078, {0xc000432000?, 0x0?, 0x800010601?})
	/usr/lib/go/src/os/file.go:118 +0x52
github.com/docker/cli/cli/connhelper/commandconn.(*commandConn).Read(0xc000034000, {0xc000432000?, 0xc0007321c0?, 0xc0000afd90?})
	/home/jedevc/go/pkg/mod/github.com/docker/cli@v26.1.0+incompatible/cli/connhelper/commandconn/commandconn.go:160 +0x29
bufio.(*Reader).Read(0xc0006925a0, {0xc0001a4040, 0x9, 0xc000580808?})
	/usr/lib/go/src/bufio/bufio.go:241 +0x197
io.ReadAtLeast({0x14357e0, 0xc0006925a0}, {0xc0001a4040, 0x9, 0x9}, 0x9)
	/usr/lib/go/src/io/io.go:335 +0x90
io.ReadFull(...)
	/usr/lib/go/src/io/io.go:354
golang.org/x/net/http2.readFrameHeader({0xc0001a4040, 0x9, 0xc0006e80d8?}, {0x14357e0?, 0xc0006925a0?})
	/home/jedevc/go/pkg/mod/golang.org/x/net@v0.24.0/http2/frame.go:237 +0x65
golang.org/x/net/http2.(*Framer).ReadFrame(0xc0001a4000)
	/home/jedevc/go/pkg/mod/golang.org/x/net@v0.24.0/http2/frame.go:498 +0x85
google.golang.org/grpc/internal/transport.(*http2Client).reader(0xc000416008, 0xc000692600)
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/transport/http2_client.go:1602 +0x22d
created by google.golang.org/grpc/internal/transport.newHTTP2Client in goroutine 86
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/transport/http2_client.go:409 +0x1e79

goroutine 59 [IO wait]:
internal/poll.runtime_pollWait(0x7e2f794ab6c8, 0x72)
	/usr/lib/go/src/runtime/netpoll.go:345 +0x85
internal/poll.(*pollDesc).wait(0xc00001a300?, 0xc0003da000?, 0x0)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:84 +0x27
internal/poll.(*pollDesc).waitRead(...)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:89
internal/poll.(*FD).Read(0xc00001a300, {0xc0003da000, 0xc80, 0xc80})
	/usr/lib/go/src/internal/poll/fd_unix.go:164 +0x27a
net.(*netFD).Read(0xc00001a300, {0xc0003da000?, 0x7e2f70f631b8?, 0xc000596240?})
	/usr/lib/go/src/net/fd_posix.go:55 +0x25
net.(*conn).Read(0xc00022a008, {0xc0003da000?, 0xc00026f918?, 0x4139db?})
	/usr/lib/go/src/net/net.go:179 +0x45
crypto/tls.(*atLeastReader).Read(0xc000596240, {0xc0003da000?, 0x0?, 0xc000596240?})
	/usr/lib/go/src/crypto/tls/conn.go:806 +0x3b
bytes.(*Buffer).ReadFrom(0xc000266d30, {0x1437c80, 0xc000596240})
	/usr/lib/go/src/bytes/buffer.go:211 +0x98
crypto/tls.(*Conn).readFromUntil(0xc000266a88, {0x1436080, 0xc00022a008}, 0xc00026f960?)
	/usr/lib/go/src/crypto/tls/conn.go:828 +0xde
crypto/tls.(*Conn).readRecordOrCCS(0xc000266a88, 0x0)
	/usr/lib/go/src/crypto/tls/conn.go:626 +0x3cf
crypto/tls.(*Conn).readRecord(...)
	/usr/lib/go/src/crypto/tls/conn.go:588
crypto/tls.(*Conn).Read(0xc000266a88, {0xc0004b6000, 0x1000, 0xc0005828c0?})
	/usr/lib/go/src/crypto/tls/conn.go:1370 +0x156
bufio.(*Reader).Read(0xc00010ecc0, {0xc0005b6120, 0x9, 0x1d57f40?})
	/usr/lib/go/src/bufio/bufio.go:241 +0x197
io.ReadAtLeast({0x14357e0, 0xc00010ecc0}, {0xc0005b6120, 0x9, 0x9}, 0x9)
	/usr/lib/go/src/io/io.go:335 +0x90
io.ReadFull(...)
	/usr/lib/go/src/io/io.go:354
net/http.http2readFrameHeader({0xc0005b6120, 0x9, 0xc0003b29f0?}, {0x14357e0?, 0xc00010ecc0?})
	/usr/lib/go/src/net/http/h2_bundle.go:1638 +0x65
net/http.(*http2Framer).ReadFrame(0xc0005b60e0)
	/usr/lib/go/src/net/http/h2_bundle.go:1905 +0x85
net/http.(*http2clientConnReadLoop).run(0xc00026ffa8)
	/usr/lib/go/src/net/http/h2_bundle.go:9342 +0x12c
net/http.(*http2ClientConn).readLoop(0xc000218600)
	/usr/lib/go/src/net/http/h2_bundle.go:9237 +0x65
created by net/http.(*http2Transport).newClientConn in goroutine 58
	/usr/lib/go/src/net/http/h2_bundle.go:7887 +0xca6

goroutine 30 [select]:
google.golang.org/grpc.newClientStreamWithParams.func4()
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/stream.go:392 +0x8c
created by google.golang.org/grpc.newClientStreamWithParams in goroutine 1
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/stream.go:391 +0xe08

goroutine 44 [IO wait]:
internal/poll.runtime_pollWait(0x7e2f794ab1f0, 0x72)
	/usr/lib/go/src/runtime/netpoll.go:345 +0x85
internal/poll.(*pollDesc).wait(0xc000692420?, 0xc000330000?, 0x1)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:84 +0x27
internal/poll.(*pollDesc).waitRead(...)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:89
internal/poll.(*FD).Read(0xc000692420, {0xc000330000, 0x8000, 0x8000})
	/usr/lib/go/src/internal/poll/fd_unix.go:164 +0x27a
os.(*File).read(...)
	/usr/lib/go/src/os/file_posix.go:29
os.(*File).Read(0xc000694088, {0xc000330000?, 0x1048020?, 0x10dc801?})
	/usr/lib/go/src/os/file.go:118 +0x52
io.copyBuffer({0x143a7c0, 0xc0000bc3a0}, {0x1435ae0, 0xc0000a0560}, {0x0, 0x0, 0x0})
	/usr/lib/go/src/io/io.go:429 +0x191
io.Copy(...)
	/usr/lib/go/src/io/io.go:388
os.genericWriteTo(0xc000694088?, {0x143a7c0, 0xc0000bc3a0})
	/usr/lib/go/src/os/file.go:269 +0x58
os.(*File).WriteTo(0xc000694088, {0x143a7c0, 0xc0000bc3a0})
	/usr/lib/go/src/os/file.go:247 +0x9c
io.copyBuffer({0x143a7c0, 0xc0000bc3a0}, {0x14356e0, 0xc000694088}, {0x0, 0x0, 0x0})
	/usr/lib/go/src/io/io.go:411 +0x9d
io.Copy(...)
	/usr/lib/go/src/io/io.go:388
os/exec.(*Cmd).writerDescriptor.func1()
	/usr/lib/go/src/os/exec/exec.go:577 +0x34
os/exec.(*Cmd).Start.func2(0xc000134308?)
	/usr/lib/go/src/os/exec/exec.go:724 +0x2c
created by os/exec.(*Cmd).Start in goroutine 86
	/usr/lib/go/src/os/exec/exec.go:723 +0x9ab

goroutine 83 [select]:
google.golang.org/grpc/internal/grpcsync.(*CallbackSerializer).run(0xc0002f32d0, {0x1444e78, 0xc0005ae4b0})
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/grpcsync/callback_serializer.go:76 +0x115
created by google.golang.org/grpc/internal/grpcsync.NewCallbackSerializer in goroutine 1
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/grpcsync/callback_serializer.go:52 +0x11a

goroutine 84 [select]:
google.golang.org/grpc/internal/grpcsync.(*CallbackSerializer).run(0xc0002f3300, {0x1444e78, 0xc0005ae500})
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/grpcsync/callback_serializer.go:76 +0x115
created by google.golang.org/grpc/internal/grpcsync.NewCallbackSerializer in goroutine 1
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/grpcsync/callback_serializer.go:52 +0x11a

goroutine 85 [select]:
google.golang.org/grpc/internal/grpcsync.(*CallbackSerializer).run(0xc0002f3330, {0x1444e78, 0xc0005ae550})
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/grpcsync/callback_serializer.go:76 +0x115
created by google.golang.org/grpc/internal/grpcsync.NewCallbackSerializer in goroutine 1
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/grpcsync/callback_serializer.go:52 +0x11a

goroutine 28 [select]:
google.golang.org/grpc.newClientStreamWithParams.func4()
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/stream.go:392 +0x8c
created by google.golang.org/grpc.newClientStreamWithParams in goroutine 1
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/stream.go:391 +0xe08

goroutine 47 [select]:
google.golang.org/grpc/internal/transport.(*controlBuffer).get(0xc00068e410, 0x1)
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/transport/controlbuf.go:418 +0x113
google.golang.org/grpc/internal/transport.(*loopyWriter).run(0xc000198070)
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/transport/controlbuf.go:551 +0x86
google.golang.org/grpc/internal/transport.newHTTP2Client.func6()
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/transport/http2_client.go:463 +0x85
created by google.golang.org/grpc/internal/transport.newHTTP2Client in goroutine 86
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/transport/http2_client.go:461 +0x242b

goroutine 31 [select]:
google.golang.org/grpc/internal/transport.(*recvBufferReader).readClient(0xc0002f8280, {0xc00021f4e0, 0x5, 0x5})
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/transport/transport.go:193 +0x95
google.golang.org/grpc/internal/transport.(*recvBufferReader).Read(0xc0002f8280, {0xc00021f4e0?, 0xc00029c048?, 0xc0004db928?})
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/transport/transport.go:173 +0x12d
google.golang.org/grpc/internal/transport.(*transportReader).Read(0xc00021f470, {0xc00021f4e0?, 0xc0004db9a0?, 0xa51d65?})
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/transport/transport.go:525 +0x2c
io.ReadAtLeast({0x1439b80, 0xc00021f470}, {0xc00021f4e0, 0x5, 0x5}, 0x5)
	/usr/lib/go/src/io/io.go:335 +0x90
io.ReadFull(...)
	/usr/lib/go/src/io/io.go:354
google.golang.org/grpc/internal/transport.(*Stream).Read(0xc0002905a0, {0xc00021f4e0, 0x5, 0x5})
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/internal/transport/transport.go:509 +0x96
google.golang.org/grpc.(*parser).recvMsg(0xc00021f4d0, 0x1000000)
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/rpc_util.go:614 +0x46
google.golang.org/grpc.recvAndDecompress(0xc00021f4d0, 0xc0002905a0, {0x0, 0x0}, 0x1000000, 0x0, {0x0, 0x0})
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/rpc_util.go:753 +0x85
google.golang.org/grpc.recv(0x1?, {0x7e2f70f6dbb0, 0x1df1b00}, 0x4d424d56457a5556?, {0x0?, 0x0?}, {0x116f9c0, 0xc000338480}, 0x0?, 0x0, ...)
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/rpc_util.go:833 +0x7d
google.golang.org/grpc.(*csAttempt).recvMsg(0xc0003d84e0, {0x116f9c0, 0xc000338480}, 0x1df1b00?)
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/stream.go:1085 +0x289
google.golang.org/grpc.(*clientStream).RecvMsg.func1(0x40?)
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/stream.go:928 +0x1f
google.golang.org/grpc.(*clientStream).withRetry(0xc000290360, 0xc0004dbdc8, 0xc0004dbe10)
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/stream.go:761 +0x3ae
google.golang.org/grpc.(*clientStream).RecvMsg(0xc000290360, {0x116f9c0?, 0xc000338480?})
	/home/jedevc/go/pkg/mod/google.golang.org/grpc@v1.63.2/stream.go:927 +0x110
github.com/dagger/dagger/telemetry.(*logsSourceSubscribeClient).Recv(0xc0004d6550)
	/home/jedevc/Documents/Projects/dagger/telemetry/servers_grpc.pb.go:186 +0x46
github.com/dagger/dagger/engine/client.(*Client).exportLogs.func1()
	/home/jedevc/Documents/Projects/dagger/engine/client/client.go:541 +0xf8
golang.org/x/sync/errgroup.(*Group).Go.func1()
	/home/jedevc/go/pkg/mod/golang.org/x/sync@v0.7.0/errgroup/errgroup.go:78 +0x56
created by golang.org/x/sync/errgroup.(*Group).Go in goroutine 1
	/home/jedevc/go/pkg/mod/golang.org/x/sync@v0.7.0/errgroup/errgroup.go:75 +0x96

goroutine 93 [IO wait]:
internal/poll.runtime_pollWait(0x7e2f794aae10, 0x72)
	/usr/lib/go/src/runtime/netpoll.go:345 +0x85
internal/poll.(*pollDesc).wait(0xc0001c8880?, 0xc000518000?, 0x0)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:84 +0x27
internal/poll.(*pollDesc).waitRead(...)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:89
internal/poll.(*FD).Read(0xc0001c8880, {0xc000518000, 0x1000, 0x1000})
	/usr/lib/go/src/internal/poll/fd_unix.go:164 +0x27a
net.(*netFD).Read(0xc0001c8880, {0xc000518000?, 0xc0004d9a98?, 0x4f7845?})
	/usr/lib/go/src/net/fd_posix.go:55 +0x25
net.(*conn).Read(0xc00022a558, {0xc000518000?, 0x0?, 0xc0000dea88?})
	/usr/lib/go/src/net/net.go:179 +0x45
net/http.(*connReader).Read(0xc0000dea80, {0xc000518000, 0x1000, 0x1000})
	/usr/lib/go/src/net/http/server.go:789 +0x14b
bufio.(*Reader).fill(0xc000040a20)
	/usr/lib/go/src/bufio/bufio.go:110 +0x103
bufio.(*Reader).Peek(0xc000040a20, 0x4)
	/usr/lib/go/src/bufio/bufio.go:148 +0x53
net/http.(*conn).serve(0xc000411560, {0x1444e40, 0xc000493020})
	/usr/lib/go/src/net/http/server.go:2074 +0x749
created by net/http.(*Server).Serve in goroutine 1
	/usr/lib/go/src/net/http/server.go:3285 +0x4b4
```