# vkignite

How to run Kubernetes with Ignite support:

`$ docker-compose -f vkignite.yaml up`

Wait until `kubeconfig.yaml` is written into the current dir.

`$ export KUBECONFIG=$PWD/kubeconfig.yaml`

Run get nodes and you'll see:

```
$ kubectl get nodes
NAME              STATUS   ROLES    AGE     VERSION
k3s               Ready    master   6d22h   v1.14.3-k3s.2
virtual-kubelet   Ready    agent    4d      v1.13.7-vk-v0.11.1-1-gf350cdf9-dev
```
Then create a VM defined as a Pod inside `pod_ignite.yaml`.
```
$ kubectl apply -f pod_ignite.yaml
pod/my-vm created
```

Then you now will be able to use `sudo ignite ssh default-my-vm` to get into the VM.

To delete it, do `$ kubectl delete pod/my-vm`.
