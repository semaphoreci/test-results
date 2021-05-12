## ID generation

This PR introduces `id` generator for tests results, test suites and tests.

`id`'s are being generated in form of UUID strings.

In order to generate consistent `id`'s between builds following method is implemented for all parsers:

- ID generation for `Test results`(top level element)

  1. If element has an ID, generate UUID based on that ID
  2. If element doesn't have an ID - generate UUID based on the `name` attribute
  3. Otherwise, generate uuid based on empty string `""`

- ID generation for `Suites`

  Same rules apply as for `Test results`, however every `Suite ID` is namespaced by `Test results` ID

- ID generation for `Tests`

  Same rules apply as for `Test results`, however every `Test ID` is namespaced by `Suite` ID

For generating IDs we're using [UUID v3 generator](https://pkg.go.dev/github.com/google/uuid#NewMD5).