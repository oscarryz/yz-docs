#example

From:

https://gist.github.com/rexim/6326e2d762283af715b7cdb69239bd65

Video

```js
map: std.collections.map

'derive: [Show]'
Group: {
  group_size Int
  group_name String
}

double_group #(Group, Group) {
    g Group
    Group(2 * g.size, g.name)
}

names: ['Sheldon', 'Leonard', 'Penny', 'Rajesh', 'Howard']

groups #([String], [Group]) {
    names.map(new_group) ++ groups.map(double_group)
}

new_group #(String, Group) {
    name String
    Groups(1, name)
}

nth #(Int, [Group], String) {
  n Int
  gs [Groups]
  match {
    gs.is_empty()   => strings.error
  }, {
    n < gs[n].size  => gs[n].name
  }, {
    nth(n - gs[n].size, gs.rest())
  }
}

```