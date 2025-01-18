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
					Name:    conf.Require("network:vpcName"),
					Region:  conf.Require("gcp:region"),
					Project: conf.Require("gcp:project"),
				},
			},
			Subnet: models.Subnet{
				Input: models.SubnetInput{
					Name:            conf.Require("wordpress:subnetName"),
					Region:          conf.Require("gcp:region"),
					IpCidrRange:     conf.Require("wordpress:subnetCidrRange"),
					EnablePrivateIP: conf.RequireBool("wordpress:enablePrivateIP"),
				},
			},
			CloudSQL: models.CloudSQL{
				Input: models.CloudSQLInput{
					Name:            conf.Require("db:name"),
					DatabaseVersion: conf.Require("db:databaseVersion"),
					Region:          conf.Require("gcp:region"),
					Tier:            conf.Require("db:tier"),
					DiskSizeGb:      conf.RequireInt("db:diskSizeGb"),
					BackupEnabled:   conf.RequireBool("db:backupEnabled"),
					DatabaseName:    "wordpress",
					UserName:        conf.Require("db:userName"),
					Password:        conf.RequireSecret("db:password"),
				},
			},
			Filestore: models.Filestore{
				Input: models.FilestoreInput{
					Name:       conf.Require("sharedStrg:name"),
					Location:   conf.Require("gcp:region") + "-a",
					Tier:       conf.Require("sharedStrg:tier"),
					CapacityGb: conf.RequireInt("sharedStrg:capacityGb"),
					Network:    conf.Require("network:vpcName"),
					Modes:      []string{"MODE_IPV4"},
				},
			},
			InstanceGroup: models.WordpressInstanceGroup{
				Input: models.WordpressInstanceGroupInput{
					InstanceTemplateName: conf.Require("wordpress:instanceTemplateName"),
					MachineType:          conf.Require("wordpress:machineType"),
					SourceImageID:        conf.Require("wordpress:sourceImageID"),
					EnableConfidentialVM: conf.RequireBool("wordpress:enableConfidentialVM"),
					AutoscalerName:       conf.Require("wordpress:autoscalerName"),
					MinReplicas:          conf.RequireInt("wordpress:minReplicas"),
					MaxReplicas:          conf.RequireInt("wordpress:maxReplicas"),
					CpuUtilizationTarget: conf.RequireFloat64("wordpress:cpuUtilizationTarget"),
					Zone:                 conf.Require("wordpress:zone"),
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
