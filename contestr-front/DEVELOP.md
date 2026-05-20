# Разработка contestr-front

Практические инструкции. Обзор проекта и структура — в [README.md](README.md).

## Данные и API

1. Схема: [`../contestr/api/openapi.yaml`](../contestr/api/openapi.yaml).
2. После изменений бэкенда:

```bash
npm run generate-client
```

3. Типы и клиент — только из `src/client/`. Ручные правки перезапишутся.
4. **Зритель:** опции React Query из `@/client/@tanstack/react-query.gen` (например `getRegattaContestStandingsOptions`); обёртка в хуке фичи — `features/contest/hooks/useContestStandings.ts`.
5. **Админка:** `adminAuthHeaders()` из `@/features/admin/auth/adminAuth`, `fetch` на `/api/admin/...` — см. `features/admin/timetable/useContestTimetable.ts`.

## Стили

1. **Глобальные** (сайдбар, сетка) — `app/styles/App.css`, `app/styles/index.css`.
2. **Компонент** — CSS module рядом с `.tsx`, один модуль на логический блок (не один файл на всё приложение).

```tsx
import styles from "./event-log.module.css";

<section className={styles.eventLog}>…</section>
```

## Добавить экран зрителя

1. Страница: `features/contest/pages/MyPage.tsx` (суффикс `Page`).
2. Данные: хук в `features/contest/hooks/` при необходимости.
3. Подключение: `app/App.tsx` — сейчас `/admin` → админка, иначе `ContestStandingsPage`; для нового маршрута расширить `AppShell`.

## Добавить раздел админки

1. UI: `features/admin/<раздел>/` или `features/admin/pages/`.
2. Вкладка: `features/admin/pages/AdminConsole.tsx`.
3. Авторизация уже в `AdminLogin` и `AdminSessionContext`.

## Статика

| Нужно | Папка | Пример |
|-------|--------|--------|
| Логотип, картинки в компонентах | `src/assets/icons/`, `src/assets/images/` | `import logo from "@/assets/icons/logo.svg"` |
| Файл по фиксированному URL | `public/` | `/file.svg` (редко) |

Эталон медали: `assets/icons/sport-winner.svg`. В UI — `shared/icons/MedalFirstIcon.tsx` (inline SVG, `currentColor`).

## Иконки

- **react-icons** / **lucide-react** — готовые пиктограммы.
- **Свой SVG** — компонент в `shared/icons/` (как `MedalFirstIcon`), если нужны `currentColor` и читаемость на ~1em.

## Импорты (примеры)

```ts
import { useContestStandings } from "@/features/contest/hooks/useContestStandings";
import { formatGroupCode } from "@/shared/utils/groupCode";
import type { RegattaEvent } from "@/client";
```

- Между фичами — только `@/…`, не `../../../`.
- `contest` ↔ `admin` напрямую не импортировать; общее — `shared/` или `app/providers/`.
- Не импортировать из `features/legacy/`.

## Проверка перед коммитом

```bash
npm run build
npm run lint
```

Если `dist/` создан Docker’ом и `build` падает с `EACCES` — удалить `dist` с нужными правами или пересобрать образ.

## Чего избегать

- Править `src/client/` вручную.
- Рабочий код в `features/legacy/`.
- Дублировать `public/` внутри `src/` — импортируемое в `assets/`.
- Импорт `admin` ↔ `contest` без выноса в `shared/`.
- Один CSS module на несвязанные экраны.
