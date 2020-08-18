---
title: Bind a Redis ServiceInstance to a microservice
type: Getting Started
---

## Bind a Redis ServiceInstance to the microservice

<div tabs name="bind-redis-to-microservice" group="bind-redis-to-microservice">
  <details>
  <summary label="cli">
  CLI
  </summary>

1. Create a ServiceBinding CR that points to the newly created ServiceInstance in the **spec.instanceRef** field:

   ```yaml
   cat <<EOF | kubectl apply -f -
   apiVersion: servicecatalog.k8s.io/v1beta1
   kind: ServiceBinding
   metadata:
     name: orders-service
     namespace: orders-service
   spec:
     instanceRef:
       name: redis-instance
   EOF
   ```

2. Check if the ServiceBinding CR was created successfully. The last condition in the CR status should state `Ready True`:

   ```bash
   kubectl get servicebinding orders-service -n orders-service -o=jsonpath="{range .status.conditions[*]}{.type}{'\t'}{.status}{'\n'}{end}"
   ```

3. Create a ServiceBindingUsage CR:

   ```yaml
   cat <<EOF | kubectl apply -f -
   apiVersion: servicecatalog.kyma-project.io/v1alpha1
   kind: ServiceBindingUsage
   metadata:
     name: orders-service
     namespace: orders-service
   spec:
     serviceBindingRef:
       name: orders-service
     usedBy:
       kind: deployment
       name: orders-service
     parameters:
       envPrefix:
         name: "REDIS_"
   EOF
   ```

   - The **spec.serviceBindingRef** and **spec.usedBy** fields are required. **spec.serviceBindingRef** points to the ServiceBinding you have just created and **spec.usedBy** points to the `orders-service` Deployment. More specifically, **spec.usedBy** refers to the name of the Deployment and the cluster-specific [UsageKind CR](/components/service-catalog/#custom-resource-usage-kind) (`kind: deployment`) that defines how Secrets should be injected to `orders-service` microservice when creating a ServiceBinding.

   - The **spec.parameters.envPrefix.name** field is optional. It adds a prefix to all environment variables injected in a Secret to the microservice when creating a ServiceBinding. In our example, **envPrefix** is `REDIS_`, so all environmental variables will follow the `REDIS_{env}` naming pattern.

     > **TIP:** It is considered good practice to use **envPrefix**. In some cases, a microservice/Function must use several instances of a given ServiceClass. Prefixes allow you to distinguish between instances and make sure that one Secret does not overwrite another one.

4. Check if the ServiceBindingUsage CR was created successfully. The last condition in the CR status should state `Ready True`:

   ```bash
   kubectl get servicebindingusage orders-service -n orders-service -o=jsonpath="{range .status.conditions[*]}{.type}{'\t'}{.status}{'\n'}{end}"
   ```

5. If you want see the Secret details from the ServiceBinding execute:

    ```bash
    kubectl get secret orders-service -n orders-service -o go-template='{{range $k,$v := .data}}{{printf "%s: " $k}}{{if not $v}}{{$v}}{{else}}{{$v | base64decode}}{{end}}{{"\n"}}{{end}}'
    ```

    You should get a result similar to the following details:

    ```bash
    HOST: hb-redis-micro-0e965585-9699-443f-b987-38bc6af0e416-redis.orders-service.svc.cluster.local
    PORT: 6379
    REDIS_PASSWORD: 1tvDcINZvp
    ```

    > **NOTE:** In step 3, we defined **envPrefix** as `REDIS_`, so all variables will start with it. For example, the **PORT** variable will take the form of **REDIS_PORT**.

  </details>
  <details>
  <summary label="console-ui">
  Console UI
  </summary>

1. Go to the **Catalog Management** > **Instances** view in the left navigation panel in `orders-service` Namespace.

2. Switch to the **Add-Ons** tab.

3. Select `redis-service` item on the list. You will redirect to details view of `redis-service` Redis instance.

4. Switch to the **Bound Applications** tab.

5. Select **Bind Application**, select from **Select Application** dropdown list **Deployment** > `redis-service`, select **Set prefix for injected variables** and write `REDIS_` in the form field that appears.

   > **NOTE:** The **Prefix for injected variables** field is optional. It adds a prefix to all environment variables injected in a Secret to the Function when creating a ServiceBinding. In our example, the prefix is set to `REDIS_`, so all environmental variables will follow the `REDIS_{ENVIRONMENT_VARIABLE}` naming pattern.

   > **TIP:** It is considered good practice to use prefixes for environment variables. In some cases, a Function must use several instances of a given ServiceClass. Prefixes allow you to distinguish between instances and make sure that one Secret does not overwrite another one.

6. Select **Bind Application** to confirm changes and wait for status `READY` of created ServiceBindingUsage CR.

  </details>
</div>

## Call and test the microservice

> **CAUTION:** If you have a Minikube cluster, you must first add the IP address of exposed Service to the `hosts` file on your machine:

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

4. If we call again the `https://$SERVICE_DOMAIN/orders` URL, then the system should return a response similar to the following:

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

5. Like in the previous tutorial (when we exposed `Orders Service` microservice) remove the Pod created by `orders-service` Deployment, execute command and wait for successful deletion and starting the new one:

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

   [{"orderCode":"762727210","consignmentCode":"76272725","consignmentStatus":"PICKUP_COMPLETE"}]
   ```

   As we can see, new instance of microservice has saved order created in previous steps. In the `Expose a microservice` tutorial, we used the in-memory storage, so in every time when you deleted a microservice's Pod or changed a Deployment definition, the orders details were lost. Using binding to Redis instance, details are stored outside the `orders-service` microservice, so the data persistence will be preserved.

// dopisać do czego mogą słuzyć bindingi
// dac zajawke do nastepnego etapu tutka