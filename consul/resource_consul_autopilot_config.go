package consul

import (
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"

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
	client := getClient(meta)
	operator := client.Operator()

	dc, err := getDC(d, client, meta)
	if err != nil {
		return err
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
	client := getClient(meta)
	operator := client.Operator()

	dc, err := getDC(d, client, meta)
	if err != nil {
		return err
	}

	qOpts := &consulapi.QueryOptions{
		Datacenter: dc,
	}
	config, err := operator.AutopilotGetConfiguration(qOpts)
	if err != nil {
		return fmt.Errorf("Failed to fetch autopilot configuration: %v", err)
	}

	d.SetId(fmt.Sprintf("consul-autopilot-%s", dc))

	if err = d.Set("cleanup_dead_servers", config.CleanupDeadServers); err != nil {
		return errwrap.Wrapf("Unable to store cleanup_dead_servers: {{err}}", err)
	}
	if err = d.Set("last_contact_threshold", config.LastContactThreshold.String()); err != nil {
		return errwrap.Wrapf("Unable to store last_contact_threshold: {{err}}", err)
	}
	if err = d.Set("max_trailing_logs", config.MaxTrailingLogs); err != nil {
		return errwrap.Wrapf("Unable to store max_trailing_logs: {{err}}", err)
	}
	if err = d.Set("server_stabilization_time", config.ServerStabilizationTime.String()); err != nil {
		return errwrap.Wrapf("Unable to store server_stabilization_time: {{err}}", err)
	}
	if err = d.Set("redundancy_zone_tag", config.RedundancyZoneTag); err != nil {
		return errwrap.Wrapf("Unable to store redundancy_zone_tag: {{err}}", err)
	}
	if err = d.Set("disable_upgrade_migration", config.DisableUpgradeMigration); err != nil {
		return errwrap.Wrapf("Unable to store disable_upgrade_migration: {{err}}", err)
	}
	if err = d.Set("upgrade_version_tag", config.UpgradeVersionTag); err != nil {
		return errwrap.Wrapf("Unable to store upgrade_version_tag: {{err}}", err)
	}

	return nil
}

func resourceConsulAutopilotConfigDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
