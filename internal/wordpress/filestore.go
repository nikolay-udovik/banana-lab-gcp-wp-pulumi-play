package wordpress

import (
	"banana-lab-gcp-wp-pulumi-play/models"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/filestore"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NewFilestore creates a new Filestore instance based on the provided input and returns its output.
func NewFilestore(ctx *pulumi.Context, input *models.FilestoreInput) (*models.Filestore, error) {

	// Convert []string to pulumi.StringArray
	modes := make(pulumi.StringArray, len(input.Modes))
	for i, mode := range input.Modes {
		modes[i] = pulumi.String(mode)
	}

	instance, err := filestore.NewInstance(ctx, input.Name, &filestore.InstanceArgs{
		Name:     pulumi.String(input.Name),
		Location: pulumi.String(input.Location),
		Tier:     pulumi.String(input.Tier),
		FileShares: &filestore.InstanceFileSharesArgs{
			CapacityGb: pulumi.Int(input.CapacityGb),
			Name:       pulumi.String(input.ShareName),
		},
		Networks: filestore.InstanceNetworkArray{
			&filestore.InstanceNetworkArgs{
				Network: pulumi.String(input.Network),
				Modes:   modes,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Return the full Filestore model
	return &models.Filestore{
		Input: *input,
		Output: models.FilestoreOutput{
			Instance: instance,
		},
	}, nil
}
