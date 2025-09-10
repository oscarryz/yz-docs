#example
From [2048.v](https://github.com/vlang/v/blob/master/examples/2048/2048.v#:~:text=%7D-,fn%20(app%20%26App)%20draw_tiles()%20%7B,app,-.theme.tile_colors.last)

```js

	App : {
		draw_title: {
			xstart: app.ui.x_padding + app.ui.border_size
			ystart: app.ui.y_padding + app.ui.border.size + app.ui.header_size
			toffset: app.ui.title_size + app.ui.padding_size

			0.to(4,  { y Int
				0.to(4, { x Int 
					tttidx: app.gborer.field[y][x]
					tile_color String
					tidx < app.theme.tile_colors.len ? {
						tiele_color = app.theme.tile_colors[tidx]
					} ,  {
						tiele_color = app.theme.tile_colors.last
					}
				
				})
					
			})
		}
	}
```

