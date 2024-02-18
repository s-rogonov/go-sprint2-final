# go-sprint2-final

1. Схема представлена в ПДФ-ке
2. Успел только оркестрартор и агента (не ругайте за поздний коммит)
    - Оркестратор был на 100% рабочий до 23:00 18 февраля (commit 84c11feb5348646ceab5d5619563c06c6602e6e5)
    - Агента закончил к 01:05, сами судите учитывать его или нет
3. Оркестратор полностью пассивен, все его ручки описаны ниже
4. Агент периодически забирает таски от оркетсратора и раздаёт их своим горутинам
5. Оркестратор не проверяет живость агентов, а ориентируется на критерий `factor`.
Если Агент не вернул ответ на задачу за время `duration*factor` - оркестратор выдаст задачу повторно
6. Оркестратор "не знает" кому из агентов достанется задача, для оркестратора агенты остаются анонимны (помним, оркестратор пассивен и не ведёт ни регистрацию ни учёт живости агентов)

По сущностям:
1. Оркестратор по ручке содаёт сущность `Query` и если получилось, сразу россыпь связанных с ней `Task`
2. Агент забирает доступные таски по ручке, на которой оркестратор создаёт сущности `Worker` - так он отличает `Task` взятые в работу
3. Завершая задачу, агент стучится оркестратору и тот помечает `Task` как готовый, и по цепочке зависящий от него `Task` может стать готовым к работе, если все подзадачи последнего готовы
4. Когда завершается самый верхнеуровневый `Task`, его `Query` помечается готовым


Парсинг выражений поддерживает:

- целые, не отрицателные числа
- 4 вида операций (+,-,*,/)
- скобки
- учитывает приоритет операций и скобок

Из его недостатков:

- порядок операций фиксирован на этапе разбора
- иногда получаются параллельные задачи, но агрессивных оптимизаций, меняющих порядок вычислений на лету ради максимизации параллельных вычислений не предусмотрено

## Agent

> Multiple agents can be run simultaneously

Runs with

```bash
go run agent/main.go
```

Env vars:

- `MASTER` - master (orchestrator) endpoint, default=`localhost:8181`
- `NWORKERS` - number of goroutines
- `BATCH` - max amount of tasks, polled from orchestrator
- `DELAY` - in seconds, period of new tasks polling

## Orchestrator (GET)

> Single orchestrator per single sqlite DB (synchronous db access)

Runs with

```bash
go run orchestrator/main.go
```

Env vars:

- `PORT` - master (orchestrator) port to listen, default=`8181`
- `DBNAME` - name of sqlite db file

---

### 1. `GET /timings`

*Return timings; `factor` is a threshold to identify workers with `status=timeout`*

Example:

```bash
curl -X GET localhost:8181/timings
```

Answer:

```json
{
  "factor": 2,
  "add": 10,
  "sub": 100,
  "mul": 20,
  "div": 200
}
```

---

### 2. `GET /queries`

*Return queries; By default, `limit=3` asc; Specify `{id}` to select specific query*

Example (`/queries/{id}`):

```bash
curl -X GET 'localhost:8181/queries/1'
```

Answer:

```json
{
  "expr": "1+2",
  "id": 1,
  "result": 3,
  "status": "finished"
}
```

Example (`/queries?limit={n}`):

```bash
curl -X GET 'localhost:8181/queries?limit=3'
```

Answer:

```json
[
  {
    "expr": "1+2+3",
    "id": 5,
    "progress": "0/2",
    "status": "in-progress"
  },
  {
    "expr": "18+3*(4-5)",
    "id": 4,
    "progress": "0/3",
    "status": "in-progress"
  },
  {
    "expr": "3*4-5",
    "id": 3,
    "progress": "0/2",
    "status": "in-progress"
  }
]
```

---

### 3. `GET /workers`

*Return workers status; Basically workers are kinda computing resources; By default, `limit=3` asc; Specify `{id}` to
select specific worker*

Example (`/workers/{id}`):

```bash
curl -X GET localhost:8181/workers/12
```

Answer:

```json
{
  "expr": "1 * 2",
  "id": 12,
  "left": 89.11674,
  "status": "computing"
}
```

Example (`/workers?limit={n}`):

```bash
curl -X GET localhost:8181/workers
```

Answer:

```json
[
  {
    "expr": "1 * 2",
    "id": 12,
    "left": 99.63454,
    "status": "computing"
  },
  {
    "deadline": "2024-02-18T23:32:51.984627591+03:00",
    "expr": "1 + 2",
    "id": 11,
    "status": "timeout"
  },
  {
    "expr": "4 - 5",
    "id": 10,
    "left": 99.6339,
    "status": "retrieving"
  }
]
```

## Orchestrator (PUT/POST)

---

### 1. `POST /query`

*Creates new query*

Example:

```bash
curl -X POST localhost:8181/query -d "1+2"
```

Answer:

```json
{
  "errorMsg": "",
  "hasError": false,
  "id": 1
}
```

---

### 2. `PUT /query`

*Updates existing query; Updated query should have parse error, else changes will be rejected*

Example:

```bash
curl -X PUT localhost:8181/query -d '{"id":4,"expr":"1+2"}'
```

Answer:

```json
{
  "errorMsg": "",
  "hasError": false,
  "id": 4
}
```

---

### 3. `POST /tasks`

*Returns operations to process; Amount of operations are limited by argument*

Example (asks 4 operations maximum):

```bash
curl -X POST localhost:8181/tasks -d '4'
```

Answer (there is only 2 operations available):

```json
[
  {
    "id": 1,
    "op": "+",
    "time": 1000000000,
    "args": [
      1,
      2
    ]
  },
  {
    "id": 2,
    "op": "*",
    "time": 1000000000,
    "args": [
      3,
      4
    ]
  }
]
```

---

### 4. `PUT /result`

*Set operation result*

Example:

```bash
curl -X PUT localhost:8181/result -d '{"id":1,"result":3.0}'
```

Answer:

```
OK
```

### 5. `PUT /timings`

*Set operator types duration*

Example (add/sub/mul/div in seconds):

```bash
curl -X PUT localhost:8181/timings -d '{"factor":2,"add":10,"sub":100,"mul":20,"div":200}'
```

Answer:

```
OK
```

---
