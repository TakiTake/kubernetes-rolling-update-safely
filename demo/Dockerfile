FROM openjdk:10-jdk-slim AS builder
ENV GRADLE_VERSION=4.9
RUN apt update -y && apt install -y wget unzip
RUN wget https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip
RUN mkdir -p /opt/gradle
RUN unzip -d /opt/gradle gradle-${GRADLE_VERSION}-bin.zip
COPY . /usr/src/myapp
WORKDIR /usr/src/myapp
RUN /opt/gradle/gradle-${GRADLE_VERSION}/bin/gradle build

FROM openjdk:10-jre-slim
COPY --from=builder /usr/src/myapp/build/libs/demo-*.jar demo.jar
CMD ["java", "-jar", "demo.jar"]
