
# interfaces

This module contains the interfaces which are required to access the code in other packages of this module.


## How To Use

Here are the interfaces and their use case:

- ```ZKComparable```: This interface is used to compare two values of a type. It contains only 1 method:
    - ```Equals```: This method returns a bool based on whether the value of receiver type is equal to the value in parameter.


- ```DbArgs```: This interface is used to provide the columns to which data should be inserted to or retrieved to. It contains only 1 method:
    - ```GetAllColumns```: This method returns a slice where each element is the db column to which data is transferred to/from.
```
go get github.com/zerok-ai/zk-utils-go/interfaces
```

To import this package do:
```
interfaces "github.com/zerok-ai/zk-utils-go/interfaces"
```



    