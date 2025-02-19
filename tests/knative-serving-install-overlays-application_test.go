package tests_test

import (
  "sigs.k8s.io/kustomize/v3/k8sdeps/kunstruct"
  "sigs.k8s.io/kustomize/v3/k8sdeps/transformer"
  "sigs.k8s.io/kustomize/v3/pkg/fs"
  "sigs.k8s.io/kustomize/v3/pkg/loader"
  "sigs.k8s.io/kustomize/v3/pkg/plugins"
  "sigs.k8s.io/kustomize/v3/pkg/resmap"
  "sigs.k8s.io/kustomize/v3/pkg/resource"
  "sigs.k8s.io/kustomize/v3/pkg/target"
  "sigs.k8s.io/kustomize/v3/pkg/validators"
  "testing"
)

func writeKnativeServingInstallOverlaysApplication(th *KustTestHarness) {
  th.writeF("/manifests/knative/knative-serving-install/overlays/application/application.yaml", `
apiVersion: app.k8s.io/v1beta1
kind: Application
metadata:
  name: knative-serving-install
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: knative-serving-install
      app.kubernetes.io/instance: knative-serving-install-v0.8.0
      app.kubernetes.io/managed-by: kfctl
      app.kubernetes.io/component: knative-serving-install
      app.kubernetes.io/part-of: kubeflow
      app.kubernetes.io/version: v0.8.0
  componentKinds:
  - group: core
    kind: ConfigMap
  - group: apps
    kind: Deployment
  descriptor:
    type: knative-serving-install
    version: v1beta1
    description: ""
    maintainers: []
    owners: []
    keywords:
     - knative-serving-install
     - kubeflow
    links:
    - description: About
      url: ""
  addOwnerRef: true
`)
  th.writeK("/manifests/knative/knative-serving-install/overlays/application", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
- ../../base
resources:
- application.yaml
commonLabels:
  app.kubernetes.io/name: knative-serving-install
  app.kubernetes.io/instance: knative-serving-install-v0.8.0
  app.kubernetes.io/managed-by: kfctl
  app.kubernetes.io/component: knative-serving-install
  app.kubernetes.io/part-of: kubeflow
  app.kubernetes.io/version: v0.8.0
`)
  th.writeF("/manifests/knative/knative-serving-install/base/gateway.yaml", `
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  labels:
    networking.knative.dev/ingress-provider: istio
    serving.knative.dev/release: "v0.8.0"
  name: knative-ingress-gateway
  namespace: knative-serving
spec:
  selector:
    istio: ingressgateway
  servers:
  - hosts:
    - '*'
    port:
      name: http
      number: 80
      protocol: HTTP
  - hosts:
    - '*'
    port:
      name: https
      number: 443
      protocol: HTTPS
    tls:
      mode: PASSTHROUGH

---
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  labels:
    networking.knative.dev/ingress-provider: istio
    serving.knative.dev/release: "v0.8.0"
  name: cluster-local-gateway
  namespace: knative-serving
spec:
  selector:
    istio: cluster-local-gateway
  servers:
  - hosts:
    - '*'
    port:
      name: http
      number: 80
      protocol: HTTP

`)
  th.writeF("/manifests/knative/knative-serving-install/base/cluster-role.yaml", `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    networking.knative.dev/ingress-provider: istio
    serving.knative.dev/controller: "true"
    serving.knative.dev/release: "v0.8.0"
  name: knative-serving-istio
rules:
- apiGroups:
  - networking.istio.io
  resources:
  - virtualservices
  - gateways
  verbs:
  - get
  - list
  - create
  - update
  - delete
  - patch
  - watch

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    autoscaling.knative.dev/metric-provider: custom-metrics
    serving.knative.dev/release: "v0.8.0"
  name: custom-metrics-server-resources
rules:
- apiGroups:
  - custom.metrics.k8s.io
  resources:
  - '*'
  verbs:
  - '*'

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    rbac.authorization.k8s.io/aggregate-to-admin: "true"
    serving.knative.dev/release: "v0.8.0"
  name: knative-serving-namespaced-admin
rules:
- apiGroups:
  - serving.knative.dev
  - networking.internal.knative.dev
  - autoscaling.internal.knative.dev
  resources:
  - '*'
  verbs:
  - '*'

---
aggregationRule:
  clusterRoleSelectors:
  - matchLabels:
      serving.knative.dev/controller: "true"
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    serving.knative.dev/release: "v0.8.0"
  name: knative-serving-admin
rules: []
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    serving.knative.dev/controller: "true"
    serving.knative.dev/release: "v0.8.0"
  name: knative-serving-core
rules:
- apiGroups:
  - ""
  resources:
  - pods
  - namespaces
  - secrets
  - configmaps
  - endpoints
  - services
  - events
  - serviceaccounts
  verbs:
  - get
  - list
  - create
  - update
  - delete
  - patch
  - watch
- apiGroups:
  - ""
  resources:
  - endpoints/restricted
  verbs:
  - create
- apiGroups:
  - apps
  resources:
  - deployments
  - deployments/finalizers
  verbs:
  - get
  - list
  - create
  - update
  - delete
  - patch
  - watch
- apiGroups:
  - admissionregistration.k8s.io
  resources:
  - mutatingwebhookconfigurations
  verbs:
  - get
  - list
  - create
  - update
  - delete
  - patch
  - watch
- apiGroups:
  - apiextensions.k8s.io
  resources:
  - customresourcedefinitions
  verbs:
  - get
  - list
  - create
  - update
  - delete
  - patch
  - watch
- apiGroups:
  - autoscaling
  resources:
  - horizontalpodautoscalers
  verbs:
  - get
  - list
  - create
  - update
  - delete
  - patch
  - watch
- apiGroups:
  - serving.knative.dev
  - autoscaling.internal.knative.dev
  - networking.internal.knative.dev
  resources:
  - '*'
  - '*/status'
  - '*/finalizers'
  verbs:
  - get
  - list
  - create
  - update
  - delete
  - deletecollection
  - patch
  - watch
- apiGroups:
  - caching.internal.knative.dev
  resources:
  - images
  verbs:
  - get
  - list
  - create
  - update
  - delete
  - patch
  - watch

---
`)
  th.writeF("/manifests/knative/knative-serving-install/base/cluster-role-binding.yaml", `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  labels:
    autoscaling.knative.dev/metric-provider: custom-metrics
    serving.knative.dev/release: "v0.8.0"
  name: custom-metrics:system:auth-delegator
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:auth-delegator
subjects:
- kind: ServiceAccount
  name: controller
  namespace: knative-serving

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  labels:
    autoscaling.knative.dev/metric-provider: custom-metrics
    serving.knative.dev/release: "v0.8.0"
  name: hpa-controller-custom-metrics
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: custom-metrics-server-resources
subjects:
- kind: ServiceAccount
  name: horizontal-pod-autoscaler
  namespace: kube-system

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  labels:
    serving.knative.dev/release: "v0.8.0"
  name: knative-serving-controller-admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: knative-serving-admin
subjects:
- kind: ServiceAccount
  name: controller
  namespace: knative-serving

---
`)
  th.writeF("/manifests/knative/knative-serving-install/base/service-role.yaml", `
apiVersion: rbac.istio.io/v1alpha1
kind: ServiceRole
metadata:
  name: istio-service-role
  namespace: knative-serving
spec:
  rules:
  - methods:
    - '*'
    services:
    - '*'


`)
  th.writeF("/manifests/knative/knative-serving-install/base/service-role-binding.yaml", `
apiVersion: rbac.istio.io/v1alpha1
kind: ServiceRoleBinding
metadata:
  name: istio-service-role-binding
  namespace: knative-serving
spec:
  roleRef:
    kind: ServiceRole
    name: istio-service-role
  subjects:
  - user: '*'
`)
  th.writeF("/manifests/knative/knative-serving-install/base/role-binding.yaml", `
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  labels:
    autoscaling.knative.dev/metric-provider: custom-metrics
    serving.knative.dev/release: "v0.8.0"
  name: custom-metrics-auth-reader
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: extension-apiserver-authentication-reader
subjects:
- kind: ServiceAccount
  name: controller
  namespace: knative-serving

`)
  th.writeF("/manifests/knative/knative-serving-install/base/config-map.yaml", `
apiVersion: v1
data:
  _example: |
    ################################
    #                              #
    #    EXAMPLE CONFIGURATION     #
    #                              #
    ################################

    # This block is not actually functional configuration,
    # but serves to illustrate the available configuration
    # options and document them in a way that is accessible
    # to users that kubectl edit this config map.
    #
    # These sample configuration options may be copied out of
    # this example block and unindented to be in the data block
    # to actually change the configuration.

    # The Revision ContainerConcurrency field specifies the maximum number
    # of requests the Container can handle at once. Container concurrency
    # target percentage is how much of that maximum to use in a stable
    # state. E.g. if a Revision specifies ContainerConcurrency of 10, then
    # the Autoscaler will try to maintain 7 concurrent connections per pod
    # on average.
    # Note: this limit will be applied to container concurrency set at every
    # level (ConfigMap, Revision Spec or Annotation).
    # For legacy and backwards compatibility reasons, this value also accepts
    # fractional values in (0, 1] interval (i.e. 0.7 ⇒ 70%).
    # Thus minimal percentage value must be greater than 1.0, or it will be
    # treated as a fraction.
    container-concurrency-target-percentage: "70"

    # The container concurrency target default is what the Autoscaler will
    # try to maintain when the Revision specifies unlimited concurrency.
    # Even when specifying unlimited concurrency, the autoscaler will
    # horizontally scale the application based on this target concurrency.
    container-concurrency-target-default: "100"

    # The target burst capacity specifies the size of burst in concurrent
    # requests that the system operator expects the system will receive.
    # Autoscaler will try to protect the system from queueing by introducing
    # Activator in the request path if the current spare capacity of the
    # service is less than this setting.
    # If this setting is 0, then Activator will be in the request path only
    # when the revision is scaled to 0.
    # If this setting is > 0 and container-concurrency-target-percentage is
    # 100% or 1.0, then activator will always be in the request path.
    # -1 denotes unlimited target-burst-capacity and activator will always
    # be in the request path.
    # Other negative values are invalid.
    target-burst-capacity: "0"

    # When operating in a stable mode, the autoscaler operates on the
    # average concurrency over the stable window.
    stable-window: "60s"

    # When observed average concurrency during the panic window reaches
    # panic-threshold-percentage the target concurrency, the autoscaler
    # enters panic mode. When operating in panic mode, the autoscaler
    # scales on the average concurrency over the panic window which is
    # panic-window-percentage of the stable-window.
    panic-window-percentage: "10.0"

    # Absolute panic window duration.
    # Deprecated in favor of panic-window-percentage.
    # Existing revisions will continue to scale based on panic-window
    # but new revisions will default to panic-window-percentage.
    panic-window: "6s"

    # The percentage of the container concurrency target at which to
    # enter panic mode when reached within the panic window.
    panic-threshold-percentage: "200.0"

    # Max scale up rate limits the rate at which the autoscaler will
    # increase pod count. It is the maximum ratio of desired pods versus
    # observed pods.
    max-scale-up-rate: "1000.0"

    # Scale to zero feature flag
    enable-scale-to-zero: "true"

    # Tick interval is the time between autoscaling calculations.
    tick-interval: "2s"

    # Dynamic parameters (take effect when config map is updated):

    # Scale to zero grace period is the time an inactive revision is left
    # running before it is scaled to zero (min: 30s).
    scale-to-zero-grace-period: "30s"
kind: ConfigMap
metadata:
  labels:
    serving.knative.dev/release: "v0.8.0"
  name: config-autoscaler
  namespace: knative-serving

---

apiVersion: v1
data:
  _example: |
    ################################
    #                              #
    #    EXAMPLE CONFIGURATION     #
    #                              #
    ################################

    # This block is not actually functional configuration,
    # but serves to illustrate the available configuration
    # options and document them in a way that is accessible
    # to users that kubectl edit this config map.
    #
    # These sample configuration options may be copied out of
    # this example block and unindented to be in the data block
    # to actually change the configuration.

    # revision-timeout-seconds contains the default number of
    # seconds to use for the revision's per-request timeout, if
    # none is specified.
    revision-timeout-seconds: "300"  # 5 minutes

    # max-revision-timeout-seconds contains the maximum number of
    # seconds that can be used for revision-timeout-seconds.
    # This value must be greater than or equal to revision-timeout-seconds.
    # If omitted, the system default is used (600 seconds).
    max-revision-timeout-seconds: "600"  # 10 minutes

    # revision-cpu-request contains the cpu allocation to assign
    # to revisions by default.  If omitted, no value is specified
    # and the system default is used.
    revision-cpu-request: "400m"  # 0.4 of a CPU (aka 400 milli-CPU)

    # revision-memory-request contains the memory allocation to assign
    # to revisions by default.  If omitted, no value is specified
    # and the system default is used.
    revision-memory-request: "100M"  # 100 megabytes of memory

    # revision-cpu-limit contains the cpu allocation to limit
    # revisions to by default.  If omitted, no value is specified
    # and the system default is used.
    revision-cpu-limit: "1000m"  # 1 CPU (aka 1000 milli-CPU)

    # revision-memory-limit contains the memory allocation to limit
    # revisions to by default.  If omitted, no value is specified
    # and the system default is used.
    revision-memory-limit: "200M"  # 200 megabytes of memory

    # container-name-template contains a template for the default
    # container name, if none is specified.  This field supports
    # Go templating and is supplied with the ObjectMeta of the
    # enclosing Service or Configuration, so values such as
    # {{.Name}} are also valid.
    container-name-template: "user-container"
kind: ConfigMap
metadata:
  labels:
    serving.knative.dev/release: "v0.8.0"
  name: config-defaults
  namespace: knative-serving

---
apiVersion: v1
data:
  _example: |
    ################################
    #                              #
    #    EXAMPLE CONFIGURATION     #
    #                              #
    ################################

    # This block is not actually functional configuration,
    # but serves to illustrate the available configuration
    # options and document them in a way that is accessible
    # to users that kubectl edit this config map.
    #
    # These sample configuration options may be copied out of
    # this example block and unindented to be in the data block
    # to actually change the configuration.

    # List of repositories for which tag to digest resolving should be skipped
    registriesSkippingTagResolving: "ko.local,dev.local"
  queueSidecarImage: gcr.io/knative-releases/knative.dev/serving/cmd/queue@sha256:e0654305370cf3bbbd0f56f97789c92cf5215f752b70902eba5d5fc0e88c5aca
kind: ConfigMap
metadata:
  labels:
    serving.knative.dev/release: "v0.8.0"
  name: config-deployment
  namespace: knative-serving

---
apiVersion: v1
data:
  _example: |
    ################################
    #                              #
    #    EXAMPLE CONFIGURATION     #
    #                              #
    ################################

    # This block is not actually functional configuration,
    # but serves to illustrate the available configuration
    # options and document them in a way that is accessible
    # to users that kubectl edit this config map.
    #
    # These sample configuration options may be copied out of
    # this example block and unindented to be in the data block
    # to actually change the configuration.

    # Default value for domain.
    # Although it will match all routes, it is the least-specific rule so it
    # will only be used if no other domain matches.
    example.com: |

    # These are example settings of domain.
    # example.org will be used for routes having app=nonprofit.
    example.org: |
      selector:
        app: nonprofit

    # Routes having domain suffix of 'svc.cluster.local' will not be exposed
    # through Ingress. You can define your own label selector to assign that
    # domain suffix to your Route here, or you can set the label
    #    "serving.knative.dev/visibility=cluster-local"
    # to achieve the same effect.  This shows how to make routes having
    # the label app=secret only exposed to the local cluster.
    svc.cluster.local: |
      selector:
        app: secret
kind: ConfigMap
metadata:
  labels:
    serving.knative.dev/release: "v0.8.0"
  name: config-domain
  namespace: knative-serving

---
apiVersion: v1
data:
  _example: |
    ################################
    #                              #
    #    EXAMPLE CONFIGURATION     #
    #                              #
    ################################

    # This block is not actually functional configuration,
    # but serves to illustrate the available configuration
    # options and document them in a way that is accessible
    # to users that kubectl edit this config map.
    #
    # These sample configuration options may be copied out of
    # this example block and unindented to be in the data block
    # to actually change the configuration.

    # Delay after revision creation before considering it for GC
    stale-revision-create-delay: "24h"

    # Duration since a route has been pointed at a revision before it should be GC'd
    # This minus lastpinned-debounce be longer than the controller resync period (10 hours)
    stale-revision-timeout: "15h"

    # Minimum number of generations of revisions to keep before considering for GC
    stale-revision-minimum-generations: "1"

    # To avoid constant updates, we allow an existing annotation to be stale by this
    # amount before we update the timestamp
    stale-revision-lastpinned-debounce: "5h"
kind: ConfigMap
metadata:
  labels:
    serving.knative.dev/release: "v0.8.0"
  name: config-gc
  namespace: knative-serving

---
apiVersion: v1
data:
  _example: |
    ################################
    #                              #
    #    EXAMPLE CONFIGURATION     #
    #                              #
    ################################

    # This block is not actually functional configuration,
    # but serves to illustrate the available configuration
    # options and document them in a way that is accessible
    # to users that kubectl edit this config map.
    #
    # These sample configuration options may be copied out of
    # this example block and unindented to be in the data block
    # to actually change the configuration.

    # Default Knative Gateway after v0.3. It points to the Istio
    # standard istio-ingressgateway, instead of a custom one that we
    # used pre-0.3.
    gateway.knative-ingress-gateway: "istio-ingressgateway.istio-system.svc.cluster.local"

    # A cluster local gateway to allow pods outside of the mesh to access
    # Services and Routes not exposing through an ingress.  If the users
    # do have a service mesh setup, this isn't required and can be removed.
    #
    # An example use case is when users want to use Istio without any
    # sidecar injection (like Knative's istio-lean.yaml).  Since every pod
    # is outside of the service mesh in that case, a cluster-local  service
    # will need to be exposed to a cluster-local gateway to be accessible.
    local-gateway.cluster-local-gateway: "cluster-local-gateway.istio-system.svc.cluster.local"

    # To use only Istio service mesh and no cluster-local-gateway, replace
    # all local-gateway.* entries the following entry.
    local-gateway.mesh: "mesh"

    # Feature flag to enable reconciling external Istio Gateways.
    # When auto TLS feature is turned on, reconcileExternalGateway will be automatically enforced.
    # 1. true: enabling reconciling external gateways.
    # 2. false: disabling reconciling external gateways.
    reconcileExternalGateway: "false"
kind: ConfigMap
metadata:
  labels:
    networking.knative.dev/ingress-provider: istio
    serving.knative.dev/release: "v0.8.0"
  name: config-istio
  namespace: knative-serving

---
apiVersion: v1
data:
  _example: |
    ################################
    #                              #
    #    EXAMPLE CONFIGURATION     #
    #                              #
    ################################

    # This block is not actually functional configuration,
    # but serves to illustrate the available configuration
    # options and document them in a way that is accessible
    # to users that kubectl edit this config map.
    #
    # These sample configuration options may be copied out of
    # this example block and unindented to be in the data block
    # to actually change the configuration.

    # Common configuration for all Knative codebase
    zap-logger-config: |
      {
        "level": "info",
        "development": false,
        "outputPaths": ["stdout"],
        "errorOutputPaths": ["stderr"],
        "encoding": "json",
        "encoderConfig": {
          "timeKey": "ts",
          "levelKey": "level",
          "nameKey": "logger",
          "callerKey": "caller",
          "messageKey": "msg",
          "stacktraceKey": "stacktrace",
          "lineEnding": "",
          "levelEncoder": "",
          "timeEncoder": "iso8601",
          "durationEncoder": "",
          "callerEncoder": ""
        }
      }

    # Log level overrides
    # For all components except the autoscaler and queue proxy,
    # changes are be picked up immediately.
    # For autoscaler and queue proxy, changes require recreation of the pods.
    loglevel.controller: "info"
    loglevel.autoscaler: "info"
    loglevel.queueproxy: "info"
    loglevel.webhook: "info"
    loglevel.activator: "info"
kind: ConfigMap
metadata:
  labels:
    serving.knative.dev/release: "v0.8.0"
  name: config-logging
  namespace: knative-serving

---
apiVersion: v1
data:
  _example: |
    ################################
    #                              #
    #    EXAMPLE CONFIGURATION     #
    #                              #
    ################################

    # This block is not actually functional configuration,
    # but serves to illustrate the available configuration
    # options and document them in a way that is accessible
    # to users that kubectl edit this config map.
    #
    # These sample configuration options may be copied out of
    # this example block and unindented to be in the data block
    # to actually change the configuration.

    # istio.sidecar.includeOutboundIPRanges specifies the IP ranges that Istio sidecar
    # will intercept.
    #
    # Replace this with the IP ranges of your cluster (see below for some examples).
    # Separate multiple entries with a comma.
    # Example: "10.4.0.0/14,10.7.240.0/20"
    #
    # If set to "*" Istio will intercept all traffic within
    # the cluster as well as traffic that is going outside the cluster.
    # Traffic going outside the cluster will be blocked unless
    # necessary egress rules are created.
    #
    # If omitted or set to "", value of global.proxy.includeIPRanges
    # provided at Istio deployment time is used. In default Knative serving
    # deployment, global.proxy.includeIPRanges value is set to "*".
    #
    # If an invalid value is passed, "" is used instead.
    #
    # If valid set of IP address ranges are put into this value,
    # Istio will no longer intercept traffic going to IP addresses
    # outside the provided ranges and there is no need to specify
    # egress rules.
    #
    # To determine the IP ranges of your cluster:
    #   IBM Cloud Private: cat cluster/config.yaml | grep service_cluster_ip_range
    #   IBM Cloud Kubernetes Service: "172.30.0.0/16,172.20.0.0/16,10.10.10.0/24"
    #   Google Container Engine (GKE): gcloud container clusters describe XXXXXXX --zone=XXXXXX | grep -e clusterIpv4Cidr -e servicesIpv4Cidr
    #   Azure Kubernetes Service (AKS): "10.0.0.0/16"
    #   Azure Container Service (ACS; deprecated): "10.244.0.0/16,10.240.0.0/16"
    #   Azure Container Service Engine (ACS-Engine; OSS): Configurable, but defaults to "10.0.0.0/16"
    #   Minikube: "10.0.0.1/24"
    #
    # For more information, visit
    # https://istio.io/docs/tasks/traffic-management/egress/
    #
    istio.sidecar.includeOutboundIPRanges: "*"

    # clusteringress.class specifies the default cluster ingress class
    # to use when not dictated by Route annotation.
    #
    # If not specified, will use the Istio ingress.
    #
    # Note that changing the ClusterIngress class of an existing Route
    # will result in undefined behavior.  Therefore it is best to only
    # update this value during the setup of Knative, to avoid getting
    # undefined behavior.
    clusteringress.class: "istio.ingress.networking.knative.dev"

    # certificate.class specifies the default Certificate class
    # to use when not dictated by Route annotation.
    #
    # If not specified, will use the Cert-Manager Certificate.
    #
    # Note that changing the Certificate class of an existing Route
    # will result in undefined behavior.  Therefore it is best to only
    # update this value during the setup of Knative, to avoid getting
    # undefined behavior.
    certificate.class: "cert-manager.certificate.networking.internal.knative.dev"

    # domainTemplate specifies the golang text template string to use
    # when constructing the Knative service's DNS name. The default
    # value is "{{.Name}}.{{.Namespace}}.{{.Domain}}". And those three
    # values (Name, Namespace, Domain) are the only variables defined.
    #
    # Changing this value might be necessary when the extra levels in
    # the domain name generated is problematic for wildcard certificates
    # that only support a single level of domain name added to the
    # certificate's domain. In those cases you might consider using a value
    # of "{{.Name}}-{{.Namespace}}.{{.Domain}}", or removing the Namespace
    # entirely from the template. When choosing a new value be thoughtful
    # of the potential for conflicts - for example, when users choose to use
    # characters such as - in their service, or namespace, names.
    # {{.Annotations}} can be used for any customization in the go template if needed.
    # We strongly recommend keeping namespace part of the template to avoid domain name clashes
    # Example '{{.Name}}-{{.Namespace}}.{{ index .Annotations "sub"}}.{{.Domain}}'
    # and you have an annotation {"sub":"foo"}, then the generated template would be {Name}-{Namespace}.foo.{Domain}
    domainTemplate: "{{.Name}}.{{.Namespace}}.{{.Domain}}"

    # tagTemplate specifies the golang text template string to use
    # when constructing the DNS name for "tags" within the traffic blocks
    # of Routes and Configuration.  This is used in conjunction with the
    # domainTemplate above to determine the full URL for the tag.
    tagTemplate: "{{.Name}}-{{.Tag}}"

    # Controls whether TLS certificates are automatically provisioned and
    # installed in the Knative ingress to terminate external TLS connection.
    # 1. Enabled: enabling auto-TLS feature.
    # 2. Disabled: disabling auto-TLS feature.
    autoTLS: "Disabled"

    # Controls the behavior of the HTTP endpoint for the Knative ingress.
    # It requires autoTLS to be enabled.
    # 1. Enabled: The Knative ingress will be able to serve HTTP connection.
    # 2. Disabled: The Knative ingress ter will reject HTTP traffic.
    # 3. Redirected: The Knative ingress will send a 302 redirect for all
    # http connections, asking the clients to use HTTPS
    httpProtocol: "Enabled"
kind: ConfigMap
metadata:
  labels:
    serving.knative.dev/release: "v0.8.0"
  name: config-network
  namespace: knative-serving

---
apiVersion: v1
data:
  _example: |
    ################################
    #                              #
    #    EXAMPLE CONFIGURATION     #
    #                              #
    ################################

    # This block is not actually functional configuration,
    # but serves to illustrate the available configuration
    # options and document them in a way that is accessible
    # to users that kubectl edit this config map.
    #
    # These sample configuration options may be copied out of
    # this example block and unindented to be in the data block
    # to actually change the configuration.

    # logging.enable-var-log-collection defaults to false.
    # The fluentd daemon set will be set up to collect /var/log if
    # this flag is true.
    logging.enable-var-log-collection: false

    # logging.revision-url-template provides a template to use for producing the
    # logging URL that is injected into the status of each Revision.
    # This value is what you might use the the Knative monitoring bundle, and provides
    # access to Kibana after setting up kubectl proxy.
    logging.revision-url-template: |
      http://localhost:8001/api/v1/namespaces/knative-monitoring/services/kibana-logging/proxy/app/kibana#/discover?_a=(query:(match:(kubernetes.labels.serving-knative-dev%2FrevisionUID:(query:'${REVISION_UID}',type:phrase))))

    # If non-empty, this enables queue proxy writing request logs to stdout.
    # The value determines the shape of the request logs and it must be a valid go text/template.
    # It is important to keep this as a single line. Multiple lines are parsed as separate entities
    # by most collection agents and will split the request logs into multiple records.
    #
    # The following fields and functions are available to the template:
    #
    # Request: An http.Request (see https://golang.org/pkg/net/http/#Request)
    # representing an HTTP request received by the server.
    #
    # Response:
    # struct {
    #   Code    int       // HTTP status code (see https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml)
    #   Size    int       // An int representing the size of the response.
    #   Latency float64   // A float64 representing the latency of the response in seconds.
    # }
    #
    # Revision:
    # struct {
    #   Name          string  // Knative revision name
    #   Namespace     string  // Knative revision namespace
    #   Service       string  // Knative service name
    #   Configuration string  // Knative configuration name
    #   PodName       string  // Name of the pod hosting the revision
    #   PodIP         string  // IP of the pod hosting the revision
    # }
    #
    logging.request-log-template: '{"httpRequest": {"requestMethod": "{{.Request.Method}}", "requestUrl": "{{js .Request.RequestURI}}", "requestSize": "{{.Request.ContentLength}}", "status": {{.Response.Code}}, "responseSize": "{{.Response.Size}}", "userAgent": "{{js .Request.UserAgent}}", "remoteIp": "{{js .Request.RemoteAddr}}", "serverIp": "{{.Revision.PodIP}}", "referer": "{{js .Request.Referer}}", "latency": "{{.Response.Latency}}s", "protocol": "{{.Request.Proto}}"}, "traceId": "{{index .Request.Header "X-B3-Traceid"}}"}'

    # metrics.backend-destination field specifies the system metrics destination.
    # It supports either prometheus (the default) or stackdriver.
    # Note: Using stackdriver will incur additional charges
    metrics.backend-destination: prometheus

    # metrics.request-metrics-backend-destination specifies the request metrics
    # destination. If non-empty, it enables queue proxy to send request metrics.
    # Currently supported values: prometheus, stackdriver.
    metrics.request-metrics-backend-destination: prometheus

    # metrics.stackdriver-project-id field specifies the stackdriver project ID. This
    # field is optional. When running on GCE, application default credentials will be
    # used if this field is not provided.
    metrics.stackdriver-project-id: "<your stackdriver project id>"

    # metrics.allow-stackdriver-custom-metrics indicates whether it is allowed to send metrics to
    # Stackdriver using "global" resource type and custom metric type if the
    # metrics are not supported by "knative_revision" resource type. Setting this
    # flag to "true" could cause extra Stackdriver charge.
    # If metrics.backend-destination is not Stackdriver, this is ignored.
    metrics.allow-stackdriver-custom-metrics: "false"
kind: ConfigMap
metadata:
  labels:
    serving.knative.dev/release: "v0.8.0"
  name: config-observability
  namespace: knative-serving

---
apiVersion: v1
data:
  _example: |
    ################################
    #                              #
    #    EXAMPLE CONFIGURATION     #
    #                              #
    ################################

    # This block is not actually functional configuration,
    # but serves to illustrate the available configuration
    # options and document them in a way that is accessible
    # to users that kubectl edit this config map.
    #
    # These sample configuration options may be copied out of
    # this example block and unindented to be in the data block
    # to actually change the configuration.
    #
    # If true we enable adding spans within our applications.
    enable: "false"

    # URL to zipkin collector where traces are sent.
    zipkin-endpoint: "http://zipkin.istio-system.svc.cluster.local:9411/api/v2/spans"

    # Enable zipkin debug mode. This allows all spans to be sent to the server
    # bypassing sampling.
    debug: "false"

    # Percentage (0-1) of requests to trace
    sample-rate: "0.1"
kind: ConfigMap
metadata:
  labels:
    serving.knative.dev/release: "v0.8.0"
  name: config-tracing
  namespace: knative-serving

---
`)
  th.writeF("/manifests/knative/knative-serving-install/base/deployment.yaml", `
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    serving.knative.dev/release: "v0.8.0"
  name: activator
  namespace: knative-serving
spec:
  selector:
    matchLabels:
      app: activator
      role: activator
  template:
    metadata:
      annotations:
        cluster-autoscaler.kubernetes.io/safe-to-evict: "false"
        sidecar.istio.io/inject: "true"
      labels:
        app: activator
        role: activator
        serving.knative.dev/release: "v0.8.0"
    spec:
      containers:
      - args:
        - -logtostderr=false
        - -stderrthreshold=FATAL
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: SYSTEM_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: CONFIG_LOGGING_NAME
          value: config-logging
        - name: CONFIG_OBSERVABILITY_NAME
          value: config-observability
        - name: METRICS_DOMAIN
          value: knative.dev/serving
        image: gcr.io/knative-releases/knative.dev/serving/cmd/activator@sha256:88d864eb3c47881cf7ac058479d1c735cc3cf4f07a11aad0621cd36dcd9ae3c6
        livenessProbe:
          httpGet:
            httpHeaders:
            - name: k-kubelet-probe
              value: activator
            path: /healthz
            port: 8012
        name: activator
        ports:
        - containerPort: 8012
          name: http1-port
        - containerPort: 8013
          name: h2c-port
        - containerPort: 9090
          name: metrics-port
        readinessProbe:
          httpGet:
            httpHeaders:
            - name: k-kubelet-probe
              value: activator
            path: /healthz
            port: 8012
        resources:
          limits:
            cpu: 1000m
            memory: 600Mi
          requests:
            cpu: 300m
            memory: 60Mi
        securityContext:
          allowPrivilegeEscalation: false
      serviceAccountName: controller
      terminationGracePeriodSeconds: 300
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    autoscaling.knative.dev/autoscaler-provider: hpa
    serving.knative.dev/release: "v0.8.0"
  name: autoscaler-hpa
  namespace: knative-serving
spec:
  replicas: 1
  selector:
    matchLabels:
      app: autoscaler-hpa
  template:
    metadata:
      annotations:
        sidecar.istio.io/inject: "false"
      labels:
        app: autoscaler-hpa
        serving.knative.dev/release: "v0.8.0"
    spec:
      containers:
      - env:
        - name: SYSTEM_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: CONFIG_LOGGING_NAME
          value: config-logging
        - name: CONFIG_OBSERVABILITY_NAME
          value: config-observability
        - name: METRICS_DOMAIN
          value: knative.dev/serving
        image: gcr.io/knative-releases/knative.dev/serving/cmd/autoscaler-hpa@sha256:a7801c3cf4edecfa51b7bd2068f97941f6714f7922cb4806245377c2b336b723
        name: autoscaler-hpa
        ports:
        - containerPort: 9090
          name: metrics
        resources:
          limits:
            cpu: 1000m
            memory: 1000Mi
          requests:
            cpu: 100m
            memory: 100Mi
        securityContext:
          allowPrivilegeEscalation: false
      serviceAccountName: controller

---

apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    serving.knative.dev/release: "v0.8.0"
  name: autoscaler
  namespace: knative-serving
spec:
  replicas: 1
  selector:
    matchLabels:
      app: autoscaler
  template:
    metadata:
      annotations:
        cluster-autoscaler.kubernetes.io/safe-to-evict: "false"
        sidecar.istio.io/inject: "true"
        traffic.sidecar.istio.io/includeInboundPorts: 8080,9090
      labels:
        app: autoscaler
        serving.knative.dev/release: "v0.8.0"
    spec:
      containers:
      - args:
        - --secure-port=8443
        - --cert-dir=/tmp
        env:
        - name: SYSTEM_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: CONFIG_LOGGING_NAME
          value: config-logging
        - name: CONFIG_OBSERVABILITY_NAME
          value: config-observability
        - name: METRICS_DOMAIN
          value: knative.dev/serving
        image: gcr.io/knative-releases/knative.dev/serving/cmd/autoscaler@sha256:aeaacec4feedee309293ac21da13e71a05a2ad84b1d5fcc01ffecfa6cfbb2870
        livenessProbe:
          httpGet:
            httpHeaders:
            - name: k-kubelet-probe
              value: autoscaler
            path: /healthz
            port: 8080
        name: autoscaler
        ports:
        - containerPort: 8080
          name: websocket
        - containerPort: 9090
          name: metrics
        - containerPort: 8443
          name: custom-metrics
        readinessProbe:
          httpGet:
            httpHeaders:
            - name: k-kubelet-probe
              value: autoscaler
            path: /healthz
            port: 8080
        resources:
          limits:
            cpu: 300m
            memory: 400Mi
          requests:
            cpu: 30m
            memory: 40Mi
        securityContext:
          allowPrivilegeEscalation: false
      serviceAccountName: controller

---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    networking.knative.dev/ingress-provider: istio
    serving.knative.dev/release: "v0.8.0"
  name: networking-istio
  namespace: knative-serving
spec:
  replicas: 1
  selector:
    matchLabels:
      app: networking-istio
  template:
    metadata:
      annotations:
        sidecar.istio.io/inject: "false"
      labels:
        app: networking-istio
    spec:
      containers:
      - env:
        - name: SYSTEM_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: CONFIG_LOGGING_NAME
          value: config-logging
        - name: CONFIG_OBSERVABILITY_NAME
          value: config-observability
        - name: METRICS_DOMAIN
          value: knative.dev/serving
        image: gcr.io/knative-releases/knative.dev/serving/cmd/networking/istio@sha256:057c999bccfe32e9889616b571dc8d389c742ff66f0b5516bad651f05459b7bc
        name: networking-istio
        ports:
        - containerPort: 9090
          name: metrics
        resources:
          limits:
            cpu: 1000m
            memory: 1000Mi
          requests:
            cpu: 100m
            memory: 100Mi
        securityContext:
          allowPrivilegeEscalation: false
      serviceAccountName: controller

---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    serving.knative.dev/release: "v0.8.0"
  name: webhook
  namespace: knative-serving
spec:
  replicas: 1
  selector:
    matchLabels:
      app: webhook
      role: webhook
  template:
    metadata:
      annotations:
        cluster-autoscaler.kubernetes.io/safe-to-evict: "false"
        sidecar.istio.io/inject: "false"
      labels:
        app: webhook
        role: webhook
        serving.knative.dev/release: "v0.8.0"
    spec:
      containers:
      - env:
        - name: SYSTEM_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: CONFIG_LOGGING_NAME
          value: config-logging
        - name: CONFIG_OBSERVABILITY_NAME
          value: config-observability
        - name: METRICS_DOMAIN
          value: knative.dev/serving
        image: gcr.io/knative-releases/knative.dev/serving/cmd/webhook@sha256:c2076674618933df53e90cf9ddd17f5ddbad513b8c95e955e45e37be7ca9e0e8
        name: webhook
        ports:
        - containerPort: 9090
          name: metrics-port
        resources:
          limits:
            cpu: 200m
            memory: 200Mi
          requests:
            cpu: 20m
            memory: 20Mi
        securityContext:
          allowPrivilegeEscalation: false
      serviceAccountName: controller

---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    serving.knative.dev/release: "v0.8.0"
  name: controller
  namespace: knative-serving
spec:
  replicas: 1
  selector:
    matchLabels:
      app: controller
  template:
    metadata:
      annotations:
        sidecar.istio.io/inject: "false"
      labels:
        app: controller
        serving.knative.dev/release: "v0.8.0"
    spec:
      containers:
      - env:
        - name: SYSTEM_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: CONFIG_LOGGING_NAME
          value: config-logging
        - name: CONFIG_OBSERVABILITY_NAME
          value: config-observability
        - name: METRICS_DOMAIN
          value: knative.dev/serving
        image: gcr.io/knative-releases/knative.dev/serving/cmd/controller@sha256:3b096e55fa907cff53d37dadc5d20c29cea9bb18ed9e921a588fee17beb937df
        name: controller
        ports:
        - containerPort: 9090
          name: metrics
        resources:
          limits:
            cpu: 1000m
            memory: 1000Mi
          requests:
            cpu: 100m
            memory: 100Mi
        securityContext:
          allowPrivilegeEscalation: false
      serviceAccountName: controller

---

`)
  th.writeF("/manifests/knative/knative-serving-install/base/service-account.yaml", `
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    serving.knative.dev/release: "v0.8.0"
  name: controller
  namespace: knative-serving

`)
  th.writeF("/manifests/knative/knative-serving-install/base/service.yaml", `
apiVersion: v1
kind: Service
metadata:
  labels:
    app: activator
    serving.knative.dev/release: "v0.8.0"
  name: activator-service
  namespace: knative-serving
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8012
  - name: http2
    port: 81
    protocol: TCP
    targetPort: 8013
  - name: metrics
    port: 9090
    protocol: TCP
    targetPort: 9090
  selector:
    app: activator
  type: ClusterIP

---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: controller
    serving.knative.dev/release: "v0.8.0"
  name: controller
  namespace: knative-serving
spec:
  ports:
  - name: metrics
    port: 9090
    protocol: TCP
    targetPort: 9090
  selector:
    app: controller

---
apiVersion: v1
kind: Service
metadata:
  labels:
    role: webhook
    serving.knative.dev/release: "v0.8.0"
  name: webhook
  namespace: knative-serving
spec:
  ports:
  - port: 443
    targetPort: 8443
  selector:
    role: webhook

---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: autoscaler
    serving.knative.dev/release: "v0.8.0"
  name: autoscaler
  namespace: knative-serving
spec:
  ports:
  - name: http
    port: 8080
    protocol: TCP
    targetPort: 8080
  - name: metrics
    port: 9090
    protocol: TCP
    targetPort: 9090
  - name: custom-metrics
    port: 443
    protocol: TCP
    targetPort: 8443
  selector:
    app: autoscaler

---
`)
  th.writeF("/manifests/knative/knative-serving-install/base/apiservice.yaml", `
apiVersion: apiregistration.k8s.io/v1beta1
kind: APIService
metadata:
  labels:
    autoscaling.knative.dev/metric-provider: custom-metrics
    serving.knative.dev/release: "v0.8.0"
  name: v1beta1.custom.metrics.k8s.io
spec:
  group: custom.metrics.k8s.io
  groupPriorityMinimum: 100
  insecureSkipTLSVerify: true
  service:
    name: autoscaler
    namespace: knative-serving
  version: v1beta1
  versionPriority: 100

`)
  th.writeF("/manifests/knative/knative-serving-install/base/image.yaml", `
apiVersion: caching.internal.knative.dev/v1alpha1
kind: Image
metadata:
  labels:
    serving.knative.dev/release: "v0.8.0"
  name: queue-proxy
  namespace: knative-serving
spec:
  image: gcr.io/knative-releases/knative.dev/serving/cmd/queue@sha256:e0654305370cf3bbbd0f56f97789c92cf5215f752b70902eba5d5fc0e88c5aca

`)
  th.writeF("/manifests/knative/knative-serving-install/base/hpa.yaml", `
apiVersion: autoscaling/v2beta1
kind: HorizontalPodAutoscaler
metadata:
  name: activator
  namespace: knative-serving
spec:
  maxReplicas: 20
  metrics:
  - resource:
      name: cpu
      targetAverageUtilization: 100
    type: Resource
  minReplicas: 1
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: activator

`)
  th.writeK("/manifests/knative/knative-serving-install/base", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: knative-serving
resources:
- gateway.yaml
- cluster-role.yaml
- cluster-role-binding.yaml
- service-role.yaml
- service-role-binding.yaml
- role-binding.yaml
- config-map.yaml
- deployment.yaml
- service-account.yaml
- service.yaml
- apiservice.yaml
- image.yaml
- hpa.yaml
commonLabels:
  kustomize.component: knative
images:
- name: gcr.io/knative-releases/knative.dev/serving/cmd/activator
  newName: gcr.io/knative-releases/knative.dev/serving/cmd/activator
  digest: sha256:88d864eb3c47881cf7ac058479d1c735cc3cf4f07a11aad0621cd36dcd9ae3c6
- name: gcr.io/knative-releases/knative.dev/serving/cmd/autoscaler-hpa
  newName: gcr.io/knative-releases/knative.dev/serving/cmd/autoscaler-hpa
  digest: sha256:a7801c3cf4edecfa51b7bd2068f97941f6714f7922cb4806245377c2b336b723
- name: gcr.io/knative-releases/knative.dev/serving/cmd/autoscaler
  newName: gcr.io/knative-releases/knative.dev/serving/cmd/autoscaler
  digest: sha256:aeaacec4feedee309293ac21da13e71a05a2ad84b1d5fcc01ffecfa6cfbb2870
- name: gcr.io/knative-releases/knative.dev/serving/cmd/networking/istio
  newName: gcr.io/knative-releases/knative.dev/serving/cmd/networking/istio
  digest: sha256:057c999bccfe32e9889616b571dc8d389c742ff66f0b5516bad651f05459b7bc
- name: gcr.io/knative-releases/knative.dev/serving/cmd/webhook
  newName: gcr.io/knative-releases/knative.dev/serving/cmd/webhook
  digest: sha256:c2076674618933df53e90cf9ddd17f5ddbad513b8c95e955e45e37be7ca9e0e8
- name: gcr.io/knative-releases/knative.dev/serving/cmd/controller
  newName: gcr.io/knative-releases/knative.dev/serving/cmd/controller
  digest: sha256:3b096e55fa907cff53d37dadc5d20c29cea9bb18ed9e921a588fee17beb937df
`)
}

func TestKnativeServingInstallOverlaysApplication(t *testing.T) {
  th := NewKustTestHarness(t, "/manifests/knative/knative-serving-install/overlays/application")
  writeKnativeServingInstallOverlaysApplication(th)
  m, err := th.makeKustTarget().MakeCustomizedResMap()
  if err != nil {
    t.Fatalf("Err: %v", err)
  }
  expected, err := m.AsYaml()
  if err != nil {
    t.Fatalf("Err: %v", err)
  }
  targetPath := "../knative/knative-serving-install/overlays/application"
  fsys := fs.MakeRealFS()
  lrc := loader.RestrictionRootOnly
  _loader, loaderErr := loader.NewLoader(lrc, validators.MakeFakeValidator(), targetPath, fsys)
  if loaderErr != nil {
    t.Fatalf("could not load kustomize loader: %v", loaderErr)
  }
  rf := resmap.NewFactory(resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl()), transformer.NewFactoryImpl())
  pc := plugins.DefaultPluginConfig()
  kt, err := target.NewKustTarget(_loader, rf, transformer.NewFactoryImpl(), plugins.NewLoader(pc, rf))
  if err != nil {
    th.t.Fatalf("Unexpected construction error %v", err)
  }
  actual, err := kt.MakeCustomizedResMap()
  if err != nil {
    t.Fatalf("Err: %v", err)
  }
  th.assertActualEqualsExpected(actual, string(expected))
}
