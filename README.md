### Building
The operator image is already uploaded to docker hub at `njhughes/go-angi:latest`.
In case you want to upload the operator to a different registry, build the docker container from the root directory:
```docker build . -t <image name>``` 

### Deploying
NOTE: Ensure port `9898` is available, or update the deploy script to use another port.

To deploy the CRD (ant its RBAC), operator and test CRD object, run this command from the root directory:
```./scripts/deploy.sh```

### Connecting
- To connecto to a podinfo pod (via its service), simply navigate to `http://localhost:8989` in a web browser
- To write a new key to the podinfo database, run ```curl -X POST -H "Content-Type: application/json" -d '{"value":"world"}' http://localhost:9898/cache/hello```
- To retrieve that same key, run ```curl http://localhost:9898/cache/hello```

### Testing:
For a production operator, I would take a test-driven approach to development, writing unit and integration tests incrementally as I go.
However, given the time constraints and objective of this project (i.e. demonstrating my knowledge of the Kubernetes reconcilliation controller pattern), I did not include tests in my submission.

### Project structure:
```
angi/
├── main.go
├── pkg/
│   ├── controller/
│   │   ├── controller.go
├── api/
│   └── myapigroup/
│       └── v1alpha1/
│           └── types.go
├── config
│   ├── crd
│   │   └── crd.yaml
│   ├── rbac
│   │   ├── clusterrole.yaml
│   ├── deployment
│   │   ├── myapp-operator-deployment.yaml
│   └── samples
│       └── test-crd.yaml
└── scripts/
    └── deploy.sh
```

### Future improvements:

- Testing: In a production system, I would ensure that there are unit tests as well as canary or integration tests that deploy to a Kuberentes cluster and validate a suite of functionality against a live deployment.

- Pod Configuration: In a production-grade operator, we'd likely want to have more configuration for the pods (such as liveness/readiness checks, volume mounts, resource limits, etc.), but for a basic example, this is sufficient.

- Concurrency Issues: If there are multiple changes in rapid succession, it's possible for the operator to get into a weird state. For example, if a user creates, deletes, and re-creates a MyAppResource object very quickly, the operator might start creating resources for the second creation before it's finished cleaning up from the deletion. This is a complex issue that's out of scope for a simple operator, but worth thinking about for a production use case.

- ConfigMap for Shared Configuration: If the application pods (PodInfo) and the Redis pod share some common configuration, it might be useful to store those in a ConfigMap and reference them in your pods.





