#example

[promises-Promise](https://stdlib.ponylang.io/promises-Promise/)

```js
Github : {
	personal_access_token String
	get_repo #( String, Promise(Repository))
}

Repository : {
	get_issue #(number Int, Promise(Issue))
}

Issue : {
	title #(String)
}

main: {
	repo Promise(Repository) = Github("my token").get_repo("ponylang/ponyc")
	//
	// Do something with the repo once the promise is fulfulled
	// in our case, get the issue
	//
	repo.then { r Repository
		fetch_issue( 123, r ).then { i
			print_issue_title(issue)
		}
	}
}
fetch_issue #(number Int, repo Repository, Promise(Issue)) = {
	repo.get_issue(number)
}
print_issue_title #(issue Promise(Issue)) = {
	issue.then {
		i Issue
		print("Issue: `i.title()`")
	}
}
...
Promise : {
  T
  next
  flatten_next
  add

}
```
