package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
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

func resourceConsulKeyPrefixCreateFile(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	kv := client.KV()
	token := d.Get("token").(string)
	dc, err := getDC(d, client)
	if err != nil {
		return err
	}

	keyClient := newKeyClient(kv, dc, token)

	subKeys := map[string]string{
		"a": "do",
		"b": "re",
		"c": "mi",
	}
	pathPrefix := d.Get("path_prefix").(string)

	// subKeys := map[string]string{}
	// for k, vI := range d.Get("subkeys").(map[string]interface{}) {
	// 	subKeys[k] = vI.(string)
	// }
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

	// Ideally we'd use d.Partial(true) here so we can correctly record
	// a partial write, but that mechanism doesn't work for individual map
	// members, so we record that the resource was created before we
	// do anything and that way we can recover from errors by doing an
	// Update on subsequent runs, rather than re-attempting Create with
	// some keys possibly already present.
	d.SetId(pathPrefix)

	// Store the datacenter on this resource, which can be helpful for reference
	// in case it was read from the provider
	d.Set("datacenter", dc)
	//ConsulKeyPrefix := Provider().(*schema.Provider)

	//ConsulProvider := Provider().(*schema.Provider)
	//ConsulProviders := map[string]terraform.ResourceProvider{
	//"consul": ConsulProvider,
	//}
	//s := &terraform.InstanceState{
	//	ID:         "yo",
	//	Attributes: map[string]string{},
	//}
	//d = resourceConsulKeyPrefixFromFile().Data(s)
	//dummy map/var to test

	d.Set("path_prefix", pathPrefix)
	d.Set("subkeys", subKeys)
	d.SetId("consul")
	// resourceConsulKeyPrefixCreate(d, meta)
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
