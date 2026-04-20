#answered  Use it

There is an [interesting suggestion](https://www.reddit.com/r/ProgrammingLanguages/comments/1m4sse2/comment/n51yxkh/) about how Yz could handle concurrency: All the functions to be async, but the value returned (not a promise) to be "thunk" that can be passed around until it is used. 

```js
load_profile: {
   // asynchronously 
   // loads the user
   // and returns a 
   // lazy `usr`
   usr: fetch_user('Mr Pink')
   // async loads an image
   // and returns a lazy `img`
   img: fetch_image('123-123')

   // do something else 
   p: Profile(usr, img) // returs p 
}
// call the function
profile: load_profile()
// at this point the execution can continue because we're not using
// `profile`
name : profile.name()
/// We are using profile, then we block until both `usr` and `img` complete

```

~~We might not even need the structured concurrency :/~~

Update: One difficulty with this approach is the lack of explicit control, for instance, a thunk is resolved when needed but it might not be obvious when that happens. As an alternative we're looking at having a value referece/dereference syntax using a `*` just like a pointer. 


```js
// p is a pointer to Profile, a value that might still
// not loaded
p *Profile = load_profile()
// can be passed
do_something(p)
q : p // aliased
arr : [p, q] // stored

// And when finally needed dereferenced
name : *p.name // using `*` _varname_ will lock the execution until the value is loaded
```

Added [more here][1]:

[1]: https://www.reddit.com/r/ProgrammingLanguages/comments/1n0zd8o/lazyish_evaluation_with_pointerish_syntax_idea/

It seems using the pointer wouldn't be needed and might just complicate things, for instance `String` and `*String` would be different data types, and while this might add clarity it could also be overkill. 



#TODO 
- Update [Concurrency](../../../Features/Concurrency.md)
- Go through all the concurrency examples and answer questions such as blocking, signaling, structured concurrency control etc. 

