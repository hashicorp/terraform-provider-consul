package consul

import (
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	consulAutopilotConfigDatacenter              = "datacenter"
	consulAutopilotConfigCleanupDeadServers      = "cleanup_dead_servers"
	consulAutopilotConfigLastContactThreshold    = "last_contact_threshold"
	consulAutopilotConfigMaxTrailingLogs         = "max_trailing_logs"
	consulAutopilotConfigServerStabilizationTime = "server_stabilization_time"
	consulAutopilotConfigRedundancyZoneTag       = "redundancy_zone_tag"
	consulAutopilotConfigDisableUpgradeMigration = "disable_upgrade_migration"
	consulAutopilotConfigUpgradeVersionTag       = "upgrade_version_tag"
)

func resourceConsulAutopilotConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulAutopilotConfigCreate,
		Update: resourceConsulAutopilotConfigUpdate,
		Read:   resourceConsulAutopilotConfigRead,
		Delete: resourceConsulAutopilotConfigDelete,

		Schema: map[string]*schema.Schema{
			consulAutopilotConfigDatacenter: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			consulAutopilotConfigCleanupDeadServers: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			consulAutopilotConfigLastContactThreshold: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "200ms",
			},
			consulAutopilotConfigMaxTrailingLogs: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  250,
			},
			consulAutopilotConfigServerStabilizationTime: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "10s",
			},
			consulAutopilotConfigRedundancyZoneTag: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			consulAutopilotConfigDisableUpgradeMigration: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			consulAutopilotConfigUpgradeVersionTag: {
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

	lastContactThreshold, err := time.ParseDuration(d.Get(consulAutopilotConfigLastContactThreshold).(string))
	if err != nil {
		return fmt.Errorf("Could not parse '%v': %v", consulAutopilotConfigLastContactThreshold, err)
	}
	serverStabilizationTime, err := time.ParseDuration(d.Get(consulAutopilotConfigServerStabilizationTime).(string))
	if err != nil {
		return fmt.Errorf("Could not parse '%v': %v", consulAutopilotConfigServerStabilizationTime, err)
	}

	config := &consulapi.AutopilotConfiguration{
		CleanupDeadServers:      d.Get(consulAutopilotConfigCleanupDeadServers).(bool),
		LastContactThreshold:    consulapi.NewReadableDuration(lastContactThreshold),
		MaxTrailingLogs:         uint64(d.Get(consulAutopilotConfigMaxTrailingLogs).(int)),
		ServerStabilizationTime: consulapi.NewReadableDuration(serverStabilizationTime),
		RedundancyZoneTag:       d.Get(consulAutopilotConfigRedundancyZoneTag).(string),
		DisableUpgradeMigration: d.Get(consulAutopilotConfigDisableUpgradeMigration).(bool),
		UpgradeVersionTag:       d.Get(consulAutopilotConfigUpgradeVersionTag).(string),
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

	d.Set(consulAutopilotConfigCleanupDeadServers, config.CleanupDeadServers)
	d.Set(consulAutopilotConfigLastContactThreshold, config.LastContactThreshold)
	d.Set(consulAutopilotConfigMaxTrailingLogs, config.MaxTrailingLogs)
	d.Set(consulAutopilotConfigServerStabilizationTime, config.ServerStabilizationTime)
	d.Set(consulAutopilotConfigRedundancyZoneTag, config.RedundancyZoneTag)
	d.Set(consulAutopilotConfigDisableUpgradeMigration, config.DisableUpgradeMigration)
	d.Set(consulAutopilotConfigUpgradeVersionTag, config.UpgradeVersionTag)

	return nil
}

func resourceConsulAutopilotConfigDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
