FROM python:3.9-slim
WORKDIR /app
COPY app.py .
RUN pip install flask redis opentelemetry-distro opentelemetry-exporter-otlp
RUN opentelemetry-bootstrap -a install
CMD ["opentelemetry-instrument", "python", "app.py"]