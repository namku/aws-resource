/*
Copyright © 2024 Isaac Lopez syak7771@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/namku/aws-resource/pkg"
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "The aws-resource describe",
	Long:  `Usage: aws-ssm describe [args].`,
	Run: func(cmd *cobra.Command, args []string) {
		// flags for custom aws config
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")

		describeTargetHealth(profile, region)

	},
}

func describeTargetHealth(profile string, region string) {
	elbClient := pkg.Newelb(profile, region)

	// to test
	tgArn := "arn:aws:elasticloadbalancing:eu-central-1:373072521776:targetgroup/test/d06a0e8740c461e2"

	result, _ := elbClient.DescribeTargetHealth(context.TODO(), &elasticloadbalancingv2.DescribeTargetHealthInput{
		TargetGroupArn: &tgArn,
	})

	fmt.Println(result)

}

func init() {
	rootCmd.AddCommand(describeCmd)
}
