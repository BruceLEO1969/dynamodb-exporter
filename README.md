# AWS DynamoDB Prometheus Exporter

A Prometheus metrics exporter for AWS DynamoDB.

~~Further function is  yet to come.~~

## Metrics

| Metric  | Labels | Description |
| ------  | ------ | ----------- |
| dynamodb\_tablesCount | aws_region | The dynamodb count of current region. |

For more information see the [AWS DynamoDB Documentation](https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_Operations_Amazon_DynamoDB.html)

## Configuration

Credentials to AWS are provided in the following order:

- Environment variables (AWS\_ACCESS\_KEY\_ID and AWS\_SECRET\_ACCESS\_KEY)
- Shared credentials file (~/.aws/credentials)
- IAM role for Amazon EC2

For more information see the [AWS SDK Documentation](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)


## Environment Variables
| Variable      | Default Value | Description                                                  |
|---------------|:---------|:-------------------------------------------------------------|
| PORT          | 9436    | The port for metrics server                                  |
| ENDPOINT      | metrics  | The metrics endpoint                                         |



## Running

`docker run -d -p 9436:9436 bruceleo1969/dynamodb-exporter`

You can provide the AWS credentials as environment variables depending upon your security rules configured in AWS;

`docker run -d -p 9436:9436 -e AWS_ACCESS_KEY_ID=<access_key> -e AWS_SECRET_ACCESS_KEY=<secret_key>  bruceleo1969/dynamodb-exporter`

