#example
https://mobile.twitter.com/v_language/status/1546061255013711875

```js

Book: {
  author: {
    name String
    age  Int
  }
  title String
}

test_anonymous_structs: { 

    empty_book: Book()
    assert empty_book.author.age == 0
    assert empty_book.author.name == ''

    book: Book (
      author: { name: 'Peter Brown'; age: 23 }
      title: 'Programming in Yz'
    )
    assert book.author.age == 23
    assert book.author.name == 'Peter Brown'

    book2: Book {
      author: { 
        name: 'Samantha Black'
        age: 24 
      }
      title: 'Anonymous blocks are cool'
    }
    assert book2.author.age == 24
    assert book2.author.name == 'Samantha Black'
}

```


https://vsql.readthedocs.io/en/latest/custom-functions.html

```js
// no_pennies will round to 0.05 denominations.
db.register_function 'no_pennies(float) float', { 
  a []vsql.Value; 
  returns vsql.Value 
  amount : math.round(a[0].f64_value / 0.05) * 0.05
  returns = vsql.new_double_precision_value(amount)
} 

db.query 'CREATE TABLE products (product_name VARCHAR(100), price FLOAT)'
db.query "INSERT INTO products (product_name, price) VALUES ('Ice Cream', 5.99)"
db.query "INSERT INTO products (product_name, price) VALUES ('Ham Sandwhich', 3.47)"
db.query "INSERT INTO products (product_name, price) VALUES ('Bagel', 1.25)"

result: db.query 'SELECT product_name, no_pennies(price) as total FROM products'
result.each({ row vsql.Row
  total: row.get_f64 'TOTAL'  
  print('${row.get_string("PRODUCT_NAME")) ?} ${total:.2f}'
}
```
