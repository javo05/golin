# golin

Golin is an approach to create a microservice that allows multiple services to use authentication without needing to implement it.

The project uses JWT as a way of securing the comunication bewteen the client and the login server. 

Interfaces and some abstractions were made to make it as generic as possible, this could be changed to use some configuration files or changing the login methods.

This microservice is fully scalable since for example it can be set aside of any API and just coordinate the validation of the tokens. Remember that all of this needs to be done via the Backend to ensure everything will work in a secure manner.

GinGonic was used to build the API since it's quick to develop and easy to maintain and debug. Also BoltDB is used to storage Tokens and is designed to have a Blacklist for tokens that are no longer valid because of logout.

Here are some useful links we read for understanding JWT and some logins procedures:
  - http://blog.brainattica.com/restful-json-api-jwt-go/
  - https://github.com/gin-gonic/contrib/blob/master/jwt/example/example.go
  - http://golangtutorials.blogspot.mx/2011/06/interfaces-in-go.html
  - https://github.com/boltdb/bolt#resources
  - http://jwt.io/
  - https://tools.ietf.org/html/rfc7519
  
