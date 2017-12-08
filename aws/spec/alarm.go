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
	"errors"
	"fmt"

	"github.com/wallix/awless/cloud/graph"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/wallix/awless/logger"
)

type CreateAlarm struct {
	_                       string `action:"create" entity:"alarm" awsAPI:"cloudwatch" awsCall:"PutMetricAlarm" awsInput:"cloudwatch.PutMetricAlarmInput" awsOutput:"cloudwatch.PutMetricAlarmOutput"`
	logger                  *logger.Logger
	graph                   cloudgraph.GraphAPI
	api                     cloudwatchiface.CloudWatchAPI
	Name                    *string   `awsName:"AlarmName" awsType:"awsstr" templateName:"name" required:""`
	Operator                *string   `awsName:"ComparisonOperator" awsType:"awsstr" templateName:"operator" required:""`
	Metric                  *string   `awsName:"MetricName" awsType:"awsstr" templateName:"metric" required:""`
	Namespace               *string   `awsName:"Namespace" awsType:"awsstr" templateName:"namespace" required:""`
	EvaluationPeriods       *int64    `awsName:"EvaluationPeriods" awsType:"awsint64" templateName:"evaluation-periods" required:""`
	Period                  *int64    `awsName:"Period" awsType:"awsint64" templateName:"period" required:""`
	StatisticFunction       *string   `awsName:"Statistic" awsType:"awsstr" templateName:"statistic-function" required:""`
	Threshold               *float64  `awsName:"Threshold" awsType:"awsfloat" templateName:"threshold" required:""`
	Enabled                 *bool     `awsName:"ActionsEnabled" awsType:"awsbool" templateName:"enabled"`
	AlarmActions            []*string `awsName:"AlarmActions" awsType:"awsstringslice" templateName:"alarm-actions"`
	InsufficientdataActions []*string `awsName:"InsufficientDataActions" awsType:"awsstringslice" templateName:"insufficientdata-actions"`
	OkActions               []*string `awsName:"OKActions" awsType:"awsstringslice" templateName:"ok-actions"`
	Description             *string   `awsName:"AlarmDescription" awsType:"awsstr" templateName:"description"`
	Dimensions              []*string `awsName:"Dimensions" awsType:"awsdimensionslice" templateName:"dimensions"`
	Unit                    *string   `awsName:"Unit" awsType:"awsstr" templateName:"unit"`
}

func (cmd *CreateAlarm) Validate_Operator() error {
	return NewEnumValidator("GreaterThanThreshold", "LessThanThreshold", "LessThanOrEqualToThreshold", "GreaterThanOrEqualToThreshold").Validate(cmd.Operator)
}

func (cmd *CreateAlarm) ExtractResult(i interface{}) string {
	return StringValue(cmd.Name)
}

type DeleteAlarm struct {
	_      string `action:"delete" entity:"alarm" awsAPI:"cloudwatch" awsCall:"DeleteAlarms" awsInput:"cloudwatch.DeleteAlarmsInput" awsOutput:"cloudwatch.DeleteAlarmsOutput"`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    cloudwatchiface.CloudWatchAPI
	Name   []*string `awsName:"AlarmNames" awsType:"awsstringslice" templateName:"name" required:""`
}

func (cmd *DeleteAlarm) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

type StartAlarm struct {
	_      string `action:"start" entity:"alarm" awsAPI:"cloudwatch" awsCall:"EnableAlarmActions" awsInput:"cloudwatch.EnableAlarmActionsInput" awsOutput:"cloudwatch.EnableAlarmActionsOutput"`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    cloudwatchiface.CloudWatchAPI
	Names  []*string `awsName:"AlarmNames" awsType:"awsstringslice" templateName:"names" required:""`
}

func (cmd *StartAlarm) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

type StopAlarm struct {
	_      string `action:"stop" entity:"alarm" awsAPI:"cloudwatch" awsCall:"DisableAlarmActions" awsInput:"cloudwatch.DisableAlarmActionsInput" awsOutput:"cloudwatch.DisableAlarmActionsOutput"`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    cloudwatchiface.CloudWatchAPI
	Names  []*string `awsName:"AlarmNames" awsType:"awsstringslice" templateName:"names" required:""`
}

func (cmd *StopAlarm) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

type AttachAlarm struct {
	_         string `action:"attach" entity:"alarm" awsAPI:"cloudwatch"`
	logger    *logger.Logger
	graph     cloudgraph.GraphAPI
	api       cloudwatchiface.CloudWatchAPI
	Name      *string `templateName:"name" required:""`
	ActionArn *string `templateName:"action-arn" required:""`
}

func (cmd *AttachAlarm) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *AttachAlarm) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	alarm, err := getAlarm(cmd.api, cmd.Name)
	if err != nil {
		return nil, err
	}
	alarm.AlarmActions = append(alarm.AlarmActions, cmd.ActionArn)

	return cmd.api.PutMetricAlarm(&cloudwatch.PutMetricAlarmInput{
		ActionsEnabled:                   alarm.ActionsEnabled,
		AlarmActions:                     alarm.AlarmActions,
		AlarmDescription:                 alarm.AlarmDescription,
		AlarmName:                        alarm.AlarmName,
		ComparisonOperator:               alarm.ComparisonOperator,
		Dimensions:                       alarm.Dimensions,
		EvaluateLowSampleCountPercentile: alarm.EvaluateLowSampleCountPercentile,
		EvaluationPeriods:                alarm.EvaluationPeriods,
		ExtendedStatistic:                alarm.ExtendedStatistic,
		InsufficientDataActions:          alarm.InsufficientDataActions,
		MetricName:                       alarm.MetricName,
		Namespace:                        alarm.Namespace,
		OKActions:                        alarm.OKActions,
		Period:                           alarm.Period,
		Statistic:                        alarm.Statistic,
		Threshold:                        alarm.Threshold,
		TreatMissingData:                 alarm.TreatMissingData,
		Unit:                             alarm.Unit,
	})
}

type DetachAlarm struct {
	_         string `action:"detach" entity:"alarm" awsAPI:"cloudwatch"`
	logger    *logger.Logger
	graph     cloudgraph.GraphAPI
	api       cloudwatchiface.CloudWatchAPI
	Name      *string `templateName:"name" required:""`
	ActionArn *string `templateName:"action-arn" required:""`
}

func (cmd *DetachAlarm) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *DetachAlarm) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	alarm, err := getAlarm(cmd.api, cmd.Name)
	if err != nil {
		return nil, err
	}
	actionArn := aws.StringValue(cmd.ActionArn)
	var found bool
	var updatedActions []*string
	for _, action := range alarm.AlarmActions {
		if aws.StringValue(action) == actionArn {
			found = true
		} else {
			updatedActions = append(updatedActions, action)
		}
	}
	if !found {
		return nil, fmt.Errorf("detach alarm: action '%s' is not attached to alarm actions of alarm %s", actionArn, aws.StringValue(alarm.AlarmName))
	}

	return cmd.api.PutMetricAlarm(&cloudwatch.PutMetricAlarmInput{
		ActionsEnabled:                   alarm.ActionsEnabled,
		AlarmActions:                     updatedActions,
		AlarmDescription:                 alarm.AlarmDescription,
		AlarmName:                        alarm.AlarmName,
		ComparisonOperator:               alarm.ComparisonOperator,
		Dimensions:                       alarm.Dimensions,
		EvaluateLowSampleCountPercentile: alarm.EvaluateLowSampleCountPercentile,
		EvaluationPeriods:                alarm.EvaluationPeriods,
		ExtendedStatistic:                alarm.ExtendedStatistic,
		InsufficientDataActions:          alarm.InsufficientDataActions,
		MetricName:                       alarm.MetricName,
		Namespace:                        alarm.Namespace,
		OKActions:                        alarm.OKActions,
		Period:                           alarm.Period,
		Statistic:                        alarm.Statistic,
		Threshold:                        alarm.Threshold,
		TreatMissingData:                 alarm.TreatMissingData,
		Unit:                             alarm.Unit,
	})
}

func getAlarm(api cloudwatchiface.CloudWatchAPI, name *string) (*cloudwatch.MetricAlarm, error) {
	if name == nil {
		return nil, errors.New("missing required params 'name'")
	}
	out, err := api.DescribeAlarms(&cloudwatch.DescribeAlarmsInput{AlarmNames: []*string{name}})
	if err != nil {
		return nil, err
	}
	if l := len(out.MetricAlarms); l == 0 {
		return nil, fmt.Errorf("alarm '%s' not found", StringValue(name))
	} else if l > 1 {
		return nil, fmt.Errorf("%d alarms found with name '%s'", l, StringValue(name))
	}
	return out.MetricAlarms[0], nil
}
