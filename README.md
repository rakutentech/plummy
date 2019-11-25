# plummy

Plummy is a CLI tool for automatically running Java image generators (such as
PlantUML) in a background daemon, without requiring to restart the JVM on every
call. This can significantly speed up static site generation when using such
Java tools to process specialized markup.

## Supported tools

* [PlantUML](http://plantuml.com/)
* [Ditaa](https://github.com/stathissideris/ditaa)

## Alternatives

### For UML diagrams

#### Running PlantUML directly

Each invocation takes at least a few seconds, which makes it painful to use
for static site generation, where you may have to process hundreds of UML
diagrams at once.

#### [PlantUML Server](https://github.com/plantuml/plantuml-server)

PlantUML server also runs a stable PlantUML instance as daemon, which cuts down
the invocation time significantly, but it has many limitations that Plummy is
trying to fix:

- PlantUML Server is built as a classic (non-embedded) Java Servlet
  application, and distributed as a WAR file, so it's far from trivial to setup
  and requires setting up your own server, complete with a few nice XML files
  and Web 1.0 era Java EE knowledge. The official docker container makes
  running PlantUML server a lot simpler, but it's still too heavyweight for
  what should have been handled by one fat jar file.
- There is no CLI wrapper, so client tools (like site-generators or editor
  plugins) have to integrate with the server themselves, and each has to be
  specially set up to know the port where PlantUML is served.
- The user needs to run the server manually and keep track of the server state.
- PlantUML Server does not support external resources very well (e.g. included
  files, fonts, images)

#### [mermaid](https://mermaidjs.github.io/)

Mermaid is a Javascript UML diagram generator. It supports a dialect that is
similar to PlantUML, but different enough to be incompatible for even the most
basic diagrams.

Since mermaid runs on Node.js, invocation start-up times are much faster than
running PlantUML directly. Unfortunately, mermaid also lacks many features
supported by PlantUML. The lack of advanced styling features, makes it harder
to generate clean-looking diagrams for more complex cases.

