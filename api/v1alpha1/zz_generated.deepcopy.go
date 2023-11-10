//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	apiv1alpha1 "github.com/spectrocloud-labs/validator/api/v1alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureValidator) DeepCopyInto(out *AzureValidator) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureValidator.
func (in *AzureValidator) DeepCopy() *AzureValidator {
	if in == nil {
		return nil
	}
	out := new(AzureValidator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AzureValidator) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureValidatorList) DeepCopyInto(out *AzureValidatorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AzureValidator, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureValidatorList.
func (in *AzureValidatorList) DeepCopy() *AzureValidatorList {
	if in == nil {
		return nil
	}
	out := new(AzureValidatorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AzureValidatorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureValidatorSpec) DeepCopyInto(out *AzureValidatorSpec) {
	*out = *in
	if in.RoleAssignmentRules != nil {
		in, out := &in.RoleAssignmentRules, &out.RoleAssignmentRules
		*out = make([]RoleAssignmentRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureValidatorSpec.
func (in *AzureValidatorSpec) DeepCopy() *AzureValidatorSpec {
	if in == nil {
		return nil
	}
	out := new(AzureValidatorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureValidatorStatus) DeepCopyInto(out *AzureValidatorStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]apiv1alpha1.ValidationCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureValidatorStatus.
func (in *AzureValidatorStatus) DeepCopy() *AzureValidatorStatus {
	if in == nil {
		return nil
	}
	out := new(AzureValidatorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Role) DeepCopyInto(out *Role) {
	*out = *in
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.RoleName != nil {
		in, out := &in.RoleName, &out.RoleName
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Role.
func (in *Role) DeepCopy() *Role {
	if in == nil {
		return nil
	}
	out := new(Role)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RoleAssignmentRule) DeepCopyInto(out *RoleAssignmentRule) {
	*out = *in
	in.Role.DeepCopyInto(&out.Role)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RoleAssignmentRule.
func (in *RoleAssignmentRule) DeepCopy() *RoleAssignmentRule {
	if in == nil {
		return nil
	}
	out := new(RoleAssignmentRule)
	in.DeepCopyInto(out)
	return out
}