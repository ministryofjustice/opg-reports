# opg-reports

This repository acts as a central hub to generate, combine and display series of key data that we report on frequently to both internal and external parties.

Please see guides for more detailed background:

- [Development enviroment](./docs/usage/development-environment.md)

## Quick startup

As the codebase will auto-generate databases if they are not present, you can get a local version up and running in docker by:

```bash
make clean && make up
```

You can then view the site at:

- [front](http://localhost:8080)
- [api](http://localhost:8081)
