package network

import (
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NewVPC creates a new VPC based on the provided input and returns its output.
func NewVPC(input *VPCInput) (*VPCOutput, error) {
	ctx := input.Ctx

	if input.Name == "" {
		input.Name = "default-vpc"
	}

	// Create the VPC resource
	vpc, err := compute.NewNetwork(ctx, input.Name, &compute.NetworkArgs{
		AutoCreateSubnetworks: pulumi.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	return &VPCOutput{
		Name: vpc.Name,
		Id:   vpc.ID(),
	}, nil
}
