package wordpress

import (
	"banana-lab-gcp-wp-pulumi-play/models"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/sql"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NewDatabaseInstance creates a new Filestore instance based on the provided input and returns its output.
func NewDatabaseInstance(ctx *pulumi.Context, input *models.CloudSQLInput) (*models.CloudSQL, error) {

	instance, err := sql.NewDatabaseInstance(ctx, "main", &sql.DatabaseInstanceArgs{
		Name:            pulumi.String(input.Name),
		DatabaseVersion: pulumi.String(input.DatabaseVersion),
		Region:          pulumi.String(input.Region),
		Settings: &sql.DatabaseInstanceSettingsArgs{
			Tier: pulumi.String(input.Tier),
		},
	})
	if err != nil {
		return nil, err
	}

	// Create the initial database
	database, err := sql.NewDatabase(ctx, input.DatabaseName, &sql.DatabaseArgs{
		Instance: instance.Name,
		Name:     pulumi.String(input.DatabaseName),
	})
	if err != nil {
		return nil, err
	}

	// Create the initial user
	user, err := sql.NewUser(ctx, input.UserName, &sql.UserArgs{
		Instance: instance.Name,
		Name:     pulumi.String(input.UserName),
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	// Return the full CloudSQL model
	return &models.CloudSQL{
		Input: *input,
		Output: models.CloudSQLOutput{
			Instance: instance,
			Database: database,
			User:     user,
		},
	}, nil
}
