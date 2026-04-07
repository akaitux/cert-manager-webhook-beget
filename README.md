<p align="center">
  <img src="https://raw.githubusercontent.com/cert-manager/cert-manager/d53c0b9270f8cd90d908460d69502694e1838f5f/logo/logo-small.png" height="32" width="32" alt="cert-manager project logo" />
</p>

# [Beget](https://beget.com/p259374) DNS01 webhook 

## Status

The module is active, but the underlying API is rarely changing, not much to update yet. Give it a star, if you're using it.

## Installation

- Read 
    - https://cert-manager.io/docs/configuration/acme/dns01/
    - https://cert-manager.io/docs/configuration/acme/

- install [cert-manager](https://github.com/cert-manager/cert-manager)
    - `kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.1/cert-manager.yaml`
- instal the issuer:
    - **NOTE**: The kubernetes resources used to install the Webhook should be deployed within the same namespace as the cert-manager ("cert-manager" by default, check ./deploy/values.yaml).
    - pull this repo
    - `helm install webhook-beget ./deploy/beget -f ./deploy/beget/values.yaml -n cert-manager`
    - default image for release `1.0.1`: `docker.io/akaitux/cert-manager-webhook-beget:1.0.1`
    - if needed, override explicitly: `--set image.repository=docker.io/akaitux/cert-manager-webhook-beget --set image.tag=1.0.1`
- create a secret for beget API
- create an issuer
- request certificates
- add the certificates to services

For webhook secret references, prefer this config shape:

```yaml
config:
  apiLoginSecretRef:
    secretName: beget-credentials
    key: login
  apiPasswdSecretRef:
    secretName: beget-credentials
    key: passwd
```

The webhook accepts both `name` and `secretName` for backward compatibility, but `secretName` is the recommended form for new manifests.

Follow ***an example*** for details: [testdata/resources](testdata/resources/README.md).

## Release 1.0.1

- Image: `docker.io/akaitux/cert-manager-webhook-beget:1.0.1`
- Helm chart version: `1.0.1`
- App version: `1.0.1`

## Tests

You can run the webhook test suite with:

```bash
$ TEST_ZONE_NAME=example.com. make test
```
