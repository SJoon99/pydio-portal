package object

import "testing"

func TestDataSourceClientConfigUsesCustomPathStyleEndpoint(t *testing.T) {
	ds := &DataSource{
		StorageType:   StorageType_S3,
		ObjectsHost:   "stale.example",
		ObjectsPort:   443,
		ObjectsSecure: true,
		ApiKey:        "access",
		ApiSecret:     "secret",
		StorageConfiguration: map[string]string{
			StorageKeyCustomEndpoint: "http://rook-ceph-rgw.e2e.svc:80",
			StorageKeyCustomRegion:   "us-east-1",
			StorageKeyPathStyle:      "true",
			StorageKeyForcePathStyle: "true",
		},
	}

	cfg := ds.ClientConfig()
	if got := cfg.Val("endpoint").String(); got != "rook-ceph-rgw.e2e.svc:80" {
		t.Fatalf("endpoint = %q, want %q", got, "rook-ceph-rgw.e2e.svc:80")
	}
	if cfg.Val("secure").Bool() {
		t.Fatal("secure = true, want false for an HTTP custom endpoint")
	}
	if !cfg.Val(StorageKeyPathStyle).Bool() {
		t.Fatal("path_style = false, want true")
	}
	if !cfg.Val(StorageKeyForcePathStyle).Bool() {
		t.Fatal("force_path_style = false, want true")
	}
}

func TestMinioClientConfigUsesCustomPathStyleEndpoint(t *testing.T) {
	config := &MinioConfig{
		StorageType:   StorageType_S3,
		RunningHost:   "stale.example",
		RunningPort:   443,
		RunningSecure: true,
		ApiKey:        "access",
		ApiSecret:     "secret",
		GatewayConfiguration: map[string]string{
			StorageKeyCustomEndpoint: "http://rook-ceph-rgw.e2e.svc:80",
			StorageKeyPathStyle:      "true",
			StorageKeyForcePathStyle: "true",
		},
	}

	cfg := config.ClientConfig()
	if got := cfg.Val("endpoint").String(); got != "rook-ceph-rgw.e2e.svc:80" {
		t.Fatalf("endpoint = %q, want %q", got, "rook-ceph-rgw.e2e.svc:80")
	}
	if cfg.Val("secure").Bool() {
		t.Fatal("secure = true, want false for an HTTP custom endpoint")
	}
	if !cfg.Val(StorageKeyPathStyle).Bool() {
		t.Fatal("path_style = false, want true")
	}
	if !cfg.Val(StorageKeyForcePathStyle).Bool() {
		t.Fatal("force_path_style = false, want true")
	}
}
