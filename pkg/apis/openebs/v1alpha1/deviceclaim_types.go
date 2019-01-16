package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DeviceClaimSpec defines the desired state of DeviceClaim
type DeviceClaimSpec struct {
	Capacity  uint64 `json:"capacity"`  // disk size in bytes
	DriveType string `json:"driveType"` // DriveType represents the type of drive like SSD, HDD etc.,
	HostName  string `json:"hostName"`  // Node name from where device has to be claimed.
}

type PoolClaimInfo struct {
	APIVersion   string    `json:"kind,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
	Kind         string    `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`
	Name         string    `json:"name,omitempty" protobuf:"bytes,3,opt,name=name"`
	PoolClaimUID types.UID `json:"poolClaimUID" protobuf: "bytes,4,opt,name=deviceClaimUID,casttype=k8s.io/apimachinery/pkg/types.UUID"`
}

// DeviceClaimPhase is a typed string for phase field of DeviceClaim.
type DeviceClaimPhase string

/*
 * DeviceClaim CR, when created pass through phased before it got some Devices Assigned.
 * Given below table, have all phases which DeviceClaim CR can go before it is marked done.
 */
const (
	// DeviceClaimStatusEmpty: DeviceClaim CR is just created.
	DeviceClaimStatusEmpty DeviceClaimPhase = ""

	// DeviceClaimStatusPending: DeviceClaim CR yet to be assigned devices. Rather
	// search is going on for matching devices.
	DeviceClaimStatusPending DeviceClaimPhase = "Pending"

	// DeviceClaimStatusDone:  DeviceClaim CR assigned backing device and ready for use.
	DeviceClaimStatusDone DeviceClaimPhase = "Bound"
)

// DeviceClaimStatus is for handling status of pool.
type DeviceClaimStatus struct {
	Phase DeviceClaimPhase `json:"phase"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DeviceClaim is the Schema for the deviceclaims API
// +k8s:openapi-gen=true
type DeviceClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeviceClaimSpec   `json:"spec,omitempty"`
	Status DeviceClaimStatus `json:"status,omitempty"`
	Claim  PoolClaimInfo     `json:"claim,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DeviceClaimList contains a list of DeviceClaim
type DeviceClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DeviceClaim `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DeviceClaim{}, &DeviceClaimList{})
}
