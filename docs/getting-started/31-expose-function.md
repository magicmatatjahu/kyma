---
title: Expose a Function
type: Getting Started
---

This tutorial shows how you can expose a Function to access it outside the cluster, through an HTTP proxy. To expose it, use an APIRule custom resource (CR) managed by the in-house API Gateway Controller. This controller reacts to an instance of the APIRule CR and, based on its details, it creates an Istio Virtual Service and Oathkeeper Access Rules that specify your permissions for the exposed Function.

When you complete this tutorial, you get a Function that:

- Is available under an unsecured endpoint (**handler** set to `noop` in the APIRule CR).
- Accepts `GET` and `POST` methods.

>**NOTE:** To learn more about securing your Function, see the [tutorial](/components/api-gateway#tutorials-expose-and-secure-a-service-deploy-expose-and-secure-the-sample-resources).

## Prerequisites

This tutorial is based on an existing Function. To create one, follow the [Create a Function](#tutorials-create-a-function) tutorial.

## Steps

Follows these steps:

<div tabs name="steps" group="expose-function">
  <details>
  <summary label="cli">
  CLI
  </summary>

// nie wiem czy wspominac o tym nocie
    >**NOTE:** Function takes the name from the Function CR name. The APIRule CR can have a different name but for the purpose of this tutorial, all related resources share a common name defined under the **NAME** variable.

1. Create an APIRule CR for your Function. It is exposed on port `80` that is the default port of the [Service](#architecture-architecture).

    ```bash
    cat <<EOF | kubectl apply -f -
    apiVersion: gateway.kyma-project.io/v1alpha1
    kind: APIRule
    metadata:
      name: orders-function
      namespace: orders-service
    spec:
      gateway: kyma-gateway.kyma-system.svc.cluster.local
      rules:
      - path: /.*
        accessStrategies:
        - config: {}
          handler: noop
        methods: ["GET","POST"]
      service:
        host: orders-function
        name: orders-function
        port: 80
    EOF
    ```

2. Check if the API Rule was created successfully and has the `OK` status:

   ```bash
   kubectl get apirules orders-function -n orders-service -o=jsonpath='{.status.APIRuleStatus.code}'
   ```

3. Access the Function's external address:

   ```bash
   curl https://orders-function.{CLUSTER_DOMAIN}
   ```

  </details>
  <details>
  <summary label="console-ui">
  Console UI
  </summary>

1. Go to the **Functions** view (under **Development** section) in the left navigation panel and select `orders-function` Function. 

   You will redirect to Function's details view.

2. Switch to **Configuration** tab, find **API Rules** section and select **Expose Function**.

   The modal box should appear. This is embedded view for creating API Rule CR from **API Rules** view. 

3. In the **General settings** section:

    - Enter the API Rule's **Name** matching the Function's name: `orders-function`.

    > **NOTE:** The APIRule CR can have a different name than the Function, but it is recommended that all related resources share a common name.

    - Enter `orders-function` as **Hostname** to indicate the host on which you want to expose your Service.

    > **NOTE**: Check that `orders-function` Service is automatically selected in **Service** dropdown

4. In the **Access strategies** section, leave the default settings, with `GET` and `POST` methods and the `noop` handler selected.

5. Select **Create** to confirm changes.

    The message appears on the screen confirming the changes were saved.

6. The modal box should close. Check if you can access the Function by selecting the HTTPS link under **Host** column of just created `orders-function` API Rule.

  </details>
</div>

## Call and test the microservice

> **CAUTION:** If you have a Minikube cluster, you must first add the IP address of exposed Service to the `hosts` file on your machine:

<!-- Improve this caution message to explain exactly why we do that-->

```bash
echo "$(minikube ip) orders-function.kyma.local" | sudo tee -a /etc/hosts
```

// mozna to pobrac z consoli (**HOST** kolumna)
1. Retrieve domain of exposed Function and save it to the environment variable:

   ```bash
   export FUNCTION_DOMAIN=$(kubectl get virtualservices -l apirule.gateway.kyma-project.io/v1alpha1=orders-function.orders-service -n orders-service -o=jsonpath='{.items[*].spec.hosts[0]}')
   ```

2. Like in tutorial of exposed microservice (jak w tutorialu do expozowania microserwisu), run this command in the terminal to call the Function:

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

   []
   ```

3. Send a `POST` request to the Function with a sample order details:

   ```bash
   curl -ikX POST "https://$FUNCTION_DOMAIN" \
     -H "Content-Type: application/json" \
     -H 'Cache-Control: no-cache' -d \
     '{
       "orderCode": "762727210",
       "consignmentCode": "76272725",
       "consignmentStatus": "PICKUP_COMPLETE"
     }'
   ```

4. Again call the Function to check the storage:

   ```bash
   curl -ik "https://$FUNCTION_DOMAIN"
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

5. Remove the [Pod](https://kubernetes.io/docs/concepts/workloads/pods/) created by `orders-function` Function, execute command and wait for successful deletion and starting the new one:

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

   []
   ```

   As you can see, Function like microservice used the in-memory storage, so in every time when you delete a Function's Pod or change a Function definition, the orders details will be lost. In the next steps, we will bind a Redis instance to Function to save orders.
