package elasticsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/pegasus-cloud/metadata_client/metadata/common"
	"github.com/pegasus-cloud/metadata_client/metadata/utility"
)

const (
	healthEnum = iota
	warnEnum
	errorEnum
	unknownEnum
)

type cluster struct {
	Nodes struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
	} `json:"_nodes"`
}

// HealthCheck returns whether the elasticsearch cluster is healthy or not.
func (p *Provider) HealthCheck() (status string, err error) {
	url := fmt.Sprintf("http://%s/_cluster/stats", p.Endpoint)
	cluster := &cluster{}

	// Get cluster status from Elasticsearch
	resp, statusCode, err := utility.SendRequest(http.MethodGet, url, headers, nil)
	if err != nil {
		return common.HealthUnknown, err
	}
	if statusCode != http.StatusOK {
		return common.HealthUnknown, errors.New(common.StatusCodeIsNotOK)
	}

	// Convert response body
	json.Unmarshal(resp, cluster)
	if cluster.Nodes.Successful == cluster.Nodes.Total {
		return common.HealthOK, nil
	} else if (cluster.Nodes.Successful < cluster.Nodes.Total) && (cluster.Nodes.Successful != 0) {
		return common.HealthWarn, nil
	} else if cluster.Nodes.Successful == 0 {
		return common.HealthError, nil
	}
	return common.HealthUnknown, nil
}
