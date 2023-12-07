From: 

https://gist.github.com/rexim/6326e2d762283af715b7cdb69239bd65

Video

```javascript
import core.collections.map

'derive: [Show]'
Group: {
  group_size Int
  greoup_name String
}

double_group { Group Group }
double_group = { g Group;  Group { 2 * g.size g.name } }

names []String
name = ['Sheldon' 'Leonard' 'Penny' 'Rajesh' 'Howard']

groups {[]String}
groups = {map new_group names ++ map double_group groups}

nth {Int []Group String}
nth = {
  n Int
  gs []Groups
  when [
    {gs.is_empty}: {strings.error}
    {n < gs[n].size}: {gs[n].name}
    {when.else}:{ nth (n - gs[n].size  gs.rest())
  ]
}

```