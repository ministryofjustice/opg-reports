# Docker build, scan & push

Uses `buildx` to create a local image, scan that image using trivy and then push the image to an AWS ECR container registry.

## Usage

Here is a short form, typical usage of the action:

```
- name: "Build and push image"
  id: build_image
  uses: ./.github/actions/docker_build_scan_push
  with:
      ecr_repository: "opg-reports/api"
      branch_name: "safe-branch-name"
      tag: "v0.11.1"
      dockerfile: './docker/api/Dockerfile'
      scan_image: true
```

Here is a more complete version of using this action with all input options setup:

```
- name: "Build and push image"
  id: build_image
  uses: ./.github/actions/docker_build_scan_push
  with:
      branch_name: "safe-branch-name"
      tag: "v1.12.1"
      dockerfile: "./docker/api/Dockerfile"
      buildx_version: "v0.15.1"
      buildx_platforms: "linux/amd64"
      ecr_repository: "opg-reports/api-test"
      ecr_registry_id: "311462405658"
      checkout_repo: false
      scan_image: false
      path_to_live: true
```
