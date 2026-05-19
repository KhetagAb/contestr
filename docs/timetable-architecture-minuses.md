# Архитектура расписания и туров

Документ описывает прежние ограничения и **текущую** модель после рефакторинга (greenfield, без миграции старых данных).

## Бывшие проблемы (до рефакторинга)

1. **Два источника правды** — план в `ToursTimetable` (`started`, `start_time`) и факт в `Tour`; админка и публичная страница могли расходиться.
2. **Админка и двойная сверка** — оператор смотрел на ось расписания и отдельно на коллекцию туров в Mongo.
3. **История и группировка** — «История» на standings группировала события по `problem_code`, разделитель тура был только между двумя последними турами.
4. Гонки ручного/авто старта, разное поведение manual vs sync, polling `timetable_sync`, статусы UI не отражали факт регаты.

## Текущая модель

| Слой | Хранение | Содержимое |
|------|----------|------------|
| **Факт** | `Tour` (Mongo `tours`) | `sequence`, `round`, `is_pause`, `duration_in_seconds`, группы/задачи — **без** `start_time` |
| **План** | `ToursTimetable` | только `pending_slots[]` (`kind`: `tour` \| `pause`, `duration`) + `auto_start_enabled` |
| **Read-model** | `GET /api/admin/timetables/{id}` | `timeline_segments[]` — прошлые Tour + pending с вычисленным `start_time` и `status` |

### Двойная нумерация

- **`sequence`** — порядок на оси (туры и паузы).
- **`round`** — номер соревновательного тура (1, 2, 3…); у паузы `round = 0`.
- Коды задач: `round` + буква (`1A`, `2B`).

### Пауза

Пауза — обычный документ `Tour` с `is_pause: true`, пустые `groups`/`problems`. Участвует в цепочке времени (сдвигает окно scoring), **не** участвует в подсчёте баллов. В плане — слот `kind: pause` в `pending_slots`.

### Время

- `segmentStart(sequence) = sum(duration)` всех документов с меньшим `sequence`.
- Scoring: окно тура = `[segmentStart, segmentStart + duration)`.
- Ручной/авто **advance**: `POST /api/admin/timetables/{contest_id}/advance` — укоротить активный сегмент, стартовать голову `pending_slots`.

### Автозапуск

`timetable_sync` вызывает `AdvanceTimetable(Auto)` для контестов с `auto_start_enabled`.

## Связанные файлы

- `contestr/pkg/regatta/timeline.go` — сегменты, pending, meta
- `contestr/pkg/regatta/tours.go`, `timetable_view.go` — модели
- `contestr/internal/services/regatta/timetable.go` — `AdvanceTimetable`, CRUD расписания
- `contestr/internal/services/regatta/start_tour.go` — создание Tour / pause
- `contestr/internal/services/timetable_sync/` — автозапуск
- `contestr-front/src/admin-timetable/` — админская ось

Перед тестом на чистой БД: `db.tours.deleteMany({})`, `db.tour_timetables.deleteMany({})`.
