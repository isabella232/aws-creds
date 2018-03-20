package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lob/aws-creds/config"
	"github.com/lob/aws-creds/test"
)

func TestExecuteRefresh(t *testing.T) {
	appPath := "/app/url"
	appSuccessResponse := test.LoadTestFile(t, "app_success_response.html")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == appPath {
			_, err := w.Write([]byte(appSuccessResponse))
			if err != nil {
				t.Fatalf("unexpected error when writing response: %s", err)
			}
			return
		}
		_, err := w.Write([]byte("{}"))
		if err != nil {
			t.Fatalf("unexpected error when writing response: %s", err)
		}
	}))
	defer srv.Close()
	conf := config.New("")
	conf.OktaHost = srv.URL
	conf.OktaAppPath = appPath
	conf.Profiles = []*config.Profile{{"staging", "arn:staging"}}
	cmd := &Cmd{
		Command: "",
		Config:  conf,
		Profile: conf.Profiles[0].Name,
		Input:   test.NewNoopInput(),
	}

	err := executeRefresh(cmd)
	if err != nil {
		t.Fatalf("unexpected error when executing refresh: %s", err)
	}

	cmd.Profile = "invalid"
	err = executeRefresh(cmd)
	if err == nil {
		t.Fatalf("expected error when executing refresh with an invalid profile: %s", err)
	}

	cmd.Profile = ""
	err = executeRefresh(cmd)
	if err == nil {
		t.Fatalf("expected error when executing refresh without a profile: %s", err)
	}
}