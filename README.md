localhost:3000/orders post
```
{
  "items": [
    { "product_id": 1 }  ]
}
```
localhost:3000/orders getall
```
[
    {
        "ID": 3,
        "CreatedAt": "2025-10-21T12:03:32.149585+05:00",
        "UpdatedAt": "2025-10-21T12:03:32.149585+05:00",
        "DeletedAt": null,
        "order_id": "ord-1761030212131561000",
        "customer_id": 1,
        "items": null
    },
    {
        "ID": 4,
        "CreatedAt": "2025-10-21T12:05:06.270914+05:00",
        "UpdatedAt": "2025-10-21T12:05:06.270914+05:00",
        "DeletedAt": null,
        "order_id": "ord-1761030306256419600",
        "customer_id": 1,
        "items": null
    }
    ]
```

# skilltestproj
localhost:3000/login POST
``` 
{
    "name":"Tem3",
    "password":"123456"
}
```

localhost:3000/signup post
```
{
    "name":"Tem4",
    "password":"123456"
}
```
das
```
localhost:3000/products or localhost:8080/products GET
    {
        "ID": 1,
        "CreatedAt": "2025-10-21T11:01:56.407958+05:00",
        "UpdatedAt": "2025-10-21T11:01:56.407958+05:00",
        "DeletedAt": null,
        "name": "new Product",
        "price": 100
    }
```
localhost:3000/products/:id or localhost:8080/products:id GET

localhost:3000/products/:id or localhost:8080/products:id delete

``` 
{
    "name":"new Product",
    "price": 100,
    "categoryID": 1 
}
localhost:3000/categories or localhost:8080/categories post

```

{
    "name": "newcategory"
}


```

localhost:3000/categories or localhost:8080/categories
```
[
    {
        "ID": 1,
        "CreatedAt": "2025-10-21T12:26:58.011084+05:00",
        "UpdatedAt": "2025-10-21T12:26:58.011084+05:00",
        "DeletedAt": null,
        "name": "newcategory"
    },
    {
        "ID": 2,
        "CreatedAt": "2025-10-21T14:24:01.670678+05:00",
        "UpdatedAt": "2025-10-21T14:24:01.670678+05:00",
        "DeletedAt": null,
        "name": "newcategory"
    }
]
```
