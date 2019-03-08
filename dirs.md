Current version 

````
pkg
├── skaffold
│   ├── apiversion
│   ├── build
│   │   ├── cache
│   │   ├── 
│   │   │   ├── bazel
│   │   │   ├── docker
│   │   │   ├── jib
│   │   │   └── kaniko
│   │   │       └── sources
│   │   ├── environments
│   │   │   └── gcb
│   │   │   └── local
│   │   │   └── incluster
│   │   ├── plugin # whatever is in shared currently + PluginBuilder implementation 
│   │        └── builder.go 
│   ├── color
│   ├── config
│   ├── constants
│   ├── deploy
│   │   ├── kubectl
│   │   └── testdata
│   │       └── foo
│   │           └── templates
│   ├── docker
│   │   └── testdata
│   ├── event
│   │   └── proto
│   ├── gcp
│   ├── kubernetes
│   │   └── context
│   ├── runner
│   ├── schema
│   │   ├── defaults
│   │   ├── latest
│   │   ├── util
│   │   ├── v1alpha1
│   │   ├── v1alpha2
│   │   ├── v1alpha3
│   │   ├── v1alpha4
│   │   ├── v1alpha5
│   │   ├── v1beta1
│   │   ├── v1beta2
│   │   ├── v1beta3
│   │   ├── v1beta4
│   │   └── v1beta5
│   ├── sources
│   ├── sync
│   │   └── kubectl
│   ├── tag
│   ├── test
│   │   └── structure
│   ├── update
│   ├── util
│   ├── version
│   ├── warnings
│   ├── watch
│   └── yamltags
└── webhook
    ├── constants
    ├── gcs
    ├── github
    ├── kubernetes
    └── labels

````