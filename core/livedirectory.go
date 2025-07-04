package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"dagger.io/dagger/telemetry"
	"github.com/containerd/continuity/fs"
	"github.com/dagger/dagger/dagql"
	bkcache "github.com/moby/buildkit/cache"
	bkclient "github.com/moby/buildkit/client"
	"github.com/vektah/gqlparser/v2/ast"
)

type LiveDirectory struct {
	Service dagql.Instance[*Service]
}

func (*LiveDirectory) Type() *ast.Type {
	return &ast.Type{
		NamedType: "LiveDirectory",
		NonNull:   true,
	}
}

func (*LiveDirectory) TypeDescription() string {
	// XXX:
	return "Live directory"
}

func (dir *LiveDirectory) Directory(ctx context.Context) (*Directory, error) {
	query, err := CurrentQuery(ctx)
	if err != nil {
		return nil, err
	}

	newRef, err := query.BuildkitCache().New(ctx, nil, nil, bkcache.WithRecordType(bkclient.UsageRecordTypeRegular))
	if err != nil {
		return nil, err
	}

	err = MountRef(ctx, newRef, nil, func(dst string) error {
		return dir.mount(ctx, func(ctx context.Context, src string) error {
			res, err := os.ReadDir(src)
			if err != nil {
				return err
			}
			fmt.Println("okay...", src)
			fmt.Println(res)
			for _, e := range res {
				fmt.Println(e.Name())
			}
			return fs.CopyDir(dst, src)
		})
	})
	if err != nil {
		return nil, err
	}

	snap, err := newRef.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &Directory{
		Result:   snap,
		Dir:      "/",
		Platform: query.Platform(),
	}, nil
}

func (dir *LiveDirectory) mount(ctx context.Context, f func(context.Context, string) error) (rerr error) {
	query, err := CurrentQuery(ctx)
	if err != nil {
		return err
	}

	services, err := query.Services(ctx)
	if err != nil {
		return err
	}
	host, err := dir.Service.Self.Hostname(ctx, dir.Service.ID())
	if err != nil {
		return err
	}
	ports, err := dir.Service.Self.Ports(ctx, dir.Service.ID())
	if err != nil {
		return err
	}
	if len(ports) == 0 {
		return fmt.Errorf("no ports available")
	}
	port := ports[0]

	detach, _, err := services.StartBindings(ctx, ServiceBindings{
		ServiceBinding{
			Hostname: host,
			Service:  dir.Service,
		},
	})
	if err != nil {
		return err
	}
	defer detach()

	// sshfs [user@]hostname:[directory] mountpoint

	mountDir, err := os.MkdirTemp("", "sshfs")
	if err != nil {
		return err
	}
	defer os.RemoveAll(mountDir)

	netConf, err := DNSConfig(ctx)
	if err != nil {
		return err
	}
	var resolvPath string
	if netConf != nil {
		resolvPath, err = mountResolv(netConf)
		if err != nil {
			return err
		}
		defer os.Remove(resolvPath)
	}

	done := make(chan error)

	ctx, cancel := context.WithCancelCause(ctx)
	go func() {
		cmd := exec.CommandContext(ctx, "sshfs", "-v", "-o", "StrictHostKeyChecking=no", "-o", "auto_unmount", "-p", strconv.Itoa(port.Port), "root@"+host+":"+".", mountDir)
		ctx, span := Tracer(ctx).Start(ctx, fmt.Sprintf("exec %s", strings.Join(cmd.Args, " ")), telemetry.Encapsulated())
		stdio := telemetry.SpanStdio(ctx, "")
		defer stdio.Close()
		cmd.Stdout, cmd.Stderr = stdio.Stdout, stdio.Stderr

		err := runWithStandardUmaskAndNetOverride(ctx, cmd, "", resolvPath)
		if errors.Is(err, context.Canceled) {
			err = nil
		}
		if err != nil {
			err = fmt.Errorf("failed to mount sshfs: %w", err)
		}
		telemetry.End(span, func() error { return err })

		cancel(err)

		done <- err
	}()

	err = f(ctx, mountDir)
	cancel(nil)
	if err != nil {
		<-done
		return err
	}
	return <-done
}
