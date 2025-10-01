#example

https://cuis-smalltalk.github.io/TheCuisBook/Block-syntax.html

```js


{ 1 + 2 }() // 3
{x Int, 1 + x}(2) // 3
{x Int, y Int, x + y}(1 2) // 3

{x Int, y Int, z: x + y,  z}(1 2) //3

x: 1
{y Int, x + y}(2) // 3
{ 1, 2, 3 }

// if `Int.to` returns a range
n: 60
m: 45
1.to(n).select { d Int, n / d == 0}
// [1 2 3 4 5 6 10 12 15 20 30 60]
1.to(m).select { d Int, m / d == 0}
// [1 3 4 9 15 45]

divisible_by: {
   d Int
   x Int
   x / d == 0
}

1 .to n .select divisible_by

```
```smalltalk
"Smalltalk"
SpaceWar>>teleport: aShip
  "Teleport a ship at a random location"
  | area randomCoordinate |
  aShip resupply.
  area := self morphLocalBounds insetBy: 20.
  randomCoordinate := [(area left to: area right) atRandom].
  aShip
    velocity: 0 @ 0;
    morphPosition: randomCoordinate value @ randomCoordinate value
```
```js
// Yz
SpaceWar: {
    teleport: { a_ship Ship
        area Area
        random_coordinate: {int.random((area.left.to(area.right))) }
      //random_coordinate: {int.random  area.left.to area.right    }
        a_ship.velocity = 0
        a_ship.morph_position( random_coordinate() random_coordinate() )
      //a_ship.morph_position random_coordinate() random_coordinate()
    }
}
```
```js
17 * 3 > 220 ? {'bigger'}, {'smaller'}

```
