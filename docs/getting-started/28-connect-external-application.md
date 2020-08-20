---
title: Connect an external application
type: Getting Started
---

// zmienic na order.deliverysent.v1.v1

// po stworzeniu microserwisu, bindingu itp dopiero tworzymy mocka

// podlinkowac tutek jak dodac eventy pod aplikacje z istniejacej apki

After creating Namespace, let's integrate an external application to Kyma. In this set of guides, we will use a mock application called [SAP Commerce Cloud - Mock](https://github.com/SAP-samples/xf-addons/tree/master/addons/commerce-mock-0.1.0) that is to simulate a monolithic application. You will learn how you can connect it to Kyma, and expose its API and events. We will subscribe to one of its events (**order.deliverysent.v1**) in other tutorials and use it to trigger the logic of a sample service and Function.

## Deploy the XF Addons and provision the Commerce mock

Commerce mock is a part of the XF Addons which give access to 3 instances of mocks that simulate external applications sending events to Kyma.

Follow these steps to deploy XF Addons and add the Commerce mock to your Namespace:

<div tabs name="provision-mock" group="connect-external-application">
  <details>
  <summary label="cli">
  CLI
  </summary>

1. Provision an Addon CR with the mocks:

   ```bash
   cat <<EOF | kubectl apply -f  -
   apiVersion: addons.kyma-project.io/v1alpha1
   kind: AddonsConfiguration
   metadata:
     name: xf-mocks
     namespace: orders-service
   spec:
     repositories:
     - url: github.com/sap/xf-addons//addons/index.yaml
   EOF
   ```

   > **NOTE**: The `index.yaml` file is an addons manifest with APIs of SAP Marketing Cloud, SAP Cloud for Customer, and SAP Commerce Cloud applications.

2. Check if the Addon CR was created successfully. The CR phase should state `Ready`:

   ```bash
   kubectl get addonsconfigurations xf-mocks -n orders-service -o=jsonpath="{.status.phase}"
   ```

3. Create a ServiceInstance CR with the mock:

   // sprawdzić to
   ```bash
   cat <<EOF | kubectl apply -f -
   apiVersion: servicecatalog.k8s.io/v1beta1
   kind: ServiceInstance
   metadata:
     name: commerce-mock
     namespace: orders-service
   spec:
     serviceClassExternalName: commerce-mock
     servicePlanExternalName: default
   EOF
   ```

4. Check if the ServiceInstance CR was created successfully. The last condition in the CR status should state `Ready True`:

   ```bash
   kubectl get serviceinstance commerce-mock -n orders-service -o=jsonpath="{range .status.conditions[*]}{.type}{'\t'}{.status}{'\n'}{end}"
   ```

  </details>
  <details>
  <summary label="console-ui">
  Console UI
  </summary>

1. If you aren't in the view of Namespace `orders-service` in the Kyma Console, select a `orders-service` Namespace from the drop-down list in the top navigation panel.

2. Go to the **Addons** view in the left navigation panel (under **Configuration** section) and select **Add New Configuration**.

3. Enter `github.com/sap/xf-addons//addons/index.yaml` in the **Urls** field. The Addon name is automatically generated.

   > **NOTE**: The `index.yaml` file is an addons manifest with APIs of SAP Marketing Cloud, SAP Cloud for Customer, and SAP Commerce Cloud applications.

4. Select **Add** to confirm changes.

5. Wait for the Addon to have the status `READY`.

6. Got to **Catalog** view (under **Service Management** group) and then to **Add-Ons** tab.

7. Select the mock you want to provision. For this example, use **[Preview] SAP Commerce Cloud - Mock**.

   > **TIP**: You can use the search in the upper right corner.

8. Click **Add once** to deploy it in your Namespace. The mock name is automatically generated.

9. Select **Create** to confirm changes.

// zostawić te generated, czy opisać, ze mozna zmienic nazwe

   You will redirect to **Catalog Management** > **Instances** > **{Generated mock name}** view.

10. Wait for the mock to have the status `RUNNING`.

  </details>
</div>

## Connect the mock application to Kyma

After provisioning the mock, connect it to Kyma.

### Create the Application and retrieve token

First create the Application CR and then retrieve token to connect the mock to an Application. Follow these steps:

<div tabs name="create-application" group="connect-external-application">
  <details>
  <summary label="cli">
  CLI
  </summary>

1. Apply an Application definition to the cluster:

   ```bash
   cat <<EOF | kubectl apply -f -
   apiVersion: applicationconnector.kyma-project.io/v1alpha1
   kind: Application
   metadata:
     name: commerce-mock
   spec:
     description: "The Application for Commerce mock"
     labels:
       app: orders-service
       example: orders-service
   EOF
   ```

2. Check if the Application CR was created successfully. The CR phase should state `deployed`:

   ```bash
   kubectl get application commerce-mock -o=jsonpath="{.status.installationStatus.status}"
   ```

3. Get a token to connect the mock to an Application. For that create a TokenRequest CR. The CR name must match the name of the Application for which you want to get the configuration details. Run:

   ```bash
   cat <<EOF | kubectl apply -f -
   apiVersion: applicationconnector.kyma-project.io/v1alpha1
   kind: TokenRequest
   metadata:
     name: commerce-mock
   EOF
   ```

4. Fetch the TokenRequest CR you created to get the token from the status section. Run:

   ```bash
   kubectl get tokenrequest commerce-mock -o=jsonpath="{.status.url}"
   ```

   > **NOTE**: If the response doesn't contain any content, wait for a few moments and run command again.

   A successful call should return a response similar to the following:

   ```bash
   https://connector-service.{CLUSTER_DOMAIN}/v1/applications/signingRequests/info?token=h31IwJiLNjnbqIwTPnzLuNmFYsCZeUtVbUvYL2hVNh6kOqFlW9zkHnzxYFCpCExBZ_voGzUo6IVS_ExlZd4muQ==
   ```

   Save this token to the clipboard, it will be needed in the next steps.

  </details>
  <details>
  <summary label="console-ui">
  Console UI
  </summary>

1. Back in the general Console UI view (clicking **Back to Namespaces**).

2. Go to the **Applications/Systems** view (under **Integration** section), click **Create Application** and set Application's name to `commerce-mock`.

3. Wait for the Application to have the status `Serving`.

4. Open the newly created application and click **Connect Application**.

5. Copy the token and select **OK** to close the pop-up box.

  </details>
</div>

### Connect events from mock to an Application 

To connect events from mock to created Application, follow these steps:  

1. Access the SAP Commerce Cloud Mock mock at `https://commerce-orders-service.{CLUSTER_DOMAIN}` or go to **API Rules** view (under **Configuration** section) in `orders-service` Namespace and select the mock, you will the direct link to the mock application under **Host** column.

2. Click **Connect**.

3. Paste the saved/copied token in previous steps and wait for the Application to connect.

4. Select **Register All** or just register **SAP Commerce Cloud - Events** to be able to send events.

   >**NOTE:** Local APIs are the ones available with the mock application. Remote APIs represent the ones registered in Kyma.

### Expose events in a Namespace

To expose events in a Namespace, first create an ApplicationMapping CR in the cluster to bind an Application to the Namespace. Then provision the Events API in the Namespace by ServiceInstance CR. Follow the instructions:

<div tabs name="expose-events-in-namespace" group="connect-external-application">
  <details>
  <summary label="cli">
  CLI
  </summary>

1. Create an ApplicationMapping CR and apply it to the cluster:

   ```bash
   cat <<EOF | kubectl apply -f -
   apiVersion: applicationconnector.kyma-project.io/v1alpha1
   kind: ApplicationMapping
   metadata:
     name: commerce-mock
     namespace: orders-service
   EOF
   ```

2. List available ServiceClass CRs in the `orders-service` Namespace and find one with the `EXTERNAL-NAME` prefix `sap-commerce-cloud-events-*`. 

   ```bash
   kubectl get serviceclasses -n orders-service
   ```

   Copy the full `EXTERNAL NAME` to environment variable like:

   ```bash
   export EVENTS_EXTERNAL_NAME="sap-commerce-cloud-events-58d21"
   ```

3. Provision the Events API in the Namespace by ServiceInstance CR:

   ```bash
   cat <<EOF | kubectl apply -f -
   apiVersion: servicecatalog.k8s.io/v1beta1
   kind: ServiceInstance
   metadata:
     name: commerce-mock-events
     namespace: orders-service
   spec:
     serviceClassExternalName: $EVENTS_EXTERNAL_NAME
     servicePlanExternalName: default
   EOF
   ```

4. Check if the ServiceInstance CR was created successfully. The last condition in the CR status should state `Ready True`:

   ```bash
   kubectl get serviceinstance commerce-mock-events -n orders-service -o=jsonpath="{range .status.conditions[*]}{.type}{'\t'}{.status}{'\n'}{end}"
   ```

  </details>
  <details>
  <summary label="console-ui">
  Console UI
  </summary>

1. Back to your Application view in the Console UI, select **Create Binding** to bind the Application to your Namespace where you will later provision the Events API provided by the mocks. Select `orders-service` Namespace and click **Create**.

2. Go to `orders-service` Namespace view, then to **Service Catalog** and open **Services** tab. Find the **SAP Commerce Cloud - Events** Service, select service and click **Add once** to add it to the Namespace.

   > **TIP**: You can use the search in the upper right corner.

// tak samo tu - opisać, ze mozna mieć inną nazwe?

   You will redirect to **Catalog Manegement** > **Instances** > **{Generated name}** view.

3. Wait for the Events API to have the status `RUNNING`.

  </details>
</div>

After all these steps, applications/Functions running in Namespace `orders-service` can consume events from mock.
