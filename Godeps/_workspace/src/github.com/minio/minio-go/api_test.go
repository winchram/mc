/*
 * Minimal object storage library (C) 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func ExampleGetPartSize() {
	fmt.Println(GetPartSize(5000000000))
	// Output: 5242880
}
func ExampleGetPartSize_second() {
	fmt.Println(GetPartSize(50000000000000000))
	// Output: 5368709120
}

func TestACLTypes(t *testing.T) {
	want := map[string]bool{
		"private":            true,
		"public-read":        true,
		"public-read-write":  true,
		"authenticated-read": true,
		"invalid":            false,
	}
	for acl, ok := range want {
		if BucketACL(acl).isValidBucketACL() != ok {
			t.Fatal("Error")
		}
	}
}

func TestMustGetEndpoint(t *testing.T) {
	conf := new(Config)
	conf.Region = "us-east-1"
	if conf.MustGetEndpoint() != "https://s3.amazonaws.com" {
		t.Fatalf("Error")
	}

	conf.Region = ""
	conf.Endpoint = "http://localhost:9000"
	if conf.MustGetEndpoint() != "http://localhost:9000" && conf.Region == "milkyway" {
		t.Fatalf("Error")
	}

	conf.Region = ""
	conf.Endpoint = "https://s3.amazonaws.com"
	if conf.MustGetEndpoint() != "https://s3.amazonaws.com" && conf.Region == "us-east-1" {
		t.Fatalf("Error")
	}

	conf.Region = "invalid"
	conf.Endpoint = "http://play.minio.io:9000"
	if conf.MustGetEndpoint() != "http://play.minio.io:9000" {
		t.Fatalf("Error")
	}

	conf.Region = ""
	conf.Endpoint = ""
	if conf.MustGetEndpoint() != "https://s3.amazonaws.com" && conf.Region == "us-east-1" {
		t.Fatalf("Error")
	}
}

func TestUserAgent(t *testing.T) {
	conf := new(Config)
	conf.AddUserAgent("minio", "1.0", "amd64")
	if !strings.Contains(conf.userAgent, "minio") {
		t.Fatalf("Error")
	}
}

func TestBucketOperations(t *testing.T) {
	bucket := bucketHandler(bucketHandler{
		resource: "/bucket",
	})
	server := httptest.NewServer(bucket)
	defer server.Close()

	a := New(&Config{Endpoint: server.URL})
	err := a.MakeBucket("bucket", "private", "")
	if err != nil {
		t.Errorf("Error")
	}

	err = a.BucketExists("bucket")
	if err != nil {
		t.Errorf("Error")
	}

	err = a.BucketExists("bucket1")
	if err == nil {
		t.Errorf("Error")
	}
	if err.Error() != "403 Forbidden" {
		t.Errorf("Error")
	}

	err = a.SetBucketACL("bucket", "public-read-write")
	if err != nil {
		t.Errorf("Error")
	}

	acl, err := a.GetBucketACL("bucket")
	if err != nil {
		t.Errorf("Error")
	}
	if !acl.isPrivate() {
		t.Fatalf("Error")
	}

	for b := range a.ListBuckets() {
		if b.Err != nil {
			t.Fatalf(b.Err.Error())
		}
		if b.Data.Name != "bucket" {
			t.Errorf("Error")
		}
	}

	for o := range a.ListObjects("bucket", "", true) {
		if o.Err != nil {
			t.Fatalf(o.Err.Error())
		}
		if o.Data.Key != "object" {
			t.Errorf("Error")
		}
	}

	err = a.RemoveBucket("bucket")
	if err != nil {
		t.Errorf("Error")
	}

	err = a.RemoveBucket("bucket1")
	if err == nil {
		t.Fatalf("Error")
	}
	if err.Error() != "404 Not Found" {
		t.Errorf("Error")
	}
}

func TestObjectOperations(t *testing.T) {
	object := objectHandler(objectHandler{
		resource: "/bucket/object",
		data:     []byte("Hello, World"),
	})
	server := httptest.NewServer(object)
	defer server.Close()

	a := New(&Config{Endpoint: server.URL})
	data := []byte("Hello, World")
	err := a.PutObject("bucket", "object", uint64(len(data)), bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Error")
	}
	metadata, err := a.StatObject("bucket", "object")
	if err != nil {
		t.Fatalf("Error")
	}
	if metadata.Key != "object" {
		t.Fatalf("Error")
	}
	if metadata.ETag != "9af2f8218b150c351ad802c6f3d66abe" {
		t.Fatalf("Error")
	}

	reader, metadata, err := a.GetObject("bucket", "object", 0, 0)
	if err != nil {
		t.Fatalf("Error")
	}
	if metadata.Key != "object" {
		t.Fatalf("Error")
	}
	if metadata.ETag != "9af2f8218b150c351ad802c6f3d66abe" {
		t.Fatalf("Error")
	}

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, reader)
	if !bytes.Equal(buffer.Bytes(), data) {
		t.Fatalf("Error")
	}

	err = a.RemoveObject("bucket", "object")
	if err != nil {
		t.Fatalf("Error")
	}
}

func TestPartSize(t *testing.T) {
	var maxPartSize uint64 = 1024 * 1024 * 1024 * 5
	partSize := GetPartSize(5000000000000000000)
	if partSize > MinimumPartSize {
		if partSize > maxPartSize {
			t.Fatal("invalid result, cannot be bigger than MaxPartSize 5GB")
		}
	}
	partSize = GetPartSize(50000000000)
	if partSize > MinimumPartSize {
		t.Fatal("invalid result, cannot be bigger than MinimumPartSize 5MB")
	}
}

func TestURLEncoding(t *testing.T) {
	type urlStrings struct {
		name        string
		encodedName string
	}

	want := []urlStrings{
		{
			name:        "bigfile-1._%",
			encodedName: "bigfile-1._%25",
		},
		{
			name:        "本語",
			encodedName: "%E6%9C%AC%E8%AA%9E",
		},
		{
			name:        "本b語.1",
			encodedName: "%E6%9C%ACb%E8%AA%9E.1",
		},
		{
			name:        ">123>3123123",
			encodedName: "%3E123%3E3123123",
		},
		{
			name:        "test 1 2.txt",
			encodedName: "test%201%202.txt",
		},
	}

	for _, u := range want {
		encodedName, err := urlEncodeName(u.name)
		if err != nil {
			t.Fatalf("Error")
		}
		if u.encodedName != encodedName {
			t.Errorf("Error")
		}
	}
}
