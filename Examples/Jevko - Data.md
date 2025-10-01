#example

From: [Jevko.org](https://jevko.org/)

## Data Interchange

```js
  person: {
      first_name: 'John'
      last_name : 'Smith'
      is_alive  : true
      age       : 27
      address   : {
        street_address: '21 2nd Street'
        city: 'New York'
        state: 'NY'
        postal_code: '10021-3100'
      }
      phone_numbers: [{
        type: 'Home'
        number: '212 555-1234'
      } {
        type: 'Office'
        number: '646 555-4567'
      }]
      children: []
      spouse: no_one
  }
  no_one: Person{}
```

## Configuration

```js
'last modified 1 April 2001 by John Doe'
owner: {
  name: 'John Doe'
  organization: 'Acme widgets Inc.'
}

database: {
  server: '192.0.0.2.62'
  port: 143
  file: 'payroll.dat'
  select_columns: [
    'name'
    'address'
    'phone number'
  ]
}
```

## Text markup

```js
html:{
  head: {
    title: 'This is a title'
  }
  body:{
    div:{
      p:'Hello world'
      abbr: {
          id:'anId'
          class:'jargon'
          style:'color:purple;'
          title:'Hypertext Markup Language'
        'HTML'
      }
      a: {
          href:'https://www.wikipedia.org'
          'A link to the Wikipedia!'
      }
      p: {
          'O well'
          span: {
              lang:'fr'
              "c'est la vie"
          }
          'as they say in france'
      }
    }
  }
}
```
