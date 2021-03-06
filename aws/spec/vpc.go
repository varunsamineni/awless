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
	"net"

	"github.com/wallix/awless/cloud/graph"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/wallix/awless/logger"
)

type CreateVpc struct {
	_      string `action:"create" entity:"vpc" awsAPI:"ec2" awsCall:"CreateVpc" awsInput:"ec2.CreateVpcInput" awsOutput:"ec2.CreateVpcOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    ec2iface.EC2API
	CIDR   *string `awsName:"CidrBlock" awsType:"awsstr" templateName:"cidr" required:""`
	Name   *string `awsName:"Name" templateName:"name"`
}

func (cmd *CreateVpc) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CreateVpc) Validate_CIDR() error {
	_, _, err := net.ParseCIDR(StringValue(cmd.CIDR))
	return err
}

func (cmd *CreateVpc) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.CreateVpcOutput).Vpc.VpcId)
}

func (cmd *CreateVpc) AfterRun(ctx map[string]interface{}, output interface{}) error {
	return createNameTag(awssdk.String(cmd.ExtractResult(output)), cmd.Name, ctx)
}

type DeleteVpc struct {
	_      string `action:"delete" entity:"vpc" awsAPI:"ec2" awsCall:"DeleteVpc" awsInput:"ec2.DeleteVpcInput" awsOutput:"ec2.DeleteVpcOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    ec2iface.EC2API
	Id     *string `awsName:"VpcId" awsType:"awsstr" templateName:"id" required:""`
}

func (cmd *DeleteVpc) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}
