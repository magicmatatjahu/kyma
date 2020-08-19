---
title: Expose a microservice
type: Getting Started
---

Now that you deployed a standalone `orders-service` application in a cluster, you can make it available outside the cluster to other resources by exposing its Kubernetes Service.

## Prerequisites

Go through the [Deploy a application](/#getting-started-deploy-a-microservice) tutorial to apply the `orders-service` application in the `orders-service` Namespace on your cluster.

## Steps

## Expose the Service

Create an APIRule resource which exposes the Kubernetes Service of the application under na unsecured endpoint (**handler** set to `noop`) and accepts the `GET` and `POST` methods.

> **TIP:** If you prefer to secure your Service, read the [tutorial](/components/api-gateway/#tutorials-expose-and-secure-a-service) to learn how to do that.

<div tabs name="create-apirule" group="expose-microservice">
  <details>
  <summary label="cli">
  CLI
  </summary>

1. Open the terminal and apply the APIRule:

   ```bash
   cat <<EOF | kubectl apply -f -
   apiVersion: gateway.kyma-project.io/v1alpha1
   kind: APIRule
   metadata:
     name: orders-service
     namespace: orders-service
     labels:
       app: orders-service
       example: orders-service
   spec:
     service:
       host: orders-service
       name: orders-service
       port: 80
     gateway: kyma-gateway.kyma-system.svc.cluster.local
     rules:
       - path: /.*
         methods: ["GET","POST"]
         accessStrategies:
           - handler: noop
         mutators: []
   EOF
   ```

2. Check if the API Rule was created successfully and has the `OK` status:

   ```bash
   kubectl get apirules orders-service -n orders-service -o=jsonpath='{.status.APIRuleStatus.code}'
   ```

  </details>
  <details>
  <summary label="console-ui">
  Console UI
  </summary>

// moze dodac tylko to co ponizej? a nie dwa tipy?
> **TIP:** You can expose a Service or Function with an API Rule from different views in the Console UI. This tutorial shows how to do that from the generic **API Rules** view.

> **TIP:** Console UI has a separate **API Rules** view from which you to create APIRules for Services and Functions. Still, you can access this view directly from their corresponding views.

1. Select the `orders-service` Namespace from the drop-down list in the top navigation panel.

2. Go to the **API Rules** view (under **Configuration** group) select **Add API Rule**.

3. In the **General settings** section:

    - Enter `orders-service` as the API Rule's **Name**.

    > **NOTE:** The APIRule CR can have a different name than the Service, but it is recommended that all related resources share a common name.

    - Enter `orders-service` as **Hostname** to indicate the host on which you want to expose your Service.

    - Select the `orders-service` Service from the drop-down list in the **Service** column.

// dodać DELETE?

4. In the **Access strategies** section, leave only the `GET` and `POST` methods marked and the `noop` handler selected.

5. Select **Create** to confirm changes.

    The message appears on the screen confirming the changes were saved.

6. In the API Rule's details view that opens up automatically, check if you can access the Service by selecting the HTTPS link under **Host**.

  </details>
</div>

## Call and test the microservice

> **CAUTION:** If you have a Minikube cluster, you must first add the IP address of exposed Service to the `hosts` file on your machine:

<!-- Improve this caution message to explain exactly why we do that-->

```bash
echo "$(minikube ip) orders-service.kyma.local" | sudo tee -a /etc/hosts
```

// mozna to pobrac z consoli (**HOST** kolumna)
1. Retrieve domain of exposed microservice and save it to the environment variable:

   ```bash
   export SERVICE_DOMAIN=$(kubectl get virtualservices -l apirule.gateway.kyma-project.io/v1alpha1=orders-service.orders-service -n orders-service -o=jsonpath='{.items[*].spec.hosts[0]}')
   ```

2. Run this command in the terminal to call the service:

   ```bash
   curl -ik "https://$SERVICE_DOMAIN/orders"
   ```

   The system should return a response similar to the following:

   ```bash
   content-length: 2
   content-type: application/json;charset=UTF-8
   date: Mon, 13 Jul 2020 13:05:33 GMT
   server: istio-envoy
   vary: Origin
   x-envoy-upstream-service-time: 37

   []
   ```

3. Send a `POST` request to the microservice with a sample order details:

   ```bash
   curl -ikX POST "https://$SERVICE_DOMAIN/orders" \
     -H 'Cache-Control: no-cache' -d \
     '{
       "orderCode": "762727210",
       "consignmentCode": "76272725",
       "consignmentStatus": PICKUP_COMPLETE
     }'
   ```

4. Again call the microservice to check the storage:

   ```bash
   curl -ik "https://$SERVICE_DOMAIN/orders"
   ```

   The system should return a response similar to the following:

   ```bash
   HTTP/2 200
   content-length: 73
   content-type: application/json;charset=UTF-8
   date: Mon, 13 Jul 2020 13:05:51 GMT
   server: istio-envoy
   vary: Origin
   x-envoy-upstream-service-time: 6

   [{"orderCode":"762727210","consignmentCode":"76272725","consignmentStatus":"PICKUP_COMPLETE"}]
   ```

   You can see the service returns the order details previously sent to it.

5. Remove the [Pod](https://kubernetes.io/docs/concepts/workloads/pods/) created by `orders-service` Deployment, execute command and wait for successful deletion and starting the new one:

   ```bash
   kubectl delete pod -n orders-service -l app=orders-service
   ```

6. Again call the microservice to check the storage:

   ```bash
   curl -ik "https://$SERVICE_DOMAIN/orders"
   ```

   The system should return a response similar to the following:

   ```bash
   content-length: 2
   content-type: application/json;charset=UTF-8
   date: Mon, 13 Jul 2020 13:05:33 GMT
   server: istio-envoy
   vary: Origin
   x-envoy-upstream-service-time: 37

   []
   ```

   As you can see, microservice used the in-memory storage, so in every time when you delete a microservice's Pod or change a Deployment definition, the orders details will be lost. We will deal with this problem in the next steps of this guide.
