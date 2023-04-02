# Online Store REST-API - Api

## Overview
A Go implemented API that an online store or eshop could use. 

tl;dr

## How to build locally

Navigate to root's project folder and execute the command below:

*requirements:* `docker-compose` and `docker` installed

```
docker-compose up -d
```

This command will use the `docker-compose.yml` file in order to create 2 running Docker containers. If everything goes fine, the command below will print the newly created running containers.

```
docker ps
```

The containers we've just creted:
- bestprice-api_db
- bestprice-api_redis

The first one serves a MySql database and the second one an empty Redis database cache. 

These 2 containers listen to the ports:
- 33066
- 63799

If ports conflicts raise, change the configuration in `docker-compose.yml` file.

Next step. At the project's root folder execute the command below:

*requirements:* `mysql` installed

```
mysql -h 0.0.0.0 -P 33066 -uuser -ppassword bestprice < bestprice_dump.sql
```
*consider to check your local ip address*

If the `docker-compose.yml` file is unchanged (contains credentials), the above command will populate the `bestprice` database with the schema and data provided from `bestprice_dump.sql` file.

If everything goes fine, the `bestprice` database will have been populated with 3 tables 
- users
- categories
- products

These 3 tables are accompanied with some dummy data, in order the project to be able to be tested easily. Also, the `bestprice` database, contains Primary/Foreign Keys, Indexes and other rules that speed up the queries applied as also as the data coherence. For further information please read the `bestprice_dump.sql` file.

## Run the Api server

Next step, is to run the application api server. Navigate from project's root folder to the folder that contains the executable file /cmd/bestprice-api
```
cd cmd/bestprice-api
```
and execute the executable file `main.go`
```
go run main.go
```
*requirements:* `go` installed

If everything goes fine, you have to see a message on terminal:
*Listening to :8080*. 
This means that the Api server is "up and running" and serves at port: 8080.

For further configauration, you can apply the global variables below, before running the application:
- BESTPRICE_MYSQL_PATH
- BESTPRICE_REDIS_PATH
- ADDRESS

These global variables modify the connection strings of MySQL and Redis servers, as also as the Api's serving address. For further information please the the `config.go` file.

## Api services provided

A Postman collection file `postman_collection.json` is provided within the project. It's suggested to be used and imported in Postam in order achieve to use the Api as fast as possible.

*requirements:* `Postman` installed

The Api returns and receives information using HTTP and JSON information protocols. 

Api routing  **`:8080/v1/`** 

The Api provides the methods below:

### **`Ping`** - Verb: GET - Path: - "/" - Checks if Api is "up and running"

Expected output:
```
{
    "Message": "I'm alive"
}
```
### **`User Signup`** - Verb: POST - Path: "/signup" - A user can signup to application
Expected input:
```
{
  "username": "johndoe",
  "password": "1234"
}
```
Expected output:
```
{
    "Message": "User created successfully"
}
```

### **`User Login`** - Verb: POST - Path: "/login" - A user can login to application and receive theirs Bearer Token for authorized access
Expected input:
```
{
  "username": "johndoe",
  "password": "1234"
}
```
Expected output:
```
{
    "Message": "Welcome johndoe",
    "Information": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9eyJ1c2VybmFtZSI6ImJheW1hbiIsImV4cCI6MTYwNDI1OTQ2MX0z5h0iQ6m5cspK_m5qxJdgH1E0P4fGhbCy6edlsRKVN4"
}
```

### **`List All Categories`** - Verb: GET - Path: "/categories" - Lists all categories
Expected output:
```
{
        "created_on": "01-11-2020 12:53:43",
        "id": 2,
        "title": "garden",
        "url_image": "http://www.archive.bp.gr/images/categories/garden.webp",
        "place": 2
    },

    ...
    
    {
        "created_on": "01-11-2020 12:53:43",
        "id": 3,
        "title": "technology",
        "url_image": "http://www.archive.bp.gr/images/categories/technology.webp",
        "place": 3
    }
```
### **`Read A Category`** - Verb: GET - Path: "/categories/{id}" - Fetches a category mathced with provided id
Expected output:
```
{
    "created_on": "01-11-2020 12:53:43",
    "updated_on": "01-11-2020 12:57:43",
    "id": 3,
    "title": "technology",
    "url_image": "http://www.archive.bp.gr/images/categories/technology.webp",
    "place": 3
}
```
### **`Create A Category`** - Verb: POST - Path: "/categories" - Creates a new category
Expected input:
```
{
    "title": "tools",
    "url_image": "http://www.archive.bp.gr/images/categories/tools.webp",
    "place": 35
}
```
Expected output:
```
{
    "Message": "Category added"
}
```
### **`Update A Category`** - Verb: PUT - Path: "/categories/{id}" - Updates a new category
Expected input:
```
{
    "title": "tools",
    "url_image": "http://www.archive.bp.gr/images/categories/tools.webp",
    "place": 18
}
```
Expected output:
```
{
    "Message": "Category updated successfully"
}
```
### **`Delete A Category`** - Verb: DELETE - Path: "/categories/{id}" - Deletes a new category
Expected output:
```
{
    "Message": "Category deleted successfully"
}
```

### **`List All Products`** - Verb: GET - Path: "/products" - Lists all products
Expected output:
```
    {
        "created_on": "01-11-2020 12:53:43",
        "category_id": 1,
        "description": "this is a fridge",
        "id": 1,
        "title": "fridge",
        "url_image": "http://www.archive.bp.gr/images/products/fridge.webp",
        "price": 80.4
    },

    ...
    
    {
        "created_on": "01-11-2020 12:53:43",
        "category_id": 1,
        "description": "this is a bed",
        "id": 2,
        "title": "bed",
        "url_image": "http://www.archive.bp.gr/images/products/bed.webp",
        "price": 61.2
    }
```
### **`Read A Product`** - Verb: GET - Path: "/products/{id}" - Fetches a product mathced with provided id
Expected output:
```
{
        "created_on": "01-11-2020 12:53:43",
        "category_id": 1,
        "description": "this is a bed",
        "id": 2,
        "title": "bed",
        "url_image": "http://www.archive.bp.gr/images/products/bed.webp",
        "price": 61.2
    }
```
### **`Create A Product`** - Verb: POST - Path: "/products" - Creates a new product
Expected input:
```
{
    "category_id": 5,
    "description": "this is a sock",
    "title": "sock",
    "url_image": "http://www.archive.bp.gr/images/products/sock.webp",
    "price": 4.2
}
```
Expected output:
```
{
    "Message": "Product added"
}
```
### **`Update A Product`** - Verb: PUT - Path: "/products/{id}" - Updates a new product
Expected input:
```
{
    "category_id": 5,
    "description": "this is a sock",
    "title": "sock",
    "url_image": "http://www.archive.bp.gr/images/products/sock.webp",
    "price": 3.2
}
```
Expected output:
```
{
    "Message": "Product updated successfully"
}
```
### **`Delete A Product`** - Verb: DELETE - Path: "/products/{id}" - Deletes a new product
Expected output:
```
{
    "Message": "Product deleted successfully"
}
```
## Information about the Api services

The services `List All Categories` and `List All Products` can be used with pagination. The pagination follows the `Limit` - `Offset` standards 
https://specs.openstack.org/openstack/api-wg/guidelines/pagination_filter_sort.html

Example:
```
http://localhost:8080/v1/categories?sort=title:desc&limit=3&offset=1
```

The services of `Create`, `Update` and `Delete` need authorization. Authorization is successfull only for the users that have created an account and are loged in. If a user is loged in, a `Bearer Token` is provided to him in order to make authorized request. The token expires within 30 minutes and can be refreshed only if the user logs in again.

If a page on application is not found, this is handled from the Api, providing the message below:
```
{
    "code": 404,
    "message": "Not Found",
    "error": "This page could not be found"
}
```
The query paramters that cannot be handled from application, there are no considered as errors.

## Implementation
### MySQL

On MySQL side, `indexes` have been applied to primary keys, as also as to the columns wich are used for the default ordering. Also, indexes have been applied to the foreign key that describe the relation of two columns of two different tables. The foreign key follows the `cascade strategy`.

`Not all columns can be NULL`, taking as consideration that are handled as keys as also as the business logic of a commercial Api.

### Redis
The Redis every time the application starts gets `flushed`. The Redis filles up with Key/Values pairs `every time the user makes a request`. Works only on read methods ofcourse.

If the result of request is not stored to Redis, the result is provided to user from MySQL and then is stored to Redis. `Every next time`, the data for the same request, are served from Redis.

When the user updates or deletes an entity, the `Redis deletes only the specific data that are linked with this request`, in order to be filled up again with an other request to the future.

The implementation doesn't use synchronous deletion to Redis, but `asynchronous` in order to `avoid problems` on heavy traffic.

### Application Entities
Some entities are designed to follow the schema of DB and some others are designed in order to be `exposed to users`. So, they are accompanied with their functions that in some cases `overwrite the default rendering functions` GO provides. This strategy works as a layer beetween application and the user and prevents some information to be provided to users.

Product and Category entities tha have "a lot" common attributes, share the same entity.

### Layers and Abstraction
In order the logic of the application to be changed, just modifying or replacing part of the code, the `layers of the application have been implemented with interfaces`. This gives us the flexibility, as an example, to change the layer which communicates with Redis, without touching other parts of the code. Also this gives us the flexibility for in depth moching and testing.

### Authentication
Authentication follows two steps. 

The first one is the user creation with Username and Password. The username is visible to the application but the `password not, because is stored hashed`.

The second one is the Bearer Token that is provided to logged in users. The token follows the standards of `RFC 7519 - JWT token` and has expiration period.

### Errors
The application provides 4 types of custom errors.
- Internal 500
- Bad Request 400
- Not Found 404
- Unauthorized 401

At this time, the errors expose the inner errors of application and database layers.

### Helper functions
In order to reduce the code, the development has been designed with a `helper` package that provides widely used functions that well described for their usage.

### Test 
For testing purposes some MySql database and Redis services have been mocked. For the time being tested functions have been created only for Api services. It would be a good idea tests to be created for all api services, especially to these with more `complex logic` like paginaton.

### SQL Injection
In order to prevent SQL injections, the application uses the `SQLX driver` for communicating with MySQL database. This `driver prevents SQL injections` when the user provides `values of variables`. In order to prevent SQL injection when the application needs to `change the structure of a query` (ex. naming the columns in order by), the application uses a `regex filtering`.

### Dependancy management
As a dependency management system it has been used the `GO Modules` standards

## Further improvements
### Cache
We may consider to use `more than one database-different clients` in Redis or `create collections` in order to organise better the stored queries results.

Now, when an entity is deleted (ex. Category), all the data in redis that refer to this entity (Category) are flushed. This needs a smarter approach.

### Test
Due to lack of time, the testing purposes have not been completed. They have to be developed further especially for `2 reasons`. For `pushing the data retrieve systems to their limits`, especially Redis, to check their behavior. And ofcourse taking care of all `possible authentication - sql injection cases`.

### Output
Some minor changes should be applied for better user experience. Example: Joining the category id of a Product with the Category name. Now we return the category id of a product and not the category name. A good idea should be both to be returned. An other example, when a entity is created, we could provide back the id of this entity.

### Errors
Better approach of error handling in order `not to expose inner errors` to users as also as to return `usefull information` in order to error to be bypassed.
