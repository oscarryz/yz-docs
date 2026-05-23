#open-question 


Define the exact way the `Compile` interface will work and how the many implementations would implement it aside from the definition in [Compile Time Bocs](docs/Features/Compile%20Time%20Bocs.md)


```js

Serializer : {
   Schema : {
	   foo Int
	   bar Int
   }
   
   run #(...)

}

```

For instance, when used for dependency manangement, to generate code from Go's stdlib ( http, json serializer, others ), serializers (json, http), configurations, etc. etc

There are many cases and more detail with concrete examples and discussion is needed for to create an effective implementation.
