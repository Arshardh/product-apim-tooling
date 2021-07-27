## apimcli reset-user

Reset user of an environment

### Synopsis



Reset user data of a particular environment (Clear the entry in env_keys_all.yaml file)

Examples:
apimcli reset-user -e dev
apimcli reset-userreset-user -e staging


```
apimcli reset-user [flags]
```

### Options

```
  -e, --environment string   Clear user details of an environment (default "default")
  -h, --help                 help for reset-user
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimcli](apimcli.md)	 - CLI for Importing and Exporting APIs and Applications

