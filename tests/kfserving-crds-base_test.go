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

func writeKfservingCrdsBase(th *KustTestHarness) {
  th.writeF("/manifests/kfserving/kfserving-crds/base/crd.yaml", `
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: inferenceservices.serving.kubeflow.org
spec:
  additionalPrinterColumns:
  - JSONPath: .status.url
    name: URL
    type: string
  - JSONPath: .status.conditions[?(@.type=='Ready')].status
    name: Ready
    type: string
  - JSONPath: .status.traffic
    name: Default Traffic
    type: integer
  - JSONPath: .status.canaryTraffic
    name: Canary Traffic
    type: integer
  - JSONPath: .metadata.creationTimestamp
    name: Age
    type: date
  group: serving.kubeflow.org
  names:
    kind: InferenceService
    plural: inferenceservices
    shortNames:
    - inferenceservice
  scope: Namespaced
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            canary:
              description: Canary defines an alternate endpoints to route a percentage
                of traffic.
              properties:
                explainer:
                  description: Explainer defines the model explanation service spec,
                    explainer service calls to predictor or transformer if it is specified.
                  properties:
                    alibi:
                      description: Spec for alibi explainer
                      properties:
                        config:
                          description: Inline custom parameter settings for explainer
                          type: object
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest Alibi Version
                          type: string
                        storageUri:
                          description: The location of a trained explanation model
                          type: string
                        type:
                          description: The type of Alibi explainer
                          type: string
                      required:
                      - type
                      type: object
                    custom:
                      description: Spec for a custom explainer
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                  type: object
                predictor:
                  description: Predictor defines the model serving spec +required
                  properties:
                    custom:
                      description: Spec for a custom predictor
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    onnx:
                      description: Spec for ONNX runtime (https://github.com/microsoft/onnxruntime)
                      properties:
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Allowed runtime versions are [v0.5.0, latest]
                            and defaults to the version specified in inferenceservice
                            config map
                          type: string
                        storageUri:
                          description: The location of the trained model
                          type: string
                      required:
                      - storageUri
                      type: object
                    pytorch:
                      description: Spec for PyTorch predictor
                      properties:
                        modelClassName:
                          description: Defaults PyTorch model class name to 'PyTorchModel'
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Allowed runtime versions are [0.2.0, latest]
                            and defaults to the version specified in inferenceservice
                            config map
                          type: string
                        storageUri:
                          description: The location of the trained model
                          type: string
                      required:
                      - storageUri
                      type: object
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                    sklearn:
                      description: Spec for SKLearn predictor
                      properties:
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Allowed runtime versions are [0.2.0, latest]
                            and defaults to the version specified in inferenceservice
                            config map
                          type: string
                        storageUri:
                          description: The location of the trained model
                          type: string
                      required:
                      - storageUri
                      type: object
                    tensorflow:
                      description: Spec for Tensorflow Serving (https://github.com/tensorflow/serving)
                      properties:
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Allowed runtime versions are [1.11.0, 1.12.0,
                            1.13.0, 1.14.0, latest] or [1.11.0-gpu, 1.12.0-gpu, 1.13.0-gpu,
                            1.14.0-gpu, latest-gpu] if gpu resource is specified and
                            defaults to the version specified in inferenceservice
                            config map.
                          type: string
                        storageUri:
                          description: The location of the trained model
                          type: string
                      required:
                      - storageUri
                      type: object
                    tensorrt:
                      description: Spec for TensorRT Inference Server (https://github.com/NVIDIA/tensorrt-inference-server)
                      properties:
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Allowed runtime versions are [19.05-py3] and
                            defaults to the version specified in inferenceservice
                            config map
                          type: string
                        storageUri:
                          description: The location of the trained model
                          type: string
                      required:
                      - storageUri
                      type: object
                    xgboost:
                      description: Spec for XGBoost predictor
                      properties:
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Allowed runtime versions are [0.2.0, latest]
                            and defaults to the version specified in inferenceservice
                            config map
                          type: string
                        storageUri:
                          description: The location of the trained model
                          type: string
                      required:
                      - storageUri
                      type: object
                  type: object
                transformer:
                  description: Transformer defines the pre/post processing before
                    and after the predictor call, transformer service calls to predictor
                    service.
                  properties:
                    custom:
                      description: Spec for a custom transformer
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                  type: object
              required:
              - predictor
              type: object
            canaryTrafficPercent:
              description: CanaryTrafficPercent defines the percentage of traffic
                going to canary InferenceService endpoints
              format: int64
              type: integer
            default:
              description: Default defines default InferenceService endpoints +required
              properties:
                explainer:
                  description: Explainer defines the model explanation service spec,
                    explainer service calls to predictor or transformer if it is specified.
                  properties:
                    alibi:
                      description: Spec for alibi explainer
                      properties:
                        config:
                          description: Inline custom parameter settings for explainer
                          type: object
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest Alibi Version
                          type: string
                        storageUri:
                          description: The location of a trained explanation model
                          type: string
                        type:
                          description: The type of Alibi explainer
                          type: string
                      required:
                      - type
                      type: object
                    custom:
                      description: Spec for a custom explainer
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                  type: object
                predictor:
                  description: Predictor defines the model serving spec +required
                  properties:
                    custom:
                      description: Spec for a custom predictor
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    onnx:
                      description: Spec for ONNX runtime (https://github.com/microsoft/onnxruntime)
                      properties:
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Allowed runtime versions are [v0.5.0, latest]
                            and defaults to the version specified in inferenceservice
                            config map
                          type: string
                        storageUri:
                          description: The location of the trained model
                          type: string
                      required:
                      - storageUri
                      type: object
                    pytorch:
                      description: Spec for PyTorch predictor
                      properties:
                        modelClassName:
                          description: Defaults PyTorch model class name to 'PyTorchModel'
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Allowed runtime versions are [0.2.0, latest]
                            and defaults to the version specified in inferenceservice
                            config map
                          type: string
                        storageUri:
                          description: The location of the trained model
                          type: string
                      required:
                      - storageUri
                      type: object
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                    sklearn:
                      description: Spec for SKLearn predictor
                      properties:
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Allowed runtime versions are [0.2.0, latest]
                            and defaults to the version specified in inferenceservice
                            config map
                          type: string
                        storageUri:
                          description: The location of the trained model
                          type: string
                      required:
                      - storageUri
                      type: object
                    tensorflow:
                      description: Spec for Tensorflow Serving (https://github.com/tensorflow/serving)
                      properties:
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Allowed runtime versions are [1.11.0, 1.12.0,
                            1.13.0, 1.14.0, latest] or [1.11.0-gpu, 1.12.0-gpu, 1.13.0-gpu,
                            1.14.0-gpu, latest-gpu] if gpu resource is specified and
                            defaults to the version specified in inferenceservice
                            config map.
                          type: string
                        storageUri:
                          description: The location of the trained model
                          type: string
                      required:
                      - storageUri
                      type: object
                    tensorrt:
                      description: Spec for TensorRT Inference Server (https://github.com/NVIDIA/tensorrt-inference-server)
                      properties:
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Allowed runtime versions are [19.05-py3] and
                            defaults to the version specified in inferenceservice
                            config map
                          type: string
                        storageUri:
                          description: The location of the trained model
                          type: string
                      required:
                      - storageUri
                      type: object
                    xgboost:
                      description: Spec for XGBoost predictor
                      properties:
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Allowed runtime versions are [0.2.0, latest]
                            and defaults to the version specified in inferenceservice
                            config map
                          type: string
                        storageUri:
                          description: The location of the trained model
                          type: string
                      required:
                      - storageUri
                      type: object
                  type: object
                transformer:
                  description: Transformer defines the pre/post processing before
                    and after the predictor call, transformer service calls to predictor
                    service.
                  properties:
                    custom:
                      description: Spec for a custom transformer
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                  type: object
              required:
              - predictor
              type: object
          required:
          - default
          type: object
        status:
          properties:
            canary:
              description: Statuses for the canary endpoints of the InferenceService
              type: object
            canaryTraffic:
              description: Traffic percentage that goes to canary services
              format: int64
              type: integer
            conditions:
              description: Conditions the latest available observations of a resource's
                current state. +patchMergeKey=type +patchStrategy=merge
              items:
                properties:
                  lastTransitionTime:
                    description: LastTransitionTime is the last time the condition
                      transitioned from one status to another. We use VolatileTime
                      in place of metav1.Time to exclude this from creating equality.Semantic
                      differences (all other things held constant).
                    type: string
                  message:
                    description: A human readable message indicating details about
                      the transition.
                    type: string
                  reason:
                    description: The reason for the condition's last transition.
                    type: string
                  severity:
                    description: Severity with which to treat failures of this type
                      of condition. When this is not specified, it defaults to Error.
                    type: string
                  status:
                    description: Status of the condition, one of True, False, Unknown.
                      +required
                    type: string
                  type:
                    description: Type of condition. +required
                    type: string
                required:
                - type
                - status
                type: object
              type: array
            default:
              description: Statuses for the default endpoints of the InferenceService
              type: object
            observedGeneration:
              description: ObservedGeneration is the 'Generation' of the Service that
                was last processed by the controller.
              format: int64
              type: integer
            traffic:
              description: Traffic percentage that goes to default services
              format: int64
              type: integer
            url:
              description: URL of the InferenceService
              type: string
          type: object
  version: v1alpha2
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)
  th.writeK("/manifests/kfserving/kfserving-crds/base", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- crd.yaml
`)
}

func TestKfservingCrdsBase(t *testing.T) {
  th := NewKustTestHarness(t, "/manifests/kfserving/kfserving-crds/base")
  writeKfservingCrdsBase(th)
  m, err := th.makeKustTarget().MakeCustomizedResMap()
  if err != nil {
    t.Fatalf("Err: %v", err)
  }
  expected, err := m.AsYaml()
  if err != nil {
    t.Fatalf("Err: %v", err)
  }
  targetPath := "../kfserving/kfserving-crds/base"
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
