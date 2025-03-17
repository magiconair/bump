# Overview

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
