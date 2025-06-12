```
./be/.env
----------
MINIO_HOST="localhost"
MINIO_PORT= 9000
MINIO_ACCESSKEY="admin"
MINIO_SECRETKEY="password"
MINIO_BUCKET="image-collection"
MINIO_BUCKET_CSV="csv-bucket"
----------
./fe/.env.local
NEXT_PUBLIC_BACKEND_URL="http://localhost:8080"
```
##### Setup minio before running fronend or backend
- docker compose up -d ./minio

