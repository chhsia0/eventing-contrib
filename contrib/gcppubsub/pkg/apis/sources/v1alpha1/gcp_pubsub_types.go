/*
Copyright 2018 The Knative Authors

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

package v1alpha1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck"
	duckv1alpha1 "knative.dev/pkg/apis/duck/v1alpha1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GcpPubSubSource is the Schema for the gcppubsubsources API.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:categories=all,knative,eventing,sources
type GcpPubSubSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GcpPubSubSourceSpec   `json:"spec,omitempty"`
	Status GcpPubSubSourceStatus `json:"status,omitempty"`
}

// Check that GcpPubSubSource can be validated and can be defaulted.
var _ runtime.Object = (*GcpPubSubSource)(nil)

// Check that GcpPubSubSource will be checked for immutable fields.
var _ apis.Immutable = (*GcpPubSubSource)(nil)

// Check that GcpPubSubSource implements the Conditions duck type.
var _ = duck.VerifyType(&GcpPubSubSource{}, &duckv1alpha1.Conditions{})

// GcpPubSubSourceSpec defines the desired state of the GcpPubSubSource.
type GcpPubSubSourceSpec struct {
	// GcpCredsSecret is the credential to use to poll the GCP PubSub Subscription. It is not used
	// to create or delete the Subscription, only to poll it. The value of the secret entry must be
	// a service account key in the JSON format
	// ( see https://cloud.google.com/iam/docs/creating-managing-service-account-keys ).
	GcpCredsSecret corev1.SecretKeySelector `json:"gcpCredsSecret,omitempty"`

	// GoogleCloudProject is the ID of the Google Cloud Project that the PubSub Topic exists in.
	GoogleCloudProject string `json:"googleCloudProject,omitempty"`

	// Topic is the ID of the GCP PubSub Topic to Subscribe to. It must be in the form of the
	// unique identifier within the project, not the entire name. E.g. it must be 'laconia', not
	// 'projects/my-gcp-project/topics/laconia'.
	Topic string `json:"topic,omitempty"`

	// Sink is a reference to an object that will resolve to a domain name to use as the sink.
	// +optional
	Sink *corev1.ObjectReference `json:"sink,omitempty"`

	// Transformer is a reference to an object that will resolve to a domain name to use as the transformer.
	// +optional
	Transformer *corev1.ObjectReference `json:"transformer,omitempty"`

	// ServiceAccoutName is the name of the ServiceAccount that will be used to run the Receive
	// Adapter Deployment.
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
}

const (
	// GcpPubSubSourceEventType is the GcpPubSub CloudEvent type, in case PubSub doesn't send a
	// CloudEvent itself.
	GcpPubSubSourceEventType = "google.pubsub.topic.publish"
)

// GcpPubSubEventSource returns the GcpPubSub CloudEvent source value.
func GcpPubSubEventSource(googleCloudProject, topic string) string {
	return fmt.Sprintf("//pubsub.googleapis.com/%s/topics/%s", googleCloudProject, topic)
}

const (
	// GcpPubSubConditionReady has status True when the GcpPubSubSource is ready to send events.
	GcpPubSubConditionReady = duckv1alpha1.ConditionReady

	// GcpPubSubConditionSinkProvided has status True when the GcpPubSubSource has been configured with a sink target.
	GcpPubSubConditionSinkProvided duckv1alpha1.ConditionType = "SinkProvided"

	// GcpPubSubConditionTransformerProvided has status True when the GcpPubSubSource has been configured with a transformer target.
	GcpPubSubConditionTransformerProvided duckv1alpha1.ConditionType = "TransformerProvided"

	// GcpPubSubConditionDeployed has status True when the GcpPubSubSource has had it's receive adapter deployment created.
	GcpPubSubConditionDeployed duckv1alpha1.ConditionType = "Deployed"

	// GcpPubSubConditionSubscribed has status True when a GCP PubSub Subscription has been created pointing at the created receive adapter deployment.
	GcpPubSubConditionSubscribed duckv1alpha1.ConditionType = "Subscribed"

	// GcpPubSubConditionEventTypesProvided has status True when the GcpPubSubSource has been configured with event types.
	GcpPubSubConditionEventTypesProvided duckv1alpha1.ConditionType = "EventTypesProvided"
)

var gcpPubSubSourceCondSet = duckv1alpha1.NewLivingConditionSet(
	GcpPubSubConditionSinkProvided,
	GcpPubSubConditionDeployed,
	GcpPubSubConditionSubscribed)

// GcpPubSubSourceStatus defines the observed state of GcpPubSubSource.
type GcpPubSubSourceStatus struct {
	// inherits duck/v1alpha1 Status, which currently provides:
	// * ObservedGeneration - the 'Generation' of the Service that was last processed by the controller.
	// * Conditions - the latest available observations of a resource's current state.
	duckv1alpha1.Status `json:",inline"`

	// SinkURI is the current active sink URI that has been configured for the GcpPubSubSource.
	// +optional
	SinkURI string `json:"sinkUri,omitempty"`

	// TransformerURI is the current active transformer URI that has been configured for the GcpPubSubSource.
	// +optional
	TransformerURI string `json:"transformerUri,omitempty"`
}

// GetCondition returns the condition currently associated with the given type, or nil.
func (s *GcpPubSubSourceStatus) GetCondition(t duckv1alpha1.ConditionType) *duckv1alpha1.Condition {
	return gcpPubSubSourceCondSet.Manage(s).GetCondition(t)
}

// IsReady returns true if the resource is ready overall.
func (s *GcpPubSubSourceStatus) IsReady() bool {
	return gcpPubSubSourceCondSet.Manage(s).IsHappy()
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (s *GcpPubSubSourceStatus) InitializeConditions() {
	gcpPubSubSourceCondSet.Manage(s).InitializeConditions()
}

// MarkSink sets the condition that the source has a sink configured.
func (s *GcpPubSubSourceStatus) MarkSink(uri string) {
	s.SinkURI = uri
	if len(uri) > 0 {
		gcpPubSubSourceCondSet.Manage(s).MarkTrue(GcpPubSubConditionSinkProvided)
	} else {
		gcpPubSubSourceCondSet.Manage(s).MarkUnknown(GcpPubSubConditionSinkProvided, "SinkEmpty", "Sink has resolved to empty.")
	}
}

// MarkSink sets the condition that the source has a transformer configured.
func (s *GcpPubSubSourceStatus) MarkTransformer(uri string) {
	s.TransformerURI = uri
	if len(uri) > 0 {
		gcpPubSubSourceCondSet.Manage(s).MarkTrue(GcpPubSubConditionTransformerProvided)
	} else {
		gcpPubSubSourceCondSet.Manage(s).MarkUnknown(GcpPubSubConditionTransformerProvided, "TransformerEmpty", "Transformer has resolved to empty.")
	}
}

// MarkNoSink sets the condition that the source does not have a sink configured.
func (s *GcpPubSubSourceStatus) MarkNoSink(reason, messageFormat string, messageA ...interface{}) {
	gcpPubSubSourceCondSet.Manage(s).MarkFalse(GcpPubSubConditionSinkProvided, reason, messageFormat, messageA...)
}

// MarkNoTransformer sets the condition that the source does not have a transformer configured.
func (s *GcpPubSubSourceStatus) MarkNoTransformer(reason, messageFormat string, messageA ...interface{}) {
	gcpPubSubSourceCondSet.Manage(s).MarkFalse(GcpPubSubConditionTransformerProvided, reason, messageFormat, messageA...)
}

// MarkDeployed sets the condition that the source has been deployed.
func (s *GcpPubSubSourceStatus) MarkDeployed() {
	gcpPubSubSourceCondSet.Manage(s).MarkTrue(GcpPubSubConditionDeployed)
}

// MarkDeploying sets the condition that the source is deploying.
func (s *GcpPubSubSourceStatus) MarkDeploying(reason, messageFormat string, messageA ...interface{}) {
	gcpPubSubSourceCondSet.Manage(s).MarkUnknown(GcpPubSubConditionDeployed, reason, messageFormat, messageA...)
}

// MarkNotDeployed sets the condition that the source has not been deployed.
func (s *GcpPubSubSourceStatus) MarkNotDeployed(reason, messageFormat string, messageA ...interface{}) {
	gcpPubSubSourceCondSet.Manage(s).MarkFalse(GcpPubSubConditionDeployed, reason, messageFormat, messageA...)
}

func (s *GcpPubSubSourceStatus) MarkSubscribed() {
	gcpPubSubSourceCondSet.Manage(s).MarkTrue(GcpPubSubConditionSubscribed)
}

// MarkEventTypes sets the condition that the source has created its event types.
func (s *GcpPubSubSourceStatus) MarkEventTypes() {
	gcpPubSubSourceCondSet.Manage(s).MarkTrue(GcpPubSubConditionEventTypesProvided)
}

// MarkNoEventTypes sets the condition that the source does not its event types configured.
func (s *GcpPubSubSourceStatus) MarkNoEventTypes(reason, messageFormat string, messageA ...interface{}) {
	gcpPubSubSourceCondSet.Manage(s).MarkFalse(GcpPubSubConditionEventTypesProvided, reason, messageFormat, messageA...)
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GcpPubSubSourceList contains a list of GcpPubSubSources.
type GcpPubSubSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GcpPubSubSource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GcpPubSubSource{}, &GcpPubSubSourceList{})
}
