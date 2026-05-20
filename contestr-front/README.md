# contestr-front

React-приложение регаты: таблица результатов для зрителей и админка (расписание туров, контесты Codeforces).

**Для разработчиков:** практические шаги и примеры — [DEVELOP.md](DEVELOP.md).

## Быстрый старт

```bash
cd contestr-front
npm install
npm run dev          # http://localhost:5173, прокси /api → бэкенд
```

Бэкенд по умолчанию: `http://127.0.0.1:8080`. Другой хост:

```bash
VITE_API_PROXY_TARGET=http://host:8080 npm run dev
```

Перед `dev` и `build` автоматически запускается `generate-client` (см. `predev` в `package.json`).

```bash
npm run build        # generate-client + tsc + vite build
npm run lint
npm run preview      # просмотр production-сборки
```

---

## Структура `src/`

```
src/
  main.tsx                 # тонкая точка входа → app/App
  app/                     # оболочка приложения
    App.tsx                # роутинг по pathname (/admin vs зритель)
    providers/             # QueryClient, AdminSessionContext
    styles/                # index.css, App.css (глобальные)
  client/                  # ⚠️ только генерация OpenAPI, не править вручную
  features/
    contest/               # зрительский UI
      pages/               # страницы целиком
      components/          # виджеты по доменам (event-log/, sidebar/)
      hooks/               # данные и логика фичи contest
    admin/                 # админка
      pages/               # AdminLogin, AdminConsole
      auth/                # токен, заголовки
      timetable/           # расписание туров
      contests/            # контесты и участники
    legacy/                # архив, не подключать к app
  shared/                  # переиспользуемое между фичами
    hooks/
    utils/
    icons/
  assets/                  # SVG/PNG для import (Vite)
    icons/
    images/
```

Отдельно: [`match-integration/`](match-integration/) — другой бандл и свой OpenAPI.

---

## Как устроен код

### Куда класть новое

| Задача | Куда |
|--------|------|
| Страница зрителя | `features/contest/pages/` — суффикс `Page` |
| Страница / раздел админки | `features/admin/pages/` или `timetable/`, `contests/` |
| UI-компонент фичи | `features/<feature>/components/<домен>/` |
| Логика и запросы фичи | `features/<feature>/hooks/` |
| Общая утилита / хук | `shared/utils/`, `shared/hooks/` |
| Иконка-компонент | `shared/icons/` |
| Картинка | `assets/icons/`, `assets/images/` |
| Провайдеры, глобальные стили | `app/providers/`, `app/styles/` |
| Старый / черновой код | `features/legacy/` |

### Нейминг

| Сущность | Стиль | Пример |
|----------|--------|--------|
| Папки | `kebab-case` | `event-log/` |
| Компоненты | `PascalCase`, файл = имя | `Sidebar.tsx` |
| Страницы | `*Page.tsx` | `ContestStandingsPage.tsx` |
| Хуки | `use` + домен | `useContestStandings` |
| CSS modules | рядом с компонентом | `event-log.module.css` |

Импорты: алиас `@/` → `src/` (`@/client`, `@/features/...`, `@/shared/...`).

---

## Точки входа

| Файл | Назначение |
|------|------------|
| [`src/main.tsx`](src/main.tsx) | `createRoot`, глобальный CSS |
| [`src/app/App.tsx`](src/app/App.tsx) | провайдеры, `Sidebar`, страница по URL |
| [`src/features/contest/pages/ContestStandingsPage.tsx`](src/features/contest/pages/ContestStandingsPage.tsx) | таблица + история |
| [`src/features/admin/pages/AdminLogin.tsx`](src/features/admin/pages/AdminLogin.tsx) | `/admin` |

URL зрителя: `/?contestId=<id>` (контест в сайдбаре).

---

## Связанные репозитории

- OpenAPI и бэкенд: [`../contestr/`](../contestr/)
- Docker/nginx: корень монорепо [`../`](../)
