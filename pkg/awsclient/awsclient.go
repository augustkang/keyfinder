package awsclient

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// GetIAMClient loads aws config then return iam.Client
func GetIAMClient() (c *iam.Client) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("[ERROR] Failed to load config")
		panic(err.Error())
	}
	return iam.NewFromConfig(cfg)
}
