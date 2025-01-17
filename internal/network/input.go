package network

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

// VPCInput defines the input arguments required to create a VPC.
type VPCInput struct {
	Ctx     *pulumi.Context // Pulumi context
	Name    string          // Name of the VPC
	Region  string          // Region where the VPC is deployed
	Project string          // GCP project ID
}

// SubnetInput defines the input arguments required to create a Subnet.
type SubnetInput struct {
	Ctx             *pulumi.Context // Pulumi context
	Name            string          // Name of the Subnet
	VPCName         string          // Name of the associated VPC
	Region          string          // Region where the Subnet is deployed
	IpCidrRange     string          // CIDR range for the Subnet
	EnablePrivateIP bool            // Whether to enable private IP access
	Project         string          // GCP project ID
}
