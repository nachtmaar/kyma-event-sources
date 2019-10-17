# Knative MQTT event source for Kyma

Knative event source for SAP S/4 applications emitting MQTT events. Meant to be used with [Kyma](https://github.com/kyma-project/kyma).

## Running locally

Create the CustomResourceDefinition for the MQTT event source type, `MQTTSource`

```console
$ kubectl create -f config/mqttsource-crd.yaml
```

Export the following mandatory environment variables:

* `KUBECONFIG`: path to a local kubeconfig file (if different from the default OS location)
* `SYSTEM_NAMESPACE`: name the namespace the controller will source its ConfigMaps from (e.g. "default")
* `METRICS_DOMAIN`: domain of the exposed Prometheus metrics. Can be set to an arbitrary value (e.g. "kyma/mqtt-event-source")

Create ConfigMaps for logging and observability inside the configured system namespace

```console
$ kubectl -n ${SYSTEM_NAMESPACE} create -f config/config-logging-dev.yaml
$ kubectl -n ${SYSTEM_NAMESPACE} create -f config/config-observability-dev.yaml
```

Build the binary

```console
$ make
```

Run the controller

```console
$ ./mqttsource-controller
```
