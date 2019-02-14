// Copyright 2018 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/google/keytransparency/core/testdata"
	"github.com/google/trillian/types"

	pb "github.com/google/keytransparency/core/api/v1/keytransparency_go_proto"
)

// Test vectors in core/testdata are generated by running
// go generate ./core/testdata
func TestVerifyGetUserResponse(t *testing.T) {
	directoryFile := "../testdata/directory.json"
	f, err := os.Open(directoryFile)
	if err != nil {
		t.Fatalf("ReadFile(%v): %v", directoryFile, err)
	}
	defer f.Close()
	var directoryPB pb.Directory
	if err := jsonpb.Unmarshal(f, &directoryPB); err != nil {
		t.Fatalf("jsonpb.Unmarshal(): %v", err)
	}
	v, err := NewVerifierFromDirectory(&directoryPB)
	if err != nil {
		t.Fatal(err)
	}

	respFile := "../testdata/getentryresponse.json"
	b, err := ioutil.ReadFile(respFile)
	if err != nil {
		t.Fatalf("ReadFile(%v): %v", respFile, err)
	}
	var getUserResponses []testdata.ResponseVector
	if err := json.Unmarshal(b, &getUserResponses); err != nil {
		t.Fatalf("Unmarshal(): %v", err)
	}

	trusted := &types.LogRootV1{}
	for _, tc := range getUserResponses {
		t.Run(tc.Desc, func(t *testing.T) {
			slr, smr, err := v.VerifyRevision(tc.GetUserResp.Revision, *trusted)
			if err != nil {
				t.Errorf("VerifyRevision(): %v", err)
			}
			if err == nil && tc.TrustNewLog {
				trusted = slr
			}
			if err := v.VerifyMapLeaf(directoryPB.DirectoryId, tc.UserIDs[0],
				tc.GetUserResp.Leaf, smr); err != nil {
				t.Errorf("VerifyMapLeaf(): %v)", err)
			}
		})
	}
}
