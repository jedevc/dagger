package core

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"dagger.io/dagger/telemetry"
	"github.com/containerd/continuity/fs"
	"github.com/dagger/dagger/dagql"
	"github.com/dagger/dagger/engine"
	"github.com/dagger/dagger/network"
	bkcache "github.com/moby/buildkit/cache"
	bkclient "github.com/moby/buildkit/client"
	"github.com/vektah/gqlparser/v2/ast"
)

type LiveDirectory struct {
	Service dagql.Instance[*Service]
	Path    string
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
	clientMetadata, err := engine.ClientMetadataFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get client metadata: %w", err)
	}

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

	sessionName := fmt.Sprintf("dagger.%s.%s", clientMetadata.ClientID, dagql.HashFrom(dir.Path))

	domain := network.SessionDomain(clientMetadata.SessionID)

	err = dir.runMutagen(ctx, resolvPath, "sync", "create", "--name="+sessionName, "/var/lib/dagger/mutagen/"+sessionName, fmt.Sprintf("root@%s:%d:%s", host+"."+domain, port.Port, dir.Path))
	if err != nil {
		return err
	}
	defer func() {
		ctx := context.WithoutCancel(ctx)
		err = dir.runMutagen(ctx, resolvPath, "sync", "terminate", sessionName)
		if err != nil && rerr == nil {
			rerr = err
		}
	}()
	err = dir.runMutagen(ctx, resolvPath, "sync", "flush", sessionName)
	if err != nil {
		return err
	}
	err = dir.runMutagen(ctx, resolvPath, "sync", "pause", sessionName)
	if err != nil {
		return err
	}

	return f(ctx, "/mnt/foobar")
}

func (dir *LiveDirectory) runMutagen(ctx context.Context, resolvPath string, args ...string) (rerr error) {
	cmd := exec.CommandContext(ctx, "mutagen", args...)
	ctx, span := Tracer(ctx).Start(ctx, fmt.Sprintf("exec %s", strings.Join(cmd.Args, " ")), telemetry.Encapsulated())
	defer telemetry.End(span, func() error { return rerr })

	stdio := telemetry.SpanStdio(ctx, "")
	defer stdio.Close()
	cmd.Stdout, cmd.Stderr = stdio.Stdout, stdio.Stderr

	return runWithStandardUmaskAndNetOverride(ctx, cmd, "", resolvPath)
}
