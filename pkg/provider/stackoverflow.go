package provider

import (
	"encoding/json"
	"fmt"
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

// EmitEvent fires an event if an answer for an error was found.
func (s *StackOverflowProvider) EmitEvent(event v1.Event) error {
	if event.Type == corev1.EventTypeWarning && event.Message != "" {
		log.Infof("found event with warning: event %s with reason %s", event.Type, event.Message)
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

		body, err := ioutil.ReadAll(resp.Body)
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
			log.Infof("Fired event for object %+v", event.InvolvedObject)
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
