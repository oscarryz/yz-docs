#example

From:

https://gist.github.com/rexim/6326e2d762283af715b7cdb69239bd65

Video

```js
map: std.collections.map

'derive: [Show]'
Group: {
  group_size Int
  greoup_name String
}

double_group #(Group, Group)
double_group = { g Group,  Group(2 * g.size, g.name) }

names [String]
name = ['Sheldon' 'Leonard' 'Penny' 'Rajesh' 'Howard']

groups #([String])
groups = { names.map(new_group) ++ groups.map(double_group) }

new_group #(String, Group)
new_group = {
    name String
     Groups(1, name)
}

nth #(Int, [Group], String)
nth = {
  n Int
  gs [Groups]
  when [
    {gs.is_empty()}: {strings.error}
    {n < gs[n].size}: {gs[n].name}
    {when.else}:{ nth (n - gs[n].size, gs.rest())
  ]
}

```
