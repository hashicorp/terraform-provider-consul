package consulyaml

import (
	"fmt"
	"io/ioutil"
	"log"
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

			/*"subkeys_file": {
				Type:     schema.TypeString,
				Required: true,
				/*StateFunc: func(v interface{}) string {
					return yamlhash(v.(string))
				},
			},*/
			"subkeys_file": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeList,
				},
				StateFunc: func(v interface{}) []string {
					yamlfile := string(v.([]string)[:0])
					v = append(v.([]string), yamlhash(yamlfile))
					return v
				},
			},
			"subkeys": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			/*"yaml_hash": {
				Type:     schema.TypeString,
				Optional: true
				Computed: true,
				StateFunc: func(v interface{}) string {
					return yamlhash(v.(string))
				},
			},*/
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
	fileContent, err := ioutil.ReadFile(d.Get("subkeys_file").(string))
	if err != nil {
		return err
	}
	d.Set("yaml_hash", yamlhash(string(fileContent)))

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
	d.Partial(true)
	if err != nil {
		return err
	}
	panic("INSIDE UPDATE")
	llaves := d.Get("subkeys").(map[string]interface{})

	for key, value := range llaves {
		log.Printf(
			"[DEBUG] !!!!!! inside UPDATE func -  Key: -> %v | Value: -> %v", key, value,
		)
	}
	currenthash := d.Get("yaml_hash")
	keyClient := newKeyClient(kv, dc, token)
	log.Printf("current hash: %v", currenthash)

	pathPrefix := d.Id()
	if d.HasChange("yaml_hash") {
		panic(fmt.Sprintf("INSIDE THE HAS CHANGE"))
		o, n := d.GetChange("subkeys")
		if o == nil {
			o = map[string]interface{}{}
		}
		if n == nil {
			n = map[string]interface{}{}
		}

		om := o.(map[string]interface{})
		nm := n.(map[string]interface{})

		for key, value := range om {
			log.Printf(
				"[DEBUG] #### inside UPDATE func - old Key: -> %v", key, "Value: -> %v", value,
			)
		}

		for key, value := range nm {
			log.Printf(
				"[DEBUG] #### inside UPDATE func - new Key: ->  %v", key, "Value: -> %v", value,
			)
		}

		// First we'll write all of the stuff in the "new map" nm,
		// and then we'll delete any keys that appear in the "old map" om
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
			if err != nil {
				return fmt.Errorf("error while writing %s: %s", fullPath, err)
			}
		}

		// Remove deleted keys
		for k, _ := range om {
			if _, exists := nm[k]; exists {
				continue
			}
			fullPath := pathPrefix + k
			err := keyClient.Delete(fullPath)
			if err != nil {
				return fmt.Errorf("error while deleting %s: %s", fullPath, err)
			}
		}

	}

	// Store the datacenter on this resource, which can be helpful for reference
	// in case it was read from the provider
	d.Set("datacenter", dc)
	return nil
}

func resourceConsulKeyPrefixReadFile(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	kv := client.KV()
	token := d.Get("token").(string)
	//yamlfile := d.Get("subkeys_file").(string)
	dc, err := getDC(d, client)
	if err != nil {
		return err
	}

	keyClient := newKeyClient(kv, dc, token)
	pathPrefix := d.Id()
	subKeys, err := keyClient.GetUnderPrefix(pathPrefix)
	if err != nil {
		return err
	}
	for key, value := range subKeys {
		log.Printf(
			"[DEBUG] #### inside READ func - Key: ->  %v Value: ->   %v", key, value,
		)
	}
	fileContent, err := ioutil.ReadFile(d.Get("subkeys_file").(string))
	if err != nil {
		return err
	}
	d.Set("yaml_hash", yamlhash(string(fileContent)))
	d.Set("subkeys", subKeys)

	// Store the datacenter on this resource, which can be helpful for reference
	// in case it was read from the provider
	d.Set("datacenter", dc)

	return nil
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
