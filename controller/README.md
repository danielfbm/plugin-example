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
3. Create a [`manager.go`](extension/manager.go) to manage loading of plugins

### 4. Create some plugins

1. Create [`plugins`](plugins) folder
2. Implement specific plugins

### 5. Change CRD

1. Change CRD [specs and status](api/v1alpha1/plugin_types.go)
2. Regenerate crd, deepcopy etc.

```
make
make manifests
```
### 6. Add a plugin loader runner

1. Add `plugin-folder` flag to [main.go](main.go) file and initiate `extension.Manager`
2. Create and implement [`plugin_loader.go`](controllers/plugin_loader.go)
3. Add loader to mgr on [main.go](main.go#L80) file

