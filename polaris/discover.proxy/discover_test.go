package discover_proxy

import (
	"context"
	"fmt"
	"testing"

	apiV1Model "git.woa.com/polaris/polaris-server-api/api/v1/model"
)

func TestDiscover(t *testing.T) {

	rsp, err := DefaultBackend.discoverWithGrpc(context.Background(), "polaris.discover",
		"lzb_test", "Test", apiV1Model.DiscoverRequest_INSTANCE)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	fmt.Printf("instances: %+v", rsp.GetInstances())
}
