package embedqueries

import (
	"bytes"
	"encoding/json"
	"net/http"
	"text/template"
	"time"

	"github.com/pkg/errors"
)

const (
	url         = "https://api.spacex.land/graphql/"
	contentType = "application/json"
)

// Client models a SpaceX API Client
type Client interface {
	// MissionByID obtains a mission by ID
	MissionByID(missionID string) ([]byte, error)

	// MissionsByManufacturer obtains all missions for the specified manufacturer
	MissionsByManufacturer(manufacturer string, limit int) ([]byte, error)

	// PastLaunches obtains the last x launches
	PastLaunches(limit int) ([]byte, error)

	// Rockets obtains rockets
	Rockets(limit int) ([]byte, error)
}

type client struct {
	queryStore map[string]*template.Template
	httpClient *http.Client
}

// NewClient constructs a new Client
func NewClient() (Client, error) {
	queryStore, err := newStore()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build query store")
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1,
		},
		Timeout: time.Second * 10,
	}

	return &client{
		queryStore: queryStore,
		httpClient: httpClient,
	}, nil
}

func (c client) MissionByID(missionID string) ([]byte, error) {
	queryTemplate := c.queryStore["mission_by_id"]

	params := struct {
		MissionID string
	}{
		MissionID: missionID,
	}

	return c.submit(queryTemplate, params)
}

func (c client) MissionsByManufacturer(manufacturer string, limit int) ([]byte, error) {
	queryTemplate := c.queryStore["missions_by_manufacturer"]

	// build the template params
	params := struct {
		Manufacturer string
		Limit        int
	}{
		Manufacturer: manufacturer,
		Limit:        limit,
	}

	return c.submit(queryTemplate, params)
}

func (c client) PastLaunches(limit int) ([]byte, error) {
	queryTemplate := c.queryStore["past_launches"]

	// build the template params
	params := struct {
		Limit int
	}{
		Limit: limit,
	}

	return c.submit(queryTemplate, params)
}

func (c client) Rockets(limit int) ([]byte, error) {
	queryTemplate := c.queryStore["rockets"]

	// build the template params
	params := struct {
		Limit int
	}{
		Limit: limit,
	}

	return c.submit(queryTemplate, params)
}

func (c client) submit(queryTemplate *template.Template, params interface{}) ([]byte, error) {

	// invoke the template
	var tmplResult bytes.Buffer
	err := queryTemplate.Execute(&tmplResult, params)
	if err != nil {
		return nil, errors.Wrap(err, "Template execution failed")
	}

	// issue the query
	return c.post(&tmplResult)
}

// graqhQLPostBody is a json struct that is necessary to wrap the graphQL query
type graqhQLPostBody struct {
	Query string `json:"query"`
}

func (c client) post(executedTemplate *bytes.Buffer) ([]byte, error) {

	// Add the json wrapper.
	// This is non-optimally writing to a string here, before serializing later again but this will do for ease of code
	graqhQLPostBody := &graqhQLPostBody{
		Query: executedTemplate.String(),
	}
	b, err := json.Marshal(graqhQLPostBody)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal body")
	}

	resp, err := c.httpClient.Post(url, contentType, bytes.NewReader(b))
	if err != nil {
		return nil, errors.Wrap(err, "POST failed")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read body failed")
	}

	return buf.Bytes(), nil
}
