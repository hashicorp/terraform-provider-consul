package consul

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/hashicorp/consul/api"
)

func waitForACLTokenReplication(acl *api.ACL, qOpts *api.QueryOptions, index uint64) error {
	for attempt := 0; attempt <= 12; attempt++ {
		rs, _, err := acl.Replication(qOpts)
		if err != nil {
			return fmt.Errorf("error fetching ACL replication status: %w", err)
		}

		if !rs.Enabled || rs.ReplicationType != "tokens" || rs.ReplicatedTokenIndex >= index {
			return nil
		}

		time.Sleep(time.Duration(math.Pow(2, float64(attempt))) * time.Millisecond)
	}

	return errors.New("timed out waiting for ACL replication")
}

func waitForACLPolicyReplication(acl *api.ACL, qOpts *api.QueryOptions, index uint64) error {
	for attempt := 0; attempt <= 12; attempt++ {
		rs, _, err := acl.Replication(qOpts)
		if err != nil {
			return fmt.Errorf("error fetching ACL replication status: %w", err)
		}

		if !rs.Enabled || rs.ReplicatedIndex >= index {
			return nil
		}

		time.Sleep(time.Duration(math.Pow(2, float64(attempt))) * time.Millisecond)
	}

	return errors.New("timed out waiting for ACL replication")
}

func waitForACLRoleReplication(acl *api.ACL, qOpts *api.QueryOptions, index uint64) error {
	for attempt := 0; attempt <= 12; attempt++ {
		rs, _, err := acl.Replication(qOpts)
		if err != nil {
			return fmt.Errorf("error fetching ACL replication status: %w", err)
		}

		if !rs.Enabled || rs.ReplicatedRoleIndex >= index {
			return nil
		}

		time.Sleep(time.Duration(math.Pow(2, float64(attempt))) * time.Millisecond)
	}

	return errors.New("timed out waiting for ACL replication")
}
