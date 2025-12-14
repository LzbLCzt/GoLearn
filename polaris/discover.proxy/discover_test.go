package discover_proxy

import (
	"context"
	"fmt"
	apiV1Model "git.woa.com/polaris/polaris-server-api/api/v1/model"
	"testing"
)

func TestDiscover(t *testing.T) {

	rsp, err := DefaultBackend.discover(context.Background(), "polaris.discover",
		"lzb_test", "Test", apiV1Model.DiscoverRequest_INSTANCE)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	fmt.Printf("instances: %+v", rsp.GetInstances())
}
