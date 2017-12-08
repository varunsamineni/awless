/* Copyright 2017 WALLIX

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

package awsspec

import (
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/wallix/awless/cloud/graph"
	"github.com/wallix/awless/logger"
)

type CreateScalingpolicy struct {
	_                   string `action:"create" entity:"scalingpolicy" awsAPI:"autoscaling" awsCall:"PutScalingPolicy" awsInput:"autoscaling.PutScalingPolicyInput" awsOutput:"autoscaling.PutScalingPolicyOutput"`
	logger              *logger.Logger
	graph               cloudgraph.GraphAPI
	api                 autoscalingiface.AutoScalingAPI
	AdjustmentType      *string `awsName:"AdjustmentType" awsType:"awsstr" templateName:"adjustment-type" required:""`
	Scalinggroup        *string `awsName:"AutoScalingGroupName" awsType:"awsstr" templateName:"scalinggroup" required:""`
	Name                *string `awsName:"PolicyName" awsType:"awsstr" templateName:"name" required:""`
	AdjustmentScaling   *int64  `awsName:"ScalingAdjustment" awsType:"awsint64" templateName:"adjustment-scaling" required:""`
	Cooldown            *int64  `awsName:"Cooldown" awsType:"awsint64" templateName:"cooldown"`
	AdjustmentMagnitude *int64  `awsName:"MinAdjustmentMagnitude" awsType:"awsint64" templateName:"adjustment-magnitude"`
}

func (cmd *CreateScalingpolicy) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*autoscaling.PutScalingPolicyOutput).PolicyARN)
}

type DeleteScalingpolicy struct {
	_      string `action:"delete" entity:"scalingpolicy" awsAPI:"autoscaling" awsCall:"DeletePolicy" awsInput:"autoscaling.DeletePolicyInput" awsOutput:"autoscaling.DeletePolicyOutput"`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    autoscalingiface.AutoScalingAPI
	Id     *string `awsName:"PolicyName" awsType:"awsstr" templateName:"id" required:""`
}

func (cmd *DeleteScalingpolicy) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}
