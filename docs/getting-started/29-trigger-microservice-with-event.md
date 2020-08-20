---
title: Trigger a microservice with an event
type: Getting Started
---

This tutorial shows how to trigger a deployed in previous tutorial `orders-service` microservice with an `order.deliverysent.v1` event from an Application connected to Kyma.

## Create the Trigger

<div tabs name="steps" group="trigger-microservice">
  <details>
  <summary label="cli">
  CLI
  </summary>

1. Create a Trigger CR for `orders-service` microservice to subscribe application to an `order.deliverysent.v1` event from `commerce-mock` Application:

   ```bash
   cat <<EOF | kubectl apply -f  -
   apiVersion: eventing.knative.dev/v1alpha1
   kind: Trigger
   metadata:
     name: orders-service
     namespace: orders-service
   spec:
     broker: default
     filter:
       attributes:
         eventtypeversion: v1
         source: commerce-mock
         type: order.deliverysent
     subscriber:
       ref:
         apiVersion: v1
         kind: Service
         name: orders-service
         namespace: orders-service
   EOF
   ```

   - **spec.filter.attributes.eventtypeversion** points to the specific event version, on our case it is `v1`.
   - **spec.filter.attributes.source** is taken from the name of the Application CR and specifies the source of events. In our example, it is created `commerce-mock` mock.
   - **spec.filter.attributes.type** points to the given event type to which you want to subscribe microservice. In our case, it is `order.deliverysent`.

2. Check if the Trigger CR was created successfully and is ready. The CR `Ready` condition should state `True`:

   ```bash
   kubectl get trigger orders-service -n orders-service -o=jsonpath="{.status.conditions[2].status}"
   ```

  </details>
  <details>
  <summary label="console-ui">
  Console UI
  </summary>

1. If you aren't in the view of Namespace `orders-service` in the Kyma Console, select a `orders-service` Namespace from the drop-down list in the top navigation panel.

2. Go to the **Services** view (under **Operation** section) in the left navigation panel and navigate to `orders-service` Service.

3. Once in the Service details view, select **Add Event Trigger** in the **Event Triggers** section.

4. Find `order.deliverysent` event with `v1` version from `commerce-mock` Application, check it and click **Add**.

   The message appears on the UI confirming that the Event Trigger was successfully created, and you will see it in the **Event Triggers** section of Service's details view.

  </details>
</div>

## Test the Trigger

To send events from mock to Orders Service application, follow these steps:  

1. Access the SAP Commerce Cloud Mock mock at `https://commerce-orders-service.{CLUSTER_DOMAIN}.` or go to **API Rules** view (under **Configuration** section) in `orders-service` Namespace and select the mock, you will the direct link to the mock application under **Host** column.

2. Switch to **Remote APIs** tab, find **SAP Commerce Cloud - Events** and click it.

3. In opened view search in dropdown list `order.deliverysent.v1` event. In pasted event change `orderCode` to `123456789` and select **Send Event**.

   The message appears on the UI confirming that the event was successfully sent.

4. For the last time call the microservice to check the storage:

   ```bash
   curl -ik "https://$SERVICE_DOMAIN/orders"
   ```

   > **NOTE**: To get the microservice domain, run:
   >
   > ```bash
   > export SERVICE_DOMAIN=$(kubectl get virtualservices -l apirule.gateway.kyma-project.io/v1alpha1=orders-service.orders-service -n orders-service -o=jsonpath='{.items[*].spec.hosts[0]}')
   > ```

   You should see a response similar to the following:

   ```bash
   content-length: 2
   content-type: application/json;charset=UTF-8
   date: Mon, 13 Jul 2020 13:05:33 GMT
   server: istio-envoy
   vary: Origin
   x-envoy-upstream-service-time: 37

   [{"orderCode":"762727210","consignmentCode":"76272725","consignmentStatus":"PICKUP_COMPLETE"}, {"orderCode":"123456789","consignmentCode":"76272725","consignmentStatus":"PICKUP_COMPLETE"}]
   ```

   The event from mock application was saved in Redis instance :)

// koniec tej trudnej przygody XD
