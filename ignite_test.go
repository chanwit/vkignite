package ignite

import (
	"os"
	"os/exec"
	"sort"
	"strings"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
	"testing"

	v1 "k8s.io/api/core/v1"
)

var fixture []string

func setup() error {
	// setup 2 VMs as fixtures
	vm0, err := exec.Command("ignite", "create", "-q", "weaveworks/ignite-ubuntu").Output()
	if err != nil {
		return err
	}

	vm1, err := exec.Command("ignite", "create", "-q", "weaveworks/ignite-ubuntu").Output()
	if err != nil {
		return err
	}

	fixture = []string{
		strings.TrimSpace(string(vm1)),
		strings.TrimSpace(string(vm0)),
	}
	sort.Strings(fixture)

	return nil
}

func shutdown() {
	// clean all VMs
	exec.Command("/bin/bash", "-c", "ignite rm -f $(ignite ps -aq)").Run()
	exec.Command("dmsetup", "remove_all").Run()
}

func TestListAllIgniteIDs(t *testing.T) {
	p, err := NewIgniteProvider("", "node1", "linux", "127.0.0.1", 0)
	assert.Check(t, err, "create provider should not fail")

	vms, err := p.listAllIgniteIDs()
	sort.Strings(vms)
	assert.Check(t, err, "list all ignite IDs should not fail")
	assert.Check(t, is.Len(vms, 2), "len should be 2")
	assert.Check(t, is.Equal(vms[0], fixture[0]), "vms[0] should be "+fixture[0])
	assert.Check(t, is.Equal(vms[1], fixture[1]), "vms[1] should be "+fixture[1])
}

func TestInspectVM(t *testing.T) {
	p, err := NewIgniteProvider("", "node1", "linux", "127.0.0.1", 0)
	assert.Check(t, err, "create provider should not fail")

	vms, err := p.listAllIgniteIDs()
	assert.Check(t, err, "list all ignite IDs should not fail")

	vm0, err := p.inspectVM(vms[0])
	assert.Check(t, err, "inspect vm should not fail")
	assert.Check(t, is.Equal(vm0.Kind, "VM"))
	assert.Check(t, is.Equal(vm0.APIVersion, "ignite.weave.works/v1alpha1"))
	assert.Check(t, is.Equal(vm0.ObjectMeta.UID.String(), vms[0]))
	assert.Check(t, is.Equal(vm0.Spec.Image.OCIClaim.Ref.String(), "weaveworks/ignite-ubuntu:latest"))
	assert.Check(t, is.Equal(string(vm0.Status.State), "Created"))
}

func TestListAllIgniteVMs(t *testing.T) {
	p, err := NewIgniteProvider("", "node1", "linux", "127.0.0.1", 0)
	assert.Check(t, err, "create provider should not fail")

	vms, err := p.listAllIgniteVMs()
	assert.Check(t, err, "list all ignite VMs should not fail")
	for _, vm := range vms {
		assert.Check(t, is.Equal(vm.Kind, "VM"))
		assert.Check(t, is.Equal(vm.APIVersion, "ignite.weave.works/v1alpha1"))
		assert.Check(t, is.Equal(vm.Spec.Image.OCIClaim.Ref.String(), "weaveworks/ignite-ubuntu:latest"))
		assert.Check(t, is.Equal(string(vm.Status.State), "Created"))
	}
}

func TestCreateVMFromPod(t *testing.T) {
	p, err := NewIgniteProvider("", "node1", "linux", "127.0.0.1", 0)
	assert.Check(t, err, "create provider should not fail")

	pod := &v1.Pod{}
	pod.ObjectMeta.Namespace = "default"
	pod.ObjectMeta.Name = "test-vm"
	pod.Annotations = map[string]string{
		"ignite.weave.works/vm": `
apiVersion: ignite.weave.works/v1alpha1
kind: VM
metadata:
  name: default__test-vm
spec:
  image:
    ociClaim:
      ref: weaveworks/ignite-ubuntu
  cpus: 1
  diskSize: 1GB
  memory: 800MB
`,
	}

	vm, err := p.createVMFromPod(pod)
	assert.Check(t, err)
	assert.Check(t, is.Equal(vm.Name, "default__test-vm"))

	defer exec.Command("/bin/bash", "-c", "ignite rm -f default__test-vm").Run()
}

func TestCreateVMFromPodWithTemplate(t *testing.T) {
	p, err := NewIgniteProvider("", "node1", "linux", "127.0.0.1", 0)
	assert.Check(t, err, "create provider should not fail")

	pod := &v1.Pod{}
	pod.ObjectMeta.Namespace = "default"
	pod.ObjectMeta.Name = "test-vm-x"
	pod.Annotations = map[string]string{
		"ignite.weave.works/vm": `
apiVersion: ignite.weave.works/v1alpha1
kind: VM
metadata:
  name: {{.POD_NAMESPACE}}__{{.POD_NAME}}
spec:
  image:
    ociClaim:
      ref: weaveworks/ignite-ubuntu
  cpus: 1
  diskSize: 1GB
  memory: 800MB
`,
	}

	vm, err := p.createVMFromPod(pod)
	assert.Check(t, err)
	assert.Check(t, is.Equal(vm.Name, "default__test-vm-x"))

	defer exec.Command("/bin/bash", "-c", "ignite rm -f default__test-vm-x").Run()
}

func TestMain(m *testing.M) {
	err := setup()
	code := 0
	if err == nil {
		code = m.Run()
	}

	shutdown()
	os.Exit(code)
}
