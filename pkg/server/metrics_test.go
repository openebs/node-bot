/*
Copyright 2018 The OpenEBS Author

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/openebs/node-disk-manager/pkg/server"
)

func TestMetricsMiddleware(t *testing.T) {
	match := []*regexp.Regexp{
		regexp.MustCompile(`# HELP ndm_uptime_seconds Uptime of node disk manager.`),
		regexp.MustCompile(`# TYPE ndm_uptime_seconds gauge`),
		regexp.MustCompile(`# HELP ndm_collector_uptime_seconds Uptime of collector.`),
		regexp.MustCompile(`# TYPE ndm_collector_uptime_seconds gauge`),
	}
	fakeHandler := server.MetricsMiddleware()
	server := httptest.NewServer(fakeHandler)
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	buf, buferr := ioutil.ReadAll(resp.Body)
	if buferr != nil {
		t.Fatal(buferr)
	}
	for _, rexp := range match {
		if !rexp.Match(buf) {
			t.Error(rexp, " Not present in response body.")
		} else {
			t.Logf(rexp.String() + " Present in response body.")
		}
	}
	defer resp.Body.Close()
	defer server.Close()
}
