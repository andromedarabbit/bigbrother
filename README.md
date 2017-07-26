[![codecov](https://codecov.io/gh/DailyHotel/bigbrother/branch/master/graph/badge.svg)](https://codecov.io/gh/DailyHotel/bigbrother)
[![Build Status](https://travis-ci.org/DailyHotel/bigbrother.svg?branch=master)](https://travis-ci.org/DailyHotel/bigbrother)
[![Docker Pulls](https://img.shields.io/docker/pulls/mashape/kong.svg)](https://hub.docker.com/r/dailyhotel/bigbrother/)
[![](https://images.microbadger.com/badges/image/dailyhotel/bigbrother.svg)](https://microbadger.com/images/dailyhotel/bigbrother "Get your own image badge on microbadger.com")
[![](https://images.microbadger.com/badges/version/dailyhotel/bigbrother.svg)](https://microbadger.com/images/dailyhotel/bigbrother "Get your own version badge on microbadger.com")

# bigbrother

Do you want to know how frequently your ELBs change their own IP addresses? **bigbrother** will help you to find out it!

## How to configure

The following environment variables can be configured:

- `SLEEP_TIME_IN_SECOND`: *not required*, every 10 minute by default
- `AWS_REGION`: *not required*
- `AWS_ACCESS_KEY_ID`: *not required*
- `AWS_SECRET_ACCESS_KEY`: *not required*
- `SLACK_TOKEN`: *required* when you want to post a notification to your Slack channel
- `SLACK_CHANNEL`: *not required*, `ops` by default
- `LOGGING_LEVEL`: *not required*, [`WARNING`](https://godoc.org/github.com/golang/glog) by default, applied to Docker container only.

## How to run on Docker

``` bash
docker run dailyhotel/bigbrother -e AWS_REGION=ap-northeast-1 -e AWS_ACCESS_KEY_ID=MY_KEY -e AWS_SECRET_ACCESS_KEY=MY_SECRET -e SLACK_TOKEN=MY_TOKEN -e SLACK_CHANNEL=ops -e LOGGING_LEVEL=INFO
```

## How to run on Kubernetes

``` yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: bigbrother
spec:
  replicas: 1
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: bigbrother
    spec:
      containers:
      - name: bigbrother
        image: dailyhotel/bigbrother:latest
        env:
        - name: SLACK_TOKEN
          value: "MY_TOKEN"
        - name: SLACK_CHANNEL
          value: ops
        - name: LOGGING_LEVEL
          value: INFO
```

## Thanks to

- [nlopes/slack](https://github.com/nlopes/slack): Slack API in Go

## See also

- [Slack Legacy Tokens](https://api.slack.com/custom-integrations/legacy-tokens)
- [Slack Message Builder](https://api.slack.com/docs/messages/builder)
- [WeltN24/metrics-discovery](https://github.com/WeltN24/metrics-discovery): can be used in a monitoring systems like nagios or zabbix to discover items on aws
- [hirose31/monitor-elb-address](https://github.com/hirose31/monitor-elb-address): check IP addresses of ELB periodically and notifies us when changed
- [doublemarket/elbipchecker](https://github.com/doublemarket/elbipchecker): ELB IP address change monitor script
