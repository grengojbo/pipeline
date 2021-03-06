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

package istiofeature

import (
	"emperror.dev/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/banzaicloud/pipeline/pkg/k8sutil"
)

func (m *MeshReconciler) ReconcileCanaryOperatorNamespace(desiredState DesiredState) error {
	m.logger.Debug("reconciling canary operator namespace")
	defer m.logger.Debug("canary operator namespace reconciled")

	client, err := m.getMasterK8sClient()
	if err != nil {
		return err
	}

	if desiredState == DesiredStatePresent {
		_, err := client.CoreV1().Namespaces().Get(canaryOperatorNamespace, metav1.GetOptions{})
		if k8serrors.IsNotFound(err) {

			err := k8sutil.EnsureNamespaceWithLabelWithRetry(client, canaryOperatorNamespace,
				map[string]string{
					"istio-injection": "enabled",
				})
			if err != nil {
				return errors.WrapIf(err, "could not create canary operator namespace")
			}
		} else if err != nil {
			return errors.WrapIf(err, "could not get canary operator namespace")
		}
	} else {
		_, err := client.CoreV1().Namespaces().Get(canaryOperatorNamespace, metav1.GetOptions{})
		if k8serrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.WrapIf(err, "could not get canary operator namespace")
		}

		err = client.CoreV1().Namespaces().Delete(canaryOperatorNamespace, &metav1.DeleteOptions{})
		if err != nil {
			return errors.WrapIf(err, "could not delete canary operator namespace")
		}

		m.logger.Debug("waiting for canary operator namespace to be deleted")
		err = m.waitForNamespaceBeDeleted(client, canaryOperatorNamespace)
		if err != nil {
			return errors.WrapIf(err, "timeout during waiting for canary operator namespace to be deleted")
		}
	}

	return nil
}
