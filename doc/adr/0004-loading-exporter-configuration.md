# 4. Loading Exporter Configuration

Date: 2023-06-24

## Status

Accepted

## Context

The application requires configuration and a way to load it that is cloud native.

For simplicity of development, a configuration file must be easy to maintain
locally.  Additionally, the PhoenixNAP Ansible and Terraform libraries already
support utilizing a yaml formatted configuration file that contains most of the
values I need to initialze a BMC API connection.  It is commonly stored at
`~/.pnap/config.yaml` and format looks like this:

```yaml
clientId: your-secret-client-id-here
clientSecret: your-secret-client-secret-here
```

The only value missing in this file is the Token URL and Scopes.  Since the scopes
are unlikely to change on a regular basis, I am not choosing to expose them at this
time.  That leaves the TokenUrl.  To capture this, I am adding a `tokenUrl` value
to the existing file.  It will look like this:

```yaml
clientId: your-secret-client-id-here
clientSecret: your-secret-client-secret-here
tokenUrl: https://auth.phoenixnap.com/auth/realms/BMC/protocol/openid-connect/token
```

Additional configuration for the exporter can be added later, while still being
compatible with the Ansible and Terraform providers.  An example might be:

```yaml
clientId: your-secret-client-id-here
clientSecret: your-secret-client-secret-here
tokenUrl: https://auth.phoenixnap.com/auth/realms/BMC/protocol/openid-connect/token
go:
  collector:
    enabled: true
bmc:
  reservations_under_contract:
    - productCode: d3.m6.xlarge
      productCategory: server
      reservationModel: ONE_MONTH_RESERVATION
      quantity: 9
    - productCode: s2.c1.small
      productCategory: server
      reservationModel: ONE_MONTH_RESERVATION
      quantity: 19
    - productCode: s2.c1.large
      productCategory: server
      reservationModel: ONE_MONTH_RESERVATION
      quantity: 5
```

## Decision

- Utilize `~/.pnap/config.yaml` to store configuration of the exporter.
- Expose via environment variable any sensitive values that should not be stored in source control.
    - These values can be injected by Kubernetes later using Secrets.

## Consequences

- Need to design tests that confirm the Ansible and Terraform modules are not impacted by the presence of additional configurationdata.
- Need to build a [config module](../../config.go) that loads the configuration.