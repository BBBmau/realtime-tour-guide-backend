import express from "express";
import dotenv from "dotenv";

dotenv.config();

const app = express();
const PORT = 8080;

// Add console log to show server is starting
console.log('Starting server...');

// An endpoint which would work with the client code above - it returns
// the contents of a REST API request to this protected endpoint
app.get("/", (req, res) => {
  res.send("Hello World");
});

app.get("/session", async (req, res) => {
  const { location = 'La Jolla, CA', destination = 'Irvine, CA', initial_desired_service = 'Coffee Shops' } = req.query;

  console.log('Received request to /session');
  try {
    const r = await fetch("https://api.openai.com/v1/realtime/sessions", {
      method: "POST",
      headers: {
        "Authorization": `Bearer ${process.env.OPENAI_API_KEY}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        model: "gpt-4o-mini-realtime-preview",
        voice: "verse",
        input_audio_transcription: {
          model: "whisper-1",
        },
        instructions: `You are a passenger in a vehicle currently on a roadtrip 
          who's currently in ${location} and is heading to ${destination}. 
          You are currently looking for ${initial_desired_service}.

          Your job is to provide recommendations for the most unique and breath-taking 
          places to visit along the way. For each recommendation, please include:
          - Name of the location
          - Detailed description
          
          Focus on hidden gems and memorable stops that would enhance the journey. 
          Please ensure that you are acting like a passenger such as a close-friend 
          or relative that joined the trip and loves adventure.
          `
      }),
    });
    const data = await r.json();
    console.log('Received response from OpenAI');
    res.send(data);
  } catch (error) {
    console.error('Error:', error);
    res.status(500).send({ error: error.message });
  }
});

// Modified listen to include callback
app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
});