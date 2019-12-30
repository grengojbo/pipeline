// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"emperror.dev/errors"
	"github.com/sirupsen/logrus"

	"github.com/banzaicloud/pipeline/internal/cloudinfo"
	"github.com/banzaicloud/pipeline/internal/global"
	pipelineContext "github.com/banzaicloud/pipeline/internal/platform/context"
	"github.com/banzaicloud/pipeline/pkg/common"
)

const labelFormatRegexp = "[^-A-Za-z0-9_.]"

type NodePoolLabels struct {
	NodePoolName string
	Existing     bool
	InstanceType string            `json:"instanceType,omitempty"`
	SpotPrice    string            `json:"spotPrice,omitempty"`
	Preemptible  bool              `json:"preemptible,omitempty"`
	CustomLabels map[string]string `json:"labels,omitempty"`
}

// GetDesiredLabelsForCluster returns desired set of labels for each node pool name, adding Banzaicloud prefixed labels like:
// head node, ondemand labels + cloudinfo to user defined labels in specified nodePools map.
// noReturnIfNoUserLabels = true, means if there are no labels specified in NodePoolStatus, no labels are returned for that node pool
// is not returned, to avoid overriding already exisisting user specified labels.
func GetDesiredLabelsForCluster(ctx context.Context, cluster CommonCluster, nodePoolLabels []NodePoolLabels) (map[string]map[string]string, error) {
	logger := pipelineContext.LoggerWithCorrelationID(ctx, log).WithFields(logrus.Fields{
		"organization": cluster.GetOrganizationId(),
		"cluster":      cluster.GetID(),
	})

	desiredLabels := make(map[string]map[string]string)

	clusterStatus, err := cluster.GetStatus()
	if err != nil {
		return desiredLabels, errors.WrapIfWithDetails(err, "failed to get cluster status", "cluster", cluster.GetName())
	}

	for _, npLabels := range nodePoolLabels {
		noReturnIfNoUserLabels := npLabels.Existing
		labelsMap := getLabelsForNodePool(logger, npLabels.NodePoolName, npLabels, noReturnIfNoUserLabels,
			clusterStatus.Cloud, clusterStatus.Distribution, clusterStatus.Region)
		if len(labelsMap) > 0 {
			desiredLabels[npLabels.NodePoolName] = labelsMap
		}
	}
	return desiredLabels, nil
}

func formatValue(value string) string {
	var re = regexp.MustCompile(labelFormatRegexp)
	norm := re.ReplaceAllString(value, "_")
	return norm
}

func getLabelsForNodePool(
	logger logrus.FieldLogger,
	nodePoolName string,
	nodePool NodePoolLabels,
	noReturnIfNoUserLabels bool,
	cloud string,
	distribution string,
	region string,
) map[string]string {

	desiredLabels := make(map[string]string)
	if len(nodePool.CustomLabels) == 0 && noReturnIfNoUserLabels {
		return desiredLabels
	}

	desiredLabels[common.LabelKey] = nodePoolName
	desiredLabels[common.OnDemandLabelKey] = getOnDemandLabel(nodePool)

	// copy user labels unless they are not reserved keys
	for labelKey, labelValue := range nodePool.CustomLabels {
		if !IsReservedDomainKey(labelKey) {
			desiredLabels[labelKey] = labelValue
		}
	}

	// get CloudInfo labels for node
	machineDetails, err := cloudinfo.GetMachineDetails(logger, cloud,
		distribution,
		region,
		nodePool.InstanceType)
	if err != nil {
		log.WithFields(logrus.Fields{
			"instance":     nodePool.InstanceType,
			"cloud":        cloud,
			"distribution": distribution,
			"region":       region,
		}).Warn(errors.Wrap(err, "failed to get instance attributes from Cloud Info"))
	} else {
		if machineDetails != nil {
			for attrKey, attrValue := range machineDetails.Attributes {
				nKey := formatValue(attrKey)
				cloudInfoAttrKey := common.CloudInfoLabelKeyPrefix + nKey
				nValue := formatValue(attrValue)
				desiredLabels[cloudInfoAttrKey] = nValue
			}
		}
	}

	return desiredLabels
}

func IsReservedDomainKey(labelKey string) bool {
	pipelineLabelDomain := global.Config.Cluster.Labels.Domain
	if strings.Contains(labelKey, pipelineLabelDomain) {
		return true
	}

	reservedNodeLabelDomains := global.Config.Cluster.Labels.ForbiddenDomains
	for _, reservedDomain := range reservedNodeLabelDomains {
		if strings.Contains(labelKey, reservedDomain) {
			return true
		}
	}
	return false
}

func getOnDemandLabel(nodePool NodePoolLabels) string {
	if p, err := strconv.ParseFloat(nodePool.SpotPrice, 64); err == nil && p > 0.0 {
		return "false"
	}
	if nodePool.Preemptible {
		return "false"
	}
	return "true"
}
