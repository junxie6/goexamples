# Go Examples
Some useful examples for Go golang

# Show how to use the third party Go middlewares and how to implement a custom Go middleware
\# cd httpmuxmiddleware 
\# openssl req -newkey rsa:2048 -nodes -subj "/C=CA/ST=British Columbia/L=Vancouver/O=My Company Name/CN=erp.local" -keyout erp.local.key -out erp.local.csr 
\# go run httpmuxmiddleware/httpmuxmiddleware.go 

# Run a web server and Transparently gzip the response body
\# go run webgzip/webgzip.go

# Run a web server and serve static files
\# go run webstaticfile/webstaticfile.go

# Convert a map to a JSON string
\# go run convmapjson/convmapjson.go

# Convert a struct to a JSON string
\# go run convstructjson/convstructjson.go

# Decode JSON string
\# go run decodejson/decodejson.go

# Two factor authentication TOTP (Time-based One Time Password) example
\# go run twofactortotpsec51/twofactortotpsec51.go

Then, go to:

http://127.0.0.1/qr

http://127.0.0.1/v

More info: https://github.com/sec51/twofactor

# Two factor authentication TOTP (Time-based One Time Password) example
\# go run twofactortotppquerna/twofactortotppquerna.go

Then, go to:

http://127.0.0.1/qr

http://127.0.0.1/v

More info: https://github.com/pquerna/otp

# JSON Web Token (JWT)
\# go run jsonwebtoken/jsonwebtoken.go

Then, go to http://127.0.0.1/

More info: https://github.com/dgrijalva/jwt-go

# Date time format conversion
\# go run datetimeconv/datetimeconv.go 127.0.0.1

# Loop through golang struct
\# go run loopthroughstruct/loopthroughstruct.go

# HTML template
\# go run htmltemplate/htmltemplate.go

# Multiple HTML templates and passing the variables to the inner templates
\# go run htmltemplatemultiple/htmltemplatemultiple.go

# Print Message
\# go run printmessage/printmessage.go

# Read arguments from console
\# go run consoleargument/consoleargument.go test.txt A B C

# Connect to MySQL
\# go run connectmysql/connectmysql.go

# Connect to MySQL (better)
\# go run connectmysqlmodel/connectmysqlmodel.go

# HTTP Post
\# go run httppost/httppost.go

# Run an external command
\# go run runexternalcommand/runexternalcommand.go
