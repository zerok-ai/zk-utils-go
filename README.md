
# zk-utils-go

This project is used as a central repository containing common models and fuctions to be used across zeroK


## How To Use


# zkcommon

This module contains the common utility functions which are useful for all the repositories.

This is a private repo and hence can only be used inside the organisation. To access this repo:

- Set GOPRIVATE environment variable

    ```
    export GOPRIVATE=github.com/zerok-ai/zk-utils-go
    ```


- If github personal access token is not configured already configre it and add to ~/.netrc file
  Ref: https://www.digitalocean.com/community/tutorials/how-to-use-a-private-go-module-in-your-own-project#providing-private-module-credentials-for-https

- Run the following command into the root dir of project where you require this library

    ```
    go get github.com/zerok-ai/zk-utils-go
    ``` 

- Add ```import "github.com/zerok-ai/zk-utils-go/<required_package>"``` to access the library


# zk-utils-go

This project is used as a central repository containing common models and functions to be used across zeroK.


## Appendix

Reference on how to use private modules in go: https://www.digitalocean.com/community/tutorials/how-to-use-a-private-go-module-in-your-own-project

We can use the command ```protoc --proto_path=. --go_out=. --go_opt=paths=source_relative ./*.proto``` to update the pb.go files. 
The command assumes the current folder contains the proto files and the pb.go files will be generated in the same folder. 