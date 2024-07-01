FROM node:alpine as builder
WORKDIR /build
RUN npx next telemetry disable
COPY ./package.json ./package-lock.json ./tsconfig.base.json ./nx.json ./
RUN npm install
COPY ./apps/client ./apps/client
RUN npx nx build client

FROM node:alpine
WORKDIR /app
COPY --from=builder /build/dist/apps/client ./
COPY --from=builder /build/node_modules ./node_modules
EXPOSE 4200
CMD ["npx", "next", "start", "-p", "4200"]
