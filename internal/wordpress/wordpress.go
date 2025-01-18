package wordpress

import (
	"banana-lab-gcp-wp-pulumi-play/models"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// DeployWordpress deploys all WordPress resources based on the provided input and returns the output.
func DeployWordpress(ctx *pulumi.Context, input *models.WordpressInput) (*models.WordpressOutput, error) {
	var output models.WordpressOutput

	// Step 1: Deploy VPC
	vpc, err := NewVPC(ctx, &input.VPC.Input)
	if err != nil {
		return nil, err
	}
	output.VPC = vpc.Output
	ctx.Export("VPCName", vpc.Output.Name)
	ctx.Export("VPCID", vpc.Output.Id)

	// Step 2: Deploy Subnet
	subnet, err := NewSubnet(ctx, &models.SubnetInput{
		Name:            input.Subnet.Input.Name,
		VPCID:           vpc.Output.Id,
		Region:          input.Subnet.Input.Region,
		IpCidrRange:     input.Subnet.Input.IpCidrRange,
		EnablePrivateIP: input.Subnet.Input.EnablePrivateIP,
	})
	if err != nil {
		return nil, err
	}
	output.Subnet = subnet.Output
	ctx.Export("SubnetName", pulumi.String(input.Subnet.Input.Name))
	ctx.Export("SubnetID", subnet.Output.ID)

	// Step 3: Deploy CloudSQL
	cloudSQL, err := NewDatabaseInstance(ctx, &input.CloudSQL.Input)
	if err != nil {
		return nil, err
	}
	output.CloudSQL = cloudSQL.Output
	ctx.Export("CloudSQLInstanceName", cloudSQL.Output.Instance.Name)
	ctx.Export("CloudSQLDatabaseName", cloudSQL.Output.Database.Name)
	ctx.Export("CloudSQLUserName", pulumi.String(input.CloudSQL.Input.UserName))

	// Step 4: Deploy Filestore
	filestore, err := NewFilestore(ctx, &input.Filestore.Input)
	if err != nil {
		return nil, err
	}
	output.Filestore = filestore.Output
	ctx.Export("FilestoreName", pulumi.String(input.Filestore.Input.Name))
	ctx.Export("FilestoreInstance", filestore.Output.Instance.Name)

	// Step 5: Deploy Instance Group (Template, Manager, and Autoscaler)
	instanceGroup, err := NewWordpressInstanceGroup(ctx, &models.WordpressInstanceGroupInput{
		InstanceTemplateName: input.InstanceGroup.Input.InstanceTemplateName,
		MachineType:          input.InstanceGroup.Input.MachineType,
		SourceImageID:        input.InstanceGroup.Input.SourceImageID,
		SubnetID:             subnet.Output.ID,
		EnableConfidentialVM: input.InstanceGroup.Input.EnableConfidentialVM,
		AutoscalerName:       input.InstanceGroup.Input.AutoscalerName,
		MinReplicas:          input.InstanceGroup.Input.MinReplicas,
		MaxReplicas:          input.InstanceGroup.Input.MaxReplicas,
		CpuUtilizationTarget: input.InstanceGroup.Input.CpuUtilizationTarget,
		Zone:                 input.InstanceGroup.Input.Zone,
	}, filestore.Output.Instance.Networks.ApplyT(func(networks []string) string {
		return networks[0]
	}).(pulumi.StringOutput))
	if err != nil {
		return nil, err
	}
	output.InstanceGroup = instanceGroup.Output
	ctx.Export("InstanceTemplateID", instanceGroup.Output.InstanceTemplateID)
	ctx.Export("InstanceGroupManagerID", instanceGroup.Output.InstanceGroupManagerID)
	ctx.Export("AutoscalerID", instanceGroup.Output.AutoscalerID)

	return &output, nil
}
