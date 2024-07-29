# Go Build

Use `make` to build the relevant go code which will also generate an a tarball of all of the binaries generated to be used as a release artifact.

the value of target must relate to a make target as underneath this calls:

```
make go-<target>
```

If you add more `go` commands to this report tool, then you will need to update the `go-*` targets in the root `Makefile`.


## Usage

Here is a short form, typical usage of the action:

```
- name: "Build go code"
  id: "go_build"
  uses: ./.github/actions/go_build
  with:
    target: "all"
```

Here is an example that only builds the api service:

```
- name: "Build go api service"
  id: "go_build"
  uses: ./.github/actions/go_build
  with:
    target: "api"

```

Here is a complete example of the options that will create a tarball as well.

```
- name: "Build go and create artifact"
  id: "go_build"
  uses: ./.github/actions/go_build
  with:
    target: "all"
    create_artifact: true
```

This action returns information about build folder locations and the tarball artficat (if `create_artifact` is true).
