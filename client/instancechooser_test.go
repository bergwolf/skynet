package client

import (
	"github.com/bketelsen/skynet"
	"testing"
)

func stubServiceInfo() (si *skynet.ServiceInfo) {
	si = &skynet.ServiceInfo{
		Config: &skynet.ServiceConfig{
			ServiceAddr: &skynet.BindAddr{},
		},
		Stats: &skynet.ServiceStatistics{},
	}

	return
}

func stubClient() (c *Client) {
	c = &Client{
		Config: &skynet.ClientConfig{},
	}

	return
}

func TestDefaultComparatorRegionFirst(t *testing.T) {
	c := stubClient()
	c.Config.Region = "A"

	s1 := stubServiceInfo()
	s1.Config.Region = "A"

	s2 := stubServiceInfo()
	s2.Config.Region = "B"

	// We should choose instances in the client's region over external
	if !defaultComparator(c, s1, s2) || defaultComparator(c, s2, s1) {
		t.Error("Failed to select instance in same region")
	}

	s2.Config.Region = "A"
	s1.Stats.LastRequest = "2012-11-14T15:04:05Z-0700"
	s2.Stats.LastRequest = "2012-11-14T15:04:10Z-0700"

	// Tie breaker LastRequest
	if !defaultComparator(c, s1, s2) || defaultComparator(c, s2, s1) {
		t.Error("Failed to select instance with older LastRequest")
	}
}

func TestDefaultComparatorHostFirst(t *testing.T) {
	c := stubClient()
	c.Config.Region = "A"
	c.Config.Host = "127.0.0.1"

	s1 := stubServiceInfo()
	s1.Config.Region = "A"
	s1.Config.ServiceAddr.IPAddress = "127.0.0.1"

	s2 := stubServiceInfo()
	s2.Config.Region = "A"
	s2.Config.ServiceAddr.IPAddress = "192.168.1.1"

	// We should choose instances on the client's host over region
	if !defaultComparator(c, s1, s2) || defaultComparator(c, s2, s1) {
		t.Error("Failed to select instance on same host")
	}

	s2.Config.ServiceAddr.IPAddress = "127.0.0.1"
	s1.Stats.LastRequest = "2012-11-14T15:04:05Z-0700"
	s2.Stats.LastRequest = "2012-11-14T15:04:10Z-0700"

	// Tie breaker LastRequest
	if !defaultComparator(c, s1, s2) || defaultComparator(c, s2, s1) {
		t.Error("Failed to select instance with older LastRequest")
	}
}
