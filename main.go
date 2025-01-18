package main

import (
	"banana-lab-gcp-wp-pulumi-play/internal/wordpress"
	"banana-lab-gcp-wp-pulumi-play/models"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Initialize configuration
		conf := config.New(ctx, "")

		// Populate WordPress Input
		wpInput := &models.WordpressInput{
			VPC: models.VPC{
				Input: models.VPCInput{
					Name:    conf.Get("network:vpcName"),
					Region:  conf.Get("gcp:region"),
					Project: conf.Get("gcp:project"),
				},
			},
			Subnet: models.Subnet{
				Input: models.SubnetInput{
					Name:            conf.Get("wordpress:subnetName"),
					Region:          conf.Get("gcp:region"),
					IpCidrRange:     conf.Get("wordpress:subnetCidrRange"),
					EnablePrivateIP: conf.GetBool("wordpress:enablePrivateIP"),
				},
			},
			CloudSQL: models.CloudSQL{
				Input: models.CloudSQLInput{
					Name:            conf.Get("db:name"),
					DatabaseVersion: conf.Get("db:databaseVersion"),
					Region:          conf.Get("gcp:region"),
					Tier:            conf.Get("db:tier"),
					DiskSizeGb:      conf.GetInt("db:diskSizeGb"),
					BackupEnabled:   conf.GetBool("db:backupEnabled"),
					DatabaseName:    "wordpress",
					UserName:        conf.Get("db:userName"),
					Password:        conf.GetSecret("db:password"),
				},
			},
			Filestore: models.Filestore{
				Input: models.FilestoreInput{
					Name:       conf.Get("sharedStrg:name"),
					Location:   conf.Get("gcp:region") + "-a",
					Tier:       conf.Get("sharedStrg:tier"),
					CapacityGb: conf.GetInt("sharedStrg:capacityGb"),
					Network:    conf.Get("network:vpcName"),
					Modes:      []string{"MODE_IPV4"},
				},
			},
			InstanceGroup: models.WordpressInstanceGroup{
				Input: models.WordpressInstanceGroupInput{
					InstanceTemplateName: conf.Get("wordpress:instanceTemplateName"),
					MachineType:          conf.Get("wordpress:machineType"),
					SourceImageID:        conf.Get("wordpress:sourceImageID"),
					EnableConfidentialVM: conf.GetBool("wordpress:enableConfidentialVM"),
					AutoscalerName:       conf.Get("wordpress:autoscalerName"),
					MinReplicas:          conf.GetInt("wordpress:minReplicas"),
					MaxReplicas:          conf.GetInt("wordpress:maxReplicas"),
					CpuUtilizationTarget: conf.GetFloat64("wordpress:cpuUtilizationTarget"),
					Zone:                 conf.Get("wordpress:zone"),
				},
			},
		}

		// Deploy WordPress resources
		_, err := wordpress.DeployWordpress(ctx, wpInput)
		if err != nil {
			return err
		}

		return nil
	})
}
