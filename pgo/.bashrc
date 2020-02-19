# User specific aliases and functions
export GOPATH=$HOME/odev
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
# NAMESPACE is the list of namespaces the Operator will watch
export NAMESPACE=pgouser1,pgouser2,pgo
export PGO_INSTALLATION_NAME=dev
# PGO_OPERATOR_NAMESPACE is the namespace the Operator is deployed into
export PGO_OPERATOR_NAMESPACE=pgo
export PGO_NAMESPACE=pgo

# PGO_CMD values are either kubectl or oc, use oc if Openshift
export PGO_CMD=kubectl

# the directory location of the Operator scripts
export PGOROOT=$GOPATH/src/github.com/crunchydata/postgres-operator

# the version of the Operator you run is set by these vars
export PGO_IMAGE_PREFIX=hub.didiyun.com/postgres
#export PGO_IMAGE_PREFIX=crunchydata
export PGO_BASEOS=centos7
export PGO_VERSION=4.1.1
export PGO_IMAGE_TAG=$PGO_BASEOS-$PGO_VERSION

# for the pgo CLI to authenticate with using TLS
export PGO_CA_CERT=/tmp/client.crt
export PGO_CLIENT_CERT=/tmp/client.crt
export PGO_CLIENT_KEY=/tmp/client.key
export PGO_NFS_IP=10.254.0.x
export PGO_APISERVER_URL=http://10.254.0.x:8443
export DISABLE_TLS=true
# common bash functions for working with the Operator
setip() 
{ 
	export PGO_APISERVER_URL=https://`$PGO_CMD -n "$PGO_OPERATOR_NAMESPACE" get service postgres-operator -o=jsonpath="{.spec.clusterIP}"`:8443 
}

alog() {
    $PGO_CMD  -n "$PGO_OPERATOR_NAMESPACE" logs `$PGO_CMD  -n "$PGO_OPERATOR_NAMESPACE" get pod --selector=name=postgres-operator -o jsonpath="{.items[0].metadata.name}"` -c apiserver
}

olog () {
    $PGO_CMD  -n "$PGO_OPERATOR_NAMESPACE" logs `$PGO_CMD  -n "$PGO_OPERATOR_NAMESPACE" get pod --selector=name=postgres-operator -o jsonpath="{.items[0].metadata.name}"` -c operator
}

slog () {
    $PGO_CMD  -n "$PGO_OPERATOR_NAMESPACE" logs `$PGO_CMD  -n "$PGO_OPERATOR_NAMESPACE" get pod --selector=name=postgres-operator -o jsonpath="{.items[0].metadata.name}"` -c scheduler
}
