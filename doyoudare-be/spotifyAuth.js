// ESM syntax for importing modules
import SpotifyWebApi from "spotify-web-api-node";
import fs from "fs";
import path from "path";
import dotenv from "dotenv";

dotenv.config();

const tokensPath = path.join(process.cwd(), "tokens.json");

// Initialize SpotifyWebApi
const spotifyApi = new SpotifyWebApi({
  redirectUri: process.env.SPOTIFY_REDIRECT_URI,
  clientId: process.env.SPOTIFY_CLIENT_ID,
  clientSecret: process.env.SPOTIFY_CLIENT_SECRET
});

const saveTokens = (tokens) => {
  const currentTime = Date.now();
  const expiresInMs = tokens.expiresIn * 1000; // Convert expiresIn from seconds to milliseconds
  const expirationTime = currentTime + expiresInMs; // Calculate the exact expiration time
  
  const tokensWithExpirationTime = {
    ...tokens,
    expirationTime, // Save the exact expiration time
  };
  
  fs.writeFileSync(tokensPath, JSON.stringify(tokensWithExpirationTime, null, 2));
};

const readTokens = () => {
  try {
    const tokensData = fs.readFileSync(tokensPath, { encoding: 'utf-8' });
    return JSON.parse(tokensData);
  } catch (error) {
    // Log removed to clean up console output; return null to indicate no tokens available
    return null;
  }
};
// Helper function to refresh access and refresh tokens
const refreshTokens = async () => {
  try {
    const data = await spotifyApi.refreshAccessToken();
    const { access_token, refresh_token, expires_in } = data.body;

    // Sometimes Spotify does not return a new refresh token, so we reuse the old one if that's the case.
    const tokens = readTokens();
    saveTokens({
      accessToken: access_token,
      refreshToken: refresh_token ? refresh_token : tokens.refreshToken,
      expiresIn: expires_in,
      tokenType: tokens.tokenType, // Reusing the existing token type
    });

    spotifyApi.setAccessToken(access_token);
    if (refresh_token) {
      spotifyApi.setRefreshToken(refresh_token);
    }

    console.log("The access token has been refreshed.");
  } catch (error) {
    console.error("Could not refresh access token, please login again.");
  }
};

// Schedule token refresh based on expiration, requires passing the refreshTokens function and the readTokens function
const scheduleTokenRefresh = async () => {
  const tokens = readTokens(); // This should now work as expected
  if (tokens && tokens.expiresIn) {
    // Set a timeout to refresh the token a minute before it expires
    setTimeout(async () => {
      await refreshTokens(); // Direct call without passing as parameter

      // Call scheduleTokenRefresh again for the next cycle
      scheduleTokenRefresh();
    }, (tokens.expiresIn - 60) * 1000); // Adjust time as needed
  }
};

export { spotifyApi, saveTokens, readTokens, refreshTokens, scheduleTokenRefresh };