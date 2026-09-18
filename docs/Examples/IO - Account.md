---
type: example
generated: { by: "oscarryz", at: 2024-03-07T16:00:59-06:00 }
---
#example

https://iolanguage.org/samples/index.html


```js
Account: { 
    balance: 0.0
    deposit : { 
        v Decimal
        balance = balance +  v
    }
    show: {
        println('Account balance ${balance}')
    }
}
account : Account()
print( "initial: " )
account.show()
println("Depositiong $10\n")
account.deposit(10.0)
print(" final: " ) 
account.show()

```