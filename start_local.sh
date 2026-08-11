echo "Creating local environment"

echo "Starting up docker..."

docker compose up -d

echo "Waiting for 15 seconds..."
sleep 15
echo "Done waiting..."

echo "Creating AWS SQS Locally..."

aws --endpoint-url="http://localhost:4566" sqs create-queue \
    --profile localstack \
    --queue-name check-booking-pending-state \
    --region="sa-east-1"
