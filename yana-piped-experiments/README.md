# Table of Contents
- [Environment](#environment)
- [Setup KubeStellar](#setup-kubeStellar)
- [Setup PipeCD and canary deployment](#setup-pipecd-and-example-scenario)
- [Cleanup](#cleanup)

## Environment

- OS: darwin
- Arch: arm64
- kubectl
    ```
    $ kubectl version
    Client Version: version.Info{Major:"1", Minor:"26", GitVersion:"v1.26.3", GitCommit:"9e644106593f3f4aa98f8a84b23db5fa378900bd", GitTreeState:"clean", BuildDate:"2023-03-15T13:40:17Z", GoVersion:"go1.19.7", Compiler:"gc", Platform:"darwin/amd64"}
    Kustomize Version: v4.5.7
    ```
- kind
    ```
    $ kind version
    kind v0.19.0 go1.20.4 darwin/arm64
    ```
- helm
    ```
    $ helm version
    version.BuildInfo{Version:"v3.13.1", GitCommit:"3547a4b5bf5edb5478ce352e18858d8a552a4110", GitTreeState:"clean", GoVersion:"go1.20.8"}
    ```
- clusteradm (https://github.com/open-cluster-management-io/clusteradm)
    ```
    $ clusteradm version
    client          version :v0.8.1-0-g3aea9c5
    Error: Get "http://localhost:8080/version?timeout=32s": dial tcp [::1]:8080: connect: connection refused
    ```
- pipectl
    ```
    $ pipectl version
    Version: v0.47.1, GitCommit: 9b136366d583c7670565092a52df6aad24e0f4cc, BuildDate: 20240422-051739
    ```
- piped
    ```
    $ piped version
    Version: v0.47.1, GitCommit: 9b136366d583c7670565092a52df6aad24e0f4cc, BuildDate: 20240422-051626
    ```
## Setup KubeStellar

Basically following to the same step as https://docs.kubestellar.io/release-0.22.0/direct/deploy-on-k3d/ but modifying it by KinD

1. Set environment variables to hold KubeStellar and OCM-status-addon desired versions: 
    ```
    export KUBESTELLAR_VERSION="0.22.0"
    export OCM_STATUS_ADDON_VERSION="0.2.0-rc8"
    export OCM_TRANSPORT_PLUGIN="0.1.7"
    ```
1. Install KubeFlex (please modify the download link for OS and Arch on your environment)
    ```
    curl -LO https://github.com/kubestellar/kubeflex/releases/download/v0.6.0/kubeflex_0.6.0_darwin_arm64.tar.gz
    tar zxvf kubeflex_0.6.0_darwin_arm64.tar.gz
    sudo cp ./bin/kflex /usr/local/bin/kflex
    ```
1. Create KinD cluster by `kflex` (this command automatically creates KinD cluster that has an Ingress controller that is listening on host port 9443 and configured with TLS passthrough)
    ```
    kflex init --create-kind
    ```
1. Create the `its1` space with OCM running in it
    ```
    kubectl apply -f https://raw.githubusercontent.com/kubestellar/kubestellar/v${KUBESTELLAR_VERSION}/config/postcreate-hooks/kubestellar.yaml
    kubectl apply -f https://raw.githubusercontent.com/kubestellar/kubestellar/v${KUBESTELLAR_VERSION}/config/postcreate-hooks/ocm.yaml
    kflex create its1 --type vcluster -p ocm
    ```
    1. Wait until ocm is ready (~1 min.)
        ```
        kubectl --context its1 api-resources | grep managedclusteraddons 
        ```
        e.g.
        ```
        $ kubectl --context its1 api-resources | grep managedclusteraddons 
        managedclusteraddons              mca,mcas                       addon.open-cluster-management.io/v1alpha1     true         ManagedClusterAddOn
        ```
1. Install OCM Status addon
    ```
    helm --kube-context its1 upgrade --install status-addon -n open-cluster-management oci://ghcr.io/kubestellar/ocm-status-addon-chart --version v${OCM_STATUS_ADDON_VERSION}
    ```
1. Create a Workload Description Space wds1 in KubeFlex. 
    ```
    kflex create wds1 -p kubestellar
    ```
1. Deploy the OCM based transport controller.
    ```
    helm --kube-context kind-kubeflex upgrade --install ocm-transport-plugin oci://ghcr.io/kubestellar/ocm-transport-plugin/chart/ocm-transport-plugin --version ${OCM_TRANSPORT_PLUGIN} \
    --set transport_cp_name=its1 \
    --set wds_cp_name=wds1 \
    -n wds1-system
    ```
1. Create the Workload Execution Cluster (WEC) `cluster1` and register it Make sure `cluster1` shares the same docker network as the kubeflex hosting cluster.
    ```
    cat << EOL > cluster1.yaml
    kind: Cluster
    apiVersion: kind.x-k8s.io/v1alpha4
    networking:
      ipFamily: dual
    nodes:
    - role: control-plane
      extraPortMappings:
      - containerPort: 80
        hostPort: 31080
        protocol: TCP
    EOL
    kind create cluster --config ./cluster1.yaml --name cluster1 --wait 5m
    kubectl config rename-context kind-cluster1 cluster1
    ```
    1. Register `cluster1`
        ```
        flags="--force-internal-endpoint-lookup"
        clusteradm --context its1 get token | grep '^clusteradm join' | sed "s/<cluster_name>/cluster1/" | awk '{print $0 " --context 'cluster1' '${flags}'"}'  | sh
        ```
    1. Wait for csr to be created (taking ~30s)
        ```
        kubectl --context its1 get csr --watch
        ```
        ``` e.g.
        $ kubectl --context its1 get csr --watch
        NAME             AGE   SIGNERNAME                            REQUESTOR                                                         REQUESTEDDURATION   CONDITION
        cluster1-9l2wz   0s    kubernetes.io/kube-apiserver-client   system:serviceaccount:open-cluster-management:cluster-bootstrap   <none>              Pending
        ```
    1. then accept pending `cluster1` cluster
        ```
        clusteradm --context its1 accept --clusters cluster1
        ```
    1. Confirm cluster1 is accepted and label it for the BindingPolicy:
        ```
        kubectl --context its1 get managedclusters
        kubectl --context its1 label managedcluster cluster1 location-group=edge name=cluster1
        ```
1. In the same way, make another WEC named "cluster2". 
    ```
    cat << EOL > cluster2.yaml
    kind: Cluster
    apiVersion: kind.x-k8s.io/v1alpha4
    networking:
      ipFamily: dual
    nodes:
    - role: control-plane
      extraPortMappings:
      - containerPort: 80
        hostPort: 31180
        protocol: TCP
    EOL
    kind create cluster --config ./cluster2.yaml --name cluster2 --wait 5m
    kubectl config rename-context kind-cluster2 cluster2
    ```
    1. Register `cluster2`
        ```
        flags="--force-internal-endpoint-lookup"
        clusteradm --context its1 get token | grep '^clusteradm join' | sed "s/<cluster_name>/cluster2/" | awk '{print $0 " --context 'cluster2' '${flags}'"}'  | sh
        ```
    1. Wait for csr to be created (taking ~30s)
        ```
        kubectl --context its1 get csr --watch
        ```
        ``` e.g.
        $ kubectl --context its1 get csr --watch
        NAME                                AGE     SIGNERNAME                            REQUESTOR                                                         REQUESTEDDURATION   CONDITION
        cluster1-9l2wz                      2m56s   kubernetes.io/kube-apiserver-client   system:serviceaccount:open-cluster-management:cluster-bootstrap   <none>              Approved,Issued
        addon-cluster1-addon-status-vc5dn   85s     kubernetes.io/kube-apiserver-client   system:open-cluster-management:cluster1:w4frl                     <none>              Approved,Issued
        cluster2-5bnzr                      0s      kubernetes.io/kube-apiserver-client   system:serviceaccount:open-cluster-management:cluster-bootstrap   <none>              Pending
        ```
    1. then accept pending `cluster2` cluster
        ```
        clusteradm --context its1 accept --clusters cluster2
        ```
    1. Confirm cluster2 is accepted and label it for the BindingPolicy:
        ```
        kubectl --context its1 get managedclusters
        kubectl --context its1 label managedcluster cluster2 location-group=edge name=cluster2
        ```
1. (optional) Check relevant deployments and statefulsets running in the hosting cluster. Expect to see the kubestellar-controller-manager in the wds1-system namespace and the statefulset vcluster in the its1-system namespace, both fully ready.
    ```
    kubectl --context kind-kubeflex get deployments,statefulsets --all-namespaces
    ```
    e.g.
    ```
    $ kubectl --context kind-kubeflex get deployments,statefulsets --all-namespaces
    NAMESPACE            NAME                                             READY   UP-TO-DATE   AVAILABLE   AGE
    ingress-nginx        deployment.apps/ingress-nginx-controller         1/1     1            1           13m
    kube-system          deployment.apps/coredns                          2/2     2            2           13m
    kubeflex-system      deployment.apps/kubeflex-controller-manager      1/1     1            1           12m
    local-path-storage   deployment.apps/local-path-provisioner           1/1     1            1           13m
    wds1-system          deployment.apps/kube-apiserver                   1/1     1            1           8m35s
    wds1-system          deployment.apps/kube-controller-manager          1/1     1            1           8m35s
    wds1-system          deployment.apps/kubestellar-controller-manager   1/1     1            1           7m13s
    wds1-system          deployment.apps/transport-controller             1/1     1            1           6m14s

    NAMESPACE         NAME                                   READY   AGE
    its1-system       statefulset.apps/vcluster              1/1     10m
    kubeflex-system   statefulset.apps/postgres-postgresql   1/1     13m
    ```
## Setup PipeCD and example scenario

1. Create a KinD cluster
    ```
    kind create cluster --name pipecd --wait 5m
    ```
1. (Optional for local environment) Update CoreDNS if you use non-publicly reachable FQDN for the cluster hostname (Or add DNS entry to /etc/hosts in host machine (not attempted))
    1. Add the mapping of your machine reachable IP address and FQDN of Workload Description Space (WDS) API server with fallthrough derective 
        ```
        { 
          IP FQDN
          fallthrough
        } 
        ```
    1. For example,
        ```
        $ kubectl --context kind-pipecd edit cm coredns -n kube-system
        ...
        apiVersion: v1
        data:
          Corefile: |
            .:53 {
                errors
                health {
                  lameduck 5s
                }
                hosts {
                  192.168.10.10 wds1.localtest.me
                  fallthrough
                }
        ```
1. Dump KubeConfig of WDS in somewhere, which will be used later for Piped to sync resources to WDS 
    ```
    kubectl --context wds1 config view --minify --raw > /tmp/kubeconfig.kflex.wds1.yaml
    ```
1. Install PipeCD on the cluster
    ```
    pipectl quickstart --version v0.47.1
    ```
1. Waiting for several minutes to finish startup, pipectl navigats you to PipeCD UI via browser
    ```
    $ pipectl quickstart --version v0.47.1

    Installing the controlplane in quickstart mode...
    Release "pipecd" does not exist. Installing it now.
    NAME: pipecd
    LAST DEPLOYED: Wed May  8 03:38:19 2024
    NAMESPACE: pipecd
    STATUS: deployed
    REVISION: 1
    TEST SUITE: None

    Intalled the controlplane successfully!

    Waiting for PipeCD control plane to be ready...
    PipeCD control plane status: Terminating 1 Running 5 Init:0/1 2
    PipeCD control plane status: Terminating 1 Running 5 Init:0/1 2

    Installing the piped for quickstart...

    Openning PipeCD control plane at http://localhost:8080/
    Please login using the following account:
    - Username: ******
    - Password: ******
    For more information refer to https://pipecd.dev/docs/quickstart/
    ```
1. Log in with the provided username and password
1. Create Piped instance
    1. Go to "setting" page
    1. Go to "Piped" tab
    1. Click "ADD"
    1. Set name and description to "kubestellar"
    1. Save
    1. You'll see Piped ID, key, and the base64 encrypted key
1. Open another terminal (or keep running PipeCD quickstart command as background job to keep port-forward PipeCD service)
1. Create piped configuration (replace `<Piped ID>`, `<The base64 encrypted key>`, `<Path to Kubeconfig file of WDS>`)
    ```
    cat << EOL > piped.yaml
    apiVersion: pipecd.dev/v1beta1
    kind: Piped
    spec:
      projectID: quickstart
      pipedID: <Piped ID>
      pipedKeyData: <The base64 encrypted key>
      # Write in a format like "host:443" because the communication is done via gRPC.
      apiAddress: localhost:8080
      repositories:
        - repoId: examples
          remote: https://github.com/yana1205/pipe-cd-examples.git
          branch: master
      syncInterval: 1m
      platformProviders:
        - name: kubestellar
          type: KUBERNETES
          config:
            masterURL: https://wds1.localtest.me:9443
            kubeConfigPath: <Path to Kubeconfig file of WDS (e.g. /tmp/kubeconfig.kflex.wds1.yaml)>
    EOL
    ```
1. Run `Piped` at local mode (or it can also run on a K8S cluster)
    ```
    piped piped --config-file=piped.yaml --insecure
    ```
1. Wait for the piped instance to be ready (On UI, it's marked as green.)
1. Create application
    1. Go to "Applications" tab
    1. Click "ADD"
    1. In step1 "Select piped and platform provider", select "kubestellar" and "kubestellar" as Piped and Platform Provider
    1. In step2 "Select application to add", select "name: canary-kubestellar, repo: examples"
    1. Keep step3 as-is
    1. SAVE
1. Open another terminal again
1. Wait for a couple of minutes to be synced
    1. Wait for PipeCD to sync resources into WDS
    ```
    watch kubectl --context wds1 get deploy,svc,bp -A
    ```
    1. Then, ensure resources labeled `workload: distributed` are distributed to each managed clusters
    ```
    kubectl --context cluster1 get deploy,pod,svc
    kubectl --context cluster2 get deploy,pod,svc
    ```
1. Try canary update
    1. Edit version tag of deployment.yaml in GitHub
    1. Merge the change into master

## Cleanup
1. Delete workloads (TODO: find a way to clean up all synced resources from PipeCD)
    ```
    kubectl --context wds1 delete deploy,svc -l workload=distributed
    kubectl --context wds1 delete bp pipecd-example-canary-bpolicy
    ```
1. Uninstall PipeCD
    1. Go to the terminal running `pipectl quickstart` and set dummy ID, Key, GitRemoteRepo to terminate
        - ID: 123e4567-e89b-12d3-a456-426614174000 
        - Key: v3f6ku5x7kujmu6f0qw0wt705ijuag4w2cxvjpf566dz0rjsoe
        - Key: abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmn
        - GitRemoteRepo: https
    1. Terminate by Ctrl+c
    1. Uninstall PipeCD
        ```
        pipectl quickstart --uninstall
        ```
1. Teardown kubestellar managed clusters and kubeflex
    ```
    kind delete cluster --name cluster1
    kind delete cluster --name cluster2
    kind delete cluster --name kubeflex
    ```