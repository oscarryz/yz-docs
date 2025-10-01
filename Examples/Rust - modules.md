#example

[Defining Modules to Control Scope and Privacy](https://doc.rust-lang.org/book/ch07-02-defining-modules-to-control-scope-and-privacy.html)

```js
// src/lib.yz
front_of_house : {
    hosting : {
        add_to_waitlist: {}

        seat_at_table: {}
    }

    serving : {
        take_order:{}

        serve_order:{}

        take_payment: {}
    }
}
eat_at_restaurant:{
   front_of_house.hosting.add_to_waitlist()
   meal: back_of_house.summer_breakfast("Rye")
   meal.toast = "Wheat"
   print("I'd like `meal.toast` please")

   order1: back_of_house.appetizer.SOUP
   order2: back_of_house.appetizer.SALAD
}
hosting: front_of_house.hosting

customer: {

   eat_at_restaurant2:{
      hosting.add_to_waitlist()
   }
}

deliver_order:{
}
back_of_house:{
   Breakfast {
      toast String
      seasonal_fruit String
   }
   appetizer{
       SOUP: {}
       SALAD: {}
    }

   summer_breakfast: {
          toast String
          Breakfast(
              toast: toast
            seasonal_fruit: "peaches"
          )
   }


   fix_incorrect_order:{
      cook_order()
      deliver_order()
    }
    cook_order:{}
}


//
boh:{
    // anenum macro could be used to declare variables
    enum!(Appetizer
          SOUP
          SALAD)
    // same as
    Appetizer{}
    SOUP: Appetizer()
    SALAD: Appetizer()
}
```
