package consul

import (
	"fmt"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	requestOptions = "request_options"

	requestOptDatacenter = "datacenter"
	requestOptToken      = "token"
)

var schemaRequestOpts = &schema.Schema{
	Optional: true,
	Type:     schema.TypeSet,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			requestOptDatacenter: {
				// Optional because we'll pull the default from the local agent if it's
				// not specified, but we can query remote data centers as a result.
				Optional: true,
				Type:     schema.TypeString,
			},
			requestOptToken: {
				Optional:  true,
				Type:      schema.TypeString,
				Sensitive: true,
			},
		},
	},
}

func getRequestOpts(d *schema.ResourceData, client *consulapi.Client) (*consulapi.WriteOptions, *consulapi.QueryOptions, error) {
	requestOpts := &consulapi.WriteOptions{}
	queryOpts := &consulapi.QueryOptions{}

	if v, ok := d.GetOk(requestOptions); ok {
		options := v.(*schema.Set).List()
		if len(options) > 0 {
			opts := options[0].(map[string]interface{})

			if token, ok := opts[requestOptToken]; ok {
				requestOpts.Token = token.(string)
				queryOpts.Token = token.(string)
			}

			if dc, ok := opts[requestOptDatacenter]; ok {
				requestOpts.Datacenter = dc.(string)
				queryOpts.Datacenter = dc.(string)
			}
		}
	}

	if requestOpts.Datacenter == "" {
		dc, err := getLocalDC(client)
		if err != nil {
			return nil, nil, err
		}
		requestOpts.Datacenter = dc
		queryOpts.Datacenter = dc
	}

	return requestOpts, queryOpts, nil
}

func getLocalDC(client *consulapi.Client) (string, error) {
	info, err := client.Agent().Self()
	if err != nil {
		// Reading can fail with `Unexpected response code: 403 (Permission denied)`
		// if the permission has not been given. Default to "" in this case.
		if strings.HasSuffix(err.Error(), "403 (Permission denied)") {
			return "", nil
		}
		return "", fmt.Errorf("Failed to get datacenter from Consul agent: %v", err)
	}
	return info["Config"]["Datacenter"].(string), nil
}
