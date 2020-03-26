package consul

import (
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

// keyClient is a wrapper around the upstream Consul client that is
// specialized for Terraform's manipulations of the key/value store.
type keyClient struct {
	client *consulapi.KV
	qOpts  *consulapi.QueryOptions
	wOpts  *consulapi.WriteOptions
}

func newKeyClient(realClient *consulapi.KV, dc, token, namespace string) *keyClient {
	qOpts := &consulapi.QueryOptions{
		Datacenter: dc,
		Token:      token,
		Namespace:  namespace,
	}
	wOpts := &consulapi.WriteOptions{
		Datacenter: dc,
		Token:      token,
		Namespace:  namespace,
	}

	return &keyClient{
		client: realClient,
		qOpts:  qOpts,
		wOpts:  wOpts,
	}
}

func (c *keyClient) Get(path string) (string, int, error) {
	log.Printf(
		"[DEBUG] Reading key '%s' in %s",
		path, c.qOpts.Datacenter,
	)
	pair, _, err := c.client.Get(path, c.qOpts)
	if err != nil {
		return "", 0, fmt.Errorf("Failed to read Consul key '%s': %s", path, err)
	}
	value := ""
	if pair != nil {
		value = string(pair.Value)
	}
	flags := 0
	if pair != nil {
		flags = int(pair.Flags)
	}
	return value, flags, nil
}

func (c *keyClient) GetUnderPrefix(pathPrefix string) (consulapi.KVPairs, error) {
	log.Printf(
		"[DEBUG] Listing keys under '%s' in %s",
		pathPrefix, c.qOpts.Datacenter,
	)
	pairs, _, err := c.client.List(pathPrefix, c.qOpts)
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to list Consul keys under prefix '%s': %s", pathPrefix, err,
		)
	}
	return pairs, nil
}

func (c *keyClient) Put(path, value string, flags int) error {
	log.Printf(
		"[DEBUG] Setting key '%s' to '%v' in %s",
		path, value, c.wOpts.Datacenter,
	)
	pair := consulapi.KVPair{Key: path, Value: []byte(value), Flags: uint64(flags)}
	if _, err := c.client.Put(&pair, c.wOpts); err != nil {
		return fmt.Errorf("Failed to write Consul key '%s': %s", path, err)
	}
	return nil
}

func (c *keyClient) Delete(path string) error {
	log.Printf(
		"[DEBUG] Deleting key '%s' in %s",
		path, c.wOpts.Datacenter,
	)
	if _, err := c.client.Delete(path, c.wOpts); err != nil {
		return fmt.Errorf("Failed to delete Consul key '%s': %s", path, err)
	}
	return nil
}

func (c *keyClient) DeleteUnderPrefix(pathPrefix string) error {
	log.Printf(
		"[DEBUG] Deleting all keys under prefix '%s' in %s",
		pathPrefix, c.wOpts.Datacenter,
	)
	if _, err := c.client.DeleteTree(pathPrefix, c.wOpts); err != nil {
		return fmt.Errorf("Failed to delete Consul keys under '%s': %s", pathPrefix, err)
	}
	return nil
}
