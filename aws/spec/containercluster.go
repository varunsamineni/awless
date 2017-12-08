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
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/wallix/awless/cloud/graph"
	"github.com/wallix/awless/logger"
)

type CreateContainercluster struct {
	_      string `action:"create" entity:"containercluster" awsAPI:"ecs" awsCall:"CreateCluster" awsInput:"ecs.CreateClusterInput" awsOutput:"ecs.CreateClusterOutput"`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    ecsiface.ECSAPI
	Name   *string `awsName:"ClusterName" awsType:"awsstr" templateName:"name" required:""`
}

func (cmd *CreateContainercluster) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ecs.CreateClusterOutput).Cluster.ClusterArn)
}

type DeleteContainercluster struct {
	_      string `action:"delete" entity:"containercluster" awsAPI:"ecs" awsCall:"DeleteCluster" awsInput:"ecs.DeleteClusterInput" awsOutput:"ecs.DeleteClusterOutput"`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    ecsiface.ECSAPI
	Id     *string `awsName:"Cluster" awsType:"awsstr" templateName:"id" required:""`
}

func (cmd *DeleteContainercluster) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}
