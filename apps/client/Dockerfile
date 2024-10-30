FROM node:alpine as builder
WORKDIR /build
ENV NX_DAEMON=false
RUN npx next telemetry disable
COPY ./package.json ./package-lock.json ./tsconfig.base.json ./nx.json ./
RUN npm install
COPY ./apps/client ./apps/client
RUN npx nx build client --skip-nx-cache

FROM node:alpine
WORKDIR /app
ENV NX_DAEMON=false
COPY --from=builder /build/dist/apps/client/package.json ./
RUN npm install --omit=dev
COPY --from=builder /build/dist/apps/client ./
EXPOSE 4200
CMD ["npx", "next", "start", "-p", "4200"]
