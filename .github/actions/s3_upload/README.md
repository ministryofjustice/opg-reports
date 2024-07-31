# S3 Upload

Uploads the content of a directory into a bucket using a bucket_path.

Please ensure any authentication (AWS creds, OIDC etc) required to upload to the bucket is done before calling the action

## Usage

Here is a typical usage of the action where a folder is uploaded to the s3 bucket.

```
- name: "Upload to s3"
  uses: ./.github/actions/s3_upload
  with:
    directory: ${{ steps.run.outputs.data_folder }}
    bucket: "report-data-development"
    bucket_path: "github/standards"

```
