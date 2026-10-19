# Печать героя

Тренажёр слепой печати с геймификацией. Хакатон «Киберрывок / Осень 2026», отборочный этап.

- **Бэкенд**: Go 1.23, chi, pgx, PostgreSQL 16, JWT (argon2id).
- **Фронтенд**: Vue 3 + Vite + TypeScript + Pinia + vue-router + Chart.js.
- **Запуск одной командой**: `docker compose up --build -d` → `http://localhost:8000`.

## Быстрый старт (Docker)

```bash
docker compose up --build -d
```

Открыть **http://localhost:8000**.

Что происходит:
- `postgres` — Postgres 16, порт `5433` на хосте (чтобы не конфликтовать с локальным Postgres на `5432`), внутри docker-сети — `postgres:5432`.
- `app` — Go-бинарник со встроенной статикой Vue (`frontend/dist`), слушает `:8000`.
- При старте `app`:
  - накатывает миграции (`goose`, `internal/store/migrations`, версия 7);
  - сидирует уровни 1..100, анлоки, 15 достижений;
  - сидирует 22 текста (отрывки `sentences` с `source`, фрагменты `code`);
  - запускает HTTP-сервер.

Остановить:

```bash
docker compose down
```

Остановить и удалить данные БД:

```bash
docker compose down -v
```

## Архитектура

```
cmd/api/                    — main: конфиг, миграции, сиды, DI, HTTP-сервер
internal/
  config/                   — env, валидация
  httpx/                    — APIError, WriteJSON/WriteError, user_id в контексте
  token/                    — JWT Claims, Parse (typ=access|refresh)
  auth/                     — register/login/refresh, argon2id, rotation refresh
  user/                     — GET/PATCH /me, merge settings
  settings/                 — Settings, Patch, Default, Apply
  progress/                 — XP, уровни, достижения, сброс прогресса
  gamification/             — формулы XP/уровня, каталог достижений, анлоки, сиды
  sessions/                 — POST/GET /sessions, метрики, звёзды узлов
  maps/                     — каталог 4 карт, звёзды, /maps
  texts/                    — /texts, генерация слов/символов, сиды
  stats/                    — summary/timeline/keys
  leaderboard/              — speed/accuracy/xp
  store/postgres/           — пул pgx
  store/migrations/         — goose-миграции 00001..00007
  httpserver/               — chi-роутер, middleware, статика + SPA-fallback
frontend/                   — Vue 3 + Vite + TS
  src/api/                  — httpApi.ts, schema.d.ts (из openapi.yaml), types.ts
  src/stores/               — Pinia
  src/views/                — экраны
  src/typing/               — логика набора
content/                    — общие файлы контента (если используются)
openapi.yaml                — контракт 0.2.1
Dockerfile                  — multi-stage: Vue → Go → alpine
docker-compose.yml          — postgres + app
```

## Эндпоинты

| Метод | Путь | Auth | Описание |
|---|---|---|---|
| POST | `/api/v1/auth/register` | — | регистрация, 201 |
| POST | `/api/v1/auth/login` | — | логин, 200 |
| POST | `/api/v1/auth/refresh` | — | rotation refresh, 200 |
| GET / PATCH | `/api/v1/me` | Bearer | профиль, merge settings |
| GET | `/api/v1/texts` | опц. | тексты, без 403 |
| POST | `/api/v1/sessions` | Bearer | сохранить сессию, `node_id` опц. |
| GET | `/api/v1/sessions` | Bearer | история, новые сверху |
| GET | `/api/v1/maps` | опц. | 4 карты, звёзды, `ready` |
| GET / DELETE | `/api/v1/me/progress` | Bearer | прогресс / сброс |
| GET | `/api/v1/achievements` | опц. | каталог 15 достижений |
| GET | `/api/v1/me/stats/summary` | Bearer | сводка |
| GET | `/api/v1/me/stats/timeline` | Bearer | по дням, `range=week\|month\|year` |
| GET | `/api/v1/me/stats/keys` | Bearer | проблемные клавиши/биграммы |
| GET | `/api/v1/leaderboard` | опц. | `metric=speed\|accuracy\|xp` |

Ошибки — единый формат:

```json
{"code":"validation_error","message":"validation error","details":{"field":"..."}}
```

Коды: `unauthorized`, `invalid_credentials`, `email_taken`, `nickname_taken`, `validation_error`, `not_found`, `rate_limited`, `internal_error`.

## Правила игры (сервер — источник истины)

- **XP за сессию** = `round(chars_correct * accuracy/100 * (1 + difficulty * 0.1))`.
- **Уровень**: с уровня `L` нужно `100 + 50*(L-1)` XP. `xp_in_level`/`xp_to_next` отдаются посчитанными.
- **Звёзды узла карты**: сколько из 3 `star_thresholds` выполнено (`wpm >= min_wpm && accuracy >= min_accuracy`). Лучший результат хранится в `map_node_results`, узел пройден при `stars >= 1`.
- **Стрик**: дни подряд с тренировкой, заканчивая сегодня или вчера, UTC.
- **Рейтинг**: `speed` — лучший WPM за период, `accuracy` — средняя точность (только пользователи с ≥3 сессиями за период), `xp` — сумма XP. Периоды `today`/`week`/`all`.
- **Проблемные клавиши**: `hits >= 3 && errors > 0`, сортировка по `error_rate` desc.
- **Достижения** (15): `first_steps`, `speed_200`, `speed_300`, `genius`, `clean_100`, `steel_finger`, `perfect`, `thousand`, `marathon`, `streak_7`, `stars_10`, `stars_30`, `combo_50`, `coder`, `polyglot`. Проверяются после каждой сессии, новые возвращаются в `SessionResult.new_achievements`.

Формулы метрик:
- `WPM = (chars_correct / 5) / (duration_ms / 60000)`
- `CPM = chars_correct / (duration_ms / 60000)`
- `accuracy = chars_correct / chars_typed * 100`

Клиент считает «живой» счёт теми же формулами; сервер пересчитывает итог по сырым данным `POST /sessions`.

## Контент

- **Слова/символы** — генерируются в Go (`internal/texts/generators.go`) по сложности 1..5.
- **`sentences`** — реальные отрывки public domain (Пушкин, Толстой, Лермонтов, Тютчев, Диккенс, Толстой/англ., Нойес, Дикинсон, техническая/научная проза). Поле `source` — автор и название. Сложность = род текста: 1 справочный, 2 научный, 3 классика, 4 стихи с ритмом, 5 сложные стихи.
- **`code`** — Python/JS/HTML/CSS, отступ 4 пробела, `\n`.
- **Карты** — 4 штуки (`journey`, `chase`, `search`, `siege`), у каждой `recommended_level`, `color`, узлы с `star_thresholds`. Все доступны всегда, `ready=false` — только предупреждение, не блокировка.

Сиды идемпотентны (`ON CONFLICT ... DO UPDATE`) и вызываются в `main.go` при старте.


В браузере: `http://localhost:8000` → регистрация → уровень → `POST /sessions` → `new_achievements`, `progress`, `node_result` → `/maps` со звёздами → `/leaderboard`.

## Известные ограничения

- **`UpsertBest` карт** и `sessions.Repo.Create` — не в одной транзакции. Если `Create` упадёт, звёзды узла уже обновятся. Помечено TODO.
- **`content/maps.json` / `content/texts.json`** — если фронт их не импортирует, контент живёт в Go-сидах и БД; фронт получает через `/maps` и `/texts`. Это поведение соответствует контракту, даже если не буквально «один файл».
- **CORS** нужен только для dev (`localhost:5173`/`127.0.0.1:5173`). В Docker фронт и API на одном origin, CORS не участвует.
