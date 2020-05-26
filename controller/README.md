# Basic plugin controller


## Intro

1. Load locally managed plugins in the controller
2. Accept plugins over the network (in-cluster)
3. Give a simple health status for each plugin
4. Make it available over CRD

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

type Bars interface {
    Bars() []string
}
```