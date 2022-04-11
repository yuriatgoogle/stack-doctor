FROM python:3.9-slim-bullseye

ENV PROJECT_ID=thegrinch-demos

COPY ./requirements.txt /app/requirements.txt
COPY ./application_default_credentials.json /app/key.json

ENV GOOGLE_APPLICATION_CREDENTIALS=/app/key.json

WORKDIR /app

RUN pip install -r requirements.txt

COPY . /app

ENTRYPOINT [ "opentelemetry-instrument" ]

CMD [ "--traces_exporter", "gcp_trace", "--exporter_gcp_trace_project_id", "$PROJECT_ID", "python", "backend.py" ]