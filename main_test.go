package main

import (
	"context"
	"net/url"
	"os"
	"testing"

	"github.com/akaitux/cert-manager-webhook-beget/begetapi"
	dns "github.com/cert-manager/cert-manager/test/acme"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

var (
	zone = os.Getenv("TEST_ZONE_NAME")
)

func TestRunsSuite(t *testing.T) {
	// The manifest path should contain a file named config.json that is a
	// snippet of valid configuration that should be included on the
	// ChallengeRequest passed as part of the test cases.
	begerURL, err := url.Parse("http://localhost:8080")
	if err != nil {
		t.FailNow()
	}

	api := begetapi.NewBegetApiMock("login", "password")
	go func() {
		api.Run(":8080")
		t.Log("run")
	}()
	go func() {
		api.RunDns("59351")
		t.Log("run dns")
	}()
	defer func() {
		api.Stop(context.TODO())
		api.StopDns(context.TODO())
		t.Log("stopped servers")
	}()

	solver := New(begerURL)
	fixture := dns.NewFixture(solver,
		dns.SetResolvedZone("example.com."),
		dns.SetManifestPath("testdata/beget"),
		dns.SetDNSServer("127.0.0.1:59351"),
		dns.SetUseAuthoritative(false),
	)

	fixture.RunConformance(t)

}

func TestLoadConfigAcceptsName(t *testing.T) {
	cfg, err := loadConfig(&extapi.JSON{Raw: []byte(`{
		"apiLoginSecretRef": {
			"name": "beget-credentials",
			"key": "login"
		},
		"apiPasswdSecretRef": {
			"name": "beget-credentials",
			"key": "passwd"
		}
	}`)})
	if err != nil {
		t.Fatalf("loadConfig returned error: %v", err)
	}

	if got := cfg.APILoginSecretRef.resourceName(); got != "beget-credentials" {
		t.Fatalf("unexpected login secret name: %q", got)
	}
	if got := cfg.APIPasswdSecretRef.resourceName(); got != "beget-credentials" {
		t.Fatalf("unexpected password secret name: %q", got)
	}
}

func TestLoadConfigAcceptsSecretName(t *testing.T) {
	cfg, err := loadConfig(&extapi.JSON{Raw: []byte(`{
		"apiLoginSecretRef": {
			"secretName": "beget-credentials",
			"key": "login"
		},
		"apiPasswdSecretRef": {
			"secretName": "beget-credentials",
			"key": "passwd"
		}
	}`)})
	if err != nil {
		t.Fatalf("loadConfig returned error: %v", err)
	}

	if got := cfg.APILoginSecretRef.resourceName(); got != "beget-credentials" {
		t.Fatalf("unexpected login secret name: %q", got)
	}
	if got := cfg.APIPasswdSecretRef.resourceName(); got != "beget-credentials" {
		t.Fatalf("unexpected password secret name: %q", got)
	}
}

func TestLoadConfigRejectsMissingSecretName(t *testing.T) {
	_, err := loadConfig(&extapi.JSON{Raw: []byte(`{
		"apiLoginSecretRef": {
			"key": "login"
		},
		"apiPasswdSecretRef": {
			"name": "beget-credentials",
			"key": "passwd"
		}
	}`)})
	if err == nil {
		t.Fatal("expected loadConfig to reject missing secret name")
	}
}
