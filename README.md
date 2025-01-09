# Meal Planner

An app to help plan meals and shops.

Add ingredients, create meals with those ingredients, and produce a shopping list for those meals, helpfully organised by category.

## Run locally
1. Install [Docker](https://www.docker.com/get-started/)
1. Run `docker compose up`
1. The application should now be running on http://localhost:4200 

## Local development
1. Install node 20 and npm. I recommend using [nvm](https://github.com/nvm-sh/nvm) for this - once installed, run `nvm install` in this repo to get the relevant version.
1. Install [Go](https://go.dev/doc/install)
1. Run `npm install` 
1. Run `npm run serve` to start the API and UI locally with hot reloading.

## Features
- Add ingredients
- Create and edit meals with ingredients in different units/quantities
- Add meals to shops
- View ingredients needed for entire shop, grouped by category
- Add ingredients to basket, to tick them off from the shopping list
- In progress: uploading meals from CSV

## Technical notes
- This repo uses [nx](https://nx.dev/) for managing the API and client apps in a monorepo
- The API uses [hallgren/event-sourcing](https://github.com/hallgren/eventsourcing), and stores events in a local SQLite database.
