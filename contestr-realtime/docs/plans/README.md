# Планы разработки

Два связанных плана по архитектуре регатты и realtime-уведомлений.

| Порядок | Файл | Содержание |
|--------|------|------------|
| 1 | [01-regatta-pipeline-snapshots.plan.md](01-regatta-pipeline-snapshots.plan.md) | Sync → compute → `regatta_snapshots` → GET API (без push) |
| 2 | [02-contestr-realtime-integration.plan.md](02-contestr-realtime-integration.plan.md) | Socket.IO, push из Go, overlay на фронте (после плана 1) |

Рекомендуется выполнять **сначала план 1**, затем план 2.
