# Knative MQTT event source for Kyma

Knative event source for SAP S/4 applications emitting MQTT events. Meant to be used with [Kyma](https://github.com/kyma-project/kyma).

## Running the controller inside a cluster

```console
$ ko apply -f config/
```

## Running the controller locally

### Environment setup (performed once)

Create the CustomResourceDefinition for the MQTT event source type, `MQTTSource`.

```console
$ kubectl create -f config/300-mqttsource-crd.yaml
```

Create the `mqtt-source` system namespace. The controller sources its ConfigMaps from it.

```console
$ kubectl create -f config/100-namespace.yaml
```

Create ConfigMaps for logging and observability inside the system namespace (`mqtt-source`)

```console
$ kubectl create -f config/400-config-logging.yaml
$ kubectl create -f config/400-config-observability.yaml
```

### Controller startup

Export the following mandatory environment variables:

* `KUBECONFIG`: path to a local kubeconfig file (if different from the default OS location)
* `SYSTEM_NAMESPACE`: set to "mqtt-source" (see above)
* `METRICS_DOMAIN`: domain of the exposed Prometheus metrics. Can be set to an arbitrary value (e.g. "kyma/mqtt-event-source")
* `ADAPTER_IMAGE`: container image of the MQTT server adapter

Build the binary

```console
$ make
```

Run the controller

```console
$ ./mqttsource-controller
```
