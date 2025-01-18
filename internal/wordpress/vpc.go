package wordpress

import (
	"banana-lab-gcp-wp-pulumi-play/models"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NewVPC creates a new VPC based on the provided input and returns its full VPC model.
func NewVPC(ctx *pulumi.Context, input *models.VPCInput) (*models.VPC, error) {

	// Set a default name if none is provided
	if input.Name == "" {
		input.Name = "default-vpc"
	}

	// Create the VPC resource
	vpcResource, err := compute.NewNetwork(ctx, input.Name, &compute.NetworkArgs{
		Name:                  pulumi.String(input.Name),
		Project:               pulumi.String(input.Project),
		AutoCreateSubnetworks: pulumi.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	// Export VPC information for debugging or other purposes
	ctx.Export("vpcName", vpcResource.Name)
	ctx.Export("vpcID", vpcResource.ID())

	// Return the full VPC model, including both input and output
	return &models.VPC{
		Input: *input,
		Output: models.VPCOutput{
			Name: vpcResource.Name,
			Id:   vpcResource.ID(),
		},
	}, nil
}
