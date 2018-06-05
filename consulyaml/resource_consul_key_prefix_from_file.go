package consulyaml

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
	yaml "gopkg.in/yaml.v2"
)

func resourceConsulKeyPrefixFromFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulKeyPrefixCreateFile,
		Read:   resourceConsulKeyPrefixReadFile,
		Update: resourceConsulKeyPrefixUpdateFile,
		Delete: resourceConsulKeyPrefixDeleteFile,

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"token": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"path_prefix": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"subkeys_file": {
				Type:     schema.TypeString,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func parseMapSlice(slice yaml.MapSlice, subKeys map[string]string, parent string) {
	for _, item := range slice {
		parseMapItem(item, subKeys, parent)
	}
}

func parseMapItem(item yaml.MapItem, subKeys map[string]string, parent string) {
	var key string

	switch k := item.Key.(type) {
	case string:
		// All the keys should be strings
		if parent != "" {
			key = fmt.Sprintf("%s/%s", parent, k)
		} else {
			key = k
		}
	default:
		panic(fmt.Sprintf("Don't know what type key is %T\n", k))
	}

	switch v := item.Value.(type) {
	case string:
		// If we see a string, just use it
		subKeys[key] = v
	case int:
		// if we see an int, convert it to a string
		subKeys[key] = strconv.Itoa(v)
	case yaml.MapSlice:
		parseMapSlice(v, subKeys, key)
	case nil:
		// if value is nil lets set an empty string
		subKeys[key] = ""
	default:
		panic(fmt.Sprintf("Found type of value we don't understand: %T\n", v))
	}
}

// getDC is used to get the datacenter of the local agent
func getDC(d *schema.ResourceData, client *consulapi.Client) (string, error) {
	if v, ok := d.GetOk("datacenter"); ok {
		return v.(string), nil
	}
	info, err := client.Agent().Self()
	if err != nil {
		return "", fmt.Errorf("Failed to get datacenter from Consul agent: %v", err)
	}
	return info["Config"]["Datacenter"].(string), nil
}

func resourceConsulKeyPrefixCreateFile(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	kv := client.KV()
	token := d.Get("token").(string)
	yamlfile := d.Get("subkeys_file").(string)
	dc, err := getDC(d, client)
	if err != nil {
		return err
	}

	keyClient := newKeyClient(kv, dc, token)
	pathPrefix := d.Get("path_prefix").(string)

	//Grab subkeys from YAML file
	f, err := os.Open(yamlfile)

	if err != nil {
		panic(err)
	}

	// A slice of yaml.MapItem structs
	var data yaml.MapSlice

	b, err := ioutil.ReadAll(f)
	err = yaml.Unmarshal(b, &data)

	subKeys := make(map[string]string)
	parseMapSlice(data, subKeys, "")

	currentSubKeys, err := keyClient.GetUnderPrefix(pathPrefix)
	if err != nil {
		return err
	}
	if len(currentSubKeys) > 0 {
		return fmt.Errorf(
			"%d keys already exist under %s; delete them before managing this prefix with Terraform",
			len(currentSubKeys), pathPrefix,
		)
	}

	d.SetId(pathPrefix)
	d.Set("datacenter", dc)
	d.Set("path_prefix", pathPrefix)
	d.Set("subkeys", subKeys)
	d.SetId("consul")

	for k, v := range subKeys {
		fullPath := pathPrefix + k
		err := keyClient.Put(fullPath, v)
		if err != nil {
			return fmt.Errorf("error while writing %s: %s", fullPath, err)
		}
	}

	return nil
}
func resourceConsulKeyPrefixReadFile(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceConsulKeyPrefixUpdateFile(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceConsulKeyPrefixDeleteFile(d *schema.ResourceData, m interface{}) error {
	return nil
}
