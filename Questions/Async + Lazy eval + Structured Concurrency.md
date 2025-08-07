#open-question 

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

We might not even need the structured concurrency :/ 