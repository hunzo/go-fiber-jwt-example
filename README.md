# go-fiber-jwt-example
* required jq
```git initsh
sudo apt install jq
```
* get access token
```sh
TOKEN=$(curl http://localhost:8080/token/testpayload | jq -r .token)
```
* show access token
```sh
echo $TOKEN
```
* call api
```sh
curl http://localhost:8080/private -H "Authorization: Bearer $TOKEN" | jq
```
