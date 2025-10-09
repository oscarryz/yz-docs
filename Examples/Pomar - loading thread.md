#example

[Pomar](https://twitter.com/pomarlang/status/1764763278624887290/photo/1)


```js
// Blocks (methods/functions/closures) are concurrent by default
// So, there's no need for channels to share data...
// one block is filling the array, another is reading from it.
asset_loading_thread: {
	// Put images here
	buffer [Image]()
	file_iterator : list_files('assets')
	file_iterator.each({
		f File
		img: raylib.load_image(f.full_name)
		buffer.push(img)
	}
}
main: {
	buffer: [Image]()
	// starts putting images in the array
	// it will go to the nextline immediately
	asset_loading_thread(buffer)
	// Righ away will start checkint for window.close and consuming images from
	// the buffer
	while { raylib.window_should_close() == false }, {
		// Gotta check if there's something in buffer first
		if buffer.is_empty() == false  {
			t: raylib.load_texture_from_image(buffer.shift())
			textures.push(t)
		} {
			// busy wait ... :/ keep trying.
			// not good  :/ :/
		}
	}
}
```
