

## Generics 

https://doc.rust-lang.org/rust-by-example/generics.html

```js
foo { arg <T> }
foo = { arg <T> }

foo ( arg T )
foo = { arg T }
```

```js
// A concreate type `A`
Aa {}

// Single using A for the first time
Single {
	data Aa
}
// Similar to SingleGeneric<A> {
//
//}
SingleGeneric {  
	A
}
main: {
    // Instantiation
	s: Single(Aa())
	// Here, `SingleGeneric` has a type parameter explicity specified
	char SingleGeneric(String) = SingleGeneric('a')

	// Implicitly specified
	t: SingleGneric(Aa())
	int: SingleGneric(6)
	str: SingleGneric('a')
}
```

## Functions

```js
Aa {}
Sa {data Aa}
SGen{ T; data T }
//SGen{ data <T> }

// takes a concrete type S
// explicit
reg_fn { s Sa } = {s Sa}
// implicit
reg_fn : { s Sa }

// Explicitly given the type `A`
gen_spec_t: { s SGen(Aa)}  
// Explicitely given the type `Int`
gen_spec_int: { s SGen(Int)}
//Nothing specified, so still generic
generic { s SGen<T> } // better 
generic { s SGen{<T>} } // ugly but correct
generic { s SGen{T} }   // ugly and confusing... :(
generic { s SGen(T) }   // Winner

main: {
	reg_fn(Sa(Aa))
	gen_spec_t(SGen(Aa()))
	gen_spec_int(SGen(6))
	generic(SGen('str'))
}

```

## Implementation

```js
Sa {}
GenericVal { 
	data T
}
GenericValIntImpl { 
	data Int
}
GenericValSImpl {
	data Sa
}

Val { 
	val Decimal
	value: {
		val
	}
}
GenVal {
	gen_val T
	value: {
		gen_val
	}
}
main: {
	x : Val(3.0)
	y : GenVal(3)
	print("$(x.value()), $(y.value())")
}

```


## Traits

(there are not traits in Yz btw, that's just the name of the doc I was copying)
```js
Empty{
	T
	double_drop: {
		_ T
	}
}
Null{}

DoubleDrop { 
	T
	double_drop { _ T}
}
main: {
	empty : Empty()
	null  : Null()
	empty.double_drop(null)
}
```

## Bounds

```js
printer : { 
	t Display
	t.fmt(/*...*/)
}
Sa {
	data Display
}
s: Sa(vec(1)) // error etc. 

```

```js
HasArea { 
	area { Decimal }
}
'#[derive(Debug)]'
Rectangle { 
	lenght Decimal
	height Decimal
	area : { lenght * height }
}
Triangle { 
	lenght Decimal
	height Decimal
}
print_debug: {
	t Display
}
area: {
	t HasArea
	t.area()
}
main: {
	rectangle : Rectangle(3.0 4.0)
	triangle : Triangle(3.0 4.0)
	print_debug(rectangle)
	print("$(area(rectangle))")

	// 
	//print_debug(triangle)
	//print("$(area(triangle))")
}


```


## Empty bounds

```js
Cardinal {}
BlueJay {}
Turkey {}

Red{}
Blue{}

red: { 
	_ Red
	"red"
}
blue: { 
	_ Red
	"blue"
}
main: {
	cardinal : Cardinal()
	blue_jay : BlueJay()
	_turkey  : Turkey()
	
	// everything will work with everything because is eveything is emp
	print("A cardinal is $(red(cardinal))")
	print("A bluejay is $(blue(blue_jay))")
	print("A turkey is $(red(turkey))")
}
```


## Multiple Bounds

Not supported, if more than one type is needed, create a new type :( (probably work for a macro)


```js
Debug {
	debugging: {print("This is debug")}
}
Display {
	displaying: {print("This is display")}
}
DebugAndDisplay {
	debugging:  {print("This is debug")}
	displaying: {print("This is display")}
}
// maybe with a macro... but that's a separate topic.
`#[Derive[Debug]]
 #[Derive[Debug]]`
DebugAndDisplay{}

compare_prints: { 
	t DebugAndDisplay
	t.debugging()
	t.displaying()
}
```

## Where clauses
N/A

## New Type Idiom
(nothing special)

```js
Years {
	data Int
	to_days: {
		data * 365
	}
}
Days {
	data Int
	to_years: {
		data / 365
	}
}
old_enough: {
	age Years
	age.data > 18
}
main: {
	age: Years(5)
	age_days = age.to_days()
	print('Old enough $(old_enough(age))')
	print('Old enough $(old_enough(age_days.to_years())')
	//print('Old enough $(old_enough(age_days)')
}
```

```js
Years {
	data Int
}
main: {
	years: Years(42)
	years_as_primitive Int = years.data
	// No destructuring
	// Years(years_as_primitive) = years
}
```

## Associated Items
n/a



# Summary

- Type
	- Declaration
	- Instantiation
	- 