package discover_proxy

import (
	"context"

	"git.code.oa.com/trpc-go/trpc-go/client"
	apiV1Model "git.woa.com/polaris/polaris-server-api/api/v1/model"
)

type IDiscover interface {
	Discover(ctx context.Context, req *apiV1Model.DiscoverRequest, opts ...client.Option) (rsp *apiV1Model.DiscoverResponse, err error)
}