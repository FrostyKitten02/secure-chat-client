# Secure chat client
## Dev
### Build
- 
### Generate client
- install tool: go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
- oapi-codegen -generate types,client -o client/client_gen.go -package client secure-chat.yaml