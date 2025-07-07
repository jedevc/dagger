package schema

import (
	"context"

	"github.com/dagger/dagger/core"
	"github.com/dagger/dagger/dagql"
)

type liveDirectorySchema struct {
	srv *dagql.Server
}

var _ SchemaResolvers = &liveDirectorySchema{}

func (s *liveDirectorySchema) Install() {
	dagql.Fields[*core.Query]{
		dagql.Func("liveDirectory", s.liveDirectory).
			Doc("Dagger liveDirectory configuration and state"),
	}.Install(s.srv)

	dagql.Fields[*core.LiveDirectory]{
		// NEVER CACHE!
		dagql.NodeFuncWithCacheKey("asDirectory", DagOpDirectoryWrapper(s.srv, s.asDirectory), dagql.CachePerClient),
	}.Install(s.srv)
}

type liveDirectoryArgs struct {
	Service core.ServiceID `name:"service"`
	Path    string         `name:"path"`
}

func (s *liveDirectorySchema) liveDirectory(ctx context.Context, parent *core.Query, args liveDirectoryArgs) (*core.LiveDirectory, error) {
	ssh, err := args.Service.Load(ctx, s.srv)
	if err != nil {
		return nil, err
	}
	// resource ids?
	return &core.LiveDirectory{
		Service: ssh,
		Path:    args.Path,
	}, nil
}

func (s *liveDirectorySchema) asDirectory(ctx context.Context, parent dagql.Instance[*core.LiveDirectory], args struct {
	DagOpInternalArgs
}) (inst dagql.Instance[*core.Directory], _ error) {
	dir, err := parent.Self.Directory(ctx)
	if err != nil {
		return inst, err
	}
	return dagql.NewInstanceForCurrentID(ctx, s.srv, parent, dir)
}
