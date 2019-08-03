# vkignite

How to run Kubernetes with Ignite support:

`$ docker-compose -f vkignite.yaml up`

Wait until `kubeconfig.yaml` is written into the current dir.

`$ export KUBECONFIG=$PWD/kubeconfig.yaml`

Run get nodes and you'll see:

```
$ kubectl get nodes
NAME     STATUS   ROLES   AGE     VERSION
ignite   Ready    agent   5h12m   v1.14.3-ignite-f5516fe-dev
```
Then create a VM defined as a Pod inside `pod_ignite.yaml`.
```
$ kubectl apply -f pod_ignite.yaml
pod/my-vm created
```

Then you now will be able to use `sudo ignite ssh default-my-vm` to get into the VM.

To delete it, do `$ kubectl delete pod/my-vm`.
