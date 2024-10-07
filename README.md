# observability-k8s
Simple observability setup for Kubernetes using Prometheus, Grafana for go api server

## Prerequisites

- Install [Minikube](https://minikube.sigs.k8s.io/docs/start/)
- Install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- Docker installed and running locally
- Install [Go](https://golang.org/doc/install)

## I. Metrics dashboards

### 1. Start Minikube

Start Minikube with the following command:

```bash
minikube start
```

### 2. Build the Docker Image

Use Minikube's Docker environment to build the image:

```bash
eval $(minikube docker-env)
docker build -t go-app:metrics .
```

### 3. Deploy the Application

Apply the deployment and service configuration files:

```bash
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```
To test the application, run the following command:

```bash
minikube service go-app-service
```
And check the ping endpoint

### 4. Deploy Prometheus and Grafana

Apply the Prometheus and Grafana configuration files:

```bash
kubectl apply -f ./prometheus/deployment.yaml
kubectl apply -f ./prometheus/service.yaml
kubectl apply -f ./prometheus/configmap.yaml
kubectl apply -f ./grafana/deployment.yaml
kubectl apply -f ./grafana/service.yaml
```

### 5. Access Grafana

To access Grafana, run the following command:

```bash
minikube service grafana-service
```

The default username and password are `admin`.

You should see the Grafana home page:

![image](img/grafana-home.jpg)

### 6. Add Prometheus as a Data Source

1. Go to `Configuration` > `Data Sources` > `Add data source`.
2. Select `Prometheus`.
3. Set the URL to `http://prometheus-service.default.svc.cluster.local:9090`.
![image](img/prometheus-input.jpg)
4. Click `Save & Test`.
![image](img/prometheus-save.jpg)

### 7. Build the Dashboard

1. Click on the `+` icon on the left sidebar.
2. Click on `New Dashboard`.
3. Click on `Add visualisation`.
4. Build the dashboard using the Prometheus data source.
5. You can use the following query to get the metrics:

```bash
up{job="go-app"}
```
It will show the uptime of the go app server.
This metric is provided by Prometheus itself and it more useful for the monitoring of the server itself.

```bash
increase(go_app_http_requests_total[5m])
```
This query shows the number of requests over the last 5 minutes. You can use labels (such as url and success) to filter
the results, allowing you to see how specific URLs or success statuses impact the total request count.

## II. Add Error Logs to the Dashboard

### 1. Deploy Loki and Promtail to k8s
```bash
kubectl apply -f ./loki/configmap.yaml
kubectl apply -f ./loki/deployment.yaml
kubectl apply -f ./loki/service.yaml
kubectl apply -f ./loki/persistentvolume.yaml
kubectl apply -f ./promtail/configmap.yaml
kubectl apply -f ./promtail/daemonset.yaml

kubectl apply -f ./promtail/service.yaml # optional
```

### 2. Add Loki datasource to Grafana
- Open Grafana in your browser
- Go to Configuration -> Data Sources
- Click on "Add data source"
- Choose "Loki" from the list
- Set the URL to `http://loki:3100`
- Click "Save & Test"

### 3. Add error logs to the dashboard
- Open the dashboard you want to add logs to
- Click on "Add panel"
- Choose "Logs" from the list
  ![image](img/add-logs.jpeg)
- Set the query to {job="go-app"} | json | detected_level="error"
- Modify options as needed. Example:
  ![image](img/logs-options.png)
- Click "Apply"
  ![image](img/dashboard.png)
- To generate some errors in the go-app you can use /error endpoint.
  These kinds of logs can be useful for monitoring and aggregating (and filtering by certain criteria like level) logs.
  It's possible to add new targets to the Promtail config map to collect logs from other sources.
