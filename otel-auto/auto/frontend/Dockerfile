FROM python:3.9-slim-bullseye

COPY ./application_default_credentials.json /app/key.json

ENV PROJECT_ID=thegrinch-demos
ENV backend=http://localhost:8081
ENV GOOGLE_APPLICATION_CREDENTIALS=/app/key.json

COPY ./requirements.txt /app/requirements.txt

WORKDIR /app

RUN pip install -r requirements.txt

COPY . /app

ENTRYPOINT [ "opentelemetry-instrument" ]

CMD [ "--traces_exporter", "gcp_trace", "--exporter_gcp_trace_project_id", "$PROJECT_ID", "python", "frontend.py" ]