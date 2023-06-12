
# common

This module contains the common logic to parse command line arguments and load the configuration.


## How To Use

The package contains a struct name ```Args```. This struct currently contains only one field named ```ConfigPath``` which contains the path of configuration file for the application.

There is also a function named ```ProcessArgs``` which takes in a generic type ```T```. This type ```T``` is the struct for the corresponding configuration file which is located at location pointed by ```ConfigPath```.
This function currently supports only one program argument which is the configuration file location. To pass the program arguments, use the following:

```
-c <file_location>
```
Example:
```
-c ./internal/config/config-local.yaml
```



    