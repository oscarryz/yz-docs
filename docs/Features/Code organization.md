#feature
# Code organization

## Module system invariants

The Yz module system is governed by five invariants (see `spec/09`):

1. **File = boc body** — the content of `foo/bar.yz` is the boc body of `foo.bar`
2. **Directory = namespace** — files in a directory compose sub-bocs; no conflicts possible
3. **Source root excluded from FQN** — `src/foo/bar.yz` defines `foo.bar`, not `src.foo.bar`
4. **`_name.yz` = infostring companion** — never a boc; content attached to the named boc's infostring slot
5. **File and directory coexist** — `http.yz` (own body) + `http/` (sub-bocs) both contribute to `http`

---

## Simple projects

For simple projects the `yz` build tool will compile each individual file and will create an executable if either they are named `main.yz` or have a `main` boc or free-floating code. You can also pass the filename to process a single file. If no entry point is found they are considered libraries and no executable would be created. This applies for subdirectories too.


Example: 
The directory contains 3 files `one.yz`, `two.yz` and `main.yz` all of them have free floating code and/or a main method 

```
project_name/
   one.yz
   two.yz
   main.yz   
```

When invoked without param the output would be:
```
yz 
proj_name/
   one.yz
   two.yz
   main.yz   
   one.exe
   two.exe
   main.exe
```
(`.exe` is just to denote an executable was created, but only applies to Windows platform)


## Larger projects

// TODO: The configuration structure is a `Compile` implementation 
If a `yz` file contains a `configuration` structure, then it will be used to create the executable. 
Also you can create this configuration structure by invoking the build tool `init` _`project_name`_  which can create additional folders 

```
yz init project_name
project_name/
   project_name.yz
   doc/
   src/
   lib/
   test/
   vendor/
```

In either case: manually creating the configuratin file or with the `init` command, a configuration file is a `.yz` file that  information like version, entry point and dependencies. It's tipically called the same name as the app (`project_name.yz`  in this example, but can be renamed to anything e.g. `mod.yz`, `manifest.yz`, `package.yz`, `project.yz`). There could be more than one and they can create different targets but most of the times there would be none or just one. 

### Example: a restaurant app

Let's say we want to create a `restaurant` app after invoking: 

`yz init restaurant`

The tool will create the following directories
```
restaurant/
   doc/
   src/
   lib/
   test/
   vendor/
```

Inside the configuration  is structured: 

File: `./restaurant.yz`

```js
 `
 compile-time: [Project]
 project: {
	 version: '0.1.0'
	 entry: 'main.yz'
	 src_path: ['./src/' './vendor/' './lib/']
	 vendor: ['./vendor']
	 dependencies: []
 }
 `
 restaurant: {
 ...
 }
```


Our app `restaurant` functionality would be as follows:  A restaurant opens the doors to the public, a menu is printed using a 3rd party library and start taking orders.

To add a dependency use `yz install printer`, the dependency will be added to the dependency section

```yz
dependencies: [
     printer: {version: "1.0.0"   url: 'https://example.org/print.git'} 
]
```

And available in the `vendor/` directory

Usually in the restaurant industry the restaurant is know as the house, and has a front where people eat and a back, where meal is prepared.

The following is the filesystem for this project:

```yz
restaurant/
    restaurant.yz
    src/
        main.yz
        house.yz
        house/
            front/
                Host.yz
    lib/
        some_other_lib.yz
    vendor/
        printer/Printer.yz
```


Our configuration has the `main.yz` as entry point, and `src` among the source path, thus the main entry point is  `src/main.yz`

```yz
// src/main.yz creates an implicit `main` block see filesystem resolution below
host: house.front.Host()
host.open_doors()
house.take_orders()
```
The `yzc` compiler will load all the files and dependencies defined in the project and will compile them into the `target/` directory using `main.yz` as entry point.


The rest of the app follows:

File: `src/house.yz`

```yz

// src/house.yz
facade: {
    Door:{
        open: {}
    }
}
take_orders: {
    // receive customers etc. 
}
```

File `src/house/front.yz`
```
main_door: facade.Door()
```

File: `src/house/front/Host.yz`
```yz
// src/house/front/Host.yz
menu String
open_doors: {
    printer: printer.Printer()
    menu = printer.print("Menu: Toast\nFruit\nCoffee")
    main_door.open() // Can access house.front.main_door directly because it's defined inside house.front block
}

```

The dependency `printer` was downloaded to the `vendor` directory which is listed in the `src_path` making it available for compilation. To re-download execute `yzc update`

File: `./vendor/printer.yz`
```yz
// vendor/printer.yz
print: {
    text String
    // do something with text
}

```

A block can be accessed using its fully qualified  block name: 

```javascript
//lib/some_other_lib.yz
other: {
    thing: {
        house.front.main_door.open() 
}

```


But is common to create aliases that work as "import" to avoid typing the whole name. 

```javascript
//lib/some_other_lib.yz
    other:{
        thing: {
            door: house.front.main_door // similar to "import"
            door.open()
        }
    }

```


### FQN resolution

The compiler resolves boc names using the `src_path` source roots. Each source root is a namespace anchor — its own name is NOT part of the FQN.

`src/house/front/Host.yz` → FQN `house.front.Host`
`vendor/printer.yz` → FQN `printer`

The boc's content is the file body:

```yz
// src/house/front/Host.yz — defines house.front.Host
menu String
open_doors: { ... }
```

### Summary

The boc structure:

```yz
house: {
    front: {
        Host: { name String }
    }
}
```

Is expressed in the filesystem as:

```
house/front/Host.yz
```

```yz
// house/front/Host.yz
name String
```

Or with an own body at each level:

```
house.yz          ← house body (facade, take_orders, ...)
house/
  front.yz        ← house.front body
  front/
    Host.yz       ← house.front.Host body
```
