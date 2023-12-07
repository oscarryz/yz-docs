[promises-Promise](https://stdlib.ponylang.io/promises-Promise/)

*Incomplete and miss understood* 
```js
Github {
	personal_access_token String
	get_repo  { String Promise }
}
Repository {
    get_issue { Int Promise }
}
Issue {
	title { String }
}
...

main: {
   repo Promise = Github('my token').get_repo('ponylang/ponyc')

   issue: repo.next(fetch_issue(1))
   issue.next(print_issue_title())

}
fetch_issue: {
   number Int
   repo.get_issue(number)
}
print_issue_title: {
   print...
}
```