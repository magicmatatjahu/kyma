# Rafter AsyncAPI Service

## Overview

The Rafter AsyncAPI Service is an HTTP server used to process AsyncAPI specifications. It contains the `/validate` and `/convert` HTTP endpoints which accept `multipart/form-data` forms:

- The `/validate` endpoint validates the AsyncAPI specification against the AsyncAPI schema in version 2.0.0., using the [AsyncAPI Parser](https://github.com/asyncapi/parser).
- The `/convert` endpoint converts the version and format of the AsyncAPI files.

To learn more about the Rafter AsyncAPI Service, see the [documentation](https://kyma-project.io/docs/components/rafter/#details-async-api-service).
