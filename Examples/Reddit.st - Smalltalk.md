#example

[Reddit.st](http://web.archive.org/web/20110122114420/http://homepage.mac.com/svc/Reddit.st/)


```js
RedditLink:  {
    id Int
    url String
    title String
    create Date
    points Int
    init: {
        points: 0
    }
    print_on: {
        "(`url`),`title`)"
    }
    posted: {
        date.now() - created
    }
    vote_up: {
        points = points + 1
    }
    vote_down: {
        points = points - 1
    }
}
with_url_and_title: {
    url String
    title String
    RedditLink (
       url: url
       title: title
    )
    // or RedLink(url title)
}
render_content_on: {
    html Html
    html.anchor.callback({ count = count + 1}, '++')
    html.space()
    html.anchor.callback({ count + count -1 }, '--')
}

```
