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

//var tGroupSlice []string

var loadbalancerCmd = &cobra.Command{
	Use:   "loadbalancer",
	Short: "The aws-resource loadbalancer",
	Long:  `Usage: aws-ssm loadbalancer [args].`,
	Run: func(cmd *cobra.Command, args []string) {
		// flags for custom aws config
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")

		withouttargets, _ := cmd.Flags().GetBool("without-targets")
		unhealthy, _ := cmd.Flags().GetBool("unhealthy")

		describeTargetGroups(profile, region, withouttargets, unhealthy)
		//describeTargetHealth(profile, region)

	},
}

func describeTargetGroups(profile string, region string, withouttargets bool, unhealthy bool) {
	elbClient := pkg.Newelb(profile, region)
	var lbarns []string

	result, err := elbClient.DescribeTargetGroups(context.TODO(), &elasticloadbalancingv2.DescribeTargetGroupsInput{
		//		TargetGroupArns: []string{tgArn},
	})

	if err != nil {
		fmt.Println(err)
	}

	for _, output := range result.TargetGroups {
		lbarns = nil
		tgroup := *output.TargetGroupArn
		for _, lb := range output.LoadBalancerArns {
			lbarns = []string{lb}
		}
		describeTargetHealth(profile, region, tgroup, lbarns, withouttargets, unhealthy)
	}

}

func describeTargetHealth(profile string, region string, tGroup string, lbArns []string, withouttargets bool, unhealthy bool) {
	elbClient := pkg.Newelb(profile, region)

	result, err := elbClient.DescribeTargetHealth(context.TODO(), &elasticloadbalancingv2.DescribeTargetHealthInput{
		TargetGroupArn: &tGroup,
	})

	if err != nil {
		fmt.Println(err)
	}

	if withouttargets {
		loadbalancerwithouttargets(result, tGroup, lbArns)
	}
	if unhealthy {
		loadbalancerunhealthy(result, tGroup, lbArns)
	}

}

// check target group without targets or not loadbalancer associated.
func loadbalancerwithouttargets(result *elasticloadbalancingv2.DescribeTargetHealthOutput, tGroup string, lbArns []string) {
	if len(result.TargetHealthDescriptions) == 0 {
		fmt.Println(tGroup)
		if lbArns == nil {
			fmt.Println("Target group isn't associated to a load balancer")
		} else {
			fmt.Println(lbArns)
			fmt.Println("Target group without targets")
		}
	}
}

// check target group without targets or not loadbalancer associated.
func loadbalancerunhealthy(result *elasticloadbalancingv2.DescribeTargetHealthOutput, tGroup string, lbArns []string) {
	// new target group
	newtGroup := tGroup

	for _, output := range result.TargetHealthDescriptions {
		tgrouphealth := &output.TargetHealth.State
		if *tgrouphealth != "healthy" && *tgrouphealth != "draining" && *tgrouphealth != "inital" {
			if newtGroup == tGroup {
				//lbarnsslice = append(lbarnsslice, lbArns)
				//tGroupSlice = append(tGroupSlice, tGroup)
				fmt.Println(lbArns)
				fmt.Println(tGroup)
				fmt.Println(*tgrouphealth)
				break
			}
		}
	}
}

func init() {
	loadbalancerCmd.Flags().BoolP("without-targets", "w", false, "Target groups without target or not associated to a load balancer.")
	loadbalancerCmd.Flags().BoolP("unhealthy", "u", false, "Target groups with unhealthy targets.")
	rootCmd.AddCommand(loadbalancerCmd)
}
