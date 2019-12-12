// Copyright © 2018 Andreas Fritzler <andreas.fritzler@gmail.com>
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

package provider

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	otypes "github.com/afritzler/oopsie/pkg/types"
	corev1 "k8s.io/api/core/v1"

	"github.com/prometheus/common/log"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
)

const (
	// ProviderName defines the name of the provider.
	ProviderName     = "StackOverflow"
	stackOverFlowAPI = "https://api.stackexchange.com/2.2/search"
)

// StackOverflowProvider defines the StackOverflow provider.
type StackOverflowProvider struct {
	Recorder record.EventRecorder
}

var _ Provider = &StackOverflowProvider{}

// EmitEvent fires an event if an answer for an error was found.
func (s *StackOverflowProvider) EmitEvent(event v1.Event) error {
	if event.Type == corev1.EventTypeWarning && event.Message != "" {
		log.Infof("found event with warning: event %s with reason %s", event.Type, event.Message)
		req, err := http.NewRequest("GET", stackOverFlowAPI, nil)
		if err != nil {
			log.Error(err, "failed to construct request")
			return err
		}
		q := req.URL.Query()
		q.Add("order", "desc")
		q.Add("sort", "votes")
		q.Add("intitle", event.Message)
		q.Add("site", "stackoverflow")
		req.URL.RawQuery = q.Encode()

		resp, err := http.Get(req.URL.String())
		if err != nil {
			log.Error(err, "failed to query backend")
			return err
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err, "failed to get response body")
			return err
		}

		var anyJSON otypes.StackOverflowAnswers
		json.Unmarshal(body, &anyJSON)

		if len(anyJSON.Items) > 0 && anyJSON.Items[0].Link != "" {
			link := anyJSON.Items[0].Link
			log.Info("Fired event for object %v", event.InvolvedObject)
			s.Recorder.Event(&event.InvolvedObject, v1.EventTypeNormal, "Hint", link)
		}
	}
	return nil
}