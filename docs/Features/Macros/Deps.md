#superseded
# Deps — superseded

Dependencies are not handled by a macro. See [Dependencies](../Dependencies.md) for the current design.

The original design treated `Deps` as a macro that satisfied the `Macro` interface and ran at compile time. This was cancelled (YZC-0041) because it coupled compilation to external network/filesystem operations, making builds unpredictable. Dependencies are now passive annotation metadata read by `yz fetch`, a standalone tool that runs before the compiler.
