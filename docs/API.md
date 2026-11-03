# Dashboard HTTP API v1

Base: `/api/v1`. All requests and responses are JSON. Errors are `{ "error": "human-readable message" }`; 401 means expired/invalid identity, 403 insufficient role, 409 conflicting configuration version, and 429 login throttling. Session cookies are HttpOnly/SameSite Strict; native clients use bearer tokens. Browser writes require the same origin or `PROBELUNE_PUBLIC_URL`.

## Authentication

| Method / path | Input and output |
| --- | --- |
| GET `/auth/status` | `{initialized: boolean}` |
| POST `/auth/setup` | Header `X-Setup-Token` from server private file; `{username,password}` creates the first admin only |
| POST `/auth/login` | `{username,password}` → `{token,expires,user}`; web/desktop session 8 hours, Probelune Mobile session 30 days |
| GET `/auth/me` | Current sanitized user |
| POST `/auth/logout` | Revoke current session; 204 |
| GET `/auth/sessions` | Own sessions with ID, device description and expiry |
| DELETE `/auth/sessions/:id` | Revoke an owned session |
| POST `/auth/pair` | Authenticated; `{code,expires}`, one-time, 120 seconds |
| POST `/auth/pair/exchange` | `{code}` → 30-day native session, code is consumed |

Pair QR JSON: `{ "type":"probelune-pair", "version":1, "url":"https://monitor.example.com", "code":"…" }`. Credentials stay in the QR only until expiry and must never be included in screenshots or public documentation.

## Observation and control

| Method / path | Role | Behavior |
| --- | --- | --- |
| GET `/overview` | any member | Monitors, nodes, active incidents, counts, server time |
| GET `/monitors` | any | Monitor configuration with sanitized headers plus live per-node states |
| POST `/monitors` | admin/operator | Create Monitor; server assigns ID/version/updated_at |
| PUT `/monitors/:id` | admin/operator | Full Monitor including current `version`; omit headers to preserve, `{}` clears them |
| GET `/monitors/:id/history` | any | `hours=1..8760`, `limit=1..1000`, `offset`, optional `node_id`; returns results, buckets, summary and coverage basis |
| GET `/nodes` | any | Sanitized node identity, labels, capabilities, heartbeat and policy |
| POST `/nodes` | admin | `{name,region,network,policy}` → node plus one-time registration token |
| PUT `/nodes/:id` | admin | Edit labels, enable/pause and private CIDR policy |
| POST `/nodes/:id/credentials` | admin | `{revoke:true}` revokes; otherwise rotate registration identity |
| GET `/diagnostics` | any | Recent jobs with optional result |
| POST `/diagnostics` | admin/operator | `{check,node_ids}` or `{monitor_id}`; returns jobs array |
| POST `/diagnostics/:id/cancel` | admin/operator | Cancel job; running node receives cancellation on poll |
| GET `/incidents` | any | Up to 500 newest incidents and timelines |
| POST `/incidents/:id/actions` | admin/operator | `{action:"ack"}`, `{action:"silence",minutes:1..1440}` or `{action:"note",note:"…"}` |

Monitor fields are specified in `models.go`. Continuous types: HTTP, TCP, ICMP (`ping`) and DNS. MTR/traceroute are diagnostic-only. Use `enabled:false` to pause; `archived:true` to archive without destroying history. Public display requires both instance `public_status` and monitor `public:true`.

## Administration

- GET `/settings`: any member; PUT: admin. Fields: `name`, `public_url`, IANA `timezone`, `raw_days` (1–90), `event_days` (7–365), `public_status`.
- GET `/users`, POST `/users`, PUT `/users/:id`: admin. Fields: `username`, `password` (12–72 bytes, blank preserves on edit), `role` (`admin`, `operator`, `viewer`), `disabled`. Password changes revoke other sessions; disabled users lose access immediately.
- GET/POST `/channels`, PUT `/channels/:id`, POST `/channels/:id/test`: admin. Types `webhook`, `smtp`, `push`. GET omits secrets; save returns webhook signing secret. Empty password/secret on update preserves it. Each channel has `enabled`, `name` and optional egress `policy`.
- GET `/deliveries`: admin/operator; delivery states include pending, sent, failed, suppressed and skipped. Three total attempts, bounded backoff. Delivery is at-least-once; use event IDs for receiver deduplication.
- GET `/maintenance`: any; POST/DELETE: admin/operator. Fields: `name`, `monitor_ids` (empty means all), `start`, `end`, max 30 days. Deleting an active window ends it and retains its historical boundary.
- GET `/audit`: admin. Records include `at`, `user`, `action`, `resource`.
- POST `/devices`: any member, own device `{token: "Firebase Installation ID", platform:"android"|"ios"}`. GET lists own sanitized registrations; DELETE `/devices/:id` removes an owned registration. `token` is the v1 wire field for FID; the provider request uses `message.fid`.
- GET `/status`: anonymous, only when explicitly enabled; returns public monitor names/status/coverage, never targets, headers, node addresses or credentials.
- GET `/healthz` (outside API prefix): process health/version, no credentials.

Webhook receivers verify `X-Probelune-Signature` according to the dashboard notification implementation and reject stale timestamps. SMTP requires TLS or STARTTLS. FCM uses the deployer's Firebase service account; native clients must belong to the same Firebase project. No shared hosted push service is provided.

Login/pair responses and overview include a stable `instance_id` for matching native notifications across LAN/public URL aliases. It is not a credential. Notification channels accept `reminder_minutes` (0 disables; 1–1440); maintenance/silence and confirmed recovery suppress further reminders.

## Configurable availability history and UDP

Protocol v1 adds optional `udp_payload`, `udp_encoding` (`text` or `hex`) and `udp_expected` fields. UDP requires a valid host, port and a 1–1400-byte decoded request. Success requires a non-empty reply; a configured expected byte sequence must occur in the response. A timeout means no response was observed, not proof that the port is closed. Nodes advertise `udp`; upgrade probes before assigning UDP checks. Responses larger than 4096 bytes fail validation rather than silently passing a partial assertion.

`GET /monitors/:id/uptime` and the public `GET /status` accept `hours` (1–720, default 24) and `bucket_minutes` (5–1440 in multiples of 5, default 15), bounded to 2016 cells. Each cell includes success/failure/executor-error counts, expected samples, availability and coverage. Missing history remains unknown; archived five-minute samples remain visible. Availability is the proportion of successful observed target checks, not a time-based SLA or confirmed incident duration. Coverage uses enabled configuration versions and includes offline gaps. Public responses exclude targets, node identities, raw evidence and incident notes.

`GET /monitors/:id/history` accepts `from`, `until` (Unix seconds), `status`, `stage`, `node_id`, `limit` and `offset`; `total` counts matching raw results. Raw-result filters do not change the time/node summary.

Diagnostic jobs now carry optional `batch_id` and `monitor_id`. Creating a multi-node run is atomic; query `/diagnostics?batch_id=...` to retrieve the whole run without unrelated jobs. Existing clients may ignore these additive fields.
