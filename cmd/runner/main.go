package main

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/platforms"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
)

func main() {
	ctx := context.Background()
	err := start(ctx)
	if err != nil {
		panic(err)
	}
}

func start(ctx context.Context) error {
	address := "/run/containerd/containerd.sock"
	client, err := containerd.New(address, containerd.WithDefaultNamespace("dagger"))
	if err != nil {
		return errors.Wrapf(err, "failed to connect client to %q . make sure containerd is running", address)
	}

	target := "/image.tar"
	r, err := os.Open(target)
	if err != nil {
		return err
	}
	defer r.Close()

	fmt.Printf("importing %s... ", target)
	imgs, err := client.Import(ctx, r, containerd.WithIndexName("dagger-engine"))
	closeErr := r.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	fmt.Println("done")

	for _, img := range imgs {
		image := containerd.NewImageWithPlatform(client, img, platforms.DefaultStrict())

		fmt.Printf("unpacking %s (%s)... ", img.Name, img.Target.Digest)
		err = image.Unpack(ctx, "overlayfs")
		if err != nil {
			return err
		}
		fmt.Println("done")
	}

	image := containerd.NewImage(client, imgs[0])

	hostname, err := os.Hostname()
	if err != nil {
		return errors.Wrap(err, "get hostname")
	}

	err = os.MkdirAll("/var/run/containers", 0o755)
	if err != nil {
		return err
	}

	mounts := []specs.Mount{
		{
			Destination: "/run/containerd/",
			Type:        "bind",
			Options:     []string{"rbind"},
			Source:      "/run/containerd/",
		},
		{
			Destination: "/var/lib/containerd/",
			Type:        "bind",
			Options:     []string{"rbind"},
			Source:      "/var/lib/containerd/",
		},
		{
			Destination: "/tmp",
			Type:        "bind",
			Options:     []string{"rbind"},
			Source:      "/tmp",
		},
		{
			Destination: "/var/lib/dagger/",
			Type:        "bind",
			Options:     []string{"rbind"},
			Source:      "/var/lib/dagger/",
		},
	}
	for _, mount := range mounts {
		err = syscall.Mount(mount.Source, mount.Source, "", syscall.MS_BIND, "")
		if err != nil {
			return err
		}
		err = syscall.Mount("", mount.Source, "", syscall.MS_SHARED|syscall.MS_REC, "")
		if err != nil {
			return err
		}
	}

	container, err := client.NewContainer(ctx, "hello",
		containerd.WithImage(image),
		containerd.WithImageConfigLabels(image),
		containerd.WithSnapshotter("overlayfs"),
		containerd.WithNewSnapshot("hello", image),
		containerd.WithNewSpec(
			oci.WithImageConfig(image),

			// --privileged
			oci.WithPrivileged,
			oci.WithAllDevicesAllowed,
			oci.WithHostDevices,

			// --network=host
			oci.WithHostNamespace(specs.NetworkNamespace),
			oci.WithHostHostsFile,
			oci.WithHostResolvconf,
			oci.WithEnv([]string{fmt.Sprintf("HOSTNAME=%s", hostname)}),

			oci.WithImageConfigArgs(image, []string{"--debug"}),

			oci.WithMounts(mounts),
			func(_ context.Context, _ oci.Client, _ *containers.Container, s *oci.Spec) error {
				// XXX: rshared really messes up dns resolution in buildkit
				// container for some reason
				if s.Linux != nil {
					s.Linux.RootfsPropagation = "shared"
				} else {
					s.Linux = &specs.Linux{
						RootfsPropagation: "shared",
					}
				}

				return nil
			},
		),
	)
	if err != nil {
		return err
	}

	cioOpts := []cio.Opt{cio.WithStreams(nil, os.Stdout, os.Stderr)}
	// cioOpts = append(cioOpts, cio.WithTerminal)

	err = os.Mkdir("/dagger-engine-bind", 0o755)
	if err != nil {
		return err
	}

	taskOpts := []containerd.NewTaskOpts{}
	if err != nil {
		return err
	}

	task, err := container.NewTask(ctx, cio.NewCreator(cioOpts...), taskOpts...)
	if err != nil {
		return err
	}

	status, err := task.Wait(ctx)
	if err != nil {
		return err
	}
	err = task.Start(ctx)
	if err != nil {
		return err
	}

	result := <-status
	fmt.Printf("dagger engine exited with %d\n", result.ExitCode())

	return result.Error()
}
