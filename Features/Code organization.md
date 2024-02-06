# Code organization

## Simple projects

For simple project usually placing all the files in a single folder is enough, `yzc` will look for a top level block called `main` and will use it as a entry point. If there's more than one you can to specify the filename as parameter.

```
# Both contain a `main` block
yzc one.yz
yzc two.yz
yzc main.yz // implicit main block

```

## Larger projects

For larger project `yzc` can create a project with the following structure after executing `yz init project_name`
```
project_name/
   project_name.yz
   doc/
   src/
   lib/
   test/
   vendor/
```

At the root, a single `.yz` file acts as the project meta information, it contains information like version, entry point and dependencies among other information. It's tipically called the same name as the app (`project_name.yz`  in this example, but can be renamed to anything e.g. `mod.yz`, `manifest.yz`, `package.yz`, `project.yz`). 

Inside the information is structured as Yz block: 

File: `./restaurant.yz`

```javascript
 version: '0.1.0'
 entry: 'main.yz'
 src_path: ['./src/' './vendor/' './lib/']
 dependencies: [
     printer: {version: "1.0.0"   url: 'https://example.org/print.git'} 
 ]
```

### Example: a restaurant app

Our app `restaurant`, opens the doors to the public, a menu is printed using a 3rd party library and start taking orders.
Usually in the restaurant industry the restaurant is know as the house, and has a front where people eat and a back, where meal is preprared.

The following is the filesystem for this project:

```
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



When run inside the `restaurant/` direvtory the `yzc` compiler takes the first `.yz` file that matches the [project](./project.md) structure ( `restaurant.yz` in this example ) and looks for the `entry` variable in the the `src_path` list to create the application entry point. If there is more than one `.yz` file in the root, the project file has to be specified e.g.:`yzc -project something_else.yz` 

In our example this is the file  `src/main.yz`

```javascript
// src/main.yz creates an implicit `main` block see filesystem resolution below
host: house.front.Host()
host.open_doors()
house.take_orders()
```
The `yzc` compiler will load all the files and dependencies defined in the project and will compile them into the `target/` directory using `main.yz` as entry point.


The rest of the app follows:

File: `src/house.yz`

```javascript

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
```javascript
// src/house/front/Host.yz
menu String
open_doors: {
    printer: printer.Printer()
    menu = printer.print "Menu: Toast\nFruit\nCoffee" 
    main_door.open() // Can access house.front.main_door directly because it's defined inside house.front block
}

```

The dependency `printer` was downloaded to the `vendor` directory which is listed in the `src_path` making it available for compilation. To re-download execute `yzc update`

File: `/vendor/printer.yz`
```javascript
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


### Filesystem Block name resolution

The compiler resolves block names to filesystem files using the `src_path` variable ( `src`, `lib`, `vendor` in this example ) using the following strategy:

 1. File name including subdirectories eill create a block, even if is empty excluding directories defined in the project's `src_parg` variable
 
 `src/house/front/Host.yz` will create the `house.front.Host` block
 ```
// src/house/front/Host.yz
// creates the empty type Host{}

```
`vendor/printer.yz`  creates/matches `printer`
```
//vendor/printer.yz
print:{} // creates the `printer:{ print:{}}` block


```
 


### Summary

Files and subdirectories can be used to create blocks and nested blocks for better code organization. 

For instance consider the following block structure:

```javascript
house: {
    front: {
        Host {
           name String
        }
    }
}
```

Could be defined in several ways in the file system: 

- Matching file name
```javascript
//house.yz
// Top level block is not `house` thus it gets defined inside `house`
front: {
    Host { // fqn house.front.House
        name String
    }
}

```

- Using directories as block names 
```javascript
// house/front/Host.yz
name String

```

#open-question How can we add new types and variables to an existing type if the type was defined in a file format vs inside of another, eg

```
// house/front/Host.yz
Host {
    name String
}
// vs 
// house/Front.yz
Front {
    House {
        name String
    }
}
```
When this would throw `already defined` or something like that?
_temptative answer_ yes it will be a compilation error but you can still add new members with a day fferent name.

#to-do Clarify only non instantiable blocks can be matched to filesystem directories.
