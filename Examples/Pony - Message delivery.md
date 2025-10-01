#example

```js

Customer: {
	store Store
	bank Bank

	run: {
		prince Int
		_ : bank.credit(me, price)
		store.buy(me, price)
	}
	me Customer
}
Store : {
	bank Bank
	buy: {
		cust Customer
		price Int
		bank.debit(customer, price)
	}
}
Bank : {
	balances [Customer:Int]= [Customer]Int

	credit: {
		customer Customer
		amount Int
		b : balances[customer]
		balances[customer]= b + acount
	}
	debit :{
		customer Customer
		price Int
		b : balances[customer]
		balance < price  ? {
			Err('Not enough balance')
		} {
			balances[customer]= b - price
			Ok()
		}
	}
}
main: {
	b: Bank()
	s: Store(b)
	c : Customer(s b)
	c.me = c
	c.run()
}
```
