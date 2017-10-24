package hosts

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var awsRegions = []string{
	"us-east-2",
	"us-east-1",
	"us-west-1",
	"us-west-2",
	"ca-central-1",
	"ap-south-1",
	"ap-northeast-2",
	"ap-southeast-1",
	"ap-southeast-2",
	"ap-northeast-1",
	"eu-central-1",
	"eu-west-1",
	"eu-west-2",
	"sa-east-1",
}

var autoscaleGroups = struct {
	sync.RWMutex
	groups map[string]int
}{groups: make(map[string]int)}

// AWS returns entries for each EC2 instance in AWS account
// By default it uses the credentials stored in ~/.aws/credentials
func AWS(withEmr, withAutoscale bool) ([]Entry, error) {
	var wg sync.WaitGroup

	entriesCh := make(chan []Entry)
	errCh := make(chan error)

	var entries []Entry
	var instErrs []error

	for _, region := range awsRegions {
		wg.Add(1)
		go func(region string) {
			var regionEntries []Entry

			// Create a different client for each region to improve concurrency
			client := newEC2Client(region)
			instances, err := client.getInstances()
			if err != nil {
				errCh <- err
				return
			}

			for _, inst := range instances {
				tags := getInstanceTags(inst)
				ip := inst.PublicIpAddress
				name := tags["Name"]

				_, isEmr := tags["aws:elasticmapreduce:instance-group-role"]
				_, isAutoscale := tags["aws:autoscaling:groupName"]

				if (isEmr && !withEmr) || (isAutoscale && !withAutoscale) {
					continue
				}

				if ip != nil && name != "" {
					regionEntries = append(regionEntries, Entry{
						Address:  *ip,
						Comment:  fmt.Sprintf("# AWS - %s", region),
						Hostname: name,
					})
				}
			}

			entriesCh <- regionEntries
		}(region)
	}

	go func() {
		for regionEntries := range entriesCh {
			entries = append(entries, regionEntries...)
			wg.Done()
		}
	}()

	go func() {
		for instErr := range errCh {
			instErrs = append(instErrs, instErr)
			wg.Done()
		}
	}()

	wg.Wait()

	var err error
	if len(instErrs) > 0 {
		err = instErrs[0]
	}

	return entries, err
}

type ec2RegionClient struct {
	region string
	svc    *ec2.EC2
}

func newEC2Client(region string) *ec2RegionClient {
	config := aws.NewConfig().WithCredentialsChainVerboseErrors(true).WithRegion(region)
	svc := ec2.New(session.New(), config)

	return &ec2RegionClient{
		region: region,
		svc:    svc,
	}
}

func (c *ec2RegionClient) getInstances() ([]*ec2.Instance, error) {
	var instances []*ec2.Instance

	resp, err := c.svc.DescribeInstances(nil)
	if err != nil {
		return nil, err
	}

	for _, res := range resp.Reservations {
		instances = append(instances, res.Instances...)
	}

	return instances, nil
}

func getInstanceTags(inst *ec2.Instance) map[string]string {
	var groupName string
	tags := make(map[string]string)

	for _, tag := range inst.Tags {
		key := *tag.Key
		value := *tag.Value

		if key == "aws:autoscaling:groupName" {
			re := regexp.MustCompile("-ondemand.*")
			groupName = re.ReplaceAllLiteralString(value, "")

			autoscaleGroups.Lock()
			autoscaleGroups.groups[groupName]++
			autoscaleGroups.Unlock()
		}

		tags[key] = value
	}

	if groupName != "" {
		autoscaleGroups.RLock()
		tags["Name"] = fmt.Sprintf("%s%d", groupName, autoscaleGroups.groups[groupName])
		autoscaleGroups.RUnlock()
	}

	return tags
}
