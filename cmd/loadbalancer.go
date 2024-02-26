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

var tGroupSlice []string

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "The aws-resource describe",
	Long:  `Usage: aws-ssm describe [args].`,
	Run: func(cmd *cobra.Command, args []string) {
		// flags for custom aws config
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")

		describeTargetGroups(profile, region)
		//describeTargetHealth(profile, region)

	},
}

func describeTargetGroups(profile string, region string) {
	elbClient := pkg.Newelb(profile, region)
	var lbarns []string

	result, err := elbClient.DescribeTargetGroups(context.TODO(), &elasticloadbalancingv2.DescribeTargetGroupsInput{
		//		TargetGroupArns: []string{tgArn},
	})

	if err != nil {
		fmt.Println(err)
	}

	for _, output := range result.TargetGroups {
		tgroup := *output.TargetGroupArn
		for _, lb := range output.LoadBalancerArns {
			lbarns = []string{lb}
		}
		describeTargetHealth(profile, region, tgroup, lbarns)
	}

}

func describeTargetHealth(profile string, region string, tGroup string, lbArns []string) {
	elbClient := pkg.Newelb(profile, region)

	result, err := elbClient.DescribeTargetHealth(context.TODO(), &elasticloadbalancingv2.DescribeTargetHealthInput{
		TargetGroupArn: &tGroup,
	})

	if err != nil {
		fmt.Println(err)
	}

	// new target group
	newtGroup := tGroup

	// check empty struct

	for _, output := range result.TargetHealthDescriptions {
		tgrouphealth := &output.TargetHealth.State
		//if *tgrouphealth != "healthy" && *tgrouphealth != "draining" && *tgrouphealth != "unused" {
		if *tgrouphealth != "healthy" && *tgrouphealth != "draining" {
			if newtGroup == tGroup {
				//lbarnsslice = append(lbarnsslice, lbArns)
				tGroupSlice = append(tGroupSlice, tGroup)
				fmt.Println(*tgrouphealth)
				fmt.Println(lbArns)
				fmt.Println(tGroupSlice)
				break
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
