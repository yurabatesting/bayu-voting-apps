FROM python:3.9-slim
WORKDIR /app
COPY app.py .
RUN pip install flask redis
CMD ["python", "app.py"]