FROM node:20-alpine as build
WORKDIR build
ENV NX_DAEMON=false
COPY package*.json ./
RUN npm ci
COPY apps/client apps/client
COPY tsconfig.base.json nx.json ./
RUN npx nx build client

FROM node:20-alpine
WORKDIR /app
ENV NX_DAEMON=false
COPY --from=build build/dist/apps/client/package.json ./
RUN npm install --omit=dev
COPY --from=build build/dist/apps/client ./
EXPOSE 4200
CMD ["npx", "next", "start", "-p", "4200"]
