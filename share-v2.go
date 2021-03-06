/*
 * Minio Client (C) 2014, 2015 Minio, Inc.
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

package main

import (
	"os"
	"time"

	"github.com/minio/minio-xl/pkg/probe"
	"github.com/minio/minio-xl/pkg/quick"
)

type sharedURLsV2 struct {
	Version string
	URLs    []struct {
		Date    time.Time
		Message ShareMessageV2
	}
}

func loadSharedURLsV2() (*sharedURLsV2, *probe.Error) {
	sharedURLsDataFile, err := getSharedURLsDataFile()
	if err != nil {
		return nil, err.Trace()
	}
	if _, err := os.Stat(sharedURLsDataFile); err != nil {
		return nil, probe.NewError(err)
	}

	qs, err := quick.New(newSharedURLsV2())
	if err != nil {
		return nil, err.Trace()
	}
	err = qs.Load(sharedURLsDataFile)
	if err != nil {
		return nil, err.Trace(sharedURLsDataFile)
	}
	s := qs.Data().(*sharedURLsV2)
	return s, nil
}

func saveSharedURLsV2(s *sharedURLsV2) *probe.Error {
	qs, err := quick.New(s)
	if err != nil {
		return err.Trace()
	}
	sharedURLsDataFile, err := getSharedURLsDataFile()
	if err != nil {
		return err.Trace()
	}
	return qs.Save(sharedURLsDataFile).Trace(sharedURLsDataFile)
}

func newSharedURLsV2() *sharedURLsV2 {
	var urls []struct {
		Date    time.Time
		Message ShareMessageV2
	}
	s := &sharedURLsV2{
		Version: "1.1.0",
		URLs:    urls,
	}
	return s
}
