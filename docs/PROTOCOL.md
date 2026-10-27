# Node protocol v1

All endpoints are under `/api/v1/node`. JSON timestamps are Unix seconds, latency is milliseconds, IDs are opaque strings. Transport is HTTPS. Local development may explicitly enable HTTP. Never place credentials in URLs or logs.

## Registration and polling

1. An administrator creates a node in the dashboard. The one-time registration token expires after 600 seconds.
2. `POST /register` with `{ "node_id": "…", "token": "…" }` exchanges it for a persistent node identity: `{ "node_id": "…", "token": "…" }`. Save this identity privately. Reusing a registration token fails.
3. Every two seconds, `POST /poll` sends a `Poll` body with `protocol: 1`, implementation version, platform, capabilities, local Unix time, queue size and dropped-buffer count. Include `X-Node-ID` and `Authorization: Bearer …`.
4. A `Work` response contains `assignments`, `jobs`, `cancelled` IDs and `server_time`. An assignment expires after 15 minutes; an offline node may execute only its cached unexpired assignments. A clock skew above 30 seconds suspends new assignments.
5. `POST /results` sends an array of at most 100 `Result` objects with the same identity headers. The response has `accepted` and `rejected` arrays. Remove both sets from the local spool. On transport or 5xx failure retain unacknowledged results for retry.

Revocation invalidates the token immediately at the server. A connected node receiving 401 clears its cached work and cancels running tasks. A disconnected node can only learn about revocation when it reconnects; its existing lease bounds continued execution.

## Scheduling and identity

A continuous run is aligned to `scheduled_at = floor(now / interval) * interval`. Its ID is lowercase SHA-256 hex of `monitor_id:node_id:version:scheduled_at`. A job uses its job ID as result ID. `assignment_id` must match the issued lease, monitor, node and version. Every result is idempotent by ID. A diagnostic never updates the continuous monitoring state.

The node starts no overlapping runs for a single assignment. It does not create fictional results for missed cycles. Local queues are limited to 100 MiB and 24 hours; dropped results increase a counter. Files use atomic rename and private permissions. Upload newest results first to recover fresh visibility promptly.

## Result semantics

| Status | Meaning | Included in availability denominator |
| --- | --- | --- |
| `success` | Target satisfied the configured assertion | Yes |
| `failure` | Target responded incorrectly or failed the check | Yes |
| `error` | Executor/policy/permission/configuration failure | No |
| `cancelled` | Explicitly cancelled run | No |

`stage` identifies the failed phase; `data` contains protocol-specific evidence. HTTP includes timing/response/TLS facts, ICMP packet statistics, DNS records/resolver, and route tools bounded output. No response body or request authorization header is returned. `attempts` and `first_error` retain retry evidence. A retry shares the whole check deadline.

The dashboard timestamps receipt itself. A result older than `2 * interval + timeout` is marked late, contributes to history, and cannot overwrite current state. Results from old configuration versions also cannot make a current version healthy. Missing results never imply target failure or recovery.

## Monitor and incident state

Each node has independent consecutive failure/recovery counters (defaults: 3 and 2). Gaps and executor errors reset consecutive counters. A monitor is healthy only when every assigned node has a fresh healthy result; all failed means down, a failed subset means partial, and incomplete evidence is pending/unknown. Node heartbeat loss after 60 seconds produces a separate node incident and unknown coverage.

One incident remains open per monitor until all assigned nodes confirm recovery. A failed-node quorum is configurable. Acknowledgment, bounded silence and notes append timeline events. Maintenance continues measurement while suppressing notification; confirmation counters restart when maintenance ends. Recovery is not inferred from pause, removal or node disconnection.

History exposes separate availability (`success / (success + failure)`) and coverage (`valid samples / expected enabled cycles`) with nullable values when no denominator exists. Expected cycles derive from configuration history, not successful polls. Maintenance counts toward coverage; monitor pause does not. Raw results are paginated. Five-minute aggregates last 90 days and hourly aggregates last 365 days; ranges are rounded to bucket boundaries.

## Access policy

Checks are structured data, never commands. HTTP permits GET/HEAD only, TLS verification defaults on, no proxy environment is inherited, redirects are limited and cross-origin custom headers are stripped. Resolve and validate all candidate IPs before dialing. Private, loopback, link-local, multicast and special networks are blocked unless the administrator adds explicit allowed CIDRs. CIDRs are a deliberate trust grant, not a DNS name allowlist. DNS resolvers on private networks need their own CIDR permission. Nodes should run unprivileged, with only the capabilities needed for ICMP/MTR.

## Optional certificate assertion

`certificate_min_days` defaults to 0 (disabled), or accepts 1–365. Supporting nodes advertise `certificate_expiry`. The panel does not assign this check to a node lacking that capability. An insufficient remaining lifetime returns a target assertion failure at `certificate_expiry` while preserving HTTP status and certificate evidence. Explicit TLS bypass remains visible in the result.

The local node credential file also binds its dashboard URL; changing that URL requires explicit registration rather than forwarding an old token. System DNS resolvers are trusted host configuration; custom private resolver addresses require a CIDR grant.
