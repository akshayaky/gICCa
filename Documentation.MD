# IRCcloud API in GO

The following are the packages I made to develop a CLI IRCcloud app:

### 1) Login.go

#### GetAuthToken()
Returns  the authorization token, which is needed in login.
No arguments are needed.
```
curl -X POST "https://www.irccloud.com/chat/auth-formtoken" --header "content-length: 0"
{"success":true,"token":"1397241172970.9a87f9s7fad9f7s9f8a7fa9sd77"}
```

#### GetSessionId()
Returns the session id which is stored as cookies, which will be used for all future requests.
Needs the auth token,email address and password of the IRCcloud user.

```
curl -d email=XXX -d password=XXX -d token=1397241172970.9a87f9s7fad9f7s9f8a7fa9sd77 --header "content-type: application/x-www-form-urlencoded" --header "x-auth-formtoken: 1397241172970.9a87f9s7fad9f7s9f8a7fa9sd77" "https://www.irccloud.com/chat/login"
{"success":true,"session":"e07910a06a086c83ba41827aa00b26ed","uid":1}
```


### 2) Cookie.go
#### SetCookie()
Sets the session id as cookie for a say post request.
This function needs to be modified setting any cookie for any request.
Returns the http client.
Currently only requires only the session ID as argument.


### 3) Backlog.go
#### EndpointConnection()
Connects to the end-point server. This request provides a link which is used to get the backlog.Calls the GetBacklog() function. Returns the names and Connection  ID's of the various connections the user has previously connected to.
Needs the session ID as argument.

```
curl https://www.irccloud.com/chat/stream --cookie 'session=7.0850938509385098503'
```
From this request we get the *stream Id* which is used for getting the backlog. This request should be sent within 10 seconds, as after that the stream ID expires.

####  GetBacklog()
Gets the whole backlog and filters out the names and IDs of the connections.Called by the EndpointConnection() function to get the connection IDs
Needs the stream ID and session ID as argument. Stream id is obtained from the endpoint connection .

```
curl https://www.irccloud.com/chat/backlog/<streamid> --cookie 'session=7.82039956563409384'
```

### 4) Connection.go
####  Reconnect()
Irrespective of the name this function can both make a particular connection active and inactive. A message can be sent to a recipient only if that particular connection is active.
Needs session ID, cid of a connection and a string "re"(for reconnection) or "dis" for disconnection.


### 5) Message.go
#### Say()
Sends a message to a given channel or user.
needs cid and session ID

```curl
curl -d "msg=/msg%20nickorchannel%20hello+world&to=%2A&cid=1&session=7.598435934539045739458" --cookie "session=7.598435934539045739458" https://www.irccloud.com/chat/say
{"_reqid":0,"success":true,"name":"nickorchannel","cid":1,"type":"open_buffer"}
```


The rough flow chart for the irc client:

![Flow Chart](https://drive.google.com/uc?id=1XwZMiY0cEzxfdrRQdqVW05GJWsIp60mz)
