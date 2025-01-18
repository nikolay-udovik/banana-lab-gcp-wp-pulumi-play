package wordpress

import (
	"banana-lab-gcp-wp-pulumi-play/models"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NewWordpressInstanceGroup creates the instance template, instance group manager, and autoscaler.
func NewWordpressInstanceGroup(ctx *pulumi.Context, input *models.WordpressInstanceGroupInput, filestoreIP pulumi.StringInput) (*models.WordpressInstanceGroup, error) {
	// Step 1: Create Instance Template
	instanceTemplate, err := createInstanceTemplate(ctx, input, filestoreIP)
	if err != nil {
		return nil, err
	}

	// Step 2: Create Instance Group Manager
	instanceGroupManager, err := compute.NewInstanceGroupManager(ctx, input.InstanceTemplateName+"-igm", &compute.InstanceGroupManagerArgs{
		Name: pulumi.String(input.InstanceTemplateName + "-igm"),
		Zone: pulumi.String(input.Zone),
		Versions: compute.InstanceGroupManagerVersionArray{
			&compute.InstanceGroupManagerVersionArgs{
				InstanceTemplate: instanceTemplate.ID(),
			},
		},
		BaseInstanceName: pulumi.String("instance"),
		TargetSize:       pulumi.Int(input.MinReplicas), // Set initial size
	})
	if err != nil {
		return nil, err
	}

	// Step 3: Create Autoscaler
	autoscaler, err := compute.NewAutoscaler(ctx, input.AutoscalerName, &compute.AutoscalerArgs{
		Name:   pulumi.String(input.AutoscalerName),
		Zone:   pulumi.String(input.Zone),
		Target: instanceGroupManager.ID(),
		AutoscalingPolicy: &compute.AutoscalerAutoscalingPolicyArgs{
			MaxReplicas: pulumi.Int(input.MaxReplicas),
			MinReplicas: pulumi.Int(input.MinReplicas),
			CpuUtilization: &compute.AutoscalerAutoscalingPolicyCpuUtilizationArgs{
				Target: pulumi.Float64(input.CpuUtilizationTarget / 100), // Scale based on CPU
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Return WordpressInstanceGroup model
	return &models.WordpressInstanceGroup{
		Input: *input,
		Output: models.WordpressInstanceGroupOutput{
			InstanceTemplateID:     instanceTemplate.ID(),
			InstanceGroupManagerID: instanceGroupManager.ID(),
			AutoscalerID:           autoscaler.ID(),
		},
	}, nil
}

// createInstanceTemplate creates an instance template with user data to configure Apache2 and Filestore.
func createInstanceTemplate(ctx *pulumi.Context, input *models.WordpressInstanceGroupInput, filestoreIP pulumi.StringInput) (*compute.InstanceTemplate, error) {
	// User data script for startup configuration
	userData := pulumi.Sprintf(`#!/bin/bash
apt-get update
apt-get install -y apache2 nfs-common
systemctl enable apache2
systemctl start apache2

FILESTORE_IP="%s"
MOUNT_POINT="/var/www/html"
mkdir -p $MOUNT_POINT
echo "$FILESTORE_IP:/ $MOUNT_POINT nfs defaults,_netdev 0 0" >> /etc/fstab
mount -a
chown -R www-data:www-data $MOUNT_POINT
chmod -R 755 $MOUNT_POINT
systemctl restart apache2
`, filestoreIP)

	// Create the instance template
	instanceTemplate, err := compute.NewInstanceTemplate(ctx, input.InstanceTemplateName, &compute.InstanceTemplateArgs{
		Name:                  pulumi.String(input.InstanceTemplateName),
		MachineType:           pulumi.String(input.MachineType),
		MetadataStartupScript: userData, // Pass the user data script
		Disks: compute.InstanceTemplateDiskArray{
			&compute.InstanceTemplateDiskArgs{
				SourceImage: pulumi.String(input.SourceImageID), // Use the provided image
			},
		},
		NetworkInterfaces: compute.InstanceTemplateNetworkInterfaceArray{
			&compute.InstanceTemplateNetworkInterfaceArgs{
				Subnetwork: input.SubnetID, // Use the subnet ID from the input
			},
		},
		ConfidentialInstanceConfig: &compute.InstanceTemplateConfidentialInstanceConfigArgs{
			EnableConfidentialCompute: pulumi.Bool(input.EnableConfidentialVM),
		},
	})
	if err != nil {
		return nil, err
	}

	return instanceTemplate, nil
}
