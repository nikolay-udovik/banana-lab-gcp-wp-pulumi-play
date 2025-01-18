package models

import (
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/filestore"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/sql"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// VPC represents the lifecycle of a VPC resource, with nested Input and Output structs.
type VPC struct {
	Input  VPCInput  // Input fields for creating the VPC
	Output VPCOutput // Output fields after the VPC is created
}

// VPCOutput defines the outputs for a VPC resource.
type VPCOutput struct {
	Name pulumi.StringOutput // The name of the VPC
	Id   pulumi.IDOutput     // The unique identifier of the VPC
}

// VPCInput defines the input arguments required to create a VPC.
type VPCInput struct {
	Name    string // Name of the VPC
	Region  string // Region where the VPC is deployed
	Project string // GCP project ID
}

// Filestore represents the lifecycle of a Filestore instance, with nested Input and Output structs.
type Filestore struct {
	Input  FilestoreInput  // Input fields for creating the Filestore instance
	Output FilestoreOutput // Output fields after the Filestore instance is created
}

// FilestoreInput represents the input parameters for creating a Filestore instance.
type FilestoreInput struct {
	Name       string   // Filestore instance name
	Location   string   // Region and zone (e.g., "us-central1-b")
	Tier       string   // Performance tier (e.g., "BASIC_HDD")
	CapacityGb int      // Storage capacity in GB
	ShareName  string   // Name of the file share
	Network    string   // VPC network name
	Modes      []string // Access modes (e.g., "MODE_IPV4")
}

// FilestoreOutput represents the output of the Filestore instance creation.
type FilestoreOutput struct {
	Instance *filestore.Instance
}

// CloudSQL represents the lifecycle of a Cloud SQL instance, with nested Input and Output structs.
type CloudSQL struct {
	Input  CloudSQLInput  // Input fields for creating the Cloud SQL instance
	Output CloudSQLOutput // Output fields after the Cloud SQL instance is created
}

// CloudSQLInput represents the input parameters for creating a Cloud SQL instance.
type CloudSQLInput struct {
	Name            string              // Name of the Cloud SQL instance
	DatabaseVersion string              // Database engine (e.g., "MYSQL_8_4", "MYSQL_8_0", "MYSQL_5_7")
	Region          string              // Region where the instance will be deployed
	Tier            string              // Machine type (e.g., "db-f1-micro")
	DiskSizeGb      int                 // Disk size in GB
	BackupEnabled   bool                // Whether automated backups are enabled
	IPConfig        []string            // Authorized IP addresses for connections
	DatabaseName    string              // Initial database name
	UserName        string              // Initial database username
	Password        pulumi.StringOutput // Initial database password
}

// CloudSQLOutput represents the output of the Cloud SQL instance creation.
type CloudSQLOutput struct {
	Instance *sql.DatabaseInstance // Reference to the created Cloud SQL instance
	Database *sql.Database         // Reference to the initial database
	User     *sql.User             // Reference to the initial database user
}

// Subnet represents the lifecycle of a subnet resource, with nested Input and Output structs.
type Subnet struct {
	Input  SubnetInput  // Input fields for creating the subnet
	Output SubnetOutput // Output fields after the subnet is created
}

// SubnetInput defines the input arguments for creating a subnet.
type SubnetInput struct {
	Name            string          // Name of the subnet
	VPCID           pulumi.IDOutput // ID of the VPC to which the subnet belongs
	Region          string          // Region where the subnet is deployed
	IpCidrRange     string          // CIDR range for the subnet
	EnablePrivateIP bool            // Whether to enable private IP access
}

// SubnetOutput defines the outputs of the subnet resource.
type SubnetOutput struct {
	ID pulumi.IDOutput // ID of the created subnet
}

// WordpressInstanceGroup represents the lifecycle of instance template, instance group manager, and autoscaler.
type WordpressInstanceGroup struct {
	Input  WordpressInstanceGroupInput  // Input fields for creating instance group-related resources
	Output WordpressInstanceGroupOutput // Output fields after the resources are created
}

// WordpressInstanceGroupInput defines the input for instance template, instance group manager, and autoscaler.
type WordpressInstanceGroupInput struct {
	InstanceTemplateName string          // Name of the instance template
	MachineType          string          // Machine type for the instance
	SourceImageID        string          // User-provided source image ID
	SubnetID             pulumi.IDOutput // ID of the subnet to attach the instance template
	EnableConfidentialVM bool            // Whether to enable confidential compute
	AutoscalerName       string          // Name of the autoscaler
	MinReplicas          int             // Minimum number of instances
	MaxReplicas          int             // Maximum number of instances
	CpuUtilizationTarget float64         // Target CPU utilization for scaling (percentage)
	Zone                 string          // GCP zone for the resources
}

// WordpressInstanceGroupOutput defines the output of instance group-related resources.
type WordpressInstanceGroupOutput struct {
	InstanceTemplateID     pulumi.IDOutput // ID of the created instance template
	InstanceGroupManagerID pulumi.IDOutput // ID of the created instance group manager
	AutoscalerID           pulumi.IDOutput // ID of the created autoscaler
}

// Wordpress represents the entire WordPress setup, including all resources and their dependencies.
type Wordpress struct {
	Input  WordpressInput  // Input fields for the WordPress setup
	Output WordpressOutput // Output fields after all resources are created
}

// WordpressInput defines the input arguments for the entire WordPress setup.
type WordpressInput struct {
	VPC           VPC                    // VPC resource
	Subnet        Subnet                 // Subnet resource
	InstanceGroup WordpressInstanceGroup // Instance group and related resources
	CloudSQL      CloudSQL               // CloudSQL resource
	Filestore     Filestore              // Filestore resource
}

// WordpressOutput defines the output of the WordPress setup.
type WordpressOutput struct {
	VPC           VPCOutput                    // VPC output
	Subnet        SubnetOutput                 // Subnet output
	InstanceGroup WordpressInstanceGroupOutput // Instance group output
	CloudSQL      CloudSQLOutput               // CloudSQL output
	Filestore     FilestoreOutput              // Filestore output
}
