# S3 Bucket

Create an s3 bucket that conforms to our organisations best practices.

Additional policies can be attached to the output of this module.

## Inputs

| Name                           | Description | Type     | Default | Required |
|--------------------------------|-------------|----------|---------|:--------:|
| bucket\_name                   | n/a         | `string` | n/a     |   yes    |
| force\_destroy                 | n/a         | `string` | false   |    no    |
| block\_public\_acls            | n/a         | `bool`   | true    |    no    |
| block\_public\_policy          | n/a         | `bool`   | true    |    no    |
| ignore\_public\_acls           | n/a         | `bool`   | true    |    no    |
| restrict\_public\_buckets      | n/a         | `bool`   | false   |    no    |
| kms\_key\_id                   | n/a         | `string` | n/a     |   yes    |
| enable\_lifecycle              | n/a         | `bool`   | false   |    no    |
| expiration\_days               | n/a         | `string` | `"365"` |    no    |
| non\_current\_expiration\_days | n/a         | `string` | `"14"`  |    no    |
| access\_logging_bucket         | n/a         | `string` | `""`    |    no    |

## Outputs

| Name   | Description                |
|--------|----------------------------|
| bucket | The created bucket object. |
