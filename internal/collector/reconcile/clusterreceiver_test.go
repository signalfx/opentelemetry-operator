// Copyright The OpenTelemetry Authors
// Copyright Splunk Inc
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

package reconcile

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/signalfx/splunk-otel-collector-operator/internal/collector"
)

func TestExpectedClusterReceivers(t *testing.T) {
	param := params()
	expectedDeploy := collector.ClusterReceiver(logger, param.Instance)

	t.Run("should create collector deployment", func(t *testing.T) {
		err := expectedClusterReceivers(context.Background(), param, []v1.Deployment{expectedDeploy})
		assert.NoError(t, err)

		exists, err := populateObjectIfExists(t, &v1.Deployment{}, types.NamespacedName{Namespace: "default", Name: "test-cluster-receiver"})

		assert.NoError(t, err)
		assert.True(t, exists)

	})

	t.Run("should update deployment", func(t *testing.T) {
		createObjectIfNotExists(t, "test-cluster-receiver", &expectedDeploy)
		err := expectedClusterReceivers(context.Background(), param, []v1.Deployment{expectedDeploy})
		assert.NoError(t, err)

		actual := v1.Deployment{}
		exists, err := populateObjectIfExists(t, &actual, types.NamespacedName{Namespace: "default", Name: "test-cluster-receiver"})

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, instanceUID, actual.OwnerReferences[0].UID)
		assert.Equal(t, int32(1), *actual.Spec.Replicas)
	})

	t.Run("should delete deployment", func(t *testing.T) {
		labels := map[string]string{
			"app.kubernetes.io/instance":   "default.test",
			"app.kubernetes.io/managed-by": "splunk-otel-operator",
		}
		deploy := v1.Deployment{}
		deploy.Name = "dummy"
		deploy.Namespace = "default"
		deploy.Labels = labels
		deploy.Spec = v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "dummy",
						Image: "busybox",
					}},
				},
			},
		}
		createObjectIfNotExists(t, "dummy", &deploy)

		err := deleteClusterReceivers(context.Background(), param, []v1.Deployment{expectedDeploy})
		assert.NoError(t, err)

		actual := v1.Deployment{}
		exists, _ := populateObjectIfExists(t, &actual, types.NamespacedName{Namespace: "default", Name: "dummy"})

		assert.False(t, exists)

	})

	t.Run("should not delete deployment", func(t *testing.T) {
		labels := map[string]string{
			"app.kubernetes.io/instance":   "default.test",
			"app.kubernetes.io/managed-by": "helm-splunk-otel-operator",
		}
		deploy := v1.Deployment{}
		deploy.Name = "dummy"
		deploy.Namespace = "default"
		deploy.Spec = v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "dummy",
						Image: "busybox",
					}},
				},
			},
		}
		createObjectIfNotExists(t, "dummy", &deploy)

		err := deleteClusterReceivers(context.Background(), param, []v1.Deployment{expectedDeploy})
		assert.NoError(t, err)

		actual := v1.Deployment{}
		exists, _ := populateObjectIfExists(t, &actual, types.NamespacedName{Namespace: "default", Name: "dummy"})

		assert.True(t, exists)

	})
}
