package network

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

// VPCOutput defines the outputs for a VPC resource.
type VPCOutput struct {
	Name pulumi.StringOutput // The name of the VPC
	Id   pulumi.IDOutput     // The unique identifier of the VPC
}

// SubnetOutput defines the outputs for a Subnet resource.
type SubnetOutput struct {
	Name pulumi.StringOutput // The name of the Subnet
	Id   pulumi.IDOutput     // The unique identifier of the Subnet
}
