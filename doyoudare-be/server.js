import express from "express";
import {
  spotifyApi,
  saveTokens,
  readTokens,
  refreshTokens,
  scheduleTokenRefresh
} from "./spotifyAuth.js";
import { generateRandomString } from "./utilities.js";
import setupRoutes from "./routes.js";
import dotenv from "dotenv";
import open from "open"
dotenv.config();

const app = express();
const PORT = process.env.PORT || 3000;

// Initialization logic for token state
const initializeTokenState = async () => {
  let tokens;
  try {
    tokens = readTokens();
    if (!tokens) throw new Error("No tokens found");
  } catch (error) {
    console.error("Failed to read tokens:", error.message);
    console.log("Please log in.");
    return; // Exit the function, requiring user action to log in
  }

  const currentTime = Date.now();

  // Check if the token is expired or not
  if (!tokens.expirationTime || tokens.expirationTime <= currentTime) {
    console.log("Token is expired or missing. Attempting to refresh...");

    try {
      await refreshTokens();
    } catch (refreshError) {
      console.error("Failed to refresh token:", refreshError.message);
      console.log("Please log in again.");
      // Optionally, you might want to clear the corrupted or outdated tokens.json here
      return; // Exit the function, requiring user action to log in
    }
  } else {
    // If token is not expired, set the access and refresh tokens
    spotifyApi.setAccessToken(tokens.accessToken);
    spotifyApi.setRefreshToken(tokens.refreshToken);
    console.log("Token is valid. Using existing tokens.");
  }
};

// Immediately invoke the async function to check and refresh token if necessary
initializeTokenState()
  .then(() => {
    console.log("Token initialization complete.");
  })
  .catch((error) => {
    console.error("Failed to initialize token state:", error);
  });

app.get("/login", async (req, res) => {
  const state = generateRandomString(16);
  const authorizeURL = spotifyApi.createAuthorizeURL(
    [
      "user-read-private",
      "user-read-email",
      "user-read-playback-state",
      "playlist-read-private",
    ],
    state,
    true
  );
  await open(authorizeURL); // Opens the URL in the default browser
  res.send("Login initiated. Please check your browser.");
});

// Callback route
app.get("/callback", async (req, res) => {
  const { code } = req.query;
  try {
    const data = await spotifyApi.authorizationCodeGrant(code);
    const { access_token, refresh_token, expires_in } = data.body;

    spotifyApi.setAccessToken(access_token);
    spotifyApi.setRefreshToken(refresh_token);

    // Save tokens to file
    saveTokens({
      accessToken: access_token,
      refreshToken: refresh_token,
      expiresIn: expires_in,
      tokenType: data.body.token_type,
    });

    scheduleTokenRefresh();
    res.send("Login successful! You can now close this window.");
  } catch (err) {
    console.error("Something went wrong during the callback:", err);
    res.send("Error during the callback.");
  }
});

// Middleware to check and refresh token if needed
const tokenRefreshMiddleware = async (req, res, next) => {
  const tokens = readTokens();
  if (!tokens || !tokens.accessToken) {
    res.redirect("/login");
  } else {
    spotifyApi.setAccessToken(tokens.accessToken);
    spotifyApi.setRefreshToken(tokens.refreshToken);
    next();
  }
};

app.use(tokenRefreshMiddleware);

// Setup routes
setupRoutes(app);

app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
});
