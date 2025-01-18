package wordpress

import (
	"banana-lab-gcp-wp-pulumi-play/models"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NewSubnet creates a new subnet and returns the Subnet model.
func NewSubnet(ctx *pulumi.Context, input *models.SubnetInput) (*models.Subnet, error) {
	// Create Subnet
	subnet, err := compute.NewSubnetwork(ctx, input.Name, &compute.SubnetworkArgs{
		Network:               input.VPCID,
		Region:                pulumi.String(input.Region),
		IpCidrRange:           pulumi.String(input.IpCidrRange),
		PrivateIpGoogleAccess: pulumi.Bool(input.EnablePrivateIP),
	})
	if err != nil {
		return nil, err
	}

	// Return Subnet model
	return &models.Subnet{
		Input:  *input,
		Output: models.SubnetOutput{ID: subnet.ID()},
	}, nil
}
