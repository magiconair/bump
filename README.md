<p align="center">
  <img src="logo-small.png" alt="BUMP logo" width="300">
</p>

# BUMP - A Go Version Nudger

`bump` is a tool for managing versions in git tags.

It parses all git tags for semantic version numbers, sorts
them and provides operations on them.

## Install

```
go install github.com/magiconair/bump@latest
```

## Usage

```
# print the current version number (highest tag)
bump cur

# list all version numbers sorted
bump list

# print the next version number
bump next

# print the next major/minor/patch version number
bump next major
bump next minor
bump next patch

# tag with the next version number
bump tag

# tag with the next major/minor/patch version number
bump tag major
bump tag minor
bump tag patch
```

## Service Tags

For monorepos with service-prefixed tags like `foo/v1.2.3`, use the `-s` flag
to scope operations to a specific service.

```
# print the current version for service foo
bump -s foo cur

# list all versions for service foo
bump -s foo list

# print the next minor version for service foo
bump -s foo next minor

# tag with the next minor version for service foo
bump -s foo tag minor
```

---

Copyright 2026 The Bump Authors. All rights reserved. See [LICENSE](LICENSE).
