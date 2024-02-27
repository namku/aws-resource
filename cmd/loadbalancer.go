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
	"golang.org/x/exp/slices"
)

// var tgarntest = "arn:aws:elasticloadbalancing:eu-central-1:079806680060:targetgroup/k8s-wipo-wiposerv-f816083ba0/f46063e09347d019"
var tGroupSlice []string

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
	},
}

func describeTargetGroups(profile string, region string, withouttargets bool, unhealthy bool) {
	elbClient := pkg.Newelb(profile, region)
	var lbarns []string

	result, err := elbClient.DescribeTargetGroups(context.TODO(), &elasticloadbalancingv2.DescribeTargetGroupsInput{
		//TargetGroupArns: []string{tgarntest},
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(len(result.TargetGroups))
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
		//TargetGroupArn: &tgarntest,
	})

	if err != nil {
		fmt.Println(err)
	}

	if withouttargets {
		loadbalancerWithoutTargets(result, tGroup, lbArns)
	}
	if unhealthy {
		loadbalancerUnhealthy(result, tGroup, lbArns)
	}
}

// check target group without targets or not loadbalancer associated.
func loadbalancerWithoutTargets(result *elasticloadbalancingv2.DescribeTargetHealthOutput, tGroup string, lbArns []string) {
	if len(result.TargetHealthDescriptions) == 0 {
		fmt.Println(tGroup)
		if lbArns == nil {
			fmt.Println("Target group isn't associated to a load balancer")
		} else {
			fmt.Print(lbArns)
			fmt.Println("Target group without targets")
		}
	}
}

// check target group without targets or not loadbalancer associated.
func loadbalancerUnhealthy(result *elasticloadbalancingv2.DescribeTargetHealthOutput, tGroup string, lbArns []string) {
	// new target group
	newtGroup := ""
	newStatus := ""
	var statusSlice []string

	for _, output := range result.TargetHealthDescriptions {
		tgrouphealth := &output.TargetHealth.State
		//if *tgrouphealth != "healthy" && *tgrouphealth != "draining" && *tgrouphealth != "inital" {
		if newtGroup != "" {
			if newStatus != string(*tgrouphealth) {
				if !slices.Contains(statusSlice, string(*tgrouphealth)) {
					statusSlice = append(statusSlice, string(*tgrouphealth))
					printTarget(tGroup, statusSlice)
				}
			}
			tGroupSlice = append(tGroupSlice, tGroup)
		} else {
			newtGroup = tGroup
			newStatus = string(*tgrouphealth)
			statusSlice = append(statusSlice, newStatus)
			tGroupSlice = nil
			tGroupSlice = append(tGroupSlice, tGroup)
			printTarget(tGroup, statusSlice)
		}
		//if newtGroup == tGroup {
		//lbarnsslice = append(lbarnsslice, lbArns)
		//tGroupSlice = append(tGroupSlice, tGroup)
		//fmt.Println(lbArns)
		//fmt.Println(tGroup)
		//fmt.Println(*tgrouphealth)
		//break
		//}
		//}
	}
}

func printTarget(tGroup string, statusSlice []string) {
	fmt.Print(tGroup + " ============> ")
	fmt.Println(statusSlice)
}

func init() {
	loadbalancerCmd.Flags().BoolP("without-targets", "w", false, "Target groups without target or not associated to a load balancer.")
	loadbalancerCmd.Flags().BoolP("unhealthy", "u", false, "Target groups with unhealthy targets.")
	rootCmd.AddCommand(loadbalancerCmd)
}
