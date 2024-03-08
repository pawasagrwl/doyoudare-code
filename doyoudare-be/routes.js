import { spotifyApi } from './spotifyAuth.js';

const setupRoutes = (app) => {
  app.get("/me", async (req, res) => {
    try {
      const me = await spotifyApi.getMe();
      res.json(me.body);
    } catch (error) {
      console.error("Failed to fetch user profile:", error);
      res.sendStatus(500);
    }
  });

  app.get("/currently-playing", async (req, res) => {
    try {
      const data = await spotifyApi.getMyCurrentPlayingTrack();
      res.json(data.body);
    } catch (error) {
      console.error("Error getting currently playing track:", error);
      if (error.statusCode === 401) {
        // If the token expired, refresh it and try again
        await refreshTokens();
        return res.redirect("/currently-playing");
      }
      res.status(500).json({ error: "Internal server error" });
    }
  });

  app.get("/playlists", async (req, res) => {
    try {
      const data = await spotifyApi.getUserPlaylists();
      res.json(data.body);
    } catch (error) {
      console.error("Error getting playlists:", error);
      if (error.statusCode === 401) {
        // If the token expired, refresh it and try again
        await refreshTokens();
        return res.redirect("/playlists");
      }
      res.status(500).json({ error: "Internal server error" });
    }
  });

  // You can add more route handlers here
};

export default setupRoutes;

