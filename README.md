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
Each `guess_results` object also comes with the start and finish time, plus an ID that correlates to the header `guessID` that gets sent with each request. I've omitted from this JSON for brevity. Solvo prints out the full body when you run him.
#### Request:
```json
{
    "game_id": "3bead1b6-cd41-4bd0-9ec0-9b451319efba",
    "guess_results": [
        {
            "guess": "tense",
            "result": [ 1, 0, 0, 0, 0 ]
        },
        {
            "guess": "finer",
            "result": [ 0, 1, 0, 0, 1 ]
        },
        {
            "guess": "unset",
            "result": [ 0, 0, 0, 0, 2 ]
        },
        {
            "guess": "cable",
            "result": [ 0, 0, 0, 0, 0 ]
        },
        {
            "guess": "deity",
            "result": [ 2, 0, 1, 1, 0 ]
        }
    ],
    "guess_durations_ns": [ 1002925700, 1001538400, 1000738000, 1000947200, 1000960600 ]
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
As with the request, all objects come with an ID, start, and end time that has been omitted for brevity where unnecessary. Check Solvo for the exact body. 
#### Request:
```json
{
    "player_id": "9fda863c-f303-47da-a8f0-35b0b84b1abe",
    "results": {
        "match_id": "777f785b-d2f9-4467-990e-e2f90efe3b52",
        "players": [
            {
                "player_id": "9fda863c-f303-47da-a8f0-35b0b84b1abe",
                "definition": {
                    "name": "solvo",
                    "description": "the magnificent"
                },
                "games_played": [
                    {
                        "game_id": "3bead1b6-cd41-4bd0-9ec0-9b451319efba",
                        "guess_results": [
                            {
                                "guess": "tense",
                                "result": [ 1, 0, 0, 0, 0 ]
                            },
                            {
                                "guess": "finer",
                                "result": [ 0, 1, 0, 0, 1 ]
                            },
                            {
                                "guess": "unset",
                                "result": [ 0, 0, 0, 0, 2 ]
                            },
                            {
                                "guess": "cable",
                                "result": [ 0, 0, 0, 0, 0 ]
                            },
                            {
                                "guess": "deity",
                                "result": [ 2, 0, 1, 1, 0 ]
                            },
                            {
                                "guess": "deter",
                                "result": [ 2, 0, 1, 0, 1 ]
                            }
                        ],
                        "guess_durations_ns": [ 1002925700, 1001538400, 1000738000, 1000947200, 1000960600, 1001443500 ]
                    }
                ],
            },
            {
                "player_id": "b6f3c0ed-c885-4253-8d23-725def324c55",
                "definition": {
                    "name": "schwordler",
                    "description": "the brave"
                },
                "games_played": [
                    {
                        "game_id": "3bead1b6-cd41-4bd0-9ec0-9b451319efba",
                        "guess_results": [
                            {
                                "guess": "crane",
                                "result": [ 0, 2, 0, 0, 0 ]
                            },
                            {
                                "guess": "droit",
                                "result": [ 2, 2, 2, 2, 2 ]
                            }
                        ],
                        "correct": true,
                        "guess_durations_ns": [ 2135600, 1794700 ]
                    }
                ],
            }
        ],
        "games": [ { "game_id": "3bead1b6-cd41-4bd0-9ec0-9b451319efba", "answer": "droit" } ],
        "rounds_per_game": 6,
        "letters_per_word": 5
    }
}
```

### Releasing
Any time you add a new tag, a release is automatically built and deployed to github. [goreleaser](https://goreleaser.com/) is awesome.

