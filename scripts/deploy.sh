kubectl apply -f config/rbac/clusterrole.yaml
kubectl apply -f config/crd/crd.yaml
kubectl apply -f config/deployment/myapp-operator-deployment.yaml
kubectl apply -f config/samples/test-crd.yaml

echo "Waiting 30 seconds for service to start before configuring port forwarding"
# Start a command in the background
sleep 30 &
# Wait for the background process to finish
wait

kubectl port-forward svc/my-app-resource-object-service 9898:9898 &
echo "Port forwarding set. Service is available at local port 9898"
