package config

import (
	"os"
	"testing"
)

func TestConfigAwsValues(t *testing.T) {
	region := "test-default-region"
	os.Setenv("AWS_DEFAULT_REGION", region)

	cfg, _ := New()
	cfg.Aws.Region = region

	found := cfg.Aws.GetRegion()
	if found != region {
		t.Errorf("region fetching failed [%s]", found)
	}

}
