#example
https://leetcode.com/problems/equal-row-and-column-pairs/


```js

/*
    - Iterate rows and store in hashmap (convert to hash if needed)
    - Iterate columns and lookup in hashmap
*/
equal_pairs: {
    grid [][Int]()
    count: 0
    map: add_rows(grid)
    0.to(grid.length() - 1).each({ i Int
        count = count + search(map, grid, i)
    })
    count
}
add_rows: {
    grid [][Int]()
    map: [[Int]: Int]()
    grid.each({ row [Int]()
        h: hash(row)
        map[h] = map[h] + 1
    })
    map
}
search: {
    map [[Int]: Int]
    grid [][Int]()
    i Int // column index

    column: [Int]()
    0.to(grid.length() - 1).each({ j Int
        column.push(grid[i][j])
    })
    map[hash(column)]
}
hash: {
    a [Int]()
    r: 31
    a.each({ i Int
        r = 31 * r + i
    })
    r
}


```
