How create custom compatible extensions

# COFF binary
In order to communicate the in-memory loaded dll we pass a golang callback as argument(https://stackoverflow.com/questions/37703779/pass-go-function-to-c-as-callback)
We additionally pass the arguments of our call as an array of bytes 
```golang
//implant/extension/extension.go (line 30)
func Run(extID string, funcName string, arguments []byte, callback func([]byte)) error {
	if ext, found := extensions[extID]; found {
		return ext.Call(funcName, arguments, callback)
	}
	for id, ext := range extensions {
		fmt.Printf("Extension '%s' (%s)", id, ext.GetArch())
	}
	return errors.New("extension not found")
}
```
So to create compatible exported functions we need to respect this pattern

```c
int entrypoint(char* argsBuffer, uint32_t bufferSize, goCallback callback)
{
    int cmd = -1;
    if (bufferSize < 1)
    {
        std::string msg{ "You must provide a command.\n\t0 = stop\n\t1 = start\n\t2 = get logs" };
        callback(msg.c_str(), msg.length());
    }
}
```

## json extension conf files
We need to create .json for the extensions in order that the server can manage the extension file location, parameters types, ...
!!!!! IMPORTANT !!!!!
For COFF, the user must pass arguments in the same order specified in the json, other case the beacon will crash. COFF execution in the same thread and poor error handling.
