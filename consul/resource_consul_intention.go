package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceConsulIntention() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulIntentionCreate,
		Update: resourceConsulIntentionUpdate,
		Read:   resourceConsulIntentionRead,
		Delete: resourceConsulIntentionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		DeprecationMessage: `The consul_intention resource is deprecated in favor of the consul_config_entry resource.
Please see https://registry.terraform.io/providers/hashicorp/consul/latest/docs/guides/upgrading#upgrading-to-2110 on instructions to upgrade.`,

		Schema: map[string]*schema.Schema{
			"source_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"source_namespace": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},

			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"destination_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"destination_namespace": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "default",
			},

			"action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"allow",
					"deny",
				}, true),
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"meta": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceConsulIntentionCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	connect := client.Connect()

	intention, err := getIntention(d)
	if err != nil {
		return err
	}

	id, _, err := connect.IntentionCreate(intention, wOpts)
	if err != nil {
		return fmt.Errorf("Failed to create intention (dc: '%s'): %v", wOpts.Datacenter, err)
	}

	d.SetId(id)

	return resourceConsulIntentionRead(d, meta)
}

func resourceConsulIntentionUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	connect := client.Connect()

	intention, err := getIntention(d)
	if err != nil {
		return err
	}
	intention.ID = d.Id()

	if _, err := connect.IntentionUpdate(intention, wOpts); err != nil {
		return fmt.Errorf("Failed to update intention (dc: '%s'): %v", wOpts.Datacenter, err)
	}

	return resourceConsulIntentionRead(d, meta)
}

func resourceConsulIntentionRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	connect := client.Connect()

	id := d.Id()

	intention, _, err := connect.IntentionGet(id, qOpts)
	if err != nil {
		return fmt.Errorf("Failed to retrieve intention (dc: '%s'): %v", qOpts.Datacenter, err)
	}

	if intention == nil {
		d.SetId("")
		return nil
	}

	if err = d.Set("datacenter", qOpts.Datacenter); err != nil {
		return fmt.Errorf("failed to set 'datacenter': %v", err)
	}
	if err = d.Set("source_name", intention.SourceName); err != nil {
		return fmt.Errorf("failed to set 'source_name': %v", err)
	}
	if err = d.Set("source_namespace", intention.SourceNS); err != nil {
		return fmt.Errorf("failed to set 'source_namespace': %v", err)
	}
	if err = d.Set("destination_name", intention.DestinationName); err != nil {
		return fmt.Errorf("failed to set 'destination_name': %v", err)
	}
	if err = d.Set("destination_namespace", intention.DestinationNS); err != nil {
		return fmt.Errorf("failed to set 'destination_namespace': %v", err)
	}
	if err = d.Set("description", intention.Description); err != nil {
		return fmt.Errorf("failed to set 'description': %v", err)
	}
	if err = d.Set("action", string(intention.Action)); err != nil {
		return fmt.Errorf("failed to set 'action': %v", err)
	}
	if err = d.Set("meta", intention.Meta); err != nil {
		return fmt.Errorf("failed to set 'meta': %v", err)
	}

	return nil
}

func resourceConsulIntentionDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	connect := client.Connect()
	id := d.Id()

	if _, err := connect.IntentionDelete(id, wOpts); err != nil {
		return fmt.Errorf("Failed to delete intention with id '%s' in %s: %v",
			id, wOpts.Datacenter, err)
	}

	// Clear the ID
	d.SetId("")
	return nil
}

func getIntention(d *schema.ResourceData) (*consulapi.Intention, error) {
	sourceName := d.Get("source_name").(string)
	sourceNamespace := d.Get("source_namespace").(string)
	destinationName := d.Get("destination_name").(string)
	destinationNamespace := d.Get("destination_namespace").(string)

	var intentionAction consulapi.IntentionAction
	action := d.Get("action").(string)

	if action == "allow" {
		intentionAction = consulapi.IntentionActionAllow
	} else if action == "deny" {
		intentionAction = consulapi.IntentionActionDeny
	} else {
		return nil, fmt.Errorf("Failed to create intention, action must match '%v' or '%v'", consulapi.IntentionActionAllow, consulapi.IntentionActionDeny)
	}

	intention := &consulapi.Intention{
		SourceName:      sourceName,
		DestinationName: destinationName,
		Action:          intentionAction,
		SourceNS:        sourceNamespace,
		DestinationNS:   destinationNamespace,
	}

	if description, ok := d.GetOk("description"); ok {
		intention.Description = description.(string)
	}

	if meta, ok := d.GetOk("meta"); ok {
		metas := meta.(map[string]interface{})
		newMeta := make(map[string]string)
		for k, v := range metas {
			newMeta[k] = v.(string)
		}
		intention.Meta = newMeta
	}
	return intention, nil
}
