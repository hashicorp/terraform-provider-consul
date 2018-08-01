package consulyaml

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

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
				Optional: true,
				Computed: true,
				StateFunc: func(v interface{}) string {
					value := v.(string)
					fileContent, err := ioutil.ReadFile(value)
					if err != nil {
						fmt.Printf("Error: YAML file -> %+v \n", value)
						dir, err := os.Getwd()
						if err != nil {
							log.Fatal(err)
						}
						fmt.Println(dir)
						panic("Error Reading yaml file ")
					}
					hashvalue := yamlhash(string(fileContent))
					value += string(":") + hashvalue
					return value
				},
			},
			"subkeys": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
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

func resourceConsulKeyPrefixCreateFile(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	kv := client.KV()
	token := d.Get("token").(string)
	yamlInfo := d.Get("subkeys_file").(string)
	yamlfile := strings.Split(yamlInfo, ":")[0]
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
	fileContent, err := ioutil.ReadFile(yamlfile)
	if err != nil {
		return err
	}
	value := yamlfile
	value += string(":") + yamlhash(string(fileContent))
	d.Set("subkeys_file", value)

	for k, v := range subKeys {
		fullPath := pathPrefix + k
		err := keyClient.Put(fullPath, v)
		if err != nil {
			return fmt.Errorf("error while writing %s: %s", fullPath, err)
		}
	}

	return nil
}

func resourceConsulKeyPrefixUpdateFile(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	kv := client.KV()
	token := d.Get("token").(string)
	dc, err := getDC(d, client)
	yamlInfo := d.Get("subkeys_file").(string)
	yamlfile := strings.Split(yamlInfo, ":")[0]
	if err != nil {
		return err
	}

	keyClient := newKeyClient(kv, dc, token)

	pathPrefix := d.Id()
	if d.HasChange("subkeys_file") {
		n := d.Get("subkeys")
		if n == nil {
			n = map[string]interface{}{}
		}
		nm := n.(map[string]interface{})
		// Grabing subkeys from Consul
		consulsubKeys, err := keyClient.GetUnderPrefix(pathPrefix)
		if err != nil {
			return err
		}

		// First we'll write all of the stuff in the "new map" nm,
		// and then we'll delete any keys that appear in the "consulsubKeys map"
		// and do not also appear in nm. This ordering means that if a subkey
		// name is changed we will briefly have both the old and new names in
		// Consul, as opposed to briefly having neither.

		// Again, we'd ideally use d.Partial(true) here but it doesn't work
		// for maps and so we'll just rely on a subsequent Read to tidy up
		// after a partial write.

		// Write new and changed keys
		for k, vI := range nm {
			v := vI.(string)
			fullPath := pathPrefix + k
			err := keyClient.Put(fullPath, v)
			log.Printf("Adding %s: ", fullPath)
			if err != nil {
				return fmt.Errorf("error while writing %s: %s", fullPath, err)
			}
		}

		// Remove deleted keys
		for k, _ := range consulsubKeys {
			if _, exists := nm[k]; exists {
				continue
			}
			fullPath := pathPrefix + k
			err := keyClient.Delete(fullPath)
			if err != nil {
				return fmt.Errorf("error while deleting %s: %s", fullPath, err)
			}
		}
		fileContent, err := ioutil.ReadFile(yamlfile)
		if err != nil {
			return err
		}
		value := yamlfile
		value += string(":") + yamlhash(string(fileContent))
		d.Set("subkeys_file", value)

	}

	// Store the datacenter on this resource, which can be helpful for reference
	// in case it was read from the provider
	d.Set("datacenter", dc)
	return nil
}

func resourceConsulKeyPrefixReadFile(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	yamlInfo := d.Get("subkeys_file").(string)
	yamlfile := strings.Split(yamlInfo, ":")[0]
	dc, err := getDC(d, client)
	if err != nil {
		return err
	}

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

	d.Set("subkeys", subKeys)
	for key, value := range subKeys {
		log.Printf(
			"[DEBUG] #### inside READ func - Key: ->  %v Value: ->   %v", key, value,
		)
	}
	// Store the datacenter on this resource, which can be helpful for reference
	// in case it was read from the provider
	d.Set("datacenter", dc)

	return resourceConsulKeyPrefixUpdateFile(d, meta)
}

func resourceConsulKeyPrefixDeleteFile(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	kv := client.KV()
	token := d.Get("token").(string)
	dc, err := getDC(d, client)
	if err != nil {
		return err
	}

	keyClient := newKeyClient(kv, dc, token)

	pathPrefix := d.Id()

	// Delete everything under our prefix, since the entire set of keys under
	// the given prefix is considered to be managed exclusively by Terraform.
	err = keyClient.DeleteUnderPrefix(pathPrefix)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
