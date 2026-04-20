#example
[Howto write a simple application](https://blog.vosca.dev/how-to-write-a-simple-v-application-step-by-step/)

```js
github_repositories_url: 'https://api.github.com/search/repositories?sort=stars&order=desc&q=language:v'

GitHubRepositoriesSearchAPI: {
	total_count Int
	items       [GitHubRepositoriesItem]()
}

GitHubRepositoriesItem: {
	full_name        String
	description      String
	stargazers_count Int
	html_url         String
}

main: {
    response: http.get(github_repositories_url).or_else(panic)

    repositories_result: json.decode(GitHubRepositoriesSearchAPI, response.body).or_else({ err Err
	    panic('An error occurred during JSON parsing: `err`')
    })

    print('The total repository count is `repositories_result.total_count`')

    repositories_result.items.each({ index Int; item GitHubRepositoriesItem

	   colored_description: chalk.fg(item.description, 'cyan')
	   colored_star_count: chalk.fg(item.stargazers_count.str(), 'green')

	   print('#`index + 1` `item.full_name`')
	   print('  URL: `item.html_url`')

       item.description != '' ? {
		  print('  Description: `colored_description`')
	   }, { }

	    print('  Star count: `colored_star_count`')
    })
    print(response.body)
}

```
