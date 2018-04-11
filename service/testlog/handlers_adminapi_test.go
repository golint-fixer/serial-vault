// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2018 Canonical Ltd
 * License granted by Canonical Limited
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package testlog_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/CanonicalLtd/serial-vault/datastore"
	"github.com/CanonicalLtd/serial-vault/service"
	"github.com/CanonicalLtd/serial-vault/service/response"
	check "gopkg.in/check.v1"
)

type SyncTest struct {
	Method      string
	URL         string
	Data        []byte
	Code        int
	Type        string
	Permissions int
	EnableAuth  bool
	Success     bool
	SkipJWT     bool
	MockError   bool
	Count       int
}

const exampleFile = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4NCjx0ZXN0X3JlcG9ydD4NCiAgICA8dXV0cz4NCiAgICAgICAgPHV1dD4NCiAgICAgICAgICAgIDxzdW1tYXJ5Pg0KICAgICAgICAgICAgICAgIDxwYXJ0X251bWJlcj44NjAtMDAwMTQ8L3BhcnRfbnVtYmVyPg0KICAgICAgICAgICAgICAgIDxzZXJpYWxfbnVtYmVyPmU0YmY3Y2ViLWY3YWYtNDQyZS1iNzM0LTU0MzJlZjMzMDZiNTwvc2VyaWFsX251bWJlcj4NCiAgICAgICAgICAgICAgICA8b3BlcmF0aW9uPlZhbGlkYXRpb24gVGVzdDwvb3BlcmF0aW9uPg0KICAgICAgICAgICAgICAgIDxzdGFydGVkX2F0PjIwMTgtMDQtMDlUMTc6NDQ6MTYrMDI6MDA8L3N0YXJ0ZWRfYXQ+DQogICAgICAgICAgICAgICAgPGVuZGVkX2F0PjIwMTgtMDQtMDlUMTc6NDQ6MTYrMDI6MDA8L2VuZGVkX2F0Pg0KICAgICAgICAgICAgICAgIDxzdGF0dXM+RmFpbGVkPC9zdGF0dXM+DQogICAgICAgICAgICA8L3N1bW1hcnk+DQogICAgICAgICAgICA8dGVzdHM+DQogICAgICAgICAgICAgICAgPHRlc3Q+DQogICAgICAgICAgICAgICAgICAgIDxuYW1lPmZhY3RvcnlfY3B1L2lNWDZVTEw8L25hbWU+DQogICAgICAgICAgICAgICAgICAgIDxzdGF0dXM+ZmFpbGVkPC9zdGF0dXM+DQogICAgICAgICAgICAgICAgPC90ZXN0Pg0KICAgICAgICAgICAgICAgIDx0ZXN0Pg0KICAgICAgICAgICAgICAgICAgICA8bmFtZT5mYWN0b3J5X2V0aGVybmV0L2NhcmQtZGV0ZWN0PC9uYW1lPg0KICAgICAgICAgICAgICAgICAgICA8c3RhdHVzPmZhaWxlZDwvc3RhdHVzPg0KICAgICAgICAgICAgICAgIDwvdGVzdD4NCiAgICAgICAgICAgICAgICA8dGVzdD4NCiAgICAgICAgICAgICAgICAgICAgPG5hbWU+ZmFjdG9yeV9oZGQvZHJpdmUtY291bnQ8L25hbWU+DQogICAgICAgICAgICAgICAgICAgIDxzdGF0dXM+ZmFpbGVkPC9zdGF0dXM+DQogICAgICAgICAgICAgICAgPC90ZXN0Pg0KICAgICAgICAgICAgICAgIDx0ZXN0Pg0KICAgICAgICAgICAgICAgICAgICA8bmFtZT5mYWN0b3J5X1JBTS9zaXplPC9uYW1lPg0KICAgICAgICAgICAgICAgICAgICA8c3RhdHVzPnBhc3NlZDwvc3RhdHVzPg0KICAgICAgICAgICAgICAgIDwvdGVzdD4NCiAgICAgICAgICAgICAgICA8dGVzdD4NCiAgICAgICAgICAgICAgICAgICAgPG5hbWU+ZmFjdG9yeV9SVEM8L25hbWU+DQogICAgICAgICAgICAgICAgICAgIDxzdGF0dXM+ZmFpbGVkPC9zdGF0dXM+DQogICAgICAgICAgICAgICAgPC90ZXN0Pg0KICAgICAgICAgICAgICAgIDx0ZXN0Pg0KICAgICAgICAgICAgICAgICAgICA8bmFtZT5mYWN0b3J5X3VzYi91c2IyLXJvb3QtaHViLXByZXNlbnQ8L25hbWU+DQogICAgICAgICAgICAgICAgICAgIDxzdGF0dXM+cGFzc2VkPC9zdGF0dXM+DQogICAgICAgICAgICAgICAgPC90ZXN0Pg0KICAgICAgICAgICAgPC90ZXN0cz4NCjwvdXV0Pg0KPC91dXRzPg0KPC90ZXN0X3JlcG9ydD4="

func (s *LogSuite) TestAPISyncHandler(c *check.C) {
	tests := []SyncTest{
		{"POST", "/api/testlog/1523460528_example_report.xml", []byte(""), 400, response.JSONHeader, datastore.SyncUser, false, false, false, false, 0},
		{"POST", "/api/testlog/1523460528_example_report.xml", []byte("bad"), 400, response.JSONHeader, datastore.SyncUser, false, false, false, false, 0},
		{"POST", "/api/testlog/1523460528_example_report.xml", []byte(exampleFile), 400, response.JSONHeader, 0, false, false, false, false, 0},
		{"POST", "/api/testlog/1523460528_example_report.xml", []byte(exampleFile), 400, response.JSONHeader, datastore.SyncUser, true, false, false, true, 0},
		{"POST", "/api/testlog/example_report.xml", []byte(exampleFile), 200, response.JSONHeader, datastore.SyncUser, true, true, false, false, 0},
		{"POST", "/api/testlog/1523460528_example_report.xml", []byte(exampleFile), 200, response.JSONHeader, datastore.SyncUser, true, true, false, false, 0},
		{"POST", "/api/testlog/1523460528_example_report.xml", []byte(exampleFile), 400, response.JSONHeader, datastore.Standard, true, false, false, false, 0},
		{"POST", "/api/testlog/1523460528_example_report.xml", []byte(exampleFile), 400, response.JSONHeader, 0, true, false, false, false, 0},
	}

	for _, t := range tests {
		if t.EnableAuth {
			datastore.Environ.Config.EnableUserAuth = true
		}
		if t.MockError {
			datastore.Environ.DB = &datastore.ErrorMockDB{}
		}

		w := sendAdminAPIRequest(t.Method, t.URL, bytes.NewReader(t.Data), t.Permissions, c)
		c.Assert(w.Code, check.Equals, t.Code)
		c.Assert(w.Header().Get("Content-Type"), check.Equals, t.Type)

		result, err := response.ParseStandardResponse(w)
		c.Assert(err, check.IsNil)
		c.Assert(result.Success, check.Equals, t.Success)

		datastore.Environ.Config.EnableUserAuth = false
		datastore.Environ.DB = &datastore.MockDB{}
	}
}

func sendAdminAPIRequest(method, url string, data io.Reader, permissions int, c *check.C) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(method, url, data)

	switch permissions {
	case datastore.SyncUser:
		r.Header.Set("user", "sync")
		r.Header.Set("api-key", "ValidAPIKey")
	case datastore.Standard:
		r.Header.Set("user", "user1")
		r.Header.Set("api-key", "ValidAPIKey")
	default:
		break
	}

	service.AdminRouter().ServeHTTP(w, r)

	return w
}
