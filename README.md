# Presign Service

## Configuration

### API Key Generation

Generate a secure API key for authentication:

**Option 1: Using the provided Go script**

```bash
go run scripts/generate-api-key.go
```

**Option 2: Using the shell script**

```bash
./scripts/generate-api-key.sh
```

### Environment Variables

| Variable                | Required | Default     | Description                                      |
| ----------------------- | -------- | ----------- | ------------------------------------------------ |
| `BUCKET_NAME`           | Yes      | -           | S3/R2 bucket name for presigning                 |
| `API_KEY`               | Yes      | -           | API key for authentication                       |
| `AWS_ACCESS_KEY_ID`     | Yes\*    | -           | AWS/R2 access key ID                             |
| `AWS_SECRET_ACCESS_KEY` | Yes\*    | -           | AWS/R2 secret access key                         |
| `AWS_ENDPOINT_URL`      | No       | -           | Custom endpoint URL (required for Cloudflare R2) |
| `PORT`                  | No       | `8080`      | Server port                                      |
| `AWS_REGION`            | No       | `us-east-1` | AWS/R2 region                                    |

### Cloudflare R2 Configuration

For Cloudflare R2, you'll need to:

1. Get your R2 credentials from the Cloudflare dashboard
2. Set the endpoint URL to your account's R2 endpoint: `https://<account-id>.r2.cloudflarestorage.com`
3. Use any region (R2 doesn't use regions like AWS, but the SDK requires one)

Example R2 configuration:

```bash
export AWS_ACCESS_KEY_ID="your-r2-access-key-id"
export AWS_SECRET_ACCESS_KEY="your-r2-secret-access-key"
export AWS_ENDPOINT_URL="https://your-account-id.r2.cloudflarestorage.com"
export AWS_REGION="auto"
export BUCKET_NAME="your-r2-bucket-name"
```

## API Endpoints

### POST /sign

Generates presigned URLs for S3 objects.

**Authentication**: Required (Bearer token)

**Request Body**:

```json
{
  "keys": ["s3/key/1.jpg", "s3/key/2.jpg"]
}
```

**Response**:

```json
{
  "urls": {
    "s3/key/1.jpg": "https://bucket.s3.amazonaws.com/s3/key/1.jpg?X-Amz-Algorithm=...",
    "s3/key/2.jpg": "https://bucket.s3.amazonaws.com/s3/key/2.jpg?X-Amz-Algorithm=..."
  }
}
```

**Error Responses**:

- `400 Bad Request`: Invalid JSON or empty keys array
- `401 Unauthorized`: Missing or invalid API key
- `500 Internal Server Error`: S3 signing failure

### GET /health

Health check endpoint (no authentication required).

**Response**:

```json
{
  "status": "healthy",
  "service": "presign-service"
}
```

### For Cloudflare R2

1. Set environment variables:

```bash
export BUCKET_NAME="my-r2-bucket"
export API_KEY="your-secret-api-key"
export AWS_ACCESS_KEY_ID="your-r2-access-key-id"
export AWS_SECRET_ACCESS_KEY="your-r2-secret-access-key"
export AWS_ENDPOINT_URL="https://your-account-id.r2.cloudflarestorage.com"
export AWS_REGION="auto"       # R2 uses "auto" region
export PORT="8080"             # optional
```

3. Run the service:

```bash
go run main.go
```

4. Test the service:

```bash
# Health check
curl http://localhost:8080/health

# Generate presigned URLs
curl -X POST http://localhost:8080/sign \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-secret-api-key" \
  -d '{"keys": ["path/to/file1.jpg", "path/to/file2.png"]}'
```

## Building

```bash
go build -o presign-service .
```
