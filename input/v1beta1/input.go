// Package v1beta1 contains the input type for this Function
// +kubebuilder:object:generate=true
// +groupName=template.fn.crossplane.io
// +versionName=v1beta1
package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// This isn't a custom resource, in the sense that we never install its CRD.
// It is a KRM-like object, so we generate a CRD to describe its schema.

// TODO: Add your input type here! It doesn't need to be called 'RandString', you can
// rename it to anything you like.

type Object struct {
	Name      string `json:"name"`
	FieldPath string `json:"fieldPath"`
	Prefix    string `json:"prefix,omitempty"`
}

type RandomString struct {
	Length int `json:"length"`
}

type Config struct {
	Objs    []Object     `json:"objects"`
	RandStr RandomString `json:"randomString"`
}

// RandString can be used to provide input to this Function.
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:categories=crossplane
type RandString struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// this always returns json response
	// +kubebuilder:pruning:PreserverUnknownFields
	Cfg Config `json:"config"`
}

type composedResource struct {
	Name string                `json:"name"`
	Base *runtime.RawExtension `json:"base,omitempty"`
}
