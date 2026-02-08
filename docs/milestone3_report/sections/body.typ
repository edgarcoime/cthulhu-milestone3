= Milestone 3 Summary

The primary focus for Cthulhu's Milestone 3 was the Lifecycle Microservice, the service would ensure that buckets would expire and auto delete at its stated expiration time. All backend service communication (auathentication, filemanager, and the new lifecycle service) now use a single, modern communication layer called gRPC. The development environment setup and run configurations have now been consolidated into a single configuration file that outlines all the tooling and prerequisites used to build the project. Each component/service can now be run using containers or with one command that will orchestrate and start up the entire Cthulhu platform. 

== Communication Between Services

All services now use a unified transmission layer to communicate in a consitent and effecient way. In this milestone I implemented gRPC which is a high-performance open-source Remote Procedure Call framework to request actions or data between services. Each service now has a clear contract or set of operations it can perform and consume outlined by its protobuffer files. 
Protobuf files also have the added benefit of being able to compile to any other target programming language. This gives me the added flexibility to experiment with creating services with other languages and frameworks since there is only one source of truth that I can use for its communication layer. 

== Development and Build Consistency

To prevent issues such as "it works on my machine" and to keep generated code in sync (protobufs and sqlc compilation), the project now uses a root configuration file (mise.toml) that pins the versions of the main development tools (e.g. Go, compilers, and code generators). Tasks such as code generation are also run through this setup. Each service's README also points back to this root configuration so anyone cloning the repository can install the same tools and run the same commands to build and generate code.

== File Lifecycle Service

One of the main goals of CTHULHU is time-limited file availability and to enforce this the File lifecycle service was added. It is responsible for recording when each shared bucket (and its files) should expire, and for running a periodic cleanup process. Other services can also request for a bucket's expiration to be deleted if required. 

== Containerization and Single-Command Run

Each microservice (authentication, filemanager, lifecycle, gateway, and web client) now has its own Dockerfile. Images are also build in multi-stage, secure manner, and run run as a non-root user where applicable. Services that require a database are also using a dedicated volume so that data persists between restarts and is not lost when containers are recreated.

There is also a root Docker Compose at the root of the project that orchestrates everything together. One command from the project root builds and starts all services. Service locations and other setting are set in one place and the root README explains how to setup the required environment variables and run the full stack.

= Conclusion

Milestone 3 is set up to keep Cthulhu on track and has consolidated how backend services communicate (gRPC), adds a dedicated lifecycle service to enforce automatic file and bucket expiration, and establishes a stadard development environment through a shared tool configuration. By having every service containerized I can also leverage the power of docker compose which allows the entire CTHULHU stack to be run in a single command. These changes will allow CTHULHU to be more maintainable and consistent which will allow for a faster and more efficient development process for the next two milestones in the future.