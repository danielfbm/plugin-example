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
3. Add loader to mgr on [main.go](main.go#L80) file, set default hclog
4. Compile plugins, manager, and run


### 7. Implement controller

This controller will only do one thing:

1. Check if the plugin is accessible and adds a condition
2. Check which implementation it serves and add a condition for each
3. Check the implementation on [plugin_controller.go](controllers/plugin_controller.go)

### 8. Build and deploy

1. Build controller and local plugins: `make docker-build`. Check [`Dockerfile`](Dockerfile) for specific build instructions
2. Install CRD: `make install`
3. Push you image and deploy: `make deploy`
4. verify that everything is working fine with `kubectl`

  `kubectl get pods --all-namespaces` to check if the controller is up and running
  `kubectl get plugins` to check if the local plugins are loaded and checked

5. Build the foobar plugin `CGO_ENABLED=0 GOOS=linux go build -o plugins/foobar/bin/foobar plugins/foobar/main.go`
6. Build the docker image `docker build -t danielfbm/foobarplugin -f plugins/foobar/Dockerfile plugins/foobar
7. Deploy on kubernetes using kubectl `run` and `expose`:

*PS: Depending on your kubectl version `expose` behaviour could be different, adapt accordingly

```
kubectl run foobar --image=danielfbm/foobarplugin:latest --image-pull-policy=Never --env BASIC_PLUGIN=hello --port 7000
kubectl expose deploy/foobar --port=7000 --target-port=7000
```

8. Check plugins and status: `kubectl get plugins`

:)