package network

import (
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NewSubnet creates a new Subnet within a VPC based on the provided input and returns its output.
func NewSubnet(input *SubnetInput) (*SubnetOutput, error) {
	ctx := input.Ctx

	if input.Name == "" {
		input.Name = input.VPCName + "-subnet"
	}

	// Create the Subnet resource
	subnet, err := compute.NewSubnetwork(ctx, input.Name, &compute.SubnetworkArgs{
		Network:               pulumi.String(input.VPCName),
		Region:                pulumi.String(input.Region),
		IpCidrRange:           pulumi.String(input.IpCidrRange),
		PrivateIpGoogleAccess: pulumi.Bool(input.EnablePrivateIP),
	})
	if err != nil {
		return nil, err
	}

	return &SubnetOutput{
		Name: subnet.Name,
		Id:   subnet.ID(),
	}, nil
}
