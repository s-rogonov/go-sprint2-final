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
