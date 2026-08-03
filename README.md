# DevOps-Interview

Go microservice for a DevOps interview exercise. This service exposes two POST endpoints that accept a JSON array of numbers and return a numeric result as plain text.

> NOTE: `build_run.sh` automatically removes any currently running container with the same name and replaces it.

## Build and Run

- `./build_run.sh`
- The service listens on `PORT`, defaulting to `3000` when the env var is not set.

## Endpoints

| Endpoint   | Description |
| ---------- | ----------- |
| `/mean`    | Calculates the mean of all numeric values in the request body |
| `/stddev`  | Calculates the population standard deviation of all numeric values in the request body |

## Request requirements

- HTTP method: `POST`
- Header: `Content-Type: application/json`
- Body: JSON array
- Numeric values only; non-numeric strings are ignored and the remaining values are calculated
- Payload max size: `1 MiB`

## HTTP responses

| HTTP Status | Description |
| ----------- | ----------- |
| `200 OK` | Valid JSON request accepted; numeric result returned as plain text |
| `400 Bad Request` | Invalid JSON or no valid numeric values provided |
| `405 Method Not Allowed` | Request method is not `POST` |
| `413 Payload Too Large` | Request body exceeds 1 MiB |

## Examples
```bash
curl http://localhost:3000/mean \
  -X POST \
  -H "Content-Type: application/json" \
  -d '[1,2,3,4,5]'
```
Response:
`3`

```bash
curl http://localhost:3000/mean \
  -X POST \
  -H "Content-Type: application/json" \
  -d '[1.5,3,1,7,4.2]'
```
Response:
`3.34`  

  

```bash
curl http://localhost:3000/stddev \
  -X POST \
  -H "Content-Type: application/json" \
  -d '[1,2,3,4,5]'
```
Response:
`1.414`

```bash
curl http://localhost:3000/stddev \
  -X POST \
  -H "Content-Type: application/json" \
  -d '[1.5,3,1,7,4.2]'
```
Response:
`2.15`
