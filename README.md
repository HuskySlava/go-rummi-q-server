# Go Rummi Q Server

## Idea
Using Go, create a server to handle a multiplayer tile-based number game inspired by Rummikub®.

> Rummikub® is a registered trademark of Orly Fidelman.
> This project is an independent implementation that follows similar rules but is not affiliated with or endorsed by the trademark owner.

## Features / Goals
1. Handle lobbies:
   * Create a lobby upon a request from frontend
   * Generate a lobby id that can be used to connect players via url / deeplink

2. Handle game state: 
   * Generate a game for lobby 
   * Keep state of the game
   * Handle turns
   * Send state updates to frontend (WebSocket or gRPC)
   * Validate player moves 
   * Anti-Cheat 

## Setup
  * TBD

## Roadmap / TODO
  * TBD
