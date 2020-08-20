---
title: Create a Redis service
type: Getting Started
---

// mozna napisac fajne wprowadzenie po co nam Redis -> patrz poprzedni dok, gdzie uzywalismy in-memory storage.

This tutorial shows how you can provision a sample [Redis](https://redis.io/) service using an Addon configuration linking to an example in the GitHub repository.

## Add Addon with Redis service

Follows these steps:

<div tabs name="add-addon" group="create-redis-service">
  <details>
  <summary label="cli">
  CLI
  </summary>

1. Provision an Addon CR with the Redis service:

   ```yaml
   cat <<EOF | kubectl apply -f  -
   apiVersion: addons.kyma-project.io/v1alpha1
   kind: AddonsConfiguration
   metadata:
     name: redis-addon
     namespace: orders-service
   spec:
     repositories:
     - url: https://github.com/kyma-project/addons/releases/download/0.12.0/index-testing.yaml
   EOF
   ```

2. Check if the Addon CR was created successfully. The CR phase should state `Ready`:

   ```bash
   kubectl get addonsconfigurations redis-addon -n orders-service -o=jsonpath="{.status.phase}"
   ```

  </details>
  <details>
  <summary label="console-ui">
  Console UI
  </summary>

1. If you aren't in the view of Namespace `orders-service` in the Kyma Console, select a `orders-service` Namespace from the drop-down list in the top navigation panel.

2. Go to the **Addons** view in the left navigation panel (under **Configuration** group) and select **Add New Configuration**.

3. Enter `https://github.com/kyma-project/addons/releases/download/0.12.0/index-testing.yaml` in the **Urls** field. The Addon name is automatically generated.

4. Select **Add** to confirm changes.

5. Wait for the Addon to have the status `READY`.

  </details>
</div>

// potem tworzymy instancje z addona

## Create a Redis instance

Follows these steps:

<div tabs name="create-redis-instance" group="create-redis-service">
  <details>
  <summary label="cli">
  CLI
  </summary>

1. Create a ServiceInstance CR. You will provision Redis with its `micro` plan:

   ```yaml
   cat <<EOF | kubectl apply -f -
   apiVersion: servicecatalog.k8s.io/v1beta1
   kind: ServiceInstance
   metadata:
     name: redis-service
     namespace: orders-service
   spec:
     serviceClassExternalName: redis
     servicePlanExternalName: micro
     parameters:
       imagePullPolicy: Always
   EOF
   ```

3. Check if the ServiceInstance CR was created successfully. The last condition in the CR status should state `Ready True`:

   ```bash
   kubectl get serviceinstance redis-service -n orders-service -o=jsonpath="{range .status.conditions[*]}{.type}{'\t'}{.status}{'\n'}{end}"
   ```

  </details>
  <details>
  <summary label="console-ui">
  Console UI
  </summary>

1. Go to the **Catalog Management** > **Catalog** view where you can see the list of all available **Add-Ons** and select **[Experimental] Redis**.

   > **TIP**: You can use the search in the upper right corner.

2. Select **Add** to provision the Redis ServiceClass and create its instance in your Namespace.

3. Change the **Name** to match `redis-service`, select `micro` from the **Plan** drop-down list, and set **Image pull policy** to `Always`.

  // mozna to przeredagowac
   > **NOTE:** The Service Instance, Service Binding, and Service Binding Usage can have different names than the Function, but it is recommended that all related resources share a common name.

4. Select **Create** to confirm changes.

   You will redirect to **Catalog Management** > **Instances** > **redis-service** view.

   Wait until the status of the instance changes from `PROVISIONING` to `RUNNING`.

  </details>
</div>

We have a provisioned Redis instance. In next steps we will bind this instance to deployed Orders Service.
