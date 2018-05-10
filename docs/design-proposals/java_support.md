# Skaffold for Java Developers


Status: Draft

Version: Alpha | Beta | GA

Contributors: @loosebazooka, @balopat, @dgageot, @patflynn, @coollog, @jstratchan

Owner: @balopat

## Motivation

We would like to enable Skaffold workflows for Java Developers, so that they can have a great user experience deploying to Kubernetes. 

## User Experience

* Users can leverage local caches of their build tools if they have one installed - including the first run of the skaffold command 
* Users can plugin their existing "source to image" strategy 
* Users don't have to modify their existing maven/gradle/docker setup to work with Skaffold 
* Users can deploy with skaffold even if there is no available maven/gradle/docker on the given machine

### Use Cases

Note: the following use case is supported out of the box: 

* **Use case**: Build image with (multi-stage) Dockerfile & deploy app to K8s `skaffold run/dev`

However this use case is not optimal in terms of local cache usage and is not applicable in certain target circumstances. 

The following use cases are proposed: 

* **Use case**: kaniko maven/gradle builder
    * **Given**: 
      * a maven/gradle project with no Dockerfile 
      * and no maven/gradle installation on the machine
      * and no docker installation on the machine 
    * **When**: 
      * skaffold.yaml defines `maven: kaniko` builder
      * and user runs `skaffold dev/run`
    * **Then**: 
      * `skaffold` generates a simple multi-stage Dockerfile for _kaniko based_ maven/gradle build for JAR/WAR
      * and deploys the image to K8S as usual
    
    
* **Use case**: in-docker maven/gradle builder
    * **Given**: 
      * a maven/gradle project with no Dockerfile 
      * and no maven/gradle installation on the machine
      * and Docker daemon is installed on the machine
    * **When**: 
      * skaffold.yaml defines `maven: in-docker` builder
      * and user runs `skaffold dev/run`
    * **Then**: 
      * `skaffold` generates a simple multi-stage Dockerfile for _in-docker_ maven/gradle build for JAR/WAR
      * and deploys the image to K8S as usual
    

* **Use case**:  maven/gradle project with no Dockerfile    
    * **Given**: 
      * a maven/gradle project with no Dockerfile 
      * and existing maven/gradle installation on the machine
      * and Docker daemon is installed on the machine
    * **When**: 
      * skaffold.yaml defines `maven: {}` or `gradle: {}` builder
      * and user runs `skaffold dev/run`
    * **Then**: 
      * *local maven/gradle builder* _locally_ builds executable JAR/WAR file 
      * and builds it into a distroless Java/Jetty image 
      * and deploys the image to K8S as usual


* **Use case**:  maven/gradle project with no Dockerfile   
    * **Given**: 
      * maven/gradle project with pre-defined or generated Dockerfile:
      * and existing maven/gradle installation on the machine
      * and Docker daemon is installed on the machine
    * **When**: 
      * skaffold.yaml defines `maven: {}` or `gradle: {}` builder
      * and user runs `skaffold dev/run`
    * **Then**: 
      *  maven/gradle builder _locally_ builds executable JAR/WAR 
      * and builds the user-specified docker image 
      * and deploys the image to K8S as usual

* **Use case**: Build image with generic script 
    * **Given**: 
      * maven/gradle project with with existing "source to image" script 
      * dependencies of the "source to image" script are available on the machine 
    * **When**: 
      * skaffold.yaml defines `command: <build_script>` builder
      * and user runs `skaffold dev/run`
    * **Then**: 
      * `skaffold` builds image with generic command 
      * and deploys the image to K8S as usual

* A new use case generic to all languages:
    * **Use case**: Trigger build & deploy based on *user-specified watch paths* using `skaffold dev`
    
* **Use case**: Hotswap modified classes in deployed java instances from a maven/gradle project using `skaffold debug` 


## Proposal for Java "Source to K8S" strategy

### Maven builder 


//TODO: describe the new tags (command, maven, gradle) in detail

//TODO: describe WAR->Jetty, JAR->openjdk logic  
 

 

## Alternatives considered

//TODO 

### Inferring source folders from Maven 

### Inferring artifacts from Maven

### Inferring source folders from Gradle

### Inferring artifacts from Gradle

### Skaffold config in Maven/Gradle project definition 

// TODO: explain: https://github.com/GoogleContainerTools/skaffold/issues/526

