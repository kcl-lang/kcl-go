import k8s.core

schema Deployment:
    metadata: str
    podSpec: core.PodSpec
