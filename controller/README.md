# Basic plugin controller


## Intro

1. Load locally managed plugins in the controller
2. Accept plugins over the network (in-cluster)
3. Give a simple health status for each plugin
4. Make it available over CRD
5. Allow different implementations for the same plugin type

## Steps

### 1. Init 

```
kubebuilder init --domain danielfbm.github.com --repo github.com/danielfbm/plugin-example/controller
```

### 2. Create api

```
kubebuilder create api --group plugins --version v1alpha1 --kind Plugin 
```

### 3. Specify the plugins

For the sake of simplicity just two simple plugin: `Foo` and `Bar`

```
type Foo interface {
    Foos() (string, error)
}

type Bar interface {
    Bars() []string
}
```

1. Create a folder `extension` and add the plugin code for the interfaces in the respective files
2. Implement the basic interface and RPC client and server
3. Create a `manager.go` to manage loading of plugins