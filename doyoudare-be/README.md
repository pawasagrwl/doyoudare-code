# Spotify Authentication Server README

## Overview

This project is a Node.js server that leverages the Spotify Web API to authenticate users via OAuth 2.0, allowing access to personal Spotify data such as profile information, currently playing tracks, and playlists. The server is built using Express and handles token management, including refresh mechanisms to maintain access to Spotify's resources.

## Features

- **OAuth 2.0 Authentication with Spotify**: Securely authenticate users with Spotify and gain authorization to access their Spotify data.
- **Token Management**: Automatically handle access and refresh token lifecycles, including token refresh before expiration.
- **User Profile Information**: Access and display Spotify user profile details.
- **Currently Playing Track**: Retrieve information about the user's currently playing track.
- **User Playlists**: List playlists available in the user's Spotify account.
- **Navigation Interface**: A simple HTML page providing links to various endpoints for easy navigation and testing.

## Prerequisites

Before running this project, you need the following:

- **Node.js and npm**: Ensure you have Node.js (version 14 or higher) and npm installed on your system. You can download them from [nodejs.org](https://nodejs.org/).
- **Spotify Developer Account**: You'll need a Spotify Developer account to create an app and obtain the `Client ID` and `Client Secret`. Register or log in at [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/).
- **Project Setup on Spotify**: Within the Spotify Developer Dashboard, create an app to get your `Client ID` and `Client Secret`. Set the Redirect URI to match the callback route configured in your server (e.g., `http://localhost:3000/callback`).

## Setup and Configuration

1. **Clone the Repository**: Clone this repository to your local machine or download the source code.

2. **Install Dependencies**: Navigate to the root directory of the project in your terminal and run:

```bash
npm install
```

3. **Configure Environment Variables**: Create a `.env` file in the root directory with the following content, replacing the placeholders with your actual Spotify app credentials and preferred port number:

```env
SPOTIFY_CLIENT_ID=<Your-Spotify-Client-ID>
SPOTIFY_CLIENT_SECRET=<Your-Spotify-Client-Secret>
SPOTIFY_REDIRECT_URI=http://localhost:3000/callback
```

## Running the Server

1. **Start the Server**: In the root directory of your project, run:

```bash
npm start
```

This command starts the Express server on the port specified in your `.env` file (default: 3000).

2. **Access the Main Page**: Open a web browser and go to `http://localhost:3000`. You'll see a main page with links to different endpoints:

- **Login**: Initiates the Spotify authentication flow.
- **Profile Information**: Displays the user's Spotify profile information.
- **Currently Playing Track**: Shows details about the track currently playing on the user's Spotify account.
- **My Playlists**: Lists the playlists in the user's Spotify account.

3. **Authenticate with Spotify**: Click on the **Login** link to authenticate with Spotify and authorize the application to access your Spotify data.

## After Authentication

Once authenticated, you can navigate through the provided links to view your Spotify profile information, currently playing track, and playlists. The server handles access token refresh automatically, ensuring continuous access to your Spotify data.

## Token Refresh and Handling Expiration

The server is designed to automatically refresh the access token using the refresh token before it expires. In case the server restarts, it checks the token's validity and refreshes it if necessary during startup. If automatic refresh fails or the tokens are invalid, the server will prompt for re-authentication.

## Additional Information

For more details on the Spotify Web API, OAuth 2.0 authentication flow, and available endpoints, refer to the [Spotify Web API Documentation](https://developer.spotify.com/documentation/web-api/).

---
Enjoy exploring and building with Spotify's rich set of resources and data!