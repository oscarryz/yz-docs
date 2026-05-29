#open-question 

> The following is a draft of what the stdlib could look like



Yz combines the visual clarity of Haskell signatures with structural constraints and an all-in-one conceptual block (**boc**), giving the language a clean, expressive, and powerful foundation.

Concurrency is completely transparent while explicit Rust-style sum types keep error paths visible: **business logic flows seamlessly, while failure paths remain completely visible.**

Below is the standard library API design in native Yz syntax, covering daily server and CLI workflows.

**1. fs (File System & Paths)**
**2. os (Processes & Environment)**
**3. http (Client & Server)**
**4. json (Structural Serialization)**
**5. str & regex (Text Processing)**
**6. cli (The Developer's UI)**
**7. time & crypto (Utilities)**
**8. log (Structured Logging)**
**9. db (Database)**
**10. encoding (TOML, YAML, CSV)**

---

**1. fs (File System & Paths)**

Since a boc can act as an object, module, and function simultaneously, path objects are created by invoking `Path` as a factory, returning a structural object with its own nested bocs.

```javascript
fs #(
    // Factory to create a Path object/boc
    Path #( path String, PathObj )
)
```

```javascript
// The structural interface of a PathObj
PathObj #(
    // Metadata & Checks
    exists    #( Bool ),
    is_dir    #( Bool ),
    size      #( Int ),

    // Fluent manipulation
    join      #( part String, PathObj ),
    parent    #( Option(PathObj) ),
    ext       #( Option(String) ),

    // High-level "Type Little, Do Much" I/O
    read_str  #( Result(String, FsError) ),
    read_bytes#( Result([Byte], FsError) ),
    write     #( content String, Result(FsError) ),
    append    #( content String, Result(FsError) ),

    // Traversal and cleanup
    walk      #( [PathObj] ), // Returns a lazily-evaluated array/stream
    delete    #( Result(FsError) )
)
```

---

**2. os (Processes & Environment)**

To allow experts to pipe commands fluidly like bash, Command returns a Process boc that can chain further actions.

```javascript
os #(
    env_get   #( key String, Option(String) ),
    env_set   #( key String, value String ),
    args      #( [String] ),

    // Process Spawning
    Command   #( program String, CmdObj )
)
```

```javascript
CmdObj #(
    arg       #( value String, CmdObj ),
    args      #( list [String], CmdObj ),
    env       #( vars [String:String], CmdObj ),

    // Execution (transparently scheduled by the Verona runtime)
    spawn     #( Result(ChildProcess, OsError) ),
    output    #( Result(ProcessOutput, OsError) ),
    pipe      #( next CmdObj, CmdObj )
)
```

```javascript
// Structural output record
ProcessOutput #(
    code   Int,
    stdout String,
    stderr String
)
```

---

**3. http (Client & Server)**

Using the generic array and hashmap syntax, parsing dynamic payloads is frictionless.

```javascript
http #(
    // High-level client fetch
    fetch     #( url String, options [String:String], Result(Response, HttpError) ),
    get       #( url String, Result(Response, HttpError) ),
    post      #( url String, body String, Result(Response, HttpError) ),

    // Web Server: Accepts a boc handler that maps a Request to a Response
    server    #( handler #(Request, Response), ServerObj )
)
```

```javascript
Response #(
    status    Int,
    headers   [String:String],
    text      #( Result(String, HttpError) ),
    json      #( T Serializable, Result(T, HttpError) )
)

ServerObj #(
    listen    #( port Int, Result(HttpError) )
)
```

---

**4. json (Structural Serialization)**

Since the language features structural typing, the JSON engine unpacks directly into any matching target shape without manual mapping code.

```javascript
json #(
    // T must satisfy the Serializable constraint
    parse     #( input String, T Serializable, Result(T, JsonError) ),
    stringify #( data T Serializable, Result(String, JsonError) ),

    // The "Escape Hatch" for dynamic unstructured JSON
    dynamic   #( input String, Result(JsonDynamic, JsonError) )
)
```

```javascript
JsonDynamic #(
    get       #( key String, Option(JsonDynamic) ),
    to_str    #( Option(String) ),
    to_int    #( Option(Int) )
)
```

---

**5. str & regex (Text Processing)**

In a high-density language, basic text transformations should live right on the String object/boc itself.

```javascript
// Native extensions on the String Boc
String #(
    trim      #( String ),
    split     #( sep String, [String] ),
    lines     #( [String] ),
    contains  #( pattern String, Bool )
)
```

```javascript
regex #(
    compile   #( pattern String, Result(RegexObj, RegexError) )
)

RegexObj #(
    match     #( input String, Option([String]) ), // Returns captured groups
    replace   #( input String, sub String, String )
)
```

---

**6. cli (The Developer's UI)**

An elegant syntax for terminal utilities makes scripting highly visual and fast to implement.

```javascript
cli #(
    print        #( msg String ),
    prompt       #( msg String, String ),
    confirm      #( msg String, Bool ),

    // Style and arguments
    color        #( msg String, ansi_code Int, String ),
    parse_flags  #( T Structural, Result(T, CliError) )
)
```

---

**7. time & crypto (Utilities)**

Essential operations that standard libraries often forget, forcing developers to look for external packages.

```javascript
time #(
    now       #( Timestamp ),
    sleep     #( ms Int ) // Yields the green thread transparently
)

Timestamp #(
    format    #( layout String, String ),
    diff      #( other Timestamp, Int ) // Returns milliseconds difference
)
```

```javascript
crypto #(
    sha256    #( input String, String ),
    md5       #( input String, String ),
    uuid_v4   #( String )
)
```

---

**8. log (Structured Logging)**

Basic text transformations live on `String` itself; structured logging lives in `log`. Every call appends a timestamped, level-tagged line. `with` returns a child logger that carries extra fields on every subsequent call — useful for request-scoped context.

```javascript
log #(
    debug  #( msg String ),
    info   #( msg String ),
    warn   #( msg String ),
    error  #( msg String ),

    // Returns a child logger with fixed key-value fields attached
    with   #( fields [String:String], Logger )
)
```

```javascript
Logger #(
    debug  #( msg String ),
    info   #( msg String ),
    warn   #( msg String ),
    error  #( msg String ),
    with   #( fields [String:String], Logger )
)
```

---

**9. db (Database)**

A single `connect` call returns a `Conn` boc. Queries return typed rows via structural scan, matching the same pattern as `json.parse`. Transactions are scoped to a handler boc — commit and rollback are automatic on success and error respectively.

```javascript
db #(
    connect #( dsn String, Result(Conn, DbError) )
)
```

```javascript
Conn #(
    query  #( sql String, args [Any], Result([Row], DbError) ),
    exec   #( sql String, args [Any], Result(ExecResult, DbError) ),
    tx     #( handler #(Tx, Result(DbError)), Result(DbError) ),
    close  #()
)

Tx #(
    query    #( sql String, args [Any], Result([Row], DbError) ),
    exec     #( sql String, args [Any], Result(ExecResult, DbError) ),
    commit   #( Result(DbError) ),
    rollback #( Result(DbError) )
)

Row #(
    scan #( T Structural, Result(T, DbError) )
)

ExecResult #(
    rows_affected   Int,
    last_insert_id  Int
)
```

---

**10. encoding (TOML, YAML, CSV)**

Config files are typically TOML or YAML; tabular data is typically CSV. All three use the same structural parse/stringify pattern as `json`, so any matching boc shape can be a target without manual mapping.

```javascript
toml #(
    parse     #( input String, T Structural, Result(T, EncodingError) ),
    stringify #( data T Structural, Result(String, EncodingError) )
)

yaml #(
    parse     #( input String, T Structural, Result(T, EncodingError) ),
    stringify #( data T Structural, Result(String, EncodingError) )
)

csv #(
    // Returns raw rows as arrays of strings
    parse   #( input String, Result([[String]], EncodingError) ),
    // Scans each row into a structural type by column order or header name
    records #( input String, T Structural, Result([T], EncodingError) ),
    write   #( rows [[String]], Result(String, EncodingError) )
)
```

Binary encoding helpers live here too:

```js
base64 #(
    encode #( input [Byte], String ),
    decode #( input String, Result([Byte], EncodingError) )
)

hex #(
    encode #( input [Byte], String ),
    decode #( input String, Result([Byte], EncodingError) )
)
```

---

**How an Expert Workflow Looks in Yz**

A minimal script — fetch weather data, format it, append to a log file:

```javascript
// 1. Fetch data structurally
result: http.get("https://weather.com").?({ e HttpError; return e })
            .json(WeatherReport).?({ e HttpError; return e })

// 2. Format a log entry
log_entry: "Status is ${result.status} at ${time.now().format('HH:MM')}\n"

// 3. Append to disk safely
fs.Path("/var/log/weather.log")
  .append(log_entry).?({ e FsError; return e })
```

A fuller script — load config, fetch weather, persist to a database, and emit structured logs throughout:

```javascript
// 1. Load typed config from TOML
AppConfig #(
    weather_url  String,
    db_dsn       String,
    env          String
)

config: toml.parse(
    fs.Path("config.toml").read_str().?({ e FsError; return e }),
    AppConfig
).?({ e EncodingError; return e })

// 2. Scope a structured logger to this run
logger: log.with(["service": "weather-sync", "env": config.env])
logger.info("starting fetch")

// 3. Fetch and parse weather report
report: http.get(config.weather_url).?({ e HttpError; return e })
            .json(WeatherReport).?({ e HttpError; return e })

logger.info("report received")

// 4. Persist to database inside a transaction
conn: db.connect(config.db_dsn).?({ e DbError; return e })

conn.tx({ tx Tx
    tx.exec(
        "INSERT INTO reports (status, summary, fetched_at) VALUES (?, ?, ?)",
        [report.status, report.summary, time.now().format("RFC3339")]
    ).?({ e DbError; return e })
}).?({ e DbError; return e })

logger.info("persisted report")

// 5. Append human-readable line to the on-disk log
fs.Path("/var/log/weather.log")
  .append("${report.status} at ${time.now().format('HH:MM')}\n")
  .?({ e FsError; return e })

logger.info("done")
```



