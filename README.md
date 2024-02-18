# go-sprint2-final

## Orchestrator (GET)

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
