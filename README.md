# I am no longer working on this but PRs opened will be reviewed!

This repo represents the different GCP services that are used in my realtime tour guide android app.

It's broken up into 3 scenarios that required accessing a VPC that's been provisioned on Google Cloud

`session` - allows a request for an ephemeral key to create a realtime session with openai. Sets up a signed JWT for the API Gateway to take in a received request.

`notify` - takes in a request w/ signed JWT to create a chat response for audio notification that includes recommendations for nearby areas. The response would be a small summary with the audio being a .WAV or .mp3 file for audio notification (just like how GPS notifications are)

`conversationHistory` - communicates with a Database inside a VPC network. uses Signed JWT for access along with user info (perhaps an integration of auth0)

## Setting up on GCP

Run the following command, you can also utilize `.env` or provide environment variables as part of the `app.yaml`. Ensure that you've setup your credentials properly that allows for app engine deployment.

`gcloud app deploy --set-env-vars GOOGLE_MAPS_API_KEY=your_key,OPENAI_API_KEY=your_key`

You will need to provide an `app.yaml` file, it can be as simple as:
```yaml
runtime: nodejs20
instance_class: F1
env_variables:
  NODE_ENV: production
```

You'll also want to provide API keys in `.env` for `OPENAI_API_KEY` and `GOOGLE_MAPS_API_KEY`
