---
title: Setup a namespace
type: Getting Started
---

// dodać info, ze na potrzeby guida tworzymy wszytsko w jednym namespacie, ale mozna w roznych :) Jeśli pojawi sie odniesienie do orders-service Namespacu to mozna uzyc innego gdzie dana appka/funckja istnieje :)

// zmienic envy na wartosci inline

In this guide almost every operations will be performed using Namespace scoped resources, so let's start with setup a namespace with `orders-service` name.

Follow these steps to create a namespace.

<div tabs name="setup-namespace" group="setup-namespace">
  <details>
  <summary label="cli">
  CLI
  </summary>

1. Create a namespace:

   ```bash
   kubectl create ns orders-service
   ```

2. Check if the Namespace was setup successfully. The Namespace phase should state `Active`:

   ```bash
   kubectl get ns orders-service -o=jsonpath="{.status.phase}"
   ```

  </details>
  <details>
  <summary label="console-ui">
  Console UI
  </summary>

1. Log in to Kyma Console UI.
// tutaj nie wiem co napisać, bo nie umiem znaleźć zadnego tutka do otwierania consoli :(
// moze pomyslec o pominieciu kroku
// https://kyma-project.io/docs/1.12/root/kyma/#installation-install-kyma-on-a-cluster-access-the-cluster

2. After login, in the **Namespaces** view select **Add New Namespace**.

3. Enter `orders-service` in the **Name** field.

4. Select **Create** to confirm changes.

   You will redirect to `orders-service` Namespace view.

  </details>
</div> 
