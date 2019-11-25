FROM gradle:5.6.0-jdk11 as build
WORKDIR /plummy
COPY build.gradle build.gradle
COPY settings.gradle settings.gradle
COPY src src
COPY conf conf
RUN gradle shadowJar

FROM zul/zulu-openjdk-alpine:11
WORKDIR /plummy
COPY --from=build /plummy/libs/plummy*.jar app.jar

COPY conf conf
EXPOSE 4545
CMD ["java", "-jar", "app.jar"]
