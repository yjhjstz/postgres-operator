---
title: "pgo_delete_pgorole"
---
## pgo delete pgorole

Delete a pgorole

### Synopsis

Delete a pgorole. For example:
    
    pgo delete pgorole somerole

```
pgo delete pgorole [flags]
```

### Options

```
      --all         all resources.
  -h, --help        help for pgorole
      --no-prompt   No command line confirmation.
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

* [pgo delete](/operatorcli/cli/pgo_delete/)	 - Delete an Operator resource

###### Auto generated by spf13/cobra on 4-Oct-2019
