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