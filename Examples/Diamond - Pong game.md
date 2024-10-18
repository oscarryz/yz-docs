
[pongGame.dmd](https://gist.githubusercontent.com/amzamora/ad4c84a73e1006a1d17d784c74b40a0a/raw/5f73279abaf9670545f696afee65838ac1b8d44c/pongGame.dmd)


```js
graphics: std.graphics
input: std.input
math: std.math

//-- Types
Ball: {
    x Int
    y Int
    dx Int
    dy Int
    size Int
}

Paddle: {
    x Int
    y Int
    width Int
    height Int
}

windowWidth: 800
windowHeghts: 600
  
//-- Initialize the game window and graphics context
window: graphics.create_window("Pong", windowWidth, windowHeight)
context: graphics.create_graphics_context(window)

//-- Set up game objects
ball : Ball(x: 400, y: 300, dx: 5, dy: 5, size: 10)
paddle1 : Paddle(x: 50, y: 250, width: 10, height: 100)
paddle2 : Paddle(x: 740, y: 250, width: 10, height: 100)
score1 : 0
score2 : 0

//-- Game loop
while {true} , {
    context.clear()
    context.draw_circle(ball.x, ball.y, ball.size)
    context.draw_rect(paddle1.x, paddle1.y, paddle1.width, paddle1.height)
    context.drar_rect(paddle2.x, paddle2.y, paddle2.width, paddle2.height)

    //-- Handle user input
   if_any [ 
    { input.key_down("W") }: {
        paddle1.y -= 5
    }

    { input.key_down("S") } : { 
        paddle1.y += 5
    }

    { input.key_down("Up") } : {
        paddle2.y -= 5
    }

    { input.key_down("Down") } : {
        paddle2.y += 5
    }
  ]

    // -- Move the ball and check for collisions with walls and paddles
    ball.x += ball.dx
    ball.y += ball.dy

    if_any [
      { ball.x < 0 }: {
        score2 += 1 
        reset(ball)
      }

      { ball.x > window_width } : {
        score1 += 1
        reset(ball)
      }

      { ball.y < 0 or ball.y > window_height } : {
        ball.dy = -ball.dy
      }

      { intersects(ball, paddle1) or intersects(ball, paddle2) } : {
        ball.dx = -ball.dx
      }

    //-- Check for game over condition
    { score1 >= 10 or score2 >= 10 } : {
        display_game_over(context, score1, score2)
        break
      }
    ]
    //-- Display current scores
    displays_scores(context, score1, score2)

    //-- Wait for next frame
    time.sleep(16)
}

//-- Helper functions
reset: {
    ball Ball
    ball.x : 400
    ball.y : 300
    ball.dx : -ball.dx
    ball.dy : math.random() * 10 - 5
}

intersects #(ball Ball, paddle Paddle, Bool) = {
      ball Ball
      paddle Paddle

      ball.x + ball.size < paddle.x 
  ||  ball.x > paddle.x + paddle.width
  ||  ball.y + ball.size < paddle.y
  ||  ball.y > paddle.y + paddle.height
  == false
}
  
displays_scores : {
    context graphics.Context
    score1 Int
    score2 Int

    context.drawText("Player 1: `score1`", 10, 20)
    context.drawText("Player 2: `score2`", 700, 20)
}

display_game_over: {
    context graphics.Context
    context.drawText("Game Over", 350, 300)
}

```