#example

```java
/*
 * Click `Run` to execute the snippet below!
 */

import java.util.*;

/*
 * To execute Java, please define "static void main" on a class
 * named Solution.
 *
 * If you need more classes, simply define them inline.
 */

/*
a: [1 2  8 4 2 6 7]
t: 7
m: {6:0,5:1, -1:2, 3:3, 5:4,  }
i: 2
r: 0, 5
*/
class Solution {
  public static void main(String[] args) {
    // var list = Arrays.asList( 1 ,2, 5,4,2,6,7 );
    test(Arrays.asList(1,2), calc(1,2,5,4,2,6,7));
    test(Arrays.asList(0,4), calc(4, 1, 2, 2, 3));
  }

  public static List<Integer> calc(int ... input ) {
    var list = new ArrayList<Integer>();
    for ( var i : input ) {
      list.add(i);
    }
    var map = new HashMap<Integer, Integer>();
    var result = new ArrayList<Integer>();
    var target = 7;
    for ( int i = 0 ; i < list.size() ; i++ ) {
      var n = list.get(i);
      if ( map.containsKey(n) ) {
        result.add(map.get(n));
        result.add(i);
        break;
      } else {
        map.put( target - n , i );
      }
    }
    result;

  }
  public static void test( List<Integer> expected, List<Integer> actual) {
    assert expected.equals(actual): String.format("Expected: %s, got %s", expected, actual);
    System.out.println("OK");
  }
}

```

```js
calc: {
   input []Int
   list: []Int
   input.each({ i Int
       list << i
   }
   map: [Int]Int
   result: []Int
   target: 7
   list.each_item: { item Int; index Int
       if map.contains(item) {
           result << map[item]
           result << index
           break
       } {
           map[target-n] = index
       }
   }
   results
}
```
