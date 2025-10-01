#example
```zig
const std = @import("std");
const json = std.json;
const payload =
    \\{
    \\    "vals": {
    \\        "testing": 1,
    \\        "production": 42
    \\    },
    \\    "uptime": 9999
    \\}
;
const Config = struct {
    vals: struct { testing: u8, production: u8 },
    uptime: u64,
};
const config = x: {
    var stream = json.TokenStream.init(payload);
    const res = json.parse(Config, &stream, .{});
    // Assert no error can occur since we are
    // parsing this JSON at comptime!
    break :x res catch unreachable;
};
pub fn main() !void {
    if (config.vals.production > 50) {
        @compileError("only up to 50 supported");
    }
    std.debug.print("up={d}", .{config.uptime});
}
```

```js
json: std.json
payload: `
        "vals": {
            "testing": 1,
            "production": 42
        },
        "uptime": 9999
  `
Configuration: {
   'json:vals'
   values: {
       'json:testing'
       tests Int
       
       'json:production'
       prod Int
   }

   'json:uptime'
   ut Int
}
config: {
   stream: json.TokenStream(payload)
   res : json.parse(Config,stream).or_fail!()
}

main: {
    if config.res.values.prod > 50 {
        print_err("Only up to 50 supported")
    }
    print("up=`config.ut`")
}
```

