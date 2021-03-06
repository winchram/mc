/*
 * Minio Client, (C) 2015 Minio, Inc.
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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/minio/mc/pkg/console"
	"github.com/minio/minio-xl/pkg/probe"
	"github.com/minio/minio-xl/pkg/quick"
)

type sessionV1 struct {
	Version     string          `json:"version"`
	Started     time.Time       `json:"started"`
	CommandType string          `json:"command-type"`
	SessionID   string          `json:"session-id"`
	URLs        []string        `json:"args"`
	Files       map[string]bool `json:"files"`

	Lock *sync.Mutex `json:"-"`
}

func (s sessionV1) String() string {
	message := console.Colorize("Time", fmt.Sprintf("[%s] ", s.Started.Local().Format(printDate)))
	message = message + console.Colorize("SessionID", fmt.Sprintf("%s", s.SessionID))
	message = message + console.Colorize("Command", fmt.Sprintf(" [%s %s]", s.CommandType, strings.Join(s.URLs, " ")))
	return message
}

// loadSession - reads session file if exists and re-initiates internal variables
func loadSessionV1(sid string) (*sessionV1, *probe.Error) {
	if !isSessionDirExists() {
		return nil, probe.NewError(errors.New("Session folder does not exist."))
	}

	sessionFile, err := getSessionFileV1(sid)
	if err != nil {
		return nil, err.Trace(sid)
	}

	s := new(sessionV1)
	s.Version = "1.0.0"
	// map of command and files copied
	s.URLs = nil
	s.Lock = new(sync.Mutex)
	s.Files = make(map[string]bool)
	qs, err := quick.New(s)
	if err != nil {
		return nil, err.Trace(s.Version)
	}
	err = qs.Load(sessionFile)
	if err != nil {
		return nil, err.Trace(sessionFile, s.Version)
	}
	return qs.Data().(*sessionV1), nil
}

func getSessionIDsV1() (sids []string) {
	sessionDir, err := getSessionDir()
	fatalIf(err.Trace(), "Unable to determine session folder.")

	sessionList, e := filepath.Glob(sessionDir + "/*")
	fatalIf(probe.NewError(e), "Unable to access session folder ‘"+sessionDir+"’.")

	for _, path := range sessionList {
		sidReg := regexp.MustCompile("^[a-zA-Z]{8}$")
		sid := filepath.Base(path)
		if sidReg.Match([]byte(sid)) {
			sessionV1, err := loadSessionV1(sid)
			fatalIf(err.Trace(sid), "Unable to load session ‘"+sid+"’.")
			if sessionV1.Version != "1.0.0" {
				continue
			}
			sids = append(sids, sid)
		}
	}
	return sids
}

func getSessionFileV1(sid string) (string, *probe.Error) {
	sessionDir, err := getSessionDir()
	if err != nil {
		return "", err.Trace()
	}

	sessionFile := filepath.Join(sessionDir, sid)
	if _, err := os.Stat(sessionFile); err != nil {
		return "", probe.NewError(err)
	}

	return sessionFile, nil
}
