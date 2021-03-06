/*
 * Pipeline API
 *
 * Pipeline is a feature rich application platform, built for containers on top of Kubernetes to automate the DevOps experience, continuous application development and the lifecycle of deployments. 
 *
 * API version: latest
 * Contact: info@banzaicloud.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package pipeline

type EksNodePoolAllOf struct {

	Autoscaling NodePoolAutoScaling `json:"autoscaling,omitempty"`

	// Machine instance type.
	InstanceType string `json:"instanceType"`

	// Instance AMI.
	Image string `json:"image,omitempty"`

	// The upper limit price for the requested spot instance. If this field is left empty or 0 passed in on-demand instances used instead of spot instances.
	SpotPrice string `json:"spotPrice,omitempty"`

	Subnet EksSubnet `json:"subnet,omitempty"`
}
