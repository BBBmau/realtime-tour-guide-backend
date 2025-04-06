This repo represents the different GCP services that are used in my explore.ai android app.

It's broken up into 3 scenarios that required accessing a VPC that's been provisioned on Google Cloud

`session` - allows a request for an ephemeral key to create a realtime session with openai. Sets up a signed JWT for the API Gateway to take in a received request.

`notify` - takes in a request w/ signed JWT to create a chat response for audio notification that includes recommendations for nearby areas. The response would be a small summary with the audio being a .WAV or .mp3 file for audio notification (just like how GPS notifications are)

`conversationHistory` - communicates with a Database inside a VPC network. uses Signed JWT for access along with user info (perhaps an integration of auth0)
