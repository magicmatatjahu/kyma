---
title: How to start
type: Getting Started
---

// do kazdego etapu guida powinny byc dodane linki do komponentow, ktore są uzywane w danym etapie, np connect-external-app to bedzie app-connector, service-catalog i helm-broker

This set of Getting Started guides aims to demonstrate the basic functionalities offered by Kyma its:
- **Integration and connectivity** feature brought in by [Application Connector](https://kyma-project.io/docs/components/application-connector/). With Kyma you can connect external applications and expose their API and events in Kyma.
- **Extensibility** feature provided by [Service Catalog](https://kyma-project.io/docs/components/service-catalog/). You can use its built-in portfolio of external services in your applications.
- **Application runtime** feature supported by [Serverless](https://kyma-project.io/docs/components/serverless/) where you can build microservices and functions to interact with external services and applications to perform a given business logic.

Having these features in mind, we will:

// nowa kolejność:
- wspomnieć o podpinaniu aplikacji z mockiem w dalszym etapie tutka
- tworzymy namespace
- deplojujemy microserwis + api rula + test orders service
- redis + test orders service
- mock
- trigger pod orders-service
- identyko z funkcją

1. Connect a mock application as an addon to Kyma and use one of its events sent whenever an order is created in an application by a user.

2. Provision the Redis service available through the Service Catalog to act as an external database in which the order information can be stored.

3. Create a microservice and a function to combine the previous two pieces of the puzzle. You can use both of them in Kyma to perform a certain business logic, integrate external applications and services and create meaningful flows. In our guides, you will see both of them used in the same flow - thanks to their logic, both the microservice and the function react to the event sent from the mock application and store the order data in the attached Redis database.

> **CAUTION:** These tutorials will refer to a sample `orders-service` application deployed on Kyma as a `microservice` to easily distinguish it from the external Commerce mock that represents an external monolithic application connected to Kyma. In Kyma docs, these are referred to as `application` and `Application` respectively.

All guides, whenever possible, will demonstrate the steps to perform both from kubectl and Console UI.

## Prerequisites

These guides show what you can do with Kyma running on a cluster of your choice. Before you start, follow the steps in the [installation tutorials](/root/kyma/#installation-install-kyma-on-a-cluster]) to get your Kyma cluster up and running.

// ode mnie
// trzeba dodać info o curlu
// dodać info, ze kubeconfig (z uprawnieniami do tworzenia resourcow) jest potrzebny i/lub tutorial o pobieraniu kubeconfiga + update tutoriala
If you choose to complete the guide in the terminal, you also need to have [kubectl](https://kubernetes.io/docs/reference/kubectl/kubectl/) 1.16 or greater installed.

## Main actors

Let's introduce the main actors that will lead us through the guides. These are our own examples, mock applications, and experimental services:

- [Orders service](https://github.com/kyma-project/examples/tree/master/orders-service) is a sample application (microservice) written in Go. It can expose HTTP endpoints used to create, read and delete basic order JSON entities. The service can run with either an in-memory database that is enabled by default or an external, Redis database. On the basis of it, we will show how you can use Kyma to deploy your own microservice, expose it on HTTP endpoints to make it available for external services, and bind it to an actual database service (Redis).

// ode mnie - zmienione
// funckja to jest to samo co orders service, ale bez mozliwosci usuwania - read i create
- [Function](https://github.com/kyma-project/examples/blob/order-service/orders-service/deployment/function.yaml) is a [Serverless](https://kyma-project.io/docs/components/serverless/) Function equivalent of the Orders Service with the ability to retrieve orders records and save. Like the microservice, Function can run with either an in-memory database or an Redis instance.

- [Redis addon](https://github.com/kyma-project/addons/tree/master/addons/redis-0.0.3) is basically a bundle of two Redis service available in two plans: `micro` and `enterprise`. We will connect it to Kyma thanks to Helm Broker and expose it in the Kyma cluster as an addon under Service Catalog. The Redis service represents an open source, in-memory data structure store, used as a database, cache and message broker. For the purpose of these guides, we will use the `micro` plan with the in-memory storage to demonstrate how it can replace the default storage of our microservice.

- [Commerce mock](https://github.com/SAP-samples/xf-addons/tree/master/addons/commerce-mock-0.1.0) is to act as a sample external and monolithic application which we want to extend with Kyma. It is based on the [Varkes](https://github.com/kyma-incubator/varkes) project and is also available in the form of an addon. It will simulate how you can pair an external application with Kyma and expose its APIs and Events. In our guides, we will use its **order.deliverysent.v1** event.

## Steps

The guides cover these steps:

1. Setup a namespace.

2. Deploy the `orders-service` microservice on the cluster.

3. Expose the microservice through the APIRule CR on HTTP endpoints. This way it will be available for other services outside the cluster and test it.

4. Add an addon with the Redis service and create a ServiceInstance CR for it so you can bind it later with your microservice and Function. // tworzymy addona a potem instancje

5. Bind a microservice to the Redis instance by creating ServiceBinding and ServiceBindingUsage CRs and test it. // zeby microserwis uzywal redisa. Na samym koncu napisac o mozliwosci podpinania aplikacji do konsumowania eventow przez microserwis

6. Connect the external SAP Commerce Cloud - Mock application.

7. Trigger your microservice to react to the **order.deliverysent.v1** event from the mock application. Send the event and see if the microservice reacts to it by saving its details in the Redis database.

8. Create a Function on the cluster.

9. Expose the Function through the APIRule CR on HTTP endpoints. This way it will be available for other services outside the cluster and test it.

// nie musimy tworzyc juz redis i addon wiec pomijamy kroki
10. Bind a Function to the Redis service by creating ServiceBinding and ServiceBindingUsage CRs and test it.

// aplikacja jest stworzona wczesniej
11. Trigger your Function to react to the **order.deliverysent.v1** event from the mock application. Send the event and see if the Function reacts to it by saving its details in the Redis database.

As a result, you get two scenarios of the same flow - a microservice and a Function that are triggered by new order events from the Commerce mock and send order data to the Redis database:

![Order flow](./assets/order-flow.svg)
