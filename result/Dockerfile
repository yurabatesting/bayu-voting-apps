FROM node:18-alpine
WORKDIR /app
COPY server.js .
RUN npm init -y && npm install express pg @opentelemetry/auto-instrumentations-node @opentelemetry/exporter-trace-otlp-proto
CMD ["node", "--require", "@opentelemetry/auto-instrumentations-node/register", "server.js"]