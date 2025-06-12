// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConsulAutopilotConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulAutopilotConfigCreate,
		Update: resourceConsulAutopilotConfigUpdate,
		Read:   resourceConsulAutopilotConfigRead,
		Delete: resourceConsulAutopilotConfigDelete,

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cleanup_dead_servers": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"last_contact_threshold": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "200ms",
			},
			"max_trailing_logs": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  250,
			},
			"server_stabilization_time": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "10s",
			},
			"redundancy_zone_tag": {
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
	client, _, wOpts := getClient(d, meta)
	operator := client.Operator()

	lastContactThreshold, err := time.ParseDuration(d.Get("last_contact_threshold").(string))
	if err != nil {
		return fmt.Errorf("could not parse '%v': %v", "last_contact_threshold", err)
	}
	serverStabilizationTime, err := time.ParseDuration(d.Get("server_stabilization_time").(string))
	if err != nil {
		return fmt.Errorf("could not parse '%v': %v", "server_stabilization_time", err)
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
	err = operator.AutopilotSetConfiguration(config, wOpts)
	if err != nil {
		return fmt.Errorf("failed to update autopilot configuration: %v", err)
	}

	return resourceConsulAutopilotConfigRead(d, meta)
}

func resourceConsulAutopilotConfigRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	operator := client.Operator()

	config, err := operator.AutopilotGetConfiguration(qOpts)
	if err != nil {
		return fmt.Errorf("failed to fetch autopilot configuration: %v", err)
	}

	d.SetId(fmt.Sprintf("consul-autopilot-%s", qOpts.Datacenter))

	sw := newStateWriter(d)
	sw.set("cleanup_dead_servers", config.CleanupDeadServers)
	sw.set("last_contact_threshold", config.LastContactThreshold.String())
	sw.set("max_trailing_logs", config.MaxTrailingLogs)
	sw.set("server_stabilization_time", config.ServerStabilizationTime.String())
	sw.set("redundancy_zone_tag", config.RedundancyZoneTag)
	sw.set("disable_upgrade_migration", config.DisableUpgradeMigration)
	sw.set("upgrade_version_tag", config.UpgradeVersionTag)

	return sw.error()
}

func resourceConsulAutopilotConfigDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
