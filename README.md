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
Once all players are finished, the engine will send you the results of everyone in the match. No response is required, except maybe to message your friends to brag. `player_id` represents your ID, look for the corresponding player in the `players` array to see how you went.
#### Request:
```json
{
    "player_id": "fc1ea6b0-b5c7-4844-a34c-09acaaea8865",
    "results": {
        "match_id": "2f2ab3b8-ed36-4208-b28e-9b05e3750bee",
        "players": [
            {
                "player_id": "fc1ea6b0-b5c7-4844-a34c-09acaaea8865",
                "definition": {
                    "name": "solvo",
                    "description": "the magnificent"
                },
                "games_played": [
                    {
                        "game_id": "ce593360-622f-4f22-9d32-f5b5163ef770",
                        "guess_results": [
                            {
                                "guess_id": "2cd00859-8889-4a68-a711-d283b454ce2c",
                                "guess": "gully",
                                "result": [ 0, 0, 0, 0, 0 ]
                            },
                            {
                                "guess_id": "67347ea2-388f-4204-9957-dea8ae011a99",
                                "guess": "revue",
                                "result": [ 0, 0, 0, 0, 0 ]
                            },
                            {
                                "guess_id": "065e2c3e-d7c7-456f-803d-4e8f7a5df01a",
                                "guess": "scare",
                                "result": [ 0, 1, 2, 0, 0 ]
                            },
                            {
                                "guess_id": "b1612e80-8f4b-4b7d-914f-a54c77f8689f",
                                "guess": "nerve",
                                "result": [ 0, 0, 0, 0, 0 ]
                            },
                            {
                                "guess_id": "51735551-8b36-4254-a49a-6038def72c72",
                                "guess": "creed",
                                "result": [ 2, 0, 0, 0, 0 ]
                            },
                            {
                                "guess_id": "9839c3c6-525c-475f-9acd-6c8ca2b2d803",
                                "guess": "chasm",
                                "result": [ 2, 2, 2, 0, 1 ]
                            }
                        ],
                        "guess_durations_ns": [ 1001797000, 1001040500, 1000737700, 1000751600, 1000742800, 1001051500 ]
                    }
                ]
            },
            {
                "player_id": "908f08bf-922a-4d3e-82a6-946e4548910f",
                "definition": {
                    "name": "schwordler",
                    "description": "the brave"
                },
                "games_played": [
                    {
                        "game_id": "ce593360-622f-4f22-9d32-f5b5163ef770",
                        "guess_results": [
                            {
                                "guess_id": "6cd20f1a-898c-4e11-9bc2-71eae9ae9f76",
                                "guess": "crane",
                                "result": [ 2, 0, 2, 0, 0 ]
                            },
                            {
                                "guess_id": "83116bd4-5daf-4529-988b-d226ddc5a74e",
                                "guess": "chasm",
                                "result": [ 2, 2, 2, 0, 1 ]
                            },
                            {
                                "guess_id": "9b636e82-6fae-49eb-80a7-c706dc930d58",
                                "guess": "chalk",
                                "result": [ 2, 2, 2, 0, 0 ]
                            },
                            {
                                "guess_id": "7bba4f1e-4d6e-4c8a-bdb7-ff5e32c99b96",
                                "guess": "chaff",
                                "result": [ 2, 2, 2, 0, 0 ]
                            },
                            {
                                "guess_id": "4c9fd232-453c-492d-ade6-c6b8013c3338",
                                "guess": "champ",
                                "result": [ 2, 2, 2, 2, 2 ]
                            }
                        ],
                        "correct": true,
                        "guess_durations_ns": [ 1309700, 978400, 1045300, 1093200, 1063800 ]
                    }
                ]
            }
        ],
        "games": [ { "game_id": "ce593360-622f-4f22-9d32-f5b5163ef770", "answer": "champ" } ],
        "rounds_per_game": 6,
        "letters_per_word": 5
    }
}
```

### Releasing
Any time you add a new tag, a release is automatically built and deployed to github. [goreleaser](https://goreleaser.com/) is awesome.

