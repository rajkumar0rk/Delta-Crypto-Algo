# Oracle Cloud deployment plan

1. Begin with local paper mode and a private control API (foundation binds loopback).
2. Add staged deployment IaC (least privilege systemd service, restricted network security groups, TLS proxy, automatic encrypted backups).
3. Never place API keys or production secrets in repository or GitHub Actions logs.
4. Store state in PostgreSQL, take verified backups to separate durable storage and regularly test full restore.
5. Collect service health, process memory, disk pressure, time sync, exchange data gaps, orders/fills mismatch, alerts and detailed logs without credentials.
6. Production should use redundant nodes and persistent databases when operational stakes warrant; a single Oracle 1GB VPS is NOT enterprise high availability.
7. Release: PR -> CI -> reviewed artifact -> demo soak -> manual production approval -> canary with hard risk limits. Provide rollback and incident runbook.

**No deployment has been performed during M1.**
