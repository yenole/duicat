```
pattrn [method] [(json|plain|tmpl|multi)] [(param|json)] handle

1.handle  ->  func(ResponseWriter,*Request)
2.json,handle -> func(ResponseWriter,*Request) interface
3.json,param,handle -> func(token string) interface
4.param,handle -> func(token string) interface
```