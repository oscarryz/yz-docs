#example
https://dl.acm.org/doi/pdf/10.1145/3622852

```js
Account: {
  balance Int
}
main: {
  acc1: Account()
  acc1.balance = acc1.balance - 100 // Ok 

  acc2: Account()
  acc2.balance = acc2.balance + 100 // Ok
}
```

```js
transfer: {
  src Account
  dst Account
  amount Int
  src.balance = src.balance - amount
  dst.balance = src.balance + amount
}
```


```js
transfer: {
  src Account
  dst Account
  amount Int

  src.balance.>=( amount ).?({
    src.balance = src.balance.-( amount )
    dst.balance = dst.balance.+( amount )
  })
  
}
```

```js
transfer: {
  src Account
  dst Account
  amount Int 
  if src.balance >= amount && src.frozen == false && dst.frozen == false ? {
    src.balance = src.balance.-( amount )
    dst.balance = dst.balance.+( amount )
  }
}
...
account Account
transfer(account, account, 100)
transfer(account, account, 100)

```

???
[[Transactions or atomic units of work]] 

