FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-airbyte"]
COPY baton-airbyte /