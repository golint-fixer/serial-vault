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

package testlog

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CanonicalLtd/serial-vault/datastore"
	"github.com/CanonicalLtd/serial-vault/service/auth"
	"github.com/CanonicalLtd/serial-vault/service/response"
)

// ListResponse is the JSON response from the API Test Log method
type ListResponse struct {
	Success      bool                `json:"success"`
	ErrorCode    string              `json:"error_code"`
	ErrorSubcode string              `json:"error_subcode"`
	ErrorMessage string              `json:"message"`
	TestLog      []datastore.TestLog `json:"logs"`
}

// listHandler is the API method to fetch the log records from signing
func listHandler(w http.ResponseWriter, user datastore.User, apiCall bool) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := auth.CheckUserPermissions(user, datastore.SyncUser, apiCall)
	if err != nil {
		response.FormatStandardResponse(false, "error-auth", "", "", w)
		return
	}

	logs, err := datastore.Environ.DB.ListAllowedTestLog(user)
	if err != nil {
		response.FormatStandardResponse(false, "error-fetch-signinglog", "", err.Error(), w)
		return
	}

	// Return successful JSON response with the list of models
	w.WriteHeader(http.StatusOK)
	formatListResponse(true, "", "", "", logs, w)
}

func formatListResponse(success bool, errorCode, errorSubcode, message string, logs []datastore.TestLog, w http.ResponseWriter) error {
	response := ListResponse{Success: success, ErrorCode: errorCode, ErrorSubcode: errorSubcode, ErrorMessage: message, TestLog: logs}

	// Encode the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error forming the signing log response.")
		return err
	}
	return nil
}
