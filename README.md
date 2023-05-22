
# zk-utils-go

This project is used as a central repository containing common models and fuctions to be used across zeroK


## How To Use

This is a private repo and hence can only be used inside the organisation. To access this repo:

    1. Set GOPRIVATE environment variable
        export GOPRIVATE=github.com/zerok-ai/zk-utils-go


    2. If github personal access token is not configured already configre it and add to ~/.netrc file
        Ref: https://www.digitalocean.com/community/tutorials/how-to-use-a-private-go-module-in-your-own-project#providing-private-module-credentials-for-https

    3. Run go get github.com/zerok-ai/zk-utils-go into the root dir of project where you require this library

    4. Add import "github.com/zerok-ai/zk-utils-go/<required_package>" to access the library 


# zk-utils-go

This project is used as a central repository containing common models and fuctions to be used across zeroK


## Appendix

Reference on how to use private modules in go: https://www.digitalocean.com/community/tutorials/how-to-use-a-private-go-module-in-your-own-project