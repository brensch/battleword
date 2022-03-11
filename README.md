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
    "player_id": "a89483bc-a35d-41d5-a168-31b9c5cee4a4",
    "results": {
        "uuid": "ced8aa1b-68f8-48e0-bf90-9bbdc2e9ee73",
        "players": [
            {
                "id": "a89483bc-a35d-41d5-a168-31b9c5cee4a4",
                "definition": {
                    "name": "schwordler",
                    "description": "the brave"
                },
                "games_played": [
                    {
                        "game_id": "40dd7298-c5fd-4567-b666-04d0a6cd8dc5",
                        "guess_results": [
                            {
                                "guess": "crane",
                                "result": [ 0, 1, 0, 0, 2 ]
                            },
                            {
                                "guess": "louse",
                                "result": [ 0, 0, 1, 0, 2 ]
                            },
                            {
                                "guess": "merge",
                                "result": [ 0, 0, 2, 2, 2 ]
                            },
                            {
                                "guess": "dirge",
                                "result": [ 0, 0, 2, 2, 2 ]
                            },
                            {
                                "guess": "purge",
                                "result": [ 2, 2, 2, 2, 2 ]
                            }
                        ],
                        "correct": true,
                        "guess_durations_ns": [ 1794600, 8304200, 1550200, 798000, 720800 ]
                    }
                ]
            },
            {
                "id": "32ad3a51-eb04-44e7-ae42-1bbefc3bb080",
                "definition": {
                    "name": "solvo",
                    "description": "the magnificent"
                },
                "games_played": [
                    {
                        "game_id": "40dd7298-c5fd-4567-b666-04d0a6cd8dc5",
                        "guess_results": [
                            {
                                "guess": "aider",
                                "result": [ 0, 0, 0, 1, 1 ]
                            },
                            {
                                "guess": "stash",
                                "result": [ 0, 0, 0, 0, 0 ]
                            },
                            {
                                "guess": "guess",
                                "result": [ 1, 2, 1, 0, 0 ]
                            },
                            {
                                "guess": "quash",
                                "result": [ 0, 2, 0, 0, 0 ]
                            },
                            {
                                "guess": "talon",
                                "result": [ 0, 0, 0, 0, 0 ]
                            },
                            {
                                "guess": "poesy",
                                "result": [ 2, 0, 1, 0, 0 ]
                            }
                        ],
                        "guess_durations_ns": [ 1885100, 582700, 449700, 559800, 549800, 503400 ]
                    }
                ]
            }
        ],
        "games": [ { "id": "40dd7298-c5fd-4567-b666-04d0a6cd8dc5", "answer": "purge" } ],
        "rounds_per_game": 6,
        "letters_per_word": 5
    }
}
```

### Releasing
Any time you add a new tag, a release is automatically built and deployed to github. [goreleaser](https://goreleaser.com/) is awesome.

