---
title: Bind a Redis ServiceInstance to a Function
type: Getting Started
---

This tutorial shows how you can bind a sample instance of the Redis service to a Function. After completing all steps, you will get the Function with encoded Secrets to the service. You can use them for authentication when you connect to the service to implement custom business logic of your Function.

## Bind a Redis ServiceInstance to the Function

<div tabs name="bind-redis-to-function" group="bind-redis-to-function">
  <details>
  <summary label="cli">
  CLI
  </summary>

1. Create a ServiceBinding CR that points to the created Redis instance in previous tutorials in the **spec.instanceRef** field:

   ```yaml
   cat <<EOF | kubectl apply -f -
   apiVersion: servicecatalog.k8s.io/v1beta1
   kind: ServiceBinding
   metadata:
     name: orders-function
     namespace: orders-service
   spec:
     instanceRef:
       name: redis-service
   EOF
   ```

2. Check if the ServiceBinding CR was created successfully. The last condition in the CR status should state `Ready True`:

   ```bash
   kubectl get servicebinding orders-function -n orders-service -o=jsonpath="{range .status.conditions[*]}{.type}{'\t'}{.status}{'\n'}{end}"
   ```

3. Create a ServiceBindingUsage CR:

   ```yaml
   cat <<EOF | kubectl apply -f -
   apiVersion: servicecatalog.kyma-project.io/v1alpha1
   kind: ServiceBindingUsage
   metadata:
     name: orders-function
     namespace: orders-service
   spec:
     serviceBindingRef:
       name: orders-function
     usedBy:
       kind: serverless-function
       name: orders-function
     parameters:
       envPrefix:
         name: "REDIS_"
   EOF
   ```

   - The **spec.serviceBindingRef** and **spec.usedBy** fields are required. **spec.serviceBindingRef** points to the ServiceBinding you have just created and **spec.usedBy** points to the `orders-function` Function. More specifically, **spec.usedBy** refers to the name of the Function and the cluster-specific [UsageKind CR](/components/service-catalog/#custom-resource-usage-kind) (`kind: serverless-function`) that defines how Secrets should be injected to `orders-function` Function when creating a ServiceBinding.

   - The **spec.parameters.envPrefix.name** field is optional. It adds a prefix to all environment variables injected in a Secret to the Function when creating a ServiceBinding. In our example, **envPrefix** is `REDIS_`, so all environmental variables will follow the `REDIS_{env}` naming pattern.

     > **TIP:** It is considered good practice to use **envPrefix**. In some cases, a Function must use several instances of a given ServiceClass. Prefixes allow you to distinguish between instances and make sure that one Secret does not overwrite another one.

4. Check if the ServiceBindingUsage CR was created successfully. The last condition in the CR status should state `Ready True`:

   ```bash
   kubectl get servicebindingusage orders-function -n orders-service -o=jsonpath="{range .status.conditions[*]}{.type}{'\t'}{.status}{'\n'}{end}"
   ```

5. If you want see the Secret details from the ServiceBinding execute:

    ```bash
    kubectl get secret orders-function -n orders-service -o go-template='{{range $k,$v := .data}}{{printf "%s: " $k}}{{if not $v}}{{$v}}{{else}}{{$v | base64decode}}{{end}}{{"\n"}}{{end}}'
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

1. Go to the **Functions** view (under **Development** section) in the left navigation panel and select `orders-function` Function. 

2. Switch to **Configuration** tab, find **Service Bindings** section and select **Create Service Binding**.

3. Select from **Service Instance** dropdown list `redis-service` and enter in **Prefix for injected variables** form field `REDIS_`.

   > **NOTE:** The **Prefix for injected variables** field is optional. It adds a prefix to all environment variables injected in a Secret to the Function when creating a ServiceBinding. In our example, the prefix is set to `REDIS_`, so all environmental variables will follow the `REDIS_{ENVIRONMENT_VARIABLE}` naming pattern.

   > **TIP:** It is considered good practice to use prefixes for environment variables. In some cases, a Function must use several instances of a given ServiceClass. Prefixes allow you to distinguish between instances and make sure that one Secret does not overwrite another one.

4. Select **Create** to confirm changes.

5. If you switch to **Code** tab and find **Environment Variables** section, you should see `REDIS_PORT`, `REDIS_HOST` and `REDIS_REDIS_PASSWORD` items with the `Service Binding` type. It indicate, that environment variable is injected by ServiceBinding.

  </details>
</div>

## Call and test the Function

> **CAUTION:** If you have a Minikube cluster, you must first add the IP address of exposed Service to the `hosts` file on your machine:

1. Retrieve domain of exposed microservice and save it to the environment variable:

   ```bash
   export FUNCTION_DOMAIN=$(kubectl get virtualservices -l apirule.gateway.kyma-project.io/v1alpha1=orders-function.orders-service -n orders-service -o=jsonpath='{.items[*].spec.hosts[0]}')
   ```

2. Run this command in the terminal to call the service:

   ```bash
   curl -ik "https://$FUNCTION_DOMAIN"
   ```

   The system should return a response similar to the following:

   ```bash
   content-length: 2
   content-type: application/json;charset=UTF-8
   date: Mon, 13 Jul 2020 13:05:33 GMT
   server: istio-envoy
   vary: Origin
   x-envoy-upstream-service-time: 37

   [{"orderCode":"762727210","consignmentCode":"76272725","consignmentStatus":"PICKUP_COMPLETE"}, {"orderCode":"123456789","consignmentCode":"76272725","consignmentStatus":"PICKUP_COMPLETE"}]
   ```

   > **NOTE**: The orders are created in previous tutorials. // w tutorialach odnośnie microserwisu

3. Send a `POST` request to the microservice with a sample order details:

   ```bash
   curl -ikX POST "https://$FUNCTION_DOMAIN" \
     -H "Content-Type: application/json" \
     -H 'Cache-Control: no-cache' -d \
     '{
       "orderCode": "762727234",
       "consignmentCode": "76272725",
       "consignmentStatus": "PICKUP_COMPLETE"
     }'
   ```

4. If we call again the `https://$FUNCTION_DOMAIN` URL, then the system should return a response similar to the following:

   ```bash
   HTTP/2 200
   content-length: 73
   content-type: application/json;charset=UTF-8
   date: Mon, 13 Jul 2020 13:05:51 GMT
   server: istio-envoy
   vary: Origin
   x-envoy-upstream-service-time: 6

   [{"orderCode":"762727234","consignmentCode":"76272725","consignmentStatus":"PICKUP_COMPLETE"}, {"orderCode":"762727210","consignmentCode":"76272725","consignmentStatus":"PICKUP_COMPLETE"}, {"orderCode":"123456789","consignmentCode":"76272725","consignmentStatus":"PICKUP_COMPLETE"}]
   ```

5. Like in the previous tutorial (when we exposed Function) remove the Pod created by `orders-function` Function, execute command and wait for successful deletion and starting the new one:

   ```bash
   kubectl delete pod -n orders-service -l "serverless.kyma-project.io/function-name=orders-function"
   ```

6. Again call the Function to check the storage:

   ```bash
   curl -ik "https://$FUNCTION_DOMAIN"
   ```

   The system should return a response similar to the following:

   ```bash
   content-length: 2
   content-type: application/json;charset=UTF-8
   date: Mon, 13 Jul 2020 13:05:33 GMT
   server: istio-envoy
   vary: Origin
   x-envoy-upstream-service-time: 37

   [{"orderCode":"762727234","consignmentCode":"76272725","consignmentStatus":"PICKUP_COMPLETE"}, {"orderCode":"762727210","consignmentCode":"76272725","consignmentStatus":"PICKUP_COMPLETE"}, {"orderCode":"123456789","consignmentCode":"76272725","consignmentStatus":"PICKUP_COMPLETE"}]
   ```

   As we can see, new instance of Function has saved order created in previous steps. In the `Expose a Function` tutorial, we used the in-memory storage, so in every time when you deleted a Function's Pod or changed a Function definition, the orders details were lost. Using binding to Redis instance, details are stored outside the `orders-function` Function, so the data persistence will be preserved. Also Redis can be shared between Function and microservice.