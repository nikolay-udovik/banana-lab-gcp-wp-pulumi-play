package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nikolay-udovik/banana-lab-gcp-wp-pulumi-play/internal/network"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	destroy := len(os.Args) > 1 && os.Args[1] == "destroy"

	stackName := os.Getenv("PULUMI_STACK")
	if stackName == "" {
		stackName = "dev"
	}

	ctx := context.Background()
	projectName := "inlinePulumiApp"

	s, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, func(ctx *pulumi.Context) error {
		// Fetch configuration dynamically
		config := make(map[string]string)
		keys := []string{"network:vpcName", "network:cidrRange", "gcp:region", "gcp:project"}
		for _, key := range keys {
			value, exists := ctx.GetConfig(key)
			if !exists {
				return fmt.Errorf("missing configuration: %s", key)
			}
			config[key] = value
		}

		// Create VPC
		vpcInput := &network.VPCInput{
			Ctx:     ctx,
			Name:    config["network:vpcName"],
			Region:  config["gcp:region"],
			Project: config["gcp:project"],
		}
		vpcOutput, err := network.NewVPC(vpcInput)
		if err != nil {
			return fmt.Errorf("failed to create VPC: %v", err)
		}

		// Create Subnet
		subnetInput := &network.SubnetInput{
			Ctx:             ctx,
			Name:            config["network:vpcName"] + "-subnet",
			VPCName:         vpcOutput.Id.ApplyT(func(id string) string { return id }).(string),
			Region:          config["gcp:region"],
			IpCidrRange:     config["network:cidrRange"],
			EnablePrivateIP: true,
			Project:         config["gcp:project"],
		}
		subnetOutput, err := network.NewSubnet(subnetInput)
		if err != nil {
			return fmt.Errorf("failed to create Subnet: %v", err)
		}

		// Export outputs
		ctx.Export("vpcName", vpcOutput.Name)
		ctx.Export("subnetName", subnetOutput.Name)

		return nil
	})
	if err != nil {
		fmt.Printf("Failed to set up a workspace: %v\n", err)
		os.Exit(1)
	}

	if _, err = s.Refresh(ctx); err != nil {
		fmt.Printf("Failed to refresh stack: %v\n", err)
		os.Exit(1)
	}

	if destroy {
		if _, err = s.Destroy(ctx, optdestroy.ProgressStreams(os.Stdout)); err != nil {
			fmt.Printf("Failed to destroy stack: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Stack successfully destroyed")
		os.Exit(0)
	}

	if _, err = s.Up(ctx, optup.ProgressStreams(os.Stdout)); err != nil {
		fmt.Printf("Failed to update stack: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Stack successfully updated")
}
