package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type policyProvider struct {
	client HttpClient
	url    string
}

const (
	evaluationIdHeader = "x-evaluation-id"
	cacheTtlHeader     = "x-cache-ttl"
)

type EvaluateResult struct {
	// Arbitrary JSON response.
	Result any
	// ETag contains unique identifier of the policy evaluation and can be used to
	// later retrieve the results from Cache.
	ETag string
}

func (p *policyProvider) Evaluate(name string, data interface{}, version string, predefinedId string) (*EvaluateResult, error) {
	var res EvaluateResult
	url := strings.Join([]string{p.url, libConfig.Policy.Repository, libConfig.Policy.Group, name, version}, "/")
	body, err := json.Marshal(data)
	if err != nil {
		return nil, SdkError{
			general:  "failed to evaluate policy",
			specific: err,
		}
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, SdkError{
			general:  "failed to evaluate policy",
			specific: err,
		}
	}
	if libConfig.Policy.ExpiresInSeconds > 0 {
		req.Header.Add(cacheTtlHeader, strconv.Itoa(libConfig.Policy.ExpiresInSeconds))
	}
	if predefinedId != "" {
		req.Header.Add(evaluationIdHeader, predefinedId)
	}
	response, err := p.client.Do(req)
	if err != nil {
		return nil, SdkError{
			general:  "failed to evaluate policy",
			specific: err,
		}
	}
	if response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices {
		resBody, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, SdkError{
				general:  "failed to evaluate policy",
				specific: err,
			}
		}
		err = json.Unmarshal(resBody, &res)
		if err != nil {
			return nil, SdkError{
				general:  "failed to evaluate policy",
				specific: err,
			}
		}
		return &res, nil
	} else {
		err = fmt.Errorf("evauluation request is not successful: %s", response.Status)
		return nil, SdkError{
			general:  "failed to evaluate policy",
			specific: err,
		}
	}
}

func getPolicyProvider(client HttpClient) *policyProvider {
	return &policyProvider{client: client, url: libConfig.Policy.Url}
}
