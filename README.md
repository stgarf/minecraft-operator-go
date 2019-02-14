# minecraft-operator-go
A Kubernetes Operator built using operator-sdk to run Minecraft servers.

This is a basic operator that just creates and cleans up servers. If a
server happens to die, the LB Service and PVC still exist so as soon as a
new pod comes back up the same game server should be online again.

I wrote this to run on top of GKE so there might be GKE-specific things
like the PVC and LB Service that don't translate 1:1 to other providers.
If this happens to be you, and you find that it doesn't work, please file
an issue.

It's hacky, but you can use `kubectl cp` to copy `/server-data` back and
forth for backups and restores of game/world data. TODO(someone?): Implement
backup/archive logic on deletion/periodically/on-demand. 

## Usage
### 1. Build the minecraft-operator-go image and push it to a registry:
```sh
$ operator-sdk build quay.io/example/minecraft-operator-go:v1.13.2
$ docker push quay.io/example/minecraft-operator-go:v1.13.2
```

### 2. Setup RBAC and deploy the minecraft-operator:

```sh
$ kubectl create -f deploy/service_account.yaml
$ kubectl create -f deploy/role.yaml
$ kubectl create -f deploy/role_binding.yaml
$ kubectl create -f deploy/operator.yaml
```

##### Verify that the minecraft-operator-go is up and running:

```sh
$ kubectl get deployment
NAME                     DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
minecraft-operator-go       1         1         1            1           1m
```

### 3. Create a Minecraft Custom Resource

Create the `Minecraft` CR that is provided at `deploy/crds/interview_v1alpha1_minecraft_cr.yaml`:

```sh
$ cat deploy/crds/interview_v1alpha1_minecraft_cr.yaml
apiVersion: "interview.example.com/v1alpha1"
kind: "Minecraft"
metadata:
  name: "minecraft-0001"
spec:
  version: "1.13.2"

$ kubectl apply -f deploy/crds/interview_v1alpha1_minecraft_cr.yaml
```

Ensure that the minecraft-operator-go creates the deployment for the Custom Resource:

```sh
$ kubectl get deployment
NAME                     DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
minecraft-operator       1         1         1            1           2m
example-minecraft        3         3         3            3           1m
```

Check the pods and CR status to confirm the status is updated with the minecraft pod names:

```sh
$ kubectl get pods
NAME                                  READY     STATUS    RESTA
minecraft-0001-6fd7c98d8-m7vn7     1/1       Running   0          1m
minecraft-operator-go-7cc7cfdf86-vvjqk   1/1       Running   0          2m
```

### 4. Get the External IP for the LB Service

```sh
$ kubectl get services minecraft-0001
NAME             TYPE           CLUSTER-IP     EXTERNAL-IP      PORT(S)           AGE
minecraft-0001   LoadBalancer   10.11.249.80   35.192.168.145   25565:30166/TCP   18m
```

That's it...! If you've got a Minecraft client you should be able to connect to the servers that this Operator produces.
