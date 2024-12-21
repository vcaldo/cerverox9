# cerverox9-bot
A Discord bot that tracks voice channel activities and user presence, storing the data in InfluxDB. A companion Telegram bot forwards real-time notifications about voice events.

## Features

- Discord bot monitoring:
    - Voice channel join/leave events
    - Stream and webcam activity
    - User online status
- Data storage in InfluxDB
- Real-time Telegram notifications for:
    - Users joining/leaving voice channels
    - Stream starts/stops
    - Webcam activation
- `/status` Telegram handler for Discord voice channel stats

## Requirements

- Docker and Docker Compose
- Discord Bot Token with intents:
        - GUILD_MEMBERS
        - GUILD_VOICE_STATES
        - GUILDS
        - GUILD_PRESENCES
- Telegram Bot Token
- Telegram Channel or Group ID

## Setup

1. Clone the repository
2. Configure environment files:
         ```bash
         cp secrets.env.sample secrets.env
         cp influx-secrets.env.sample influx-secrets.env
         ```
3. Update environment files with your configuration:
    - In `influx-secrets.env`:
      - Set InfluxDB credentials
      - Configure database parameters
    - In `secrets.env`:
      - Add your Discord bot token
      - Set your Telegram bot token
      - Configure Telegram channel/group ID
      - Set InfluxDB credentials anda parameters

## Usage

Start the application:
```bash
docker compose up --build
```
