# GophKeeper

You can use client app to create owen vault

Here is instruction

1. Sign Up

register new user with name 'testName' and password 'testPassword'

```bash
./client signUp http://localhost:8080 testName testPassword
```

4. Sign In for receiving token

login user by user name 'testName' and password 'testPassword' for receiving auth token

```bash
./client signIn http://localhost:8080 testName testPassword
```

Response
```bash
eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIxIn0.bMFEsrvtCxd5i3SMn3E_8HcRx6RzNfTX2PI1eWXJsbNUbeG_VaEpf9trTcm4KsYqYp_wpLzMYEYKQCtQykb4lQ
```

3. Create new vault with token

create vault with key 'k1' and value 'v1' with token "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIxIn0.iok4gCKCJP3d7vXMUyDFEvgZQ2-hyyk85gvHvmoGkx5-aMByqGyq8GjfNcpgY1Mc31xRn-d0BHnmy3H1kwNWXg"

```bash
./client saveVault http://localhost:8080 "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIxIn0.iok4gCKCJP3d7vXMUyDFEvgZQ2-hyyk85gvHvmoGkx5-aMByqGyq8GjfNcpgY1Mc31xRn-d0BHnmy3H1kwNWXg" k1 v1
```
4. Get vault with token

get vault by id '1' with token "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIxIn0.iok4gCKCJP3d7vXMUyDFEvgZQ2-hyyk85gvHvmoGkx5-aMByqGyq8GjfNcpgY1Mc31xRn-d0BHnmy3H1kwNWXg"

```bash
./client getVault http://localhost:8080 "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIxIn0.iok4gCKCJP3d7vXMUyDFEvgZQ2-hyyk85gvHvmoGkx5-aMByqGyq8GjfNcpgY1Mc31xRn-d0BHnmy3H1kwNWXg" 1
 1037  history
```