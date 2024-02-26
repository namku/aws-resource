package pkg

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

func Newelb(profile string, region string) *elasticloadbalancingv2.Client {
	// initialize aws session using config files
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
		config.WithRetryMaxAttempts(20),
	)

	if err != nil {
		panic(fmt.Sprintf("failed loading config, %v", err))
	}

	elbClient := elasticloadbalancingv2.NewFromConfig(cfg)
	return elbClient
}
