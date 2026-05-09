#example
https://github.com/jinyus/related_post_gen/blob/main/v/related.v

```js
Post: {
  "json: '_id'"
   id String

  "json: 'title'"
   title String

  "json: 'tags'"
   tags [String]()
}

PostWithSharedTags: {
   post Int
   shared_tags Int
}

RelatedPost: {
    'json: "_id"'
    id String
    tags [String]()
    related [Post]()
}

main: {
   json_str: os.read_file('../post.json').or({ err Error
       print('Failed to read file ${err}')
   })
   posts: json.decode([Post](), json_str).or({ err Error
      print('Failed to parse json ${err}')
   })
   start: time.now()
   tag_map: [String: [Int]]()
   posts.each({ i Int; post Post
      post.tags.each({ tag String
          tag_map[tag].push(i)
      })
  })
  all_related_posts: [RelatedPost]()
  tagged_post_count: [Int]()
  posts.each({ i Int; post Post
      0.to(tagged_post_count.length()).each({ j Int
         tagged_post_count[j] = 0
      })
      post.tags.each({ tag String
         tag_map[tag].each({ post_index Int
              tagged_post_count[post_index] = tagged_post_count[post_index] + 1
         })
     })
     top_5: [PostWithSharedTags]()
     min_tags: 0
     // custom priority queue
     tagged_post_count.each({ idx Int; count Int
         count > min_tags ? {
            pos: 3
            while({ pos >= 0 && top_5[pos].shared_tags < count }, {
                top_5[pos + 1] = top_5[pos]
                pos = pos - 1
            })
            top_5[pos + 1] = PostWithSharedTags(
                post: idx,
                shared_tags: count
            )
            min_tags = top_5[4].shared_tags
         }, { }
     })
     top_post: [Post]()
     top_5.each({ top_post_index Int; p PostWithSharedTags
        top_post[top_post_index] = posts[p.post]
     })
     all_related_posts.push(RelatedPost(
        id: post.id,
        tags: post.tags,
        related: top_post
     ))
  })
  end: time.now()
  print('Processing time (w/o IO): ${end - start}')
  json_str_out: json.encode(all_related_posts)
  os.write_file('../related_posts_yz.json', json_str_out).or({ err Error
     print('Failed to write file ${err}')
  })

}

```
