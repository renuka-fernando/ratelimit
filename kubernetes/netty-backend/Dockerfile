FROM eclipse-temurin:8u345-b01-jre-focal
ENV JAVA_OPTS=""

COPY netty-http-echo-service-* netty-http-echo-service.jar
CMD java $JAVA_OPTS -jar netty-http-echo-service.jar
