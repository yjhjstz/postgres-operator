---
title: "pgo_scale"
---
## pgo scale

Scale a PostgreSQL cluster

### Synopsis

The scale command allows you to adjust a Cluster's replica configuration. For example:

	pgo scale mycluster --replica-count=1

```
pgo scale [flags]
```

### Options

```
      --ccp-image-tag string      The CCPImageTag to use for cluster creation. If specified, overrides the .pgo.yaml setting.
  -h, --help                      help for scale
      --no-prompt                 No command line confirmation.
      --node-label string         The node label (key) to use in placing the primary database. If not set, any node is used.
      --replica-count int         The replica count to apply to the clusters. (default 1)
      --resources-config string   The name of a container resource configuration in pgo.yaml that holds CPU and memory requests and limits.
      --service-type string       The service type to use in the replica Service. If not set, the default in pgo.yaml will be used.
      --storage-config string     The name of a Storage config in pgo.yaml to use for the replica storage.
```

### Options inherited from parent commands

```
      --apiserver-url string     The URL for the PostgreSQL Operator apiserver.
      --debug                    Enable debugging when true.
  -n, --namespace string         The namespace to use for pgo requests.
      --pgo-ca-cert string       The CA Certificate file path for authenticating to the PostgreSQL Operator apiserver.
      --pgo-client-cert string   The Client Certificate file path for authenticating to the PostgreSQL Operator apiserver.
      --pgo-client-key string    The Client Key file path for authenticating to the PostgreSQL Operator apiserver.
```

### SEE ALSO

* [pgo](/operatorcli/cli/pgo/)	 - The pgo command line interface.

###### Auto generated by spf13/cobra on 31-May-2019