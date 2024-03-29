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
	"time"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/briandowns/spinner"
	"github.com/namku/aws-resource/pkg"
	progressbar "github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

// var tgarntest = "arn:aws:elasticloadbalancing:eu-central-1:079806680060:targetgroup/k8s-wipo-wiposerv-f816083ba0/f46063e09347d019"
var indicatorSpinner *spinner.Spinner

var tGroupUnused []string
var tGroupWithoutTargets []string
var lbWithoutTargets []string

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

		//startSpinner()
		describeTargetGroups(nil, profile, region, withouttargets, unhealthy)
		//indicatorSpinner.Stop()

		if withouttargets {
			fmt.Println("[ TARGET GROUPS UNUSED ]")
			fmt.Println()
			for _, tGroupun := range tGroupUnused {
				fmt.Println(tGroupun)
			}
			fmt.Println()
			fmt.Println()
			fmt.Println("[ TARGET GROUPS WITHOUT TARGETS ]")
			fmt.Println()
			for _, tGroupwith := range tGroupWithoutTargets {
				fmt.Println(tGroupwith)
			}
			fmt.Println()
			fmt.Println()
			fmt.Println("[ LOADBALANCER WITHOUT TARGETS ]")
			fmt.Println()
			for _, lbwith := range lbWithoutTargets {
				fmt.Println(lbwith)
			}

			fmt.Println()
			fmt.Println()
			fmt.Println("[ TOTAL ]")
			fmt.Println()
			fmt.Print("TargetGroups unused ===========> ")
			fmt.Println(len(tGroupUnused))
			fmt.Print("TargetGroups without targets ===========> ")
			fmt.Println(len(tGroupWithoutTargets))
			fmt.Print("LoadBalancers without targets ===========> ")
			fmt.Println(len(lbWithoutTargets))

		}
	},
}

func describeTargetGroups(nextMarker *string, profile string, region string, withouttargets bool, unhealthy bool) {
	elbClient := pkg.Newelb(profile, region)
	var lbarns []string

	result, err := elbClient.DescribeTargetGroups(context.TODO(), &elasticloadbalancingv2.DescribeTargetGroupsInput{
		Marker: nextMarker,
		//TargetGroupArns: []string{tgarntest},
	})

	if err != nil {
		fmt.Println(err)
	}

	bar := progressbar.Default(int64(len(result.TargetGroups)))
	for _, output := range result.TargetGroups {
		bar.Add(1)
		time.Sleep(40 * time.Millisecond)
		lbarns = nil
		tgroup := *output.TargetGroupArn
		for _, lb := range output.LoadBalancerArns {
			lbarns = []string{lb}
		}
		describeTargetHealth(profile, region, tgroup, lbarns, withouttargets, unhealthy)
	}

	// check if there are more target group pages.
	if result.NextMarker != nil {
		describeTargetGroups(result.NextMarker, profile, region, withouttargets, unhealthy)
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
	// Define suffix spinner
	//indicatorSpinner.Suffix = "  " + tGroup

	if len(result.TargetHealthDescriptions) == 0 {
		if lbArns == nil {
			tGroupUnused = append(tGroupUnused, tGroup)
		} else {
			tGroupWithoutTargets = append(tGroupWithoutTargets, tGroup)
			for _, lbarn := range lbArns {
				lbWithoutTargets = append(lbWithoutTargets, lbarn)
			}
		}
	}
}

// check target group without targets or not loadbalancer associated.
func loadbalancerUnhealthy(result *elasticloadbalancingv2.DescribeTargetHealthOutput, tGroup string, lbArns []string) {
	newtGroup := ""
	newStatus := ""
	var statusSlice []string

	for _, output := range result.TargetHealthDescriptions {
		tgrouphealth := &output.TargetHealth.State
		// asdfasdfa
		if *tgrouphealth != "healthy" && *tgrouphealth != "draining" && *tgrouphealth != "inital" {
			if newtGroup != "" {
				if newStatus != string(*tgrouphealth) {
					if !slices.Contains(statusSlice, string(*tgrouphealth)) {
						statusSlice = append(statusSlice, string(*tgrouphealth))
					}
				}
			} else {
				newtGroup = tGroup
				newStatus = string(*tgrouphealth)
				statusSlice = append(statusSlice, newStatus)
			}
		}
	}
	if statusSlice != nil {
		printTarget(tGroup, lbArns, statusSlice)
	}
}

func printTarget(tGroup string, lbArns []string, statusSlice []string) {
	fmt.Println(lbArns)
	fmt.Print(tGroup + " ============> ")
	fmt.Println(statusSlice)
}

func startSpinner() {
	// Start spinner
	indicatorSpinner = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	indicatorSpinner.Start()
	indicatorSpinner.Prefix = "  "
	pkg.SetupCloseHandler(indicatorSpinner)
}

func init() {
	loadbalancerCmd.Flags().BoolP("without-targets", "w", false, "Target groups without target or not associated to a load balancer.")
	loadbalancerCmd.Flags().BoolP("unhealthy", "u", false, "Target groups with unhealthy targets.")
	rootCmd.AddCommand(loadbalancerCmd)
}
