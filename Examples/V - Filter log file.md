#example
https://vosca.dev/p/e31e438324
```js
os.write_file('app.log' `
ERROR: log file not found
DEBUG: create new file
DEBUG: write text to log file
ERROR: file not writeable
`).or {
    err Error
    eprintln('failed to write the file `err`')
}

text: os.read_file 'app.log' .or { e Error
	eprintln('failed to read the file: `e`')
}

list: text.split_into_lines()

lines.each({ line String
    line.starts_with 'DEBUG' ? {
        print(line)
    }
}
```