---
name: Regatta pipeline snapshots
overview: "План 1: единый pipeline (1m) — sync raw → compute regatta → snapshot в Mongo → GET standings читает snapshot. Без realtime и без push."
todos:
  - id: compute-extract
    content: "Вынести ComputeContestResult из GetContestResult (формулы без изменений)"
    status: pending
  - id: snapshot-repo
    content: "Mongo regatta_snapshots + RegattaSnapshotRepository"
    status: pending
  - id: pipeline-service
    content: "RegattaPipelineService: sync → compute → upsert snapshot"
    status: pending
  - id: wire-pipeline
    content: "Заменить tick contest_sync на pipeline; admin SyncContest — тот же путь"
    status: pending
  - id: api-snapshot
    content: "GET standings → чтение snapshot; OpenAPI computed_at"
    status: pending
  - id: verify-pipeline
    content: "Проверка: API= snapshot; rows/events с regatta-правилами; первый tick без сюрпризов"
    status: pending
isProject: false
---

# План 1: Pipeline просчёта и snapshot регатты

**Зависимости:** нет.  
**Связанный план:** [02-contestr-realtime-integration.plan.md](02-contestr-realtime-integration.plan.md) (push-уведомления — после этого плана).

## Цель

Убрать пересчёт standings/events на каждый HTTP-запрос. Один фоновый цикл по **`contest_sync.interval: 1m`** ([`contestr/configs/config.yaml`](../../contestr/configs/config.yaml)):

1. Загрузка с CF/Ejudge → `contests` (raw submissions).
2. Просчёт regatta (`CalculateResultWithSettingsAt`, `BuildContestEventsAt`).
3. Сохранение в **`regatta_snapshots`**.
4. API отдаёт готовый snapshot.

```mermaid
flowchart LR
    Ext[CF/Ejudge] --> Sync[Upsert contests]
    Sync --> Compute[ComputeContestResult]
    Compute --> Snap[Upsert regatta_snapshots]
    Snap --> API[GET standings]
```

## Модель данных

```go
type RegattaSnapshot struct {
    ContestID  int
    ComputedAt time.Time
    Standings  regatta.ContestStandings // rows + events + metadata
}
```

- Одна метка **`ComputedAt`** на успешный проход pipeline.
- **Без** `PublishedEventKeys` — diff для push будет в плане 2.

## RegattaPipelineService

[`contestr/internal/services/regatta_pipeline/`](../../contestr/internal/services/regatta_pipeline/) (или рефакторинг `contest_sync`):

| Шаг | Действие |
|-----|----------|
| 1 | `adapter.FetchContest` → `contestRepo.Upsert` |
| 2 | `ComputeContestResult(ctx, contestID)` |
| 3 | `snapshotRepo.Upsert` |

Триггеры: тикер 1m, admin `SyncContest` / refresh.

Ошибка на шаге 1 → snapshot не трогаем. Ошибка на 2–3 → лог, старый snapshot остаётся для API.

## API

[`get_contest_standings.go`](../../contestr/internal/handlers/regatta/get_contest_standings.go):

- **Было:** `GetContestResult` (compute на каждый GET).
- **Станет:** `snapshotRepo.Get(contestID)` → JSON как сейчас + поле **`computed_at`**.
- Нет snapshot: `404` или пустой ответ (зафиксировать в OpenAPI).

**Events:** полный `events[]` в snapshot — лента [`ContestEventLog`](../../contestr-front/src/features/contest/components/event-log/ContestEventLog.tsx) без изменений контракта.

## Рефакторинг Go

- `ComputeContestResult` — internal, только pipeline.
- `GetContestResult` на HTTP handler **не** вызывается (удалить или оставить alias для тестов → snapshot).
- Wire: pipeline вместо голого `ContestSyncService.Start` loop body.

## Фронт (минимально)

- Standings poll можно оставить 60s — данные обновляются раз в 1m.
- Опционально показать `computed_at` в UI.
- **Не** добавлять socket/overlay в этом плане.

## Проверка

| Сценарий | Ожидание |
|----------|----------|
| После pipeline | GET === snapshot, очки с overtake/bonus |
| До первого pipeline | 404 / пусто |
| Admin refresh | Один contest проходит pipeline |
| CF solve | В ленте events после следующего тика (~1m) |

## Вне scope

- contestr-realtime, push, overlay.
- Ускорение interval &lt; 1m.
- Append-only коллекция events.
