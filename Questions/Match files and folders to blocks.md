In project structurenI basically equated files to blocks, let's make sure that's what I want

```
// foo.yz
bar: {
  baz:{
      hi String
  }
}
```

Will be the same as

```
foo: {
   bar: {
      baz:{
         hi String
      }
   }
}
```
(Although now there wont be a way to define that block)

Also will be this
```
foo/bar.yz
baz:{
      hi String
}
```
and
```
// foo/bar/baz.yz

hi String
```
Furthermore, it would allow to add members in different ways

```
./foo.yz
./foo/bar.yz
./foo/bar/baz.yz
```

Wich should either not compile or allow to redefine blocks

#answered It won't compile and you'll be only allowed to define things in one place. Using folders however would be as if the curly braces hand close
