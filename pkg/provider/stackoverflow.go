// Copyright 2022 Andreas Fritzler
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	otypes "github.com/afritzler/oopsie/pkg/types"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
)

const (
	ProviderName     = "StackOverflow"
	stackOverFlowAPI = "https://api.stackexchange.com/2.2/search"
)

// StackOverflowProvider defines the StackOverflow provider.
type StackOverflowProvider struct {
	Recorder record.EventRecorder
	Log      logr.Logger
}

func (s *StackOverflowProvider) GetName() string {
	return ProviderName
}

// EmitEvent fires an event if an answer for an error was found.
func (s *StackOverflowProvider) EmitEvent(event v1.Event) error {
	if event.Type == corev1.EventTypeWarning && event.Message != "" {
		s.Log.Info("found warning event event", "Event", event.Type, "Message", event.Message)
		req, err := http.NewRequest("GET", stackOverFlowAPI, nil)
		if err != nil {
			return fmt.Errorf("failed to construct request: %s", err)
		}
		req.URL.RawQuery = s.constructQuery(event.Message, req)

		resp, err := http.Get(req.URL.String())
		if err != nil {
			return fmt.Errorf("failed to query backend: %s", err)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to get response body: %s", err)
		}

		var anyJSON otypes.StackOverflowAnswers
		if err := json.Unmarshal(body, &anyJSON); err != nil {
			return fmt.Errorf("failed to unmarshal response json: %s", err)
		}

		if len(anyJSON.Items) > 0 && anyJSON.Items[0].Link != "" {
			link := anyJSON.Items[0].Link
			errorHint := fmt.Sprintf("For error '%s' I found something here -> %s", event.Message, link)
			s.Log.Info("fired event for object", "Object", event.InvolvedObject)
			s.Recorder.Event(&event.InvolvedObject, v1.EventTypeNormal, "Hint", errorHint)
		}
	}
	return nil
}

// constructQuery contructs a RAW request query string
func (s *StackOverflowProvider) constructQuery(message string, req *http.Request) string {
	q := req.URL.Query()
	q.Add("order", "desc")
	q.Add("sort", "votes")
	q.Add("intitle", message)
	q.Add("site", "stackoverflow")
	return q.Encode()
}
