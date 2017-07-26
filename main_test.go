package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostToSlack(t *testing.T) {
	elbName := "Test Load Balancer"
	dnsName := "test.dailyhotel.com"
	old := []string{"1.2.3.4", "5.6.7.8"}
	new := []string{"2.3.4.5", "5.6.7.8"}
	assert.Nilf(t, postToSlack(dnsName, elbName, old, new), "Posting to Slack failed")
}

func TestProcessLoadBalancer(t *testing.T) {
	elbAddrs := make(map[string][]string)
	elbName := "elb1"
	dnsName := "elb1.example.com"
	addresses := []string{"1.2.3.4", "2.3.4.5"}

	processLoadBalancer(elbAddrs, elbName, dnsName, addresses)
	processLoadBalancer(elbAddrs, elbName, dnsName, addresses)
	processLoadBalancer(elbAddrs, elbName, dnsName, addresses)

	assert.Equal(t, 1, len(elbAddrs))
	assert.Equal(t, elbAddrs[dnsName], addresses)
}

func TestProcessLoadBalancerWhenIpAddressIsChanged(t *testing.T) {
	elbAddrs := make(map[string][]string)
	elbName := "elb1"
	dnsName := "elb1.example.com"
	addresses := []string{"1.2.3.4", "2.3.4.5"}

	processLoadBalancer(elbAddrs, elbName, dnsName, addresses)

	newAddresses := []string{"1.2.3.4", "5.6.7.8"}
	processLoadBalancer(elbAddrs, elbName, dnsName, newAddresses)

	assert.Equal(t, 1, len(elbAddrs))
	assert.Equal(t, elbAddrs[dnsName], newAddresses)
}
