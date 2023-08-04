The operator image is already uploaded to docker hub at njhughes/go-angi:latest.
In case you want to upload the operator to a different registry, build the docker container from the root directory:
```docker build . -t <image name>``` 

NOTE: Ensure port 9898 is available, or update the deploy script to use another port.

To deploy the CRD (ant its RBAC), operator and test CRD object, run this command from the root directory:
```./scripts/deploy.sh```

Testing:
For a production operator, I would take a test-driven approach to development, writing unit and integration tests incrementally as I go.
However, given the time constraints and objective of this project (i.e. demonstrating my knowledge of the Kubernetes reconcilliation controller pattern),
I did not include tests in my submission.

Project structure:
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

Notes on future improvements:

Error Handling: Right now, if an error occurs during the creation of resources (pods, services etc.), the program prints an error message and continues executing. It might be more robust to exit the function or return an error when a serious error occurs (like failing to create a pod).

Pod Configuration: In a production-grade operator, you'd likely want to have more configuration for your pods (such as liveness/readiness checks, volume mounts, resource limits, etc.), but for a basic example, what you have is sufficient.

Update Mechanism: Currently, your operator handles resource addition, but not update or deletion. This might be out of the scope for your current work, but keep in mind that a fully functional operator would also handle updates and deletions to the Custom Resource.

Concurrency Issues: If there are multiple changes in rapid succession, it's possible for the operator to get into a weird state. For example, if a user creates, deletes, and re-creates a MyAppResource object very quickly, the operator might start creating resources for the second creation before it's finished cleaning up from the deletion. This is a complex issue that's out of scope for a simple operator, but worth thinking about for a production use case.

ConfigMap for Shared Configuration: If your application pods (PodInfo) and the Redis pod share some common configuration, it might be useful to store those in a ConfigMap and reference them in your pods.





