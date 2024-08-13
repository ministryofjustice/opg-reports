# Go Build

Build all go code. Please make sure to update the steps when more go is created.

## Usage

Here is a short form, typical usage of the action:

```
- name: "Build go code"
  id: "go_build"
  uses: ./.github/actions/go_build
```


Here is a complete example of the options that will create a tarball as well.

```
- name: "Build go and create artifact"
  id: "go_build"
  uses: ./.github/actions/go_build
  with:
    create_artifact: true
```

This action returns information about build folder locations and the tarball artficat (if `create_artifact` is true).
