package consul

import (
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceConsulAutopilotConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulAutopilotConfigCreate,
		Update: resourceConsulAutopilotConfigUpdate,
		Read:   resourceConsulAutopilotConfigRead,
		Delete: resourceConsulAutopilotConfigDelete,

		Schema: map[string]*schema.Schema{
			"datacenter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"cleanup_dead_servers": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"last_contact_threshold": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "200ms",
			},
			"max_trailing_logs": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  250,
			},
			"server_stabilization_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "10s",
			},
			"redundancy_zone_tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"disable_upgrade_migration": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"upgrade_version_tag": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceConsulAutopilotConfigCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceConsulAutopilotConfigUpdate(d, meta)
}

func resourceConsulAutopilotConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	operator := client.Operator()

	var dc string
	if v, ok := d.GetOk("datacenter"); ok {
		dc = v.(string)
	} else {
		var err error
		dc, err = getDC(d, client)
		if err != nil {
			return err
		}
	}

	lastContactThreshold, err := time.ParseDuration(d.Get("last_contact_threshold").(string))
	if err != nil {
		return fmt.Errorf("Could not parse '%v': %v", "last_contact_threshold", err)
	}
	serverStabilizationTime, err := time.ParseDuration(d.Get("server_stabilization_time").(string))
	if err != nil {
		return fmt.Errorf("Could not parse '%v': %v", "server_stabilization_time", err)
	}

	config := &consulapi.AutopilotConfiguration{
		CleanupDeadServers:      d.Get("cleanup_dead_servers").(bool),
		LastContactThreshold:    consulapi.NewReadableDuration(lastContactThreshold),
		MaxTrailingLogs:         uint64(d.Get("max_trailing_logs").(int)),
		ServerStabilizationTime: consulapi.NewReadableDuration(serverStabilizationTime),
		RedundancyZoneTag:       d.Get("redundancy_zone_tag").(string),
		DisableUpgradeMigration: d.Get("disable_upgrade_migration").(bool),
		UpgradeVersionTag:       d.Get("upgrade_version_tag").(string),
	}
	wOpts := &consulapi.WriteOptions{
		Datacenter: dc,
	}
	err = operator.AutopilotSetConfiguration(config, wOpts)
	if err != nil {
		return fmt.Errorf("Failed to update autopilot configuration: %v", err)
	}

	return resourceConsulAutopilotConfigRead(d, meta)
}

func resourceConsulAutopilotConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	operator := client.Operator()

	dc := "dc1"

	qOpts := &consulapi.QueryOptions{
		Datacenter: dc,
	}
	config, err := operator.AutopilotGetConfiguration(qOpts)
	if err != nil {
		return fmt.Errorf("Failed to fetch autopilot configuration: %v", err)
	}

	d.SetId(fmt.Sprintf("consul-autopilot-%s", dc))

	d.Set("cleanup_dead_servers", config.CleanupDeadServers)
	d.Set("last_contact_threshold", config.LastContactThreshold)
	d.Set("max_trailing_logs", config.MaxTrailingLogs)
	d.Set("server_stabilization_time", config.ServerStabilizationTime)
	d.Set("redundancy_zone_tag", config.RedundancyZoneTag)
	d.Set("disable_upgrade_migration", config.DisableUpgradeMigration)
	d.Set("upgrade_version_tag", config.UpgradeVersionTag)

	return nil
}

func resourceConsulAutopilotConfigDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
