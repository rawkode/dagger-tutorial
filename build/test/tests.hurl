GET http://backend:8080/ping
HTTP 200

POST http://backend:8080/ask
```json
{
  "question": "Translate \"Hello\" from English to German"
}
```

HTTP 200
[Asserts]
jsonpath "$.result" contains "Rawkode"
