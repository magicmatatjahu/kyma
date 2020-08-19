---
title: Create a Function
type: Getting Started
---

// dodać jakiś fajny wstęp 

This tutorial shows how you can create a simple "Hello World!" Function.

## Steps

Follows these steps:

<div tabs name="steps" group="create-function">
  <details>
  <summary label="cli">
  CLI
  </summary>

1. Create a Function CR that specifies the Function's logic:

   ```yaml
   apiVersion: serverless.kyma-project.io/v1alpha1
   kind: Function
   metadata:
     name: orders-function
     namespace: orders-service
     labels:
       app: orders-function
       example: orders-function
   spec:
     maxReplicas: 1
     minReplicas: 1
     resources:
       limits:
         cpu: 20m
         memory: 32Mi
       requests:
         cpu: 10m
         memory: 16Mi
     env:
       - name: APP_REDIS_PREFIX
         value: "REDIS_"
     deps: |-
       { 
         "name": "orders-function",
         "version": "1.0.0",
         "dependencies": {
           "redis": "3.0.2"
         }
       }
     source: |-
      const redis = require("redis");
      const { promisify } = require("util");

      let storage = undefined;
      const ordersPrefix = "orders";
      const errors = {
        codeRequired: new Error("orderCode is required"),
        alreadyExists: new Error("object already exists"),
      }

      module.exports = { 
        main: async function (event, _) {
          const storage = getStorage();

          if (!event.data || !Object.keys(event.data).length) {
            return await onList(storage, event);
          }

          const { orderCode, consignmentCode, consignmentStatus } = event.data;
          if (orderCode && consignmentCode && consignmentStatus) {
            return await onCreate(storage, event);
          }

          event.extensions.response.status(500);
        }
      }

      async function onList(storage, event) {
        try {
          return await storage.getAll();
        } catch(err) {
          event.extensions.response.status(500);
          return;
        }
      }

      async function onCreate(storage, event) {
        try {
          await storage.set(event.data);
        } catch(err) {
          let status = 500;
          switch (err) {
            case errors.codeRequired: {
              status = 400;
              break;
            };
            case errors.alreadyExists: {
              status = 409;
              break;
            };
          }
          event.extensions.response.status(status);
        }
      }

      class RedisStorage {
        storage = undefined;
        asyncGet = void 0;
        asyncKeys = void 0;
        asyncSet = void 0;

        constructor(options) {
          this.storage = redis.createClient(options);
          this.asyncGet = promisify(this.storage.get).bind(this.storage);
          this.asyncKeys = promisify(this.storage.keys).bind(this.storage);
          this.asyncSet = promisify(this.storage.set).bind(this.storage);
        }

        async getAll() {
          let values = [];

          const keys = await this.asyncKeys(`${ordersPrefix}:*`);
          for (const key of keys) {
            const value = await this.asyncGet(key);
            values.push(JSON.parse(value));
          }

          return values;
        }

        async set(order = {}) {
          if (!order.orderCode) {
            throw errors.codeRequired;
          }
          const value = await this.asyncGet(`${ordersPrefix}:${order.orderCode}`);
          if (value) {
            throw errors.alreadyExists;
          }
          await this.asyncSet(`${ordersPrefix}:${order.orderCode}`, JSON.stringify(order));
        }
      }

      class InMemoryStorage {
        storage = new Map();

        getAll() {
          return Array.from(this.storage)
            .filter(([orderCode, _]) => orderCode.startsWith(`${ordersPrefix}:`))
            .map(([_, order]) => order)
        }

        set(order = {}) {
          if (!order.orderCode) {
            throw errors.codeRequired;
          }
          if (this.storage.get(`${ordersPrefix}:${order.orderCode}`)) {
            throw errors.alreadyExists;
          }
          return this.storage.set(`${ordersPrefix}:${order.orderCode}`, order);
        }
      }

      function readEnv(env = "") {
        return process.env[env] || undefined;
      }

      function createStorage() {
        let redisPrefix = readEnv("APP_REDIS_PREFIX");
        if (!redisPrefix) {
          redisPrefix = "REDIS_";
        }
        const port = readEnv(`${redisPrefix}PORT`);
        const host = readEnv(`${redisPrefix}HOST`);
        const password = readEnv(`${redisPrefix}REDIS_PASSWORD`);

        if (host && port && password) {
          return new RedisStorage({ host, port, password });
        }
        return new InMemoryStorage();
      }

      function getStorage() {
        if (!storage) {
          storage = createStorage();
        }
        return storage;
      }
   EOF
   ```

2. Check if your Function was created successfully and all conditions are set to `True`:

    ```bash
    kubectl get functions orders-function -n orders-service
    ```

    You should get a result similar to the following example:

    ```bash
    NAME                CONFIGURED   BUILT   RUNNING   VERSION   AGE
    orders-function     True         True    True      1         18m
    ```

    </details>
    <details>
    <summary label="console-ui">
    Console UI
    </summary>

1. If you aren't in the view of Namespace `orders-service` in the Kyma Console, select a `orders-service` Namespace from the drop-down list in the top navigation panel.

2. Go to the **Functions** view (under **Development** section) in the left navigation panel and select **Create Function**.

// dodac jeszcze labelke `app: orders-function` i `example: orders-function`
3. In the pop-up box, provide the `orders-function` name and select **Create** to confirm changes.

   The pop-up box closes and the message appears on the screen after a while, confirming that the Function was created successfully.

4. In the Function details view that opens up automatically, enter the Function's code in the **Source** tab:

    ```js
    const redis = require("redis");
    const { promisify } = require("util");

    let storage = undefined;
    const ordersPrefix = "orders";
    const errors = {
      codeRequired: new Error("orderCode is required"),
      alreadyExists: new Error("object already exists"),
    }

    module.exports = { 
      main: async function (event, _) {
        const storage = getStorage();

        if (!event.data || !Object.keys(event.data).length) {
          return await onList(storage, event);
        }

        const { orderCode, consignmentCode, consignmentStatus } = event.data;
        if (orderCode && consignmentCode && consignmentStatus) {
          return await onCreate(storage, event);
        }

        event.extensions.response.status(500);
      }
    }

    async function onList(storage, event) {
      try {
        return await storage.getAll();
      } catch(err) {
        event.extensions.response.status(500);
        return;
      }
    }

    async function onCreate(storage, event) {
      try {
        await storage.set(event.data);
      } catch(err) {
        let status = 500;
        switch (err) {
          case errors.codeRequired: {
            status = 400;
            break;
          };
          case errors.alreadyExists: {
            status = 409;
            break;
          };
        }
        event.extensions.response.status(status);
      }
    }

    class RedisStorage {
      storage = undefined;
      asyncGet = void 0;
      asyncKeys = void 0;
      asyncSet = void 0;

      constructor(options) {
        this.storage = redis.createClient(options);
        this.asyncGet = promisify(this.storage.get).bind(this.storage);
        this.asyncKeys = promisify(this.storage.keys).bind(this.storage);
        this.asyncSet = promisify(this.storage.set).bind(this.storage);
      }

      async getAll() {
        let values = [];

        const keys = await this.asyncKeys(`${ordersPrefix}:*`);
        for (const key of keys) {
          const value = await this.asyncGet(key);
          values.push(JSON.parse(value));
        }

        return values;
      }

      async set(order = {}) {
        if (!order.orderCode) {
          throw errors.codeRequired;
        }
        const value = await this.asyncGet(`${ordersPrefix}:${order.orderCode}`);
        if (value) {
          throw errors.alreadyExists;
        }
        await this.asyncSet(`${ordersPrefix}:${order.orderCode}`, JSON.stringify(order));
      }
    }

    class InMemoryStorage {
      storage = new Map();

      getAll() {
        return Array.from(this.storage)
          .filter(([orderCode, _]) => orderCode.startsWith(`${ordersPrefix}:`))
          .map(([_, order]) => order)
      }

      set(order = {}) {
        if (!order.orderCode) {
          throw errors.codeRequired;
        }
        if (this.storage.get(`${ordersPrefix}:${order.orderCode}`)) {
          throw errors.alreadyExists;
        }
        return this.storage.set(`${ordersPrefix}:${order.orderCode}`, order);
      }
    }

    function readEnv(env = "") {
      return process.env[env] || undefined;
    }

    function createStorage() {
      let redisPrefix = readEnv("APP_REDIS_PREFIX");
      if (!redisPrefix) {
        redisPrefix = "REDIS_";
      }
      const port = readEnv(`${redisPrefix}PORT`);
      const host = readEnv(`${redisPrefix}HOST`);
      const password = readEnv(`${redisPrefix}REDIS_PASSWORD`);

      if (host && port && password) {
        return new RedisStorage({ host, port, password });
      }
      return new InMemoryStorage();
    }

    function getStorage() {
      if (!storage) {
        storage = createStorage();
      }
      return storage;
    }
    ```

5. In the **Dependencies** tab enter:

   ```js
   { 
     "name": "orders-function",
     "version": "1.0.0",
     "dependencies": {
       "redis": "3.0.2"
     }
   }
   ```

6. Select **Save** to confirm changes.

   You will see the message confirming the changes were saved. Once deployed, the new Function should have the `RUNNING` status.

  </details>
</div>
