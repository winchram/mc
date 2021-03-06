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
	"errors"

	"github.com/minio/minio-xl/pkg/probe"
)

var (
	errDummy = func() *probe.Error {
		return probe.NewError(errors.New("")).Untrace()
	}

	errInvalidArgument = func() *probe.Error {
		return probe.NewError(errors.New("Invalid arguments provided, cannot proceed.")).Untrace()
	}

	errSourceListEmpty = func() *probe.Error {
		return probe.NewError(errors.New("Source argument list is empty.")).Untrace()
	}

	errInvalidGlobURL = func(glob, request string) *probe.Error {
		return probe.NewError(errors.New("Error reading glob URL ‘" + glob + "’ while comparing with ‘" + request + "’.")).Untrace()
	}

	errNoMatchingHost = func(URL string) *probe.Error {
		return probe.NewError(errors.New("No matching host found for the given URL ‘" + URL + "’.")).Untrace()
	}

	errInitClient = func(URL string) *probe.Error {
		return probe.NewError(errors.New("Unable to initialize client for URL ‘" + URL + "’.")).Untrace()
	}

	errInvalidSource = func(URL string) *probe.Error {
		return probe.NewError(errors.New("Invalid source ‘" + URL + "’.")).Untrace()
	}

	errInvalidTarget = func(URL string) *probe.Error {
		return probe.NewError(errors.New("Invalid target ‘" + URL + "’.")).Untrace()
	}

	errSourceNotRecursive = func(URL string) *probe.Error {
		return probe.NewError(errors.New("Source ‘" + URL + "’ is not recursive.")).Untrace()
	}

	errSourceIsNotDir = func(URL string) *probe.Error {
		return probe.NewError(errors.New("Source ‘" + URL + "’ is not a folder.")).Untrace()
	}
	errSourceIsDir = func(URL string) *probe.Error {
		return probe.NewError(errors.New("Source ‘" + URL + "’ is a folder.")).Untrace()
	}
	errNotAnObject = func(URL string) *probe.Error {
		return probe.NewError(errors.New("‘" + URL + "’ is not an object.")).Untrace()
	}
)
