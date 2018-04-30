GOOS=linux GOARCH=amd64 go build -o handler .
zip step-001.zip handler
rm handler
aws s3 mv step-001.zip s3://grpc-interactive/step-001/handler.zip
aws lambda update-function-code --function-name grpc-interactive-step-001 --region us-west-2 --s3-bucket grpc-interactive --s3-key step-001/handler.zip
