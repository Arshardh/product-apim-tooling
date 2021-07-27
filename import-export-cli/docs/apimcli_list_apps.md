## apimcli list apps

Display a list of Applications in an environment specific to an owner

### Synopsis



Display a list of Applications of the user in the environment specified by the flag --environment, -e

apimcli list apps -e dev
apimcli list apps -e dev -o sampleUser
apimcli list apps -e prod -o sampleUser -u admin
apimcli list apps -e staging -o sampleUser -u admin -p admin


```
apimcli list apps [flags]
```

### Options

```
  -e, --environment string   Environment to be searched (default "default")
  -h, --help                 help for apps
  -o, --owner string         Owner of the Application
  -p, --password string      Password
  -u, --username string      Username
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimcli list](apimcli_list.md)	 - List APIs/Applications in an environment or List the environments

