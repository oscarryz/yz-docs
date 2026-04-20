#example


https://github.com/chshersh/ccl

```js
// : "This is (not) a ccl document"
// import ommited
title: "CCL Example"
database: {
	enabled: true
	ports: [
		8000,
		8001,
		8002
	]
	limits: {
		cpu: Mi(1500)
		memory: Gb(10)
	}
	user: {
		guest_id: 42
		login: something
		created_at: Date(2024, 12, 31)
	}
}
``` 


