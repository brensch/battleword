# Battleword
Wordle is cool right now

## What is this
Battleword is a competition to see who can come up with the fastest/most accurate/shoutiest wordle solver.

Players host an api, then the battleword engine will make a `POST` request to their api with the state of a wordle (starting empty).

The player's api should then respond to the state of the game in the body of the post with their best guess. As soon as the battleword engine hears back from them, it will send the results of their guess in a new request. It will do that until the player's api guesses correctly, or they reach the guess limit.

## Quickstart

1. Download the latest release for your OS and unpack
2. Run `solvo` (double click) - this starts solvo the solver. He will listen for game states from engine.
3. Run `engine` - this starts sending game states to solvo. With every guess solvo makes, engine will send a new request to solvo with the results of his previous guess. Solvo will ignore those results and choose a completely random word to send next. Your solver should do better than solvo.

## Setup
To test your own guesser against the engine, create an api that implements the schema below. Once you've done that, run the engine against the api location of your solver like so:

```
./engine --apis http://localhost:8081
```

You can specify multiple solvers to compete against each other:
```
./engine --apis http://localhost:8081,http://localhost:8080
```

NB these commands are executed in a command line of your choice. Exact syntax may change based on your OS.

## API
This is what all solvers need to implement.
### /guess
The engine will hit your api here with the previous results of a game. You are expected to respond with your best guess.
#### Request:
```json
{
    "game_id": "c975280b-5ff8-47e1-a384-843e803dceef",
    "guess_results": [
        {
            "guess": "crane",
            "result": [ 0, 1, 1, 0, 0 ]
        },
        {
            "guess": "solar",
            "result": [ 0, 0, 0, 1, 1 ]
        },
        {
            "guess": "party",
            "result": [ 0, 1, 1, 0, 0 ]
        },
        {
            "guess": "guava",
            "result": [ 0, 2, 0, 0, 2 ]
        }
    ],
    "guess_durations_ns": [ 880900, 67876700, 2913600, 1261100 ]
}
```
#### Response:
```json
{
	"guess": "rumba",
	"shout": "why is everybody shouting"
}
```
Shouts server no purpose except to intimidate your opponents.

### /ping
In order to get the definition of your character, the engine will ping you. This is also run at the start of each match up to 10 times in order to wake up your server if you're hosting it in serverless land where everything is slightly less reliable.

#### Request:
GET request - no payload
#### Response:
```json
{
	"name": "solvo",
	"description": "the magnificent"
}
```
There will be more things here in the future. stay posted.

### /results
Once all players are finished, the engine will send you the results of everyone in the match. No response is required, except maybe to message your friends to brag. The below body is just one player, if there were multiple you'd see more objects in the `players` array.
#### Request:
```json
{
    "uuid": "e4d817ef-778c-422f-ad50-fd3b749eaefa",
    "players": [
        {
            "id": "da6c26ce-0234-40f0-8d79-3d03249bd770",
            "definition": {
                "name": "schwordler",
                "description": "the brave"
            },
            "games_played": [
                {
                    "game_id": "8ab55e9b-04ef-4d86-adc1-38f522ad2b68",
                    "guess_results": [
                        {
                            "guess": "crane",
                            "result": [ 0, 0, 0, 0, 0 ]
                        },
                        {
                            "guess": "lousy",
                            "result": [ 1, 1, 0, 0, 0 ]
                        },
                        {
                            "guess": "fight",
                            "result": [ 0, 2, 0, 0, 2 ]
                        },
                        {
                            "guess": "pivot",
                            "result": [ 2, 2, 0, 2, 2 ]
                        },
                        {
                            "guess": "pilot",
                            "result": [ 2, 2, 2, 2, 2 ]
                        }
                    ],
                    "guess_durations_ns": [ 610200, 27952000, 2574200, 1087400, 1776400 ]
                },
                {
                    "game_id": "fed87c2e-a2e1-4f2b-842e-8c7e2990f2e9",
                    "guess_results": [
                        {
                            "guess": "crane",
                            "result": [ 1, 0, 0, 0, 0 ]
                        },
                        {
                            "guess": "lousy",
                            "result": [ 0, 1, 0, 1, 0 ]
                        },
                        {
                            "guess": "smith",
                            "result": [ 2, 0, 0, 0, 0 ]
                        },
                        {
                            "guess": "spook",
                            "result": [ 2, 1, 2, 2, 0 ]
                        },
                        {
                            "guess": "swoop",
                            "result": [ 2, 0, 2, 2, 2 ]
                        },
                        {
                            "guess": "scoop",
                            "result": [ 2, 2, 2, 2, 2 ]
                        }
                    ],
                    "guess_durations_ns": [ 1104800, 37972900, 1881100, 915000, 910100, 735900 ]
                },
                {
                    "game_id": "17fc1e3a-233e-4185-aa35-d4a138eec7f1",
                    "guess_results": [
                        {
                            "guess": "crane",
                            "result": [ 0, 0, 1, 0, 1 ]
                        },
                        {
                            "guess": "salty",
                            "result": [ 0, 1, 2, 0, 2 ]
                        },
                        {
                            "guess": "alloy",
                            "result": [ 1, 0, 2, 0, 2 ]
                        },
                        {
                            "guess": "milky",
                            "result": [ 0, 0, 2, 0, 2 ]
                        },
                        {
                            "guess": "delay",
                            "result": [ 2, 2, 2, 2, 2 ]
                        }
                    ],
                    "guess_durations_ns": [ 1642600, 67321300, 713900, 998700, 642200 ]
                }
            ]
        }
    ],
    "games": [
        {
            "id": "8ab55e9b-04ef-4d86-adc1-38f522ad2b68",
            "answer": "pilot"
        },
        {
            "id": "17fc1e3a-233e-4185-aa35-d4a138eec7f1",
            "answer": "delay"
        },
        {
            "id": "fed87c2e-a2e1-4f2b-842e-8c7e2990f2e9",
            "answer": "scoop"
        }
    ],
    "rounds_per_game": 6,
    "letters_per_word": 5
}

```

### Releasing
In order to release, run:

```
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
goreleaser release
```

To do a dry run (no upload): 
```
goreleaser release --snapshot --rm-dist
```

`GITHUB_TOKEN` is a standard PAT from github and needs to be set to upload.

### Identity Federation

To allow a github project to use gcloud resources:

Setup pool:
```bash
gcloud iam workload-identity-pools create "github-pool" \
  --project="battleword" \
  --location="global" \
  --display-name="github-pool"
```

Setup workload:
```bash
gcloud iam workload-identity-pools providers create-oidc "github-provider" \
  --project="battleword" \
  --location="global" \
  --workload-identity-pool="github-pool" \
  --display-name="github-provider" \
  --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.aud=assertion.aud,attribute.repository=assertion.repository" \
  --issuer-uri="https://token.actions.githubusercontent.com"
```

Allow the identity provider to impersonate the service account:

```bash
gcloud iam service-accounts add-iam-policy-binding "github@battleword.iam.gserviceaccount.com" \
  --project="battleword" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/339690027814/locations/global/workloadIdentityPools/github-pool/attribute.repository/brensch/battleword"
```
This is kind of magic and tbh I don't understand it well yet. It is working successfully however and this repo auto deploys to GCP on every commit to the branches `api/dev` and `api/prod`.