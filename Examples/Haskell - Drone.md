#example

```haskell
-- Haskell
import Control.Monad (replicateM)
import Data.Foldable (foldl')
import qualified System.Random.Stateful as Rand

data Drone = Drone
  { xPos :: Int
  , yPos :: Int
  , zPos :: Int
  } deriving Show

data Movement
  = Forward | Back | ToLeft | ToRight | Up | Down
  deriving (Show, Enum, Bounded)

main :: IO ()
main = do
  let initDrone = Drone { xPos = 0, yPos = 100, zPos = 0 }
  -- Generate 15 moves randomly
  randomMoves <- replicateM 15 $ Rand.uniformEnumM Rand.globalStdGen
  let resultDrone = foldl' moveDrone initDrone randomMoves
  print(resultDrone)

moveDrone :: Drone -> Movement -> Drone
moveDrone drone move =
  case move of
    Forward -> drone { zPos = zPos drone + 1 }
    Back    -> drone { zPos = zPos drone - 1 }
    ToLeft  -> drone { xPos = xPos drone - 1 }
    ToRight -> drone { xPos = xPos drone + 1 }
    Up      -> drone { yPos = yPos drone + 1 }
    Down    -> drone { yPos = yPos drone - 1 }
```

```js
// Yz attempt
Drone {
    x_pos Int
    y_pos Int
    z_pos Int
    move: {
        direction Int
        when_eq direction [
            { movement.forward }  : { z_pos = z_pos + 1 }
            { movement.back }     : { z_pos = z_pos - 1 }
            { movement.to_left }  : { x_pos = x_pos - 1 }
            { movement.to_right } : { x_pos = x_pos + 1 }
            { movement.up }       : { y_pos = y_pos + 1 }
            { movement.down }     : { y_pos = y_pos - 1 }
        ]
    }
}
movement: {
    forward  : 0
    back     : 1
    to_left  : 2
    to_right : 3
    up       : 4
    down     : 5

    values: [forward back to_left to_right up down]
    random: {
        index: int.random(values.len())
        values[index]
    }
}
main: {
    drone: Drone(0 10 0)
    _ : 15.times({
        drone.move(movement.random())
    }
    print('`drone`')
}
```
