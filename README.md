# DevOps-Excersize

Go lightweight HTTP microservice to calculate math functions from a supplied array of numbers. Done for a DevOps exercise.

## Constraints and Features

* Strict JSON validation (no coercion, no partial parsing)
* Responses rounded to **3 decimal places**
* Configurable port via environment variable (`PORT`)
* Dockerized for easy deployment

| Endpoint   | Description |
| ---------- | ----------- |
| `/mean`    | Calculates the mean of all numeric values in the JSON array |
| `/stddev`  | Calculates the population standard deviation of all numeric values in the JSON array |

## Build and Run Options

### 1. Use Provided Script

If using the included script:

```bash
chmod +x build_run.sh
./build_run.sh
```

Customizations to build params can be done in your `.env`:

```env
IMAGE_NAME=math-service
CONTAINER_NAME=math-service-container
PORT=3000
```
> NOTE: `build_run.sh` automatically removes any currently running container with the same name and replaces it.

---

### 2. Run with Docker

#### Build the image

```bash
docker build -t math-service .
```

#### Run the container

```bash
docker run -p 3000:3000 -e PORT=3000 math-service
```

---

### 3. Run Locally

Default port: `3000`

```bash
go run main.go
```

Override port:

```bash
PORT=8080 go run main.go
```


## Request Requirements

- HTTP method: `POST`
- Header: `Content-Type: application/json`
- Body: JSON array
- Numeric values only; non-numeric strings are ignored and the remaining values are calculated
- Payload max size: `1 MiB`

## HTTP Details

All requests must include:

```http
Content-Type: application/json
```

---

| HTTP Status | Description |
| ----------- | ----------- |
| `200 OK` | Valid JSON request accepted; numeric result returned as plain text |
| `400 Bad Request` | Invalid JSON or no valid numeric values provided |
| `405 Method Not Allowed` | Request method is not `POST` |
| `413 Payload Too Large` | Request body exceeds 1 MiB |

## Usage

### Request

```bash
curl -X POST http://localhost:3000/mean \
  -H "Content-Type: application/json" \
  -d '[1,2,3,4,5]'
```

### Response

#### HTTP Status
`200 OK`

```json
{
  "result": 3
}
```

---

### Request

```bash
curl -X POST http://localhost:3000/stddev \
  -H "Content-Type: application/json" \
  -d '[1,2,3,4,5]'
```

### Response

#### HTTP Status
`200 OK`

```json
{
  "result": 1.414
}
```

---

### Decimal Input

```bash
curl -X POST http://localhost:3000/mean \
  -H "Content-Type: application/json" \
  -d '[1.5,3,1,7,4.2]'
```

### Response

#### HTTP Status
`200 OK`

```json
{
  "result": 3.34
}
```

---

## Error Handling

### Invalid JSON

```bash
curl -X POST http://localhost:3000/mean \
  -H "Content-Type: application/json" \
  -d '{"numbers":[1,2,3]}'
```

### Response

#### HTTP Status
`400 Bad Request`

```json
{
  "error": "Request body must be a JSON array of numbers"
}
```

---

### Non-Numeric Values

```bash
curl -X POST http://localhost:3000/mean \
  -H "Content-Type: application/json" \
  -d '[1,"abc",3]'
```

### Response

#### HTTP Status
`400 Bad Request`

```json
{
  "error": "Request body must be a JSON array of numbers"
}
```

---

### Empty Array

```bash
curl -X POST http://localhost:3000/mean \
  -H "Content-Type: application/json" \
  -d '[]'
```

### Response

#### HTTP Status
`400 Bad Request`

```json
{
  "error": "Request array must contain at least one Numeric Value"
}
```