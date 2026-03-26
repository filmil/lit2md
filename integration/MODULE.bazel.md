# SPDX-License-Identifier: Apache-2.0
Declare the module, this seems to be a requirement.

```python
module(
    name = "lit2md_integration",
    version = "0.0.0",
)

```
Required for `write_source_files`. Refer to [`README.md`](./README.md) for
details.

```python
bazel_dep(name = "aspect_bazel_lib", version = "2.21.1")

```
This brings in `lit2md`. Be sure to pick the most recent version.

```python
bazel_dep(name = "lit2md", version = "0.2.0")

```
You can omit this part: it is only used in this integration repository so that
it would always use the most recent version of `lit2md`.

```python
local_path_override(
    module_name = "lit2md",
    path = "../",
)

```
