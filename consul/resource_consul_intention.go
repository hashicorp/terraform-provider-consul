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

		Schema: map[string]*schema.Schema{
			"source_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"destination_name": {
				Type:     schema.TypeString,
				Required: true,
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
	client := getClient(meta)
	connect := client.Connect()

	sourceName := d.Get("source_name").(string)
	destinationName := d.Get("destination_name").(string)

	var intentionAction consulapi.IntentionAction
	action := d.Get("action").(string)

	dc := ""
	if _, ok := d.GetOk("datacenter"); ok {
		dc = d.Get("datacenter").(string)
	}

	if action == "allow" {
		intentionAction = consulapi.IntentionActionAllow
	} else if action == "deny" {
		intentionAction = consulapi.IntentionActionDeny
	} else {
		return fmt.Errorf("Failed to create intention, action must match '%v' or '%v'", consulapi.IntentionActionAllow, consulapi.IntentionActionDeny)
	}

	wOpts := consulapi.WriteOptions{Datacenter: dc}

	intention := &consulapi.Intention{
		SourceName:      sourceName,
		DestinationName: destinationName,
		Action:          intentionAction,
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

	id, _, err := connect.IntentionCreate(intention, &wOpts)
	if err != nil {
		return fmt.Errorf("Failed to create intention (dc: '%s'): %v", dc, err)
	}

	d.SetId(id)

	return resourceConsulIntentionRead(d, meta)
}

func resourceConsulIntentionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)
	connect := client.Connect()

	sourceName := d.Get("source_name").(string)
	destinationName := d.Get("destination_name").(string)

	var intentionAction consulapi.IntentionAction
	action := d.Get("action").(string)

	if action == "allow" {
		intentionAction = consulapi.IntentionActionAllow
	} else if action == "deny" {
		intentionAction = consulapi.IntentionActionDeny
	} else {
		return fmt.Errorf("Failed to create intention, action must match '%v' or '%v'", consulapi.IntentionActionAllow, consulapi.IntentionActionDeny)
	}

	dc := ""
	if _, ok := d.GetOk("datacenter"); ok {
		dc = d.Get("datacenter").(string)
	}

	intention := &consulapi.Intention{
		SourceName:      sourceName,
		DestinationName: destinationName,
		Action:          intentionAction,
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

	// Setup the operations using the datacenter
	wOpts := consulapi.WriteOptions{Datacenter: dc}

	if _, err := connect.IntentionUpdate(intention, &wOpts); err != nil {
		return fmt.Errorf("Failed to update intention (dc: '%s'): %v", dc, err)
	}

	return resourceConsulIntentionRead(d, meta)
}

func resourceConsulIntentionRead(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)
	connect := client.Connect()

	dc := ""
	if _, ok := d.GetOk("datacenter"); ok {
		dc = d.Get("datacenter").(string)
	}

	id := d.Id()

	qOpts := consulapi.QueryOptions{Datacenter: dc}
	intention, _, err := connect.IntentionGet(id, &qOpts)
	if err != nil {
		return fmt.Errorf("Failed to retrieve intention (dc: '%s'): %v", dc, err)
	}

	if intention == nil {
		d.SetId("")
		return nil
	}

	d.Set("source_name", intention.SourceName)
	d.Set("destination_name", intention.DestinationName)
	d.Set("description", intention.Description)
	d.Set("action", string(intention.Action))
	d.Set("meta", intention.Meta)

	return nil
}

func resourceConsulIntentionDelete(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)
	connect := client.Connect()
	id := d.Id()

	dc := ""
	if _, ok := d.GetOk("datacenter"); ok {
		dc = d.Get("datacenter").(string)
	}

	// Setup the operations using the datacenter
	wOpts := consulapi.WriteOptions{Datacenter: dc}

	if _, err := connect.IntentionDelete(id, &wOpts); err != nil {
		return fmt.Errorf("Failed to delete intention with id '%s' in %s: %v",
			id, dc, err)
	}

	// Clear the ID
	d.SetId("")
	return nil
}
