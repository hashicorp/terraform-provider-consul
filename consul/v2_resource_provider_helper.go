// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package consul

import (
	"fmt"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/proto-public/pbresource"
)

type GVK struct {
	Group   string
	Version string
	Kind    string
}

type V2WriteRequest struct {
	Metadata map[string]string `json:"metadata"`
	Data     map[string]any    `json:"data"`
	Owner    *pbresource.ID    `json:"owner"`
}

type V2WriteResponse struct {
	Metadata   map[string]string `json:"metadata"`
	Data       map[string]any    `json:"data"`
	Owner      *pbresource.ID    `json:"owner,omitempty"`
	ID         *pbresource.ID    `json:"id"`
	Version    string            `json:"version"`
	Generation string            `json:"generation"`
	Status     map[string]any    `json:"status"`
}

func v2MulticlusterRead(client *api.Client, gvk *GVK, resourceName string, q *api.QueryOptions) (map[string]interface{}, error) {
	endpoint := strings.ToLower(fmt.Sprintf("/api/%s/%s/%s/%s", gvk.Group, gvk.Version, gvk.Kind, resourceName))
	var out map[string]interface{}
	_, err := client.Raw().Query(endpoint, &out, q)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func v2MulticlusterDelete(client *api.Client, gvk *GVK, resourceName string, q *api.QueryOptions) error {
	endpoint := strings.ToLower(fmt.Sprintf("/api/%s/%s/%s/%s", gvk.Group, gvk.Version, gvk.Kind, resourceName))
	_, err := client.Raw().Delete(endpoint, q)
	if err != nil {
		return err
	}
	return nil
}

func v2MulticlusterApply(client *api.Client, gvk *GVK, resourceName string, w *api.WriteOptions, payload *V2WriteRequest) (*V2WriteResponse, *api.WriteMeta, error) {
	endpoint := strings.ToLower(fmt.Sprintf("/api/%s/%s/%s/%s", gvk.Group, gvk.Version, gvk.Kind, resourceName))
	out := &V2WriteResponse{}
	wm, err := client.Raw().Write(endpoint, payload, out, w)
	if err != nil {
		return nil, nil, err
	}
	return out, wm, nil
}
