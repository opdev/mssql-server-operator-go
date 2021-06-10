# MSSQL Server Operator
This operator is currently in development and using [existing helm chart](https://github.com/microsoft/mssql-docker/tree/master/linux/sample-helm-chart-statefulset-deployment) from [MSSQL Docker](https://github.com/microsoft/mssql-docker) repository. This implementation using the StatefulSet helm chart sample which is where the link will lead you.  

## Installation

In order to get this working currently you have to install the configmap, secret, service account, and role binding in the config/samples directory. Once installed the operator can run which will create a statefulset that will deploy successfully with the previously mentioned manifests deployed. Namespace/project mssql was used.

### OpenShift

### Kubernetes + OLM

## 

### Helpful Documentation Links
[mssql.conf](https://docs.microsoft.com/en-us/sql/linux/sql-server-linux-configure-mssql-conf?view=sql-server-ver15#mssql-conf-format)  
[environment variables](https://docs.microsoft.com/en-us/sql/linux/sql-server-linux-configure-environment-variables?view=sql-server-ver15)
