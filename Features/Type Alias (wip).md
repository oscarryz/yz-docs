#rejected `::` will be used to declare the block type e.g. `a :: {}`  `a` is a variable of type `block` (empty block), so no type alias

The `::` operators is used to create an alias to a type and optionally add new behavior. 

The most common case is to define a block instance (see: [Instances](Instances.md))

```javascript
// full syntax
Person:: {name String last_name String full_name {String}} = {
    full_name: {name ++ last_name}
}
// shortcut
Person :: {
    name String
    last_name String
    full_name: {name ++ last_name}
}
```

Another use is to narrow the possible values of a type, e.g. think of a card deck where the possible suit values are: club, diamond, heart and spade. 

A possible implementation could be: 
```javascript
deck: {
    Suit :: {}
    club Suit{}
    diamond Suit{}
    heart Suit{}
    spade Suit{}
}
deck.club
deck.diamond
```

Another implementation would be to alias the type `Int` but allow only those 4 values as per the deck definition

```javascript
range: Range{0 3}
deck: {
   Suit :: Int
   club Suit = range.next()
   diamond Suit = range.next()
   heart Suit = range.next()
   spade Suit = range.next()
}
s deck.Suit // s is of type deck.Suit
s = 5 // possible but doesn't make sense, better is to choose from deck
s = deck.spade

```

A more complex example ... 

```javascript
// https://wiki.haskell.org/Type#A_simple_example
deck: {
    suit: {
        Suit::Int // Suit is an alias to Int. `type Suit Int` in other languages
        
        club    Suit = 0
        diamond Suit = 1
        heart   Suit = 2
        spade   Suit = 3
        values: [club diamond heart spade]
    }
    
    card_value: {
        CardValue::Int
        two CardValue =  2
        three CardValue = 3
        ...
        king CardValue = 13
        ace CardValue = 14
        values: [two three ... king ace]
    }
    Card: {
        value card_value.CardValue
        suit suit.Suit
        compare: {
            c1 Card
            c2 Card
            compare({c1.value c1.suit}{c2.value c2.suit})
        }
        // Satisfy the Enum block
        to_enum: Enum.to_enum
        from_enum: Enum.from_enum
   }
    Deck :: []Card = {
        self Deck
        suffle: {
            i:cards.len() -1;
            while {i > 0 } {
                j: int.random(i+1)
                self[j] self[i] = {self[i]; self[j]}()
                i = i - 1
            }
        }
    }
    new: {
        cards Deck
        suit.values.each { s Suit; 
            card_value.values.each {cv CardValue;
                 cards << Card{s cv}
             }
        }
        cards.self = cards
        cards
    }

}
cards deck.Deck = deck.new()
cards.suffle()

Person:: {name String last_name String full_name {String}} = {
    full_name: { name ++ last_name }
}
Person::{
    name String
    last_name String
    full_name: { name ++ last_name }
}
Student::Person = {
}

Summary::{
    summarize:{
        "(Read more ...)"
    }
}
NewsArticle: {
    headline String
    location String
    author   String
    content  String
    summarize: Summary.summarize
}
NewsArticle:: Summary {
    headline String
    location String
    author   String
    content  String
}
```

#rejected `::` will be used to declare the block type e.g. `a :: {}`  `a` is a variable of type `block` (empty block), so no type alias
#challenged We'll be using type aliases after all, this is because it's not possible to know when a new type is defined vs instantiated ( see [Define new types](Define%20new%20types.md) )
So, `Person::{}` defines a new type and `Person{}`  instantiates it (`Person:{}` still infers it). Thus keep things consistent `Suit :: Int` is a valid declaration (notice uppercase)


How can we add methods to `Suit :: Int` ?  You can't. Maybe a possible feature would be [[../Questions/solved/Extension methods]], but that something we should consider for a future release.

