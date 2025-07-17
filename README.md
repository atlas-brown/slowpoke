# Slowpoke: End-to-end Throughput Optimization Modeling for Microservice Applications

## Using Slowpoke

To use Slowpoke for a microservice application, changes in three places are needed

#### Application Code
The application needs to include Slowpoke's runtime as dependency, and register request handler to the Slowpoke wrapper. Currently we only support Go (See example in XX)

#### Change to Kubernetes Deployment Configuration
Slowpoke takes the Kubernetes YAML files 
