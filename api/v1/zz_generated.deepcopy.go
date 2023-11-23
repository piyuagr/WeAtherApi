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

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Weatherapi) DeepCopyInto(out *Weatherapi) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Weatherapi.
func (in *Weatherapi) DeepCopy() *Weatherapi {
	if in == nil {
		return nil
	}
	out := new(Weatherapi)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Weatherapi) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WeatherapiList) DeepCopyInto(out *WeatherapiList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Weatherapi, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WeatherapiList.
func (in *WeatherapiList) DeepCopy() *WeatherapiList {
	if in == nil {
		return nil
	}
	out := new(WeatherapiList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WeatherapiList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WeatherapiSpec) DeepCopyInto(out *WeatherapiSpec) {
	*out = *in
	out.Tempmin = in.Tempmin.DeepCopy()
	out.Tempmax = in.Tempmax.DeepCopy()
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WeatherapiSpec.
func (in *WeatherapiSpec) DeepCopy() *WeatherapiSpec {
	if in == nil {
		return nil
	}
	out := new(WeatherapiSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WeatherapiStatus) DeepCopyInto(out *WeatherapiStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WeatherapiStatus.
func (in *WeatherapiStatus) DeepCopy() *WeatherapiStatus {
	if in == nil {
		return nil
	}
	out := new(WeatherapiStatus)
	in.DeepCopyInto(out)
	return out
}
