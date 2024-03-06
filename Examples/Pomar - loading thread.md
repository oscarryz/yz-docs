[Pomar](https://twitter.com/pomarlang/status/1764763278624887290/photo/1)


```js
asset_laoding_thread: {
	buffer []Image
	file_iterator : list_files('assets')
	it File
	file.iterator.each { 
		f File
		img: raylib.load_image(f.full_name)
		buffer.push(img)
	}
}
main: {
	buffer: []Image
	// starts putting images in the array
	// it will go to the nextline immediately 
	asset_loading_thread(buffer)
	// Righ away will start checkint for window.close and consuming images from
	// the buffer
	while { raylib.window_should_close() == false } {
		t: raylib.load_texture_from_image(buffer.shift())
		textures.push(t)
	}
}
```