// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package metadata_test

import (
	"fmt"

	"github.com/siderolabs/go-pointer"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "github.com/siderolabs/sidero/app/caps-controller-manager/api/v1alpha3"
	metalv1 "github.com/siderolabs/sidero/app/sidero-controller-manager/api/v1alpha2"
)

func fixture() []client.Object {
	var objects []client.Object

	for _, fixture := range []func() []client.Object{
		fixture1,
		fixture2,
		fixture3,
		fixture4,
		fixture5,
		fixture6,
		fixture7,
		fixture8,
		fixture9,
	} {
		objects = append(objects, fixture()...)
	}

	return objects
}

// fixture1 creates a server without config patches.
func fixture1() []client.Object {
	return fixtureSimple("0000-1111-2222", 1, `
version: v1alpha1
machine:
  kubelet: {}
`)
}

// fixture2 creates a server with Server-level config patches.
func fixture2() []client.Object {
	return []client.Object{
		&infrav1.ServerBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: "1111-2222-3333",
			},
			Spec: infrav1.ServerBindingSpec{
				MetalMachineRef: corev1.ObjectReference{
					Name: "metal-machine-2",
				},
			},
		},
		&infrav1.MetalMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name: "metal-machine-2",
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: "cluster.x-k8s.io/v1beta1",
						Kind:       "Machine",
						Name:       "machine-2",
					},
				},
			},
			Spec: infrav1.MetalMachineSpec{},
		},
		&capiv1.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name: "machine-2",
			},
			Spec: capiv1.MachineSpec{
				Bootstrap: capiv1.Bootstrap{
					DataSecretName: pointer.To("bootstrap2"),
				},
			},
		},
		&metalv1.Server{
			ObjectMeta: metav1.ObjectMeta{
				Name: "1111-2222-3333",
			},
			Spec: metalv1.ServerSpec{
				ConfigPatches: []metalv1.ConfigPatches{
					{
						Op:   "add",
						Path: "/machine/network",
						Value: v1.JSON{
							Raw: []byte(`{"hostname":"example2"}`),
						},
					},
				},
			},
		},
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: "bootstrap2",
			},
			Data: map[string][]byte{
				"value": []byte(`
version: v1alpha1
machine:
  kubelet:
    extraArgs:
      foo: bar
`),
			},
		},
	}
}

// fixture3 creates a server with Server- & ServerClass-level config patches.
func fixture3() []client.Object {
	return []client.Object{
		&infrav1.ServerBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: "2222-3333-4444",
			},
			Spec: infrav1.ServerBindingSpec{
				MetalMachineRef: corev1.ObjectReference{
					Name: "metal-machine-3",
				},
				ServerClassRef: &corev1.ObjectReference{
					Name: "server-class-3",
				},
			},
		},
		&infrav1.MetalMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name: "metal-machine-3",
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: "cluster.x-k8s.io/v1beta1",
						Kind:       "Machine",
						Name:       "machine-3",
					},
				},
			},
			Spec: infrav1.MetalMachineSpec{},
		},
		&capiv1.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name: "machine-3",
			},
			Spec: capiv1.MachineSpec{
				Bootstrap: capiv1.Bootstrap{
					DataSecretName: pointer.To("bootstrap3"),
				},
			},
		},
		&metalv1.ServerClass{
			ObjectMeta: metav1.ObjectMeta{
				Name: "server-class-3",
			},
			Spec: metalv1.ServerClassSpec{
				ConfigPatches: []metalv1.ConfigPatches{
					{
						Op:   "add",
						Path: "/machine/network",
						Value: v1.JSON{
							Raw: []byte(`{"hostname":"invalid3"}`),
						},
					},
				},
			},
		},
		&metalv1.Server{
			ObjectMeta: metav1.ObjectMeta{
				Name: "2222-3333-4444",
			},
			Spec: metalv1.ServerSpec{
				ConfigPatches: []metalv1.ConfigPatches{
					{
						Op:   "replace",
						Path: "/machine/network/hostname",
						Value: v1.JSON{
							Raw: []byte(`"example3"`),
						},
					},
				},
			},
		},
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: "bootstrap3",
			},
			Data: map[string][]byte{
				"value": []byte(`
version: v1alpha1
machine:
  kubelet:
    extraArgs:
      node-labels: foo=bar
`),
			},
		},
	}
}

// fixture4 creates a server machine config without machine.kubelet.
func fixture4() []client.Object {
	return fixtureSimple("4444-5555-6666", 4, `
version: v1alpha1
machine:
`)
}

// fixture5 creates a server machine config without machine.
func fixture5() []client.Object {
	return fixtureSimple("5555-6666-7777", 5, `
version: v1alpha1
cluster: {}
`)
}

// fixture6 creates a server with Server- & ServerClass-level with strategic merge config patches.
func fixture6() []client.Object {
	oldConfigPatch := "machine:\n  network:\n    hostname: invalid6"
	newConfigPatch := "machine:\n  network:\n    hostname: example6"

	return []client.Object{
		&infrav1.ServerBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: "6666-7777-8888",
			},
			Spec: infrav1.ServerBindingSpec{
				MetalMachineRef: corev1.ObjectReference{
					Name: "metal-machine-6",
				},
				ServerClassRef: &corev1.ObjectReference{
					Name: "server-class-6",
				},
			},
		},
		&infrav1.MetalMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name: "metal-machine-6",
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: "cluster.x-k8s.io/v1beta1",
						Kind:       "Machine",
						Name:       "machine-6",
					},
				},
			},
			Spec: infrav1.MetalMachineSpec{},
		},
		&capiv1.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name: "machine-6",
			},
			Spec: capiv1.MachineSpec{
				Bootstrap: capiv1.Bootstrap{
					DataSecretName: pointer.To("bootstrap6"),
				},
			},
		},
		&metalv1.ServerClass{
			ObjectMeta: metav1.ObjectMeta{
				Name: "server-class-6",
			},
			Spec: metalv1.ServerClassSpec{
				StrategicPatches: []string{oldConfigPatch},
			},
		},
		&metalv1.Server{
			ObjectMeta: metav1.ObjectMeta{
				Name: "6666-7777-8888",
			},
			Spec: metalv1.ServerSpec{
				StrategicPatches: []string{
					newConfigPatch,
					`
apiVersion: v1alpha1
kind: ExtensionServiceConfig
name: frr
environment:
- TESTKEY=TESTVALUE`,
				},
			},
		},
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: "bootstrap6",
			},
			Data: map[string][]byte{
				"value": []byte(`
version: v1alpha1
machine:
  kubelet:
    extraArgs:
      node-labels: foo=bar
`),
			},
		},
	}
}

// fixture7 creates a server with a Talos 1.14+ machine config: kubelet config lives in a
// separate `KubeletConfig` document, and it has no `extraArgs`.
func fixture7() []client.Object {
	return fixtureSimple("7777-8888-9999", 7, `
version: v1alpha1
machine: {}
---
apiVersion: v1alpha1
kind: KubeletConfig
image: ghcr.io/siderolabs/kubelet:v1.34.0
`)
}

// fixture8 creates a server with a Talos 1.14+ machine config with `node-labels` already set
// in the `KubeletConfig` document as a string.
func fixture8() []client.Object {
	return fixtureSimple("8888-9999-0000", 8, `
version: v1alpha1
machine: {}
---
apiVersion: v1alpha1
kind: KubeletConfig
image: ghcr.io/siderolabs/kubelet:v1.34.0
extraArgs:
  node-labels: foo=bar
  feature-gates: AllBeta=true
`)
}

// fixture9 creates a server with a Talos 1.14+ machine config with `node-labels` already set
// in the `KubeletConfig` document as a list.
func fixture9() []client.Object {
	return fixtureSimple("9999-0000-1111", 9, `
version: v1alpha1
machine: {}
---
apiVersion: v1alpha1
kind: KubeletConfig
image: ghcr.io/siderolabs/kubelet:v1.34.0
extraArgs:
  node-labels:
    - foo=bar
`)
}

func fixtureSimple(uuid string, index int, config string) []client.Object {
	return []client.Object{
		&infrav1.ServerBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: uuid,
			},
			Spec: infrav1.ServerBindingSpec{
				MetalMachineRef: corev1.ObjectReference{
					Name: fmt.Sprintf("metal-machine-%d", index),
				},
			},
		},
		&infrav1.MetalMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("metal-machine-%d", index),
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: "cluster.x-k8s.io/v1beta1",
						Kind:       "Machine",
						Name:       fmt.Sprintf("machine-%d", index),
					},
				},
			},
			Spec: infrav1.MetalMachineSpec{},
		},
		&capiv1.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("machine-%d", index),
			},
			Spec: capiv1.MachineSpec{
				Bootstrap: capiv1.Bootstrap{
					DataSecretName: pointer.To(fmt.Sprintf("bootstrap%d", index)),
				},
			},
		},
		&metalv1.Server{
			ObjectMeta: metav1.ObjectMeta{
				Name: uuid,
			},
			Spec: metalv1.ServerSpec{},
		},
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("bootstrap%d", index),
			},
			Data: map[string][]byte{
				"value": []byte(config),
			},
		},
	}
}
