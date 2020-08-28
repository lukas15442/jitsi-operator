# OpenShift Jitsi Operator

## Introduction

This OpenShift Operator is developed on operator-sdk 0.19.1 with the old layout.
What this operator is doing for you:

- initializes whole jitsi infrastructure **once at operator start**, this includes:
  - prosody
  - jicofo
  - jitsi-web
  - multiple jvb's
- does NOT watch status or deployment updates for:
  - prosody
  - jicofo
  - jitsi-web
- automatically deploys and deletes *n* jvb's at live system running. *n* is specified in the custom resource over *size* in the *spec* part
- creates and deletes for each jvb deployment:
  - deployment (not deployment config)
  - three services (standard, http, tcp)

This operator was first developed to manage only multiple jvb's. Because one jvb deployment need multiple services you can not use the standard scaling from OpenShift. After that the operator was changed to initialize a whole jitsi infrastructure, but only once at start of the operator. This could be changed in future work (see *Future work*).

## Usage

### First time usage

If you are using this operator for first time, you have to add the ClusterRole and the CustomResourceDefinition in your System. For this you need to have privileges on the OpenShift system. Do not add the ClusterRole or CustomResourceDefinition in a specific namespace

```bash
# Add ClusterRole
oc create -f ClusterRole.yml

# Add CustomResourceDefinition
oc create -f deploy/crds/jitsi.fbi.h-da.de_jitsis_crd.yaml
```

### Create a project

```bash
# Create a project
oc new-project jitsi-operator

# Ensure that created project is active
oc project jitsi-operator

# Add operator essentials
oc create -f deploy/service_account.yaml
oc create -f deploy/role.yaml
oc create -f deploy/role_binding.yaml

# Add CustomResource (optional: change parameter in template file, e.g. size of jvb's)
oc create -f deploy/crds/jitsi.fbi.h-da.de_v1alpha1_jitsi_cr.yaml

# Deploy the operator (to change version, see template file)
oc create -f deploy/operator.yaml
```

## Changelog

**0.1_beta**

- *add* initializes whole jitsi infrastructure at operator startup

## Future work

- change layout to new operator-sdk version
- image versions should be in custom resource definition
- operator should watch deployment changes and handle these
