package nuclei

import (
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/projectdiscovery/nuclei/v2/pkg/model"
	"github.com/projectdiscovery/nuclei/v2/pkg/output"
	"github.com/projectdiscovery/nuclei/v2/pkg/types"
	"github.com/projectdiscovery/nuclei/v2/pkg/utils"
)

func NewWriter() (*Writer, error) {
	writer := &Writer{}
	writer.auroraColorizer = aurora.NewAurora(false)
	writer.ResultEvents = []*output.ResultEvent{}

	return writer, nil
}

type Writer struct {
	auroraColorizer aurora.Aurora
	ResultEvents    []*output.ResultEvent
	FailureEvents   []*output.ResultEvent
}

// Write writes the event to file and/or screen.
func (w *Writer) Write(event *output.ResultEvent) error {
	w.ResultEvents = append(w.ResultEvents, event)
	return nil
}

// Close closes the output writer interface
func (w *Writer) Close() {
	return
}

// Colorizer returns the colorizer instance for writer
func (w *Writer) Colorizer() aurora.Aurora {
	return w.auroraColorizer
}

// WriteFailure writes the optional failure event for template to file and/or screen.
func (w *Writer) WriteFailure(event output.InternalEvent) error {
	templatePath, templateURL := utils.TemplatePathURL(types.ToString(event["template-path"]))
	data := &output.ResultEvent{
		Template:      templatePath,
		TemplateURL:   templateURL,
		TemplateID:    types.ToString(event["template-id"]),
		TemplatePath:  types.ToString(event["template-path"]),
		Info:          event["template-info"].(model.Info),
		Type:          types.ToString(event["type"]),
		Host:          types.ToString(event["host"]),
		MatcherStatus: false,
		Timestamp:     time.Now(),
	}

	w.FailureEvents = append(w.FailureEvents, data)
	return nil
}

// Request logs a request in the trace log
func (w *Writer) Request(templateID, url, requestType string, err error) {
	return
}

// WriteStoreDebugData writes the request/response debug data to file
func (w *Writer) WriteStoreDebugData(host, templateID, eventType string, data string) {}
