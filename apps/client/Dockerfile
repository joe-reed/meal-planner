FROM node:20-alpine
WORKDIR /app
ENV NX_DAEMON=false
COPY dist/apps/client/package.json ./
RUN npm install --omit=dev
COPY dist/apps/client ./
EXPOSE 4200
CMD ["npx", "next", "start", "-p", "4200"]
