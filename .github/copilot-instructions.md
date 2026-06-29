# AI Review Guidelines — `akamaiopen-edgegrid-golang`

## Meta Information for AI Agents

**Review scope**: AI agents must always review the **entire pull request in full**,
not just the incremental diff. This means:
- Read and consider all files in the repository that are relevant to the changes,
  not only the lines added or modified.
- Evaluate the impact of changes in the context of the whole codebase.
- Ensure consistency, correctness, and adherence to project conventions across the
  full scope of the PR, including unchanged surrounding code and its dependencies.

---

This file is the source of truth used by AI review agents when reviewing pull requests
in the `akamaiopen-edgegrid-golang` repository. Every rule below MUST be checked for
the changed code (additions and modifications). Do NOT flag pre-existing issues in
code that is not part of the change unless they are directly touched by the PR.

Each rule has a stable `id` so review tools can reference it.

---

## 1. General practices (apply to every change)

- **GEN-01**: If newly added functionality is not yet Globally Available (GA), the code
  must contain a comment stating when it is expected to reach GA.
- **GEN-02**: When a change is part of a multi-repository change (e.g. requires a
  matching change in `terraform-provider-akamai` or `cli-terraform`), the branch
  name in this repo must be exactly the same as in the other repositories.
- **GEN-03**: If the underlying API is still changing, new major features may be
  introduced as Beta. While in Beta, breaking changes on that feature are allowed
  without a breaking-change release. Beta status must be explicit in code/docs.
- **GEN-04**: A PR that changes exported types, functions, methods, variables, or
  constants in a way that forces end-users to manually adjust their code is a
  breaking change (unless the feature is still in Beta). It must target the
  appropriate `sp-breaking-changes` branch. These are released roughly once a
  quarter or less often. When suspecting breaking changes in the PR, let someone
  from the DEVEXP team know to help identify the correct `sp-breaking-changes`
  branch.
- **GEN-05**: Export types, variables, methods, and constants ONLY if they should
  be exposed to the customer. Default to unexported.
- **GEN-06**: Do not include methods, structs, fields, constants, or files that
  are not used anywhere.
- **GEN-07**: If new methods are being added together with new Terraform
  resources/data sources or new cli-terraform exports, the PR description must
  state clearly that the related PRs must be released together.

## 2. Common coding practices

- **COD-01**: Any change affecting customers must be reflected in `CHANGELOG.md`.
- **COD-02**: No secrets, real contract IDs, real Akamai employee IDs, or any
  other non-public data may appear anywhere in the code (especially unit tests
  and fixture JSON).
- **COD-03**: Never skip error checking. Every returned `error` must be handled.
- **COD-04**: Do not add files that are never used (frequently happens with unit
  test fixtures).
- **COD-05**: Prefer `any` over `interface{}`.
- **COD-06**: Acronyms in exported identifiers (e.g. `IP`, `DNS`, `URL`, `ID`,
  `API`) must be in all uppercase. This rule does not apply to unexported
  identifiers.
- **COD-07**: Descriptions and comments that are valid sentences must end with a
  period.
- **COD-08**: Unit tests must cover all corner cases — especially presence and
  absence of fields. Such cases may be aggregated where it makes sense.

## 3. Changelog rules

- **CHG-01**: Entries must describe what changed from the customer perspective of
  THIS project (edgegrid-golang). Do not describe the effect on TFP or cli.
- **CHG-02**: Fixes/changes based on GitHub issues must include a link to that
  issue in the changelog entry, e.g.
  `([#436](https://github.com/akamai/terraform-provider-akamai/issues/436))`.
- **CHG-03**: Use past tense.
- **CHG-04**: Use backticks (`` ` ``) around proper names (types, fields, methods,
  packages, etc.).
- **CHG-05**: Entries should be placed in a random line position (within the
  correct section / correct release) to mitigate merge conflicts.
- **CHG-06**: Do NOT delete empty lines in the changelog — they exist to mitigate
  conflicts.

## 4. PR hygiene (informational checks)

- **PRH-01**: The PR must not be in "Ready for review" state if the build is not
  passing or if it is still WIP — should be a draft instead.
- **PRH-02**: Branch must be rebased on top of its target branch before opening
  the PR for review.
- **PRH-03**: Commits must be squashed wherever it makes sense.
- **PRH-04**: Every commit message must include the appropriate JIRA number AND
  the story title or a description of the change.

---

## 5. Methods (interface methods of each package)

- **MTH-01**: Each package interface method must be a wrapper around exactly ONE
  OpenAPI call. Combining multiple calls in a single interface method is not
  allowed. Such combined methods belong in `terraform-provider-akamai` or
  `cli-terraform`.
- **MTH-02**: Preferably ALL methods from the given OpenAPI package/family are
  covered. Missing methods should be called out.
- **MTH-03**: Each method must have a working link to public techdocs. The link
  must appear:
    - in the changelog entry, AND
    - in the doc comment immediately above the method.
  If the link is not live yet, a comment must state when it will be working.
- **MTH-04**: When introducing a whole new subprovider package, the boilerplate
  (interface, mocks, client, session) must match the current library structure.
- **MTH-05**: Each interface method must accept at most TWO arguments:
  `context.Context` and optionally a `Request` struct (which may itself contain a
  nested `Body`). It must return at most TWO values, with `error` as the second.
- **MTH-06**: Each interface method must:
    1. Log the method name first.
    2. Then run validation (if applicable).
    3. Then build the request.
    4. Then execute it.
    5. Then check the response status and return either the result or an error.
    6. The response body must be closed.
- **MTH-07**: When checking for a successful response, only the HTTP status
  code(s) that the API actually returns on success may be allowed (usually a
  single code).
- **MTH-08**: New methods must use the requests builder located in
  `akamaiopen-edgegrid-golang/internal/request/request.go`.
- **MTH-09**: GET endpoints returning a list must be named `List...`; GET
  endpoints returning a single element (even a complex one) must be named
  `Get...`.
- **MTH-10**: At least one failing test case for an API error must exist for each
  method.
- **MTH-11**: When wrapping errors, use `%w` for each wrapped error — do not mix
  `%s` and `%w`.
- **MTH-12**: If the API returns non-standard response headers (e.g. rate-limit
  metadata), expose them as dedicated fields on the response struct.
- **MTH-13**: Each new package must include a `TestClient` test (in
  `<package>_test.go`) that tests the `Client` constructor and any options it
  accepts.
- **MTH-14**: Consider defining exported sentinel errors (`ErrXxx`, e.g.
  `ErrNotFound`) for each major failure case so that callers can use
  `errors.Is` to differentiate between error types. This is preferred over
  forcing downstream projects to inspect HTTP status codes or parse error
  message strings. These sentinel errors must be defined in the `errors.go`
  file. When declaring such an error, a unit test must verify that it can be
  matched using `errors.Is`.
- **MTH-15**: Each new package must define an `errors.go` file that specifies
  how the package handles error processing.
- **MTH-16**: The `Error` struct defined in `errors.go` must declare all
  possible fields that the API can return in an error response.

## 6. Structs (requests and responses)

- **STR-01**: All fields available for an OpenAPI endpoint must be covered in
  BOTH requests and responses (where the API exposes them).
- **STR-02**: `AccountSwitchKey` handling is done at the library core
  (`pkg/edgegrid/signer.go`) and must NOT be reintroduced at the method level.
- **STR-03**: Top-level request and response structs must be named after the
  method, with the appropriate suffix (e.g. `Activate` and `ActivateResponse`).
- **STR-04**: In responses: fields returned only sometimes must be pointers
  (booleans require extra care). Fields never returned must NOT be in the
  response struct. In requests: fields never accepted by the API must NOT be in
  the request struct.
- **STR-05**: A struct may be reused between responses of multiple methods only
  if all fields are returned by all methods and no extra fields exist for either.
- **STR-06**: A struct may be reused for request + response only if every field
  is both provided and returned, nothing more can be provided, and nothing more
  can be returned.
  When in doubt whether a struct reuse satisfies STR-05 or STR-06, flag it for
  human review.
- **STR-07**: When two structs are needed (request and response), the response
  must carry the suffix (e.g. `Resp` / `Response`); the request stays bare. For
  example: `Activate` and `ActivateResponse`.
- **STR-08**: Do not use JSON tags on fields that are not being marshalled
  (especially URL query/path parameters).
- **STR-09**: Field types must match the API schema. Prefer `int64` over `int`
  (avoiding 32-bit issues). Date/time fields must be `time.Time`, not `string`.
- **STR-10**: Alias types and enums for fields may be used ONLY in requests. In
  responses, the library must accept whatever the server returns, so do not use
  them.
- **STR-11**: Struct embedding is allowed ONLY in responses, and only sparingly —
  it adds overhead for customers.
- **STR-12**: All fields and structs must be documented.
- **STR-13**: Do not use pointers to slices or maps — slices and maps can
  already hold `nil`.
- **STR-14**: Do not use anonymous structs inside other types.
- **STR-15**: For readability, separate each field (and its doc comment) with a
  newline from the next field, instead of one solid block of declarations.
- **STR-16**: If the endpoint has a known request/response schema, define the
  matching struct. Avoid `interface{}` / `any` / `json.RawMessage` when the type
  is known.
- **STR-17**: If a method has BOTH body and path/query parameters, there must be
  TWO structs: one for the Request and one for the Body, with the Body being a
  field of the Request. Structure of that Body should have `Body` suffix.

## 7. Usage of `omitempty`

- **OMT-01**: Do NOT use `omitempty` in structs that are used ONLY in responses.
- **OMT-02**: If a struct is used in BOTH request and response, prefer a single
  struct with `omitempty` over two separate structs with the same fields where
  one uses `omitempty` (for requests) and the other does not (for responses).
  Note: OMT-01 applies only to structs used *exclusively* in responses; this
  rule applies when the same struct serves both request and response.
- **OMT-03**: To distinguish `{"g": ""}` from `{}` / `{"g": null}` in a response,
  use a pointer type (`*string`, `*bool`, ...) — `omitempty` does not affect
  unmarshalling.
- **OMT-04**: For requests, use `omitempty` according to API behaviour: usually
  add it to avoid sending Go zero values. For boolean fields, if the API
  distinguishes "absent" from `false`, do NOT use plain `bool` with
  `omitempty` — consider `*bool` with `omitempty` instead.

## 8. Validations

- **VAL-01**: Every type (at any nesting level — slices, structs, etc.) used in
  a request or body that has easily verifiable values must implement a
  `Validate` method. "Easily verifiable" means: a list of allowed values,
  required-vs-optional check, format check, etc. It is NOT calling another
  endpoint to check.
- **VAL-02**: Library validation must not be stricter than API validation.
- **VAL-03**: When nested structures (e.g. Body) need validation, the nested
  struct must have its own `Validate` method.
- **VAL-04**: Top-level validation must use
  `edgegriderr.ParseValidationErrors(...)`. Nested-level validation (typically
  Body or its elements) must use `validation.Errors{}.Filter()`. This is
  required for consistent error formatting.
- **VAL-05**: For every `Validate` method there must be a negative test that
  checks the resulting error message (or its key fragment). One test may cover
  several validation checks (e.g. all required fields missing) as long as the
  exact error message is asserted.
- **VAL-06**: Validation must use the currently adopted library: `ozzo-validation`.

## 9. Tests

- **TST-01**: Each interface method must have tests that, given a request as a
  Go struct and a mocked JSON server response, verify:
    - the correct Go response struct is returned (or the correct error),
    - the correct URL and HTTP method were used,
    - for body-carrying methods (POST/PUT/PATCH), the correct JSON request body
      was sent,
    - any non-standard response headers that the method processes are checked.
- **TST-02**: One test method per interface method, using a parametrised
  (table-driven) approach. Moving validation checks into a dedicated test method
  is also acceptable. See the `cloudcertificates` package for a reference.
- **TST-03**: Avoid duplicated tests.
- **TST-04**: Validation tests are exempt from rule TST-01 (no server mock
  needed). See also VAL-05 for validation-specific test requirements.
- **TST-05**: `t.Parallel()` must be present wherever applicable.

