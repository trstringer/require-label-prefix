# GitHub Action - Require Label Prefix

[![CI](https://github.com/trstringer/require-label-prefix/actions/workflows/main.yaml/badge.svg)](https://github.com/trstringer/require-label-prefix/actions/workflows/main.yaml)

Use this GitHub action to either warn or add a default label based on prefixes.

## Usage

```yaml
steps:
  - uses: trstringer/require-label-prefix@v1
    with:
      secret: ${{ github.TOKEN }}

      # prefix is set to whatever prefix you are trying to enforce. For
      # instance, if you want to make sure size labels (e.g. "size/S", "size/L")
      # are enforced, the prefix would be "size".
      prefix: size

      # The prefix is divided by the suffix by some separator. This defaults
      # to "/" and is typically this, but it could be anything (e.g. ":").
      # labelSeparator: "/"

      # addLabel, when set to "true", will label the issue with defaultLabel if
      # the issue doesn't have a label with the prefix. If this is set to "false"
      # then a label won't be added, there will just be a comment requesting that
      # somebody adds a label with the labelPrefix.
      # Options: "true", "false" (default).
      # addLabel: false

      # If addLabel is set, defaultLabel is the label that will be added if there
      # is no label with this prefix already on the issue. E.g. "size/needed".
      # defaultLabel: "size/needed"

      # If you want to only comment on or label issues that are part of a milestone
      # then you would set this to "true". Otherwise, all issues are evaluated.
      # Options: "true", "false" (default).
      # onlyMilestone: false
```
