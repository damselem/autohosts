package hosts

import (
	"fmt"
	"log"

	"net/http"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
)

var gcpRegions []string = []string{
	"europe-west1",
	"us-central1",
	"us-east1",
}

func GCP(projectsWhitelist []string) ([]Entry, error) {
	var wg sync.WaitGroup

	entriesCh := make(chan []Entry)
	errCh := make(chan error)

	ctx := context.Background()
	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope, cloudresourcemanager.CloudPlatformScope)
	if err != nil {
		return []Entry{}, errors.Wrap(err, "could not initialize GCP client")
	}

	projects, err := getAllProjects(c, ctx, projectsWhitelist)
	if err != nil {
		return []Entry{}, err
	}

	for _, project := range projects {
		ctx := context.Background()

		c, err := google.DefaultClient(ctx, compute.CloudPlatformScope, cloudresourcemanager.CloudPlatformScope)
		if err != nil {
			errCh <- errors.Wrap(err, "could not initialize GCP client")

			continue
		}

		zones, err := getAllZonesForProject(c, ctx, project)
		if err != nil {
			msg := fmt.Sprintf("could not get the list of GCP zones for project %s", project.ProjectId)
			errCh <- errors.Wrap(err, msg)

			continue
		}

		for _, zone := range zones {
			wg.Add(1)

			go func(project cloudresourcemanager.Project, zone compute.Zone) {
				var entries []Entry
				ctx := context.Background()

				c, err := google.DefaultClient(ctx, compute.CloudPlatformScope, cloudresourcemanager.CloudPlatformScope)
				if err != nil {
					errCh <- errors.Wrap(err, "could not initialize GCP client")
					return
				}

				instances, _ := getAllInstancesForProject(c, ctx, project, zone.Name)
				for _, instance := range instances {
					entries = append(entries, Entry{
						Address:  instance.NetworkInterfaces[0].AccessConfigs[0].NatIP,
						Comment:  fmt.Sprintf("# GCP - %s", zone.Name),
						Hostname: instance.Name,
					})
				}

				entriesCh <- entries
			}(project, zone)
		}
	}

	var entries []Entry
	go func() {
		for projectEntries := range entriesCh {
			entries = append(entries, projectEntries...)
			wg.Done()
		}
	}()

	var instErrs []error
	go func() {
		for projectErr := range errCh {
			instErrs = append(instErrs, projectErr)
			wg.Done()
		}
	}()

	wg.Wait()

	if len(instErrs) > 0 {
		err = instErrs[0]
	}

	return entries, err
}

func getAllProjects(c *http.Client, ctx context.Context, projectsWhitelist []string) ([]cloudresourcemanager.Project, error) {
	projects := make([]cloudresourcemanager.Project, 0)

	service, err := cloudresourcemanager.New(c)
	if err != nil {
		log.Fatal(err)
	}

	req := service.Projects.List()
	if err := req.Pages(ctx, func(page *cloudresourcemanager.ListProjectsResponse) error {
		for _, project := range page.Projects {
			if contains(projectsWhitelist, project.ProjectId) {
				projects = append(projects, *project)
			}
		}

		return nil
	}); err != nil {
		return projects, errors.Wrap(err, "could not obtain the list of GCP projects")
	}

	return projects, nil
}

func getAllZonesForProject(c *http.Client, ctx context.Context, project cloudresourcemanager.Project) ([]compute.Zone, error) {
	zones := make([]compute.Zone, 0)

	computeService, err := compute.New(c)
	if err != nil {
		return zones, errors.Wrap(err, "could not initialize GCP Compute service")
	}

	req := computeService.Zones.List(project.ProjectId)
	if err := req.Pages(ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			zones = append(zones, *zone)
		}

		return nil
	}); err != nil {
		msg := fmt.Sprintf("could not obtain the list of zones for GCP's project %s", project.ProjectId)
		return zones, errors.Wrap(err, msg)
	}

	return zones, nil
}

func getAllInstancesForProject(c *http.Client, ctx context.Context, project cloudresourcemanager.Project, zone string) ([]compute.Instance, error) {
	instances := make([]compute.Instance, 0)

	computeService, err := compute.New(c)
	if err != nil {
		return instances, errors.Wrap(err, "could not initialize GCP Compute service")
	}

	req := computeService.Instances.List(project.ProjectId, zone)
	if err := req.Pages(ctx, func(page *compute.InstanceList) error {
		for _, instance := range page.Items {
			instances = append(instances, *instance)
		}

		return nil
	}); err != nil {
		msg := fmt.Sprintf("could not obtain the list of instances for GCP's project %s", project.ProjectId)
		return instances, errors.Wrap(err, msg)
	}

	return instances, nil
}

func contains(slc []string, val string) bool {
	for _, slcVal := range slc {
		if val == slcVal {
			return true
		}
	}

	return false
}
