# Get OPG Metadata

Fetches the latest metadata tarball release from `opg-metadata` repository to get account data for using in matrix workflows in our reports.

As `opg-metadata` is an internal repository, you will need to pass in a github token with the correct permissions to theis action to fetch the matrix data correctly.

## Usage

Here is a short form of using this as a step:

```
- name: "Get account matrix"
  id: account_matrix
  uses: "./.github/actions/get-opg-metadata"
  with:
      github_token: ${{ secrets.ORG_ACCESS_TOKEN }}
      prereleases: false
```

Here is a more complete version of how to use this action in a wider context and how to feed the output of this step into another jobs matrix:

```
jobs:
  load_matrix:
    name: "Load matrix"
    runs-on: ubuntu-latest
    outputs:
      matrix: ${{ steps.account_matrix.outputs.matrix }}
    steps:
      - name: "Checkout"
        id: checkout
        uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 #v4.2.2
      - name: "get account matrix"
        id: account_matrix
        uses: "./.github/actions/get-opg-metadata"
        with:
          github_token: ${{ secrets.GH_ORG_ACCESS_TOKEN }}
          prereleases: false

  load_matrix_test:
    name: "Test matrix load"
    runs-on: ubuntu-latest
    needs: [load_matrix]
    strategy:
      fail-fast: false
      matrix: ${{ fromJson(needs.load_matrix.outputs.matrix) }}
    steps:
      - name: "test matrix values"
        env:
          unit: ${{ matrix.billing-unit }}
          label: ${{ matrix.label }}
          name: ${{ matrix.name }}
          id: ${{ matrix.id }}
          environment: ${{ matrix.environment }}
        run: |
          echo "name: ${{ env.name }}"
          echo "unit: ${{ env.unit }}"
          echo "label: ${{ env.label }}"
          echo "environment: ${{ env.environment }}"

```
