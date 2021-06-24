# Contributing

We welcome contributions from the community. Please read the following guidelines carefully to maximize the chances of your PR being merged.

## Coding Style

The code is linted using a stringent golang-ci. To run this linter (and a few others) use run `make check`. To format your files, you can run `make format`.

## Running tests

```
# Run local tests without running envoy processes.
make test

# Run all e2e tests.
# This requires you to have Envoy binary locally.
make test.e2e

# Run e2e tests for a specific example.
# This requires you to have Envoy binary locally.
make test.e2e.single name=helloworld
```

## Code Reviews

* Indicate the priority of each comment, following this
[feedback ladder](https://www.netlify.com/blog/2020/03/05/feedback-ladders-how-we-encode-code-reviews-at-netlify/).
If none was indicated it will be treated as `[dust]`.
* A single approval is sufficient to merge, except when the change cuts
across several components; then it should be approved by at least one owner
of each component. If a reviewer asks for changes in a PR they should be
addressed before the PR is merged, even if another reviewer has already
approved the PR.
* During the review, address the comments and commit the changes _without_ squashing the commits.
This facilitates incremental reviews since the reviewer does not go through all the code again to
find out what has changed since the last review.
