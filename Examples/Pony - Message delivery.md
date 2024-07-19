```javascript

Customer: {
	store Store
	bank Bank

	run: {
		prince Int
		_ : bank.credit(me price)
		store.buy(me price)
	}
	me Customer
}
Store : {
	bank Bank
	buy: {
		cust Customer
		price Int
		bank.debit(cust price)
	}
}
Bank : {
	balances [Customer]Int = [Customer]Int()

	credit: { 
		cust Customer
		amount Int
		b : balances[cust]
		balances[customer] b + acount 
	}
	debit :{
		cust Customer
		price Int
		b : balances[cust]
		balance < price  ? {
			Error('Not enough balance')
		} {
			balances[cust] b - price
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