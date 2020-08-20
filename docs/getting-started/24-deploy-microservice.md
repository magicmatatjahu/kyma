---
title: Deploy a microservice
type: Getting Started
---

Learn how to quickly deploy a standalone [`Orders Service`](https://github.com/kyma-project/examples/blob/master/orders-service/README.md) microservice on a Kyma cluster.

You will create:

- Deployment in which you specify the application configuration
- Kubernetes Service through which your application will communicate with other resources on the Kyma cluster

## Steps

// to powinno iść do overview guida :) - powinno to byc usuniete

### Get the kubeconfig file and configure the CLI

Follow these steps to get the `kubeconfig` file and configure the CLI to connect to the cluster:

1. Access the Console UI of your Kyma cluster.
2. Click the user icon in the top right corner.
3. Select **Get Kubeconfig** from the drop-down menu to download the configuration file to a selected location on your machine.
4. Open a terminal window.
5. Export the **KUBECONFIG** environment variable to point to the downloaded `kubeconfig`. Run this command:

   ```bash
   export KUBECONFIG={KUBECONFIG_FILE_PATH}
   ```

   >**NOTE:** Drag and drop the `kubeconfig` file in the terminal to easily add the path of the file to the `export KUBECONFIG` command you run.

6. Run `kubectl cluster-info` to check if the CLI is connected to the correct cluster.

### Create a Deployment

Create a [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) that provides the application definition and enables you to run it on the cluster. The Deployment uses the `eu.gcr.io/kyma-project/pr/orders-service:PR-162` image. This Docker image exposes the `8080` port on which the related Service is listening.

<div tabs name="create-deployment" group="create-microservice">
  <details>
  <summary label="cli">
  CLI
  </summary>

1. Apply an application definition to the cluster:

   ```bash
   cat <<EOF | kubectl apply -f -
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: orders-service
     namespace: orders-service
     labels:
       app: orders-service
       example: orders-service
   spec:
     replicas: 1
     selector:
       matchLabels:
         app: orders-service
         example: orders-service
     template:
       metadata:
         labels:
           app: orders-service
           example: orders-service
       spec:
         containers:
           - name: orders-service
             image: "eu.gcr.io/kyma-project/pr/orders-service:PR-162"
             imagePullPolicy: IfNotPresent
             resources:
               limits:
                 cpu: 20m
                 memory: 32Mi
               requests:
                 cpu: 10m
                 memory: 16Mi
             env:
               - name: APP_PORT
                 value: "8080"
               - name: APP_REDIS_PREFIX
                 value: "REDIS_"
   EOF
   ```

2. Check if the Deployment was created successfully. The Deployment status should have set `readyReplicas` to 1:

   ```bash
   kubectl get deployment orders-service -n orders-service -o=jsonpath="{.status.readyReplicas}"
   ```

  </details>
  <details>
  <summary label="console-ui">
  Console UI
  </summary>

1. Create a YAML file with the Deployment definition:

   ```yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: orders-service
     namespace: orders-service
     labels:
       app: orders-service
       example: orders-service
   spec:
     replicas: 1
     selector:
       matchLabels:
         app: orders-service
         example: orders-service
     template:
       metadata:
         labels:
           app: orders-service
           example: orders-service
       spec:
         containers:
           - name: orders-service
             image: "eu.gcr.io/kyma-project/pr/orders-service:PR-162"
             imagePullPolicy: IfNotPresent
             resources:
               limits:
                 cpu: 20m
                 memory: 32Mi
               requests:
                 cpu: 10m
                 memory: 16Mi
             env:
               - name: APP_PORT
                 value: "8080"
               - name: APP_REDIS_PREFIX
                 value: "REDIS_"
   ```

2. Go to the `orders-service` Namespace view in the Console UI and select the **Deploy new resource** button.

3. Browse your Deployment file and select **Deploy** to confirm changes.

4. Go to the **Deployments** view (under **Operation** section) to make sure `orders-service` Deployment has `RUNNING` status.

  </details>
</div>

### Create the Service

Deploy the Kubernetes `orders-service` [Service](https://kubernetes.io/docs/concepts/services-networking/service/) in the `orders-service` Namespace to allow other Kubernetes resources to communicate with your application.

<div tabs name="create-service" group="create-microservice">
  <details>
  <summary label="cli">
  CLI
  </summary>

Create a Kubernetes Service in the cluster:

```bash
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Service
metadata:
  name: orders-service
  namespace: orders-service
  labels:
    app: orders-service
    example: orders-service
spec:
  type: ClusterIP
  ports:
    - name: http
      port: 80
      protocol: TCP
      targetPort: 8080
  selector:
    app: orders-service
    example: orders-service
EOF
```

  </details>
  <details>
  <summary label="console-ui">
  Console UI
  </summary>

1. Create a YAML file with the Service definition:

   ```yaml
   apiVersion: v1
   kind: Service
   metadata:
     name: orders-service
     namespace: orders-service
     labels:
       app: orders-service
       example: orders-service
   spec:
     type: ClusterIP
     ports:
       - name: http
         port: 80
         protocol: TCP
         targetPort: 8080
     selector:
       app: orders-service
       example: orders-service
   ```

2. Go to the `orders-service` Namespace view in the Console UI and select the **Deploy new resource** button.

3. Browse your Service file and select **Deploy** to confirm changes.

4. Go to the **Services** view (under **Operation** group) to make sure `orders-service` Service has `RUNNING` status.

  </details>
</div>
