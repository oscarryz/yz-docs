---
type: example
generated: { by: oscarryz, at: 2026-11-18T12:00:00+01:00 }
tags: [module]
---

From the article [ML Modules for mare mortals](https://thunderseethe.dev/posts/ml-modules-for-mere-mortals/intro/)

```js
BuildStep #( Env #(
    cache Cache
    runner Runner
   )
) 
Cache #( 
    is_cached #(input [String], output Path, Bool)
    sore #(inputs [String], output Path)
)
Runner #(
   run #(input [String], output Path)
)
BuildSet  #(...) {
   run_cached #(inputs [String], output Path, Bool) {
      env.cache.is_cached(intput, output) ? {
         env.runner.run(inputs, output)
        env.cache.sore(intput, output)
      }
   }

    compile #(src String, String) { 
       obj : path_buf.from("${src}.o")
       run_cached([src], obj)
       obj.display().to_string()
   }

   build #(srcs [String], output Path) {
      objs [String] = srcs.map( compile )
      run_cached(objs, output)
   }
}
```
...      
