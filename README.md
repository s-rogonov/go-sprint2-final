# go-sprint2-final

## Orchestrator

---

> POST /query

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

> PUT /query

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

> POST /tasks

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

> PUT /result

*Set operation result*

Example:

```bash
curl -X PUT localhost:8181/result -d '{"id":1,"result":3.0}'
```

Answer:

```
OK
```