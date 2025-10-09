#example

https://github.com/golang/go/blob/f32519e5fbcf1b12f9654a6175e5e72b09ae8f3a/src/cmd/go/internal/work/gc.go#L68

```js
Gc_tool_chain: {

    gc: {

        b Builder;
        a Action
        archive String
        importcfg [B]()yte
        embedcfg []byte
        symabis String
        asmhdr Bool
        gofiles [String]()
        ofile String
        output [B]()yte
        error Error

        p: a.package
        objdir: a.objdir

        archive != '' ? {
            ofile = archive
        } {
            out: '_go_.o'
            ofile = objdir + out
        }
         pkgpath: pg_path a
         default_gc_flags : ["-p", pkgpath]
         p.module != nil ? {
             v: p.module.goversion
             v == '' ? {
                // We started adding a 'go' directive to the go.mod file unconditionally
                // as of Go 1.12, so any module that still lacks such a directive must
                // either have been authored before then, or have a hand-edited go.mod
                // file that hasn't been updated by cmd/go since that edit.
                //
                // Unfortunately, through at least Go 1.16 we didn't add versions to
                // vendor/modules.txt. So this could also be a vendored 1.16 dependency.
                //
                // Fortunately, there were no breaking changes to the language between Go
                // 1.11 and 1.16, so if we assume Go 1.16 semantics we will not introduce
                // any spurious errors — we will only mask errors, and not particularly
                // important ones at that.
                v = "1.16"
             }
             allowed_version(v) ? {
                 default_gc_flags << '-lang=go' + v
             }
         }
         p.standard ? {
             default_gc_flags << '-std'
         }

        compiling_runtime: runtime_packages[p.import_path]
        compiling_runtime = compiling_runtime && p.standard
        compiling_runtime ? {
            // runtime compiles with a special gc flag to check for
            // memory allocations that are invalid in the runtime package,
            // and to implement some special compiler pragmas.
            default_gc_flags << '-+'
        }
        // If we're giving the compiler the entire package (no C etc files), tell it that,
        // so that it can give good error messages about forward declarations.
        // Exceptions: a few standard packages have forward declarations for
        // pieces supplied behind-the-scenes by package runtime.
        ext_file: p.cgo_files.len()
            + p.cfiles.len()
            + p.cxxfiles.len()
            + p.mfiles.len()
            + p.ffiiles.len()
            + p.sfiles.len()
            + p.sysofiles.len()
            + p.swigfiles.len()
            + p.swigcxxfile.len()
        p.standard ?  {
            ["bytes"
             "internal/poll"
             "net"
             "os"
             "runtime/metrics"
             "runtime/pprof"
             "runtime/trace"
             "sync"
             "syscall"
             "time"].contains(p.import_path) ? {
                  ext_files = ext_files + 1
              }
        }
        ext_files == 0 ? {
            default_gc_flags << '-complete'
        }
        cfg.build_context.install_fussix != '' ? {
            default_gc_flags << '-installfuffix'
            default_gc_flags << cg.build_context.install_suffix
        }
        a.build_id != '' ? {
            default_gc_flags << '-buildid'
            default_gc_flags << a.build_id
        }
        p.internal.omit_debug || { cfg.goos == 'plan9'} || {cfg.goarch == "wasm"} ? {
            default_gc_flags << '-dwarf=false'

        }
        runtime_version.has_prefix "go1" && {os.args[0].contains 'go_bootstrap'} ? {
            default_gc_flags << '-goversion'
            default_gc_flags << runtime_version
        }
        symabis != '' ? {
            default_gc_flags << '-symabis'
            default_gc_flags << symabis
        }

        gcflags:  compiling_runtime ? {
            // Remove -N, if present.
            // It is not possible to build the runtime with no optimizations,
            // because the compiler cannot eliminate enough write barriers.
            gcflags.remove('-N')
            i: 0
            while { i < gcflags.len() } {
                gcflags[i] == '-N' ? {
                    gcflags = gcflags.sub_array(0 i) ++ gcflags.subarray(i+1)
                }
                i = i + 1
            }
        }
    }
}
```
