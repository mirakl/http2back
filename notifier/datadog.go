package notifier

import (
	"errors"
	"fmt"

	"github.com/zorkian/go-datadog-api"
)

var (
	_           Notifier = Datadog{}
	defaultTags          = []string{
		"origin:http2back",
	}
)

type Datadog struct {
	ApiKey    string
	AppKey    string
	ExtraTags []string
}

func (d Datadog) Notify(event *Event) error {
	if d.ApiKey == "" {
		return errors.New("Cannot send event to Datadog, api key is not defined")
	}

	tags := append(defaultTags, d.ExtraTags...)

	client := datadog.NewClient(d.ApiKey, d.AppKey)

	_, err := client.PostEvent(&datadog.Event{
		Title:     &event.Title,
		Text:      &event.Message,
		AlertType: datadog.String("info"),
		Tags:      tags,
	})

	return err
}

func (d Datadog) String() string {
	return fmt.Sprintf("Datadog notifier - ApiKey: %s, AppKey: %s, ExtraTags: %+v", d.ApiKey, d.AppKey, d.ExtraTags)
}
