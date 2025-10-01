#example
https://leetcode.com/problems/equal-row-and-column-pairs/


```js

/*
    - Iterate rows and store in hashmap (convert to hash if needed)
    - Iterate columns and lookup in hashmap
*/
equal_pairs: {
    grid [][]Int
    count : 0 
    map: add_rows grid 
    grid.each({ i Int v Int 
        count = count + search(map, grid, i)
    }
    0.to(len - 1 ).each({ i Int
        count = count + search(map, grid, i)
    }
    count
}
add_rows: {
    grid [][]Int
    map [Int:Int] = [:]
    grid.each({ row []Int 
        h: hash(row)
        map[h] = map[h] + 1
    }
    map
}
search: {
    map [Int:Int]
    grid [][]Int
    i Int // column index

    column: []
    0.to( grid.len() -1 ).each({j Int 
        column << grid[i][j]
    }
    map[hash(colum)]
}
hash: { 
    a []Int
    r : 31
    a.each({ i Int 
        r = 31 * r + i
    }
    r
}


```