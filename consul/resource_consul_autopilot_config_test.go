package consul

import (
	"fmt"
	"regexp"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccConsulAutopilotConfig_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testFinalConfiguration,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConsulAutopilotConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_autopilot_config.config", "id", "consul-autopilot-dc1"),
					resource.TestCheckResourceAttr("consul_autopilot_config.config", "cleanup_dead_servers", "true"),
					resource.TestCheckResourceAttr("consul_autopilot_config.config", "last_contact_threshold", "200ms"),
					resource.TestCheckResourceAttr("consul_autopilot_config.config", "max_trailing_logs", "250"),
					resource.TestCheckResourceAttr("consul_autopilot_config.config", "server_stabilization_time", "10s"),
					resource.TestCheckResourceAttr("consul_autopilot_config.config", "redundancy_zone_tag", ""),
					resource.TestCheckResourceAttr("consul_autopilot_config.config", "disable_upgrade_migration", "false"),
					resource.TestCheckResourceAttr("consul_autopilot_config.config", "upgrade_version_tag", ""),
				),
			},
			resource.TestStep{
				Config: testAccConsulAutopilotConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_autopilot_config.config", "id", "consul-autopilot-dc1"),
					resource.TestCheckResourceAttr("consul_autopilot_config.config", "cleanup_dead_servers", "false"),
					resource.TestCheckResourceAttr("consul_autopilot_config.config", "last_contact_threshold", "1s"),
					resource.TestCheckResourceAttr("consul_autopilot_config.config", "max_trailing_logs", "100"),
					resource.TestCheckResourceAttr("consul_autopilot_config.config", "server_stabilization_time", "5s"),
					resource.TestCheckResourceAttr("consul_autopilot_config.config", "redundancy_zone_tag", "redundancy_tag"),
					resource.TestCheckResourceAttr("consul_autopilot_config.config", "disable_upgrade_migration", "true"),
					resource.TestCheckResourceAttr("consul_autopilot_config.config", "upgrade_version_tag", "version_tag"),
				),
			},
		},
	})
}

func TestAccConsulAutopilotConfig_parseduration(t *testing.T) {
	errorRegexp := regexp.MustCompile("Could not parse 'last_contact_threshold': time: invalid duration")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccConsulAutopilotConfigParseDuration,
				ExpectError: errorRegexp,
			},
		},
	})
}

// when destroying the consul_autopilot_config resource, the configuration
// should not be changed
func testFinalConfiguration(s *terraform.State) error {
	operator := testAccProvider.Meta().(*consulapi.Client).Operator()
	qOpts := &consulapi.QueryOptions{}
	config, err := operator.AutopilotGetConfiguration(qOpts)
	if err != nil {
		return fmt.Errorf("err: %v", err)
	}
	if config.CleanupDeadServers != false {
		return fmt.Errorf("err: cleanup_dead_servers during destroy: %v", config.CleanupDeadServers)
	}
	if config.LastContactThreshold.String() != "1s" {
		return fmt.Errorf("err: last_contact_threshold during destroy: %v", config.LastContactThreshold)
	}
	if config.MaxTrailingLogs != 100 {
		return fmt.Errorf("err: max_trailing_logs during destroy: %v", config.MaxTrailingLogs)
	}
	if config.ServerStabilizationTime.String() != "5s" {
		return fmt.Errorf("err: server_stabilization_time during destroy: %v", config.ServerStabilizationTime)
	}
	if config.RedundancyZoneTag != "redundancy_tag" {
		return fmt.Errorf("err: redundancy_zone_tag during destroy: %v", config.RedundancyZoneTag)
	}
	if config.DisableUpgradeMigration != true {
		return fmt.Errorf("err: disable_upgrade_migration during destroy: %v", config.DisableUpgradeMigration)
	}
	if config.UpgradeVersionTag != "version_tag" {
		return fmt.Errorf("err: upgrade_version_tag during destroy: %v", config.UpgradeVersionTag)
	}
	return nil
}

const testAccConsulAutopilotConfigBasic = `
resource "consul_autopilot_config" "config" {}
`

const testAccConsulAutopilotConfig = `
resource "consul_autopilot_config" "config" {
	cleanup_dead_servers      =  false
	last_contact_threshold    =  "1s"
	max_trailing_logs         =  100
	server_stabilization_time =  "5s"
	redundancy_zone_tag       =  "redundancy_tag"
	disable_upgrade_migration =  true
	upgrade_version_tag       =  "version_tag"
}`

const testAccConsulAutopilotConfigParseDuration = `
resource "consul_autopilot_config" "config" {
	last_contact_threshold = "one minute"
}
`
