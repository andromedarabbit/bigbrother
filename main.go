package main

import (
	"flag"
	"net"
	"sort"

	"reflect"

	"strings"

	"os"
	"time"

	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/golang/glog"
	"github.com/nlopes/slack"
	_ "golang.org/x/net/publicsuffix"
)

var sleepTime time.Duration

func init() {
	sleepTimeString, exists := os.LookupEnv("SLEEP_TIME_IN_SECOND")
	if exists == true {
		i64, err := strconv.ParseInt(sleepTimeString, 10, 32)
		if err != nil {
			glog.Fatal("Error while trying to parse SLEEP_TIME_IN_SECOND env var")
		}
		sleepTime = time.Duration(int32(i64))
	} else {
		sleepTime = 60 * 10 // every 10 minute by default
	}
}

func getRegion(metadata *ec2metadata.EC2Metadata) string {
	region, exists := os.LookupEnv("AWS_REGION")
	if exists == true {
		return region
	}

	region, err := metadata.Region()
	if err == nil {
		return region
	}

	glog.Warningf("Unable to retrieve the region from the EC2 instance %v\n", err)

	return "ap-northeast-1"
}

func main() {
	flag.Parse()
	glog.Info("AWS ELB IP Addresses Watcher")

	metadata := ec2metadata.New(session.New())
	region := getRegion(metadata)

	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{},
			&ec2rolecreds.EC2RoleProvider{Client: metadata},
		})

	awsConfig := aws.NewConfig()
	awsConfig.WithCredentials(creds)
	awsConfig.WithRegion(region)
	sess := session.New(awsConfig)

	elbApi := elb.New(sess)
	if elbApi == nil {
		glog.Fatal("Failed to make AWS connection")
	}

	elbAddrs := make(map[string][]string)

	for {
		list, err := elbApi.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{})
		if err != nil {
			glog.Warning(err)
			continue
		}

		for _, desc := range list.LoadBalancerDescriptions {
			if *desc.Scheme != "internet-facing" {
				break
			}

			elbName := *desc.LoadBalancerName
			dnsName := *desc.DNSName
			addresses, err := net.LookupHost(dnsName)
			if err != nil {
				glog.Warningf("DNS resolution failed against ELB %s's address `%s`", elbName, dnsName)
				continue
			}

			processLoadBalancer(elbAddrs, elbName, dnsName, addresses)
		}

		time.Sleep(sleepTime * time.Second)
	}
}

func processLoadBalancer(elbAddrs map[string][]string, elbName string, dnsName string, addresses []string) {
	sort.Strings(addresses)
	if addressesSaved, ok := elbAddrs[dnsName]; ok {
		if reflect.DeepEqual(addressesSaved, addresses) {
			return
		}

		glog.Infof("ELB: %s, Domain: %s, Old IP: %s, New IP: %s", dnsName, elbName, addressesSaved, addresses)
		postToSlack(dnsName, elbName, addressesSaved, addresses)
	} else {
		glog.Infof("New ELB: %s, Domain: %s, IP: %s", dnsName, elbName, addresses)
	}
	elbAddrs[dnsName] = addresses
}

func postToSlack(dnsName string, elbName string, oldAddresses []string, newAddresses []string) error {
	token, exists := os.LookupEnv("SLACK_TOKEN")
	if exists == false {
		glog.Info("`SLACK_TOKEN` is not defined. Skip posting a message to Slack")
		return nil
	}

	channel, exists := os.LookupEnv("SLACK_CHANNEL")
	if exists == false {
		channel = "ops"
	}

	api := slack.New(token)
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Pretext: "ELB IP Change Detected",
		Text:    elbName,
		Fields: []slack.AttachmentField{
			{
				Title: "Domain",
				Value: dnsName,
			},
			{
				Title: "Old IP",
				Value: strings.Join(oldAddresses[:], ", "),
			},
			{
				Title: "New IP",
				Value: strings.Join(newAddresses[:], ", "),
			},
		},
	}
	params.Attachments = []slack.Attachment{attachment}
	msg := ""
	channelID, timestamp, err := api.PostMessage(channel, msg, params)
	if err != nil {
		glog.Errorf("%s\n", err)
		return err
	}
	glog.Infof("Message successfully sent to channel %s at %s", channelID, timestamp)
	return nil
}
