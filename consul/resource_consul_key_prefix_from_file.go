package consul

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceConsulKeyPrefixFromFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulKeyPrefixCreateFile,
		Read:   resourceConsulKeyPrefixReadFile,
		Update: resourceConsulKeyPrefixUpdateFile,
		Delete: resourceConsulKeyPrefixDeleteFile,

		Schema: map[string]*schema.Schema{
			"datacenter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"token": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"path_prefix": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"subkeys_file": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceConsulKeyPrefixCreateFile(d *schema.ResourceData, m interface{}) error {
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
	subkeys := map[string]string{
		"a": "do",
		"b": "re",
		"c": "mi",
	}
	prefix := d.Get("path_prefix").(string)
	d.Set("path_prefix", prefix)
	d.Set("subkeys", subkeys)
	d.SetId("consul")
	resourceConsulKeyPrefixCreate(d, m)

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
