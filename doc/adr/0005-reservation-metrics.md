# 5. Reservation Metrics

Date: 2023-06-25

## Status

Accepted

## Context

Visibility on reservations is required to understand what the monthly recurring
commitment of an organization is to PhoenixNAP.

## Decision

The number and total cost metrics for reservations need to be tracked.

The following labels will be applied to these metrics:

- ProductCode
- LocationCode (PHX, ASH, SGP, NLD, CHI, SEA, AUS, )
- Status (active, renewing, expiring, expired)
- AutoRenewal (enabled, disabled)
- Assignment (assigned, unassigned)

## Consequences

Each of these metrics must be captured by `productCode` and `location` as
reservations are not transitory between locations and products.

The `productCode` dimension has 52 possible values at time of writing.

```bash
wget --quiet \
  --method GET \
  --header 'Accept: */*' \
  --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
  --output-document \
  - 'https://api.phoenixnap.com/billing/v1//products?productCategory=SERVER' | jq '. | length'
```
