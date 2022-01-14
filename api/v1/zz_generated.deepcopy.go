//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022.

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

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EurekaApplication) DeepCopyInto(out *EurekaApplication) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EurekaApplication.
func (in *EurekaApplication) DeepCopy() *EurekaApplication {
	if in == nil {
		return nil
	}
	out := new(EurekaApplication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EurekaApplication) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EurekaApplicationList) DeepCopyInto(out *EurekaApplicationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EurekaApplication, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EurekaApplicationList.
func (in *EurekaApplicationList) DeepCopy() *EurekaApplicationList {
	if in == nil {
		return nil
	}
	out := new(EurekaApplicationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EurekaApplicationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EurekaApplicationPaths) DeepCopyInto(out *EurekaApplicationPaths) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EurekaApplicationPaths.
func (in *EurekaApplicationPaths) DeepCopy() *EurekaApplicationPaths {
	if in == nil {
		return nil
	}
	out := new(EurekaApplicationPaths)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EurekaApplicationSpec) DeepCopyInto(out *EurekaApplicationSpec) {
	*out = *in
	out.Paths = in.Paths
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EurekaApplicationSpec.
func (in *EurekaApplicationSpec) DeepCopy() *EurekaApplicationSpec {
	if in == nil {
		return nil
	}
	out := new(EurekaApplicationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EurekaApplicationStatus) DeepCopyInto(out *EurekaApplicationStatus) {
	*out = *in
	if in.LastReconcileTime != nil {
		in, out := &in.LastReconcileTime, &out.LastReconcileTime
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EurekaApplicationStatus.
func (in *EurekaApplicationStatus) DeepCopy() *EurekaApplicationStatus {
	if in == nil {
		return nil
	}
	out := new(EurekaApplicationStatus)
	in.DeepCopyInto(out)
	return out
}
