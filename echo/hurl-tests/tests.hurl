PUT http://{{host}}/api/v1/numbers
Content-Type: application/json
{
  "value": 2, "index": 0
}
HTTP 201
[Asserts]
jsonpath "$.index" == 0
jsonpath "$.value" == 2

PUT http://{{host}}/api/v1/numbers
Content-Type: application/json
{
  "value": 4, "index": 1
}
HTTP 201
[Asserts]
jsonpath "$.index" == 1
jsonpath "$.value" == 4

PUT http://{{host}}/api/v1/numbers
Content-Type: application/json
{
  "value": 3, "index": 3
}
HTTP 400
[Asserts]
jsonpath "$.message" == "Invalid index"

GET http://{{host}}/api/v1/numbers/index/0
HTTP 200
[Asserts]
jsonpath "$.index" == 0
jsonpath "$.value" == 2

GET http://{{host}}/api/v1/numbers/index/2
HTTP 404
[Asserts]
jsonpath "$.message" == "Index not found"

GET http://{{host}}/api/v1/numbers/value/4
HTTP 200
[Asserts]
jsonpath "$.index" == 1
jsonpath "$.value" == 4

GET http://{{host}}/api/v1/numbers/value/6
HTTP 404
[Asserts]
jsonpath "$.message" == "Value not found"

DELETE http://{{host}}/api/v1/numbers/1
HTTP 200

GET http://{{host}}/api/v1/numbers/index/1
HTTP 404
[Asserts]
jsonpath "$.message" == "Index not found"
