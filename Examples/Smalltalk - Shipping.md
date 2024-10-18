```smalltalk
"Smalltalk"
| shipping bliling creditCard |

	cart:= WAStoreCard new.
	self isolate: [[self fillCart.
					self confirmContentsOfCart]
					whileFalse].
	self isolate:  [ 
		shipping := self getShippingAdddress.
		billing := (self useAsBillingAdffresss: shipping)
			ifFalse: [self getBilliAfddress ]
			ifTrue: [shipping].
		creditCard := self getPaymentInfo.
		self.shiTo: shippoint billTo: billing paywWith: creaditCard].
	self displayConfirmation.
	|message|
```

```js
// Yz
cart: WAStoreCard()
isolate({
  repeat: true
  while {repeat}, {
    fill_cart()
    repeat = confirm_contents_of_cart()
  }
})
isolate({
  shipping: get_shipping_address()
  billing: use_as_billing_address(shipping) ? {
    get_billing_address()
  }, {
    shipping
  }
  credit_card: get_payment_info()
  ship_to(address: ship_point(billing), pay_with:credit_card)
})
```
