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
    "game_id": "37f77f29-5865-477c-9bef-ca01d22483c2",
    "guess_results": [
        {
            "guess_id": "950b240e-7f55-4585-b51e-37bd950c590f",
            "guess": "reach",
            "result": [ 1, 0, 0, 0, 0 ]
        },
        {
            "guess_id": "e0086c47-c094-41d0-bbea-dc587b385e9e",
            "guess": "titan",
            "result": [ 1, 1, 0, 0, 0 ]
        },
        {
            "guess_id": "e2a233f8-2e1d-4d74-8ce6-df0f419053f2",
            "guess": "rajah",
            "result": [ 1, 0, 0, 0, 0 ]
        },
        {
            "guess_id": "2b84c35e-abc2-4aef-87c9-3a90eb77b976",
            "guess": "imbue",
            "result": [ 1, 0, 0, 0, 0 ]
        }
    ],
    "guess_durations_ns": [ 1001874500, 1000829600, 1000972000, 1000716200 ]
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
    "player_id": "96abf359-1c33-4293-aeff-51bef796ee1d",
    "results": {
        "match_id": "54a479dc-c3f0-461b-a97a-22d8aafcee4e",
        "players": [
            {
                "player_id": "96abf359-1c33-4293-aeff-51bef796ee1d",
                "definition": {
                    "name": "solvo",
                    "description": "the magnificent"
                },
                "games_played": [
                    {
                        "game_id": "37f77f29-5865-477c-9bef-ca01d22483c2",
                        "guess_results": [
                            {
                                "guess_id": "950b240e-7f55-4585-b51e-37bd950c590f",
                                "guess": "reach",
                                "result": [ 1, 0, 0, 0, 0 ]
                            },
                            {
                                "guess_id": "e0086c47-c094-41d0-bbea-dc587b385e9e",
                                "guess": "titan",
                                "result": [ 1, 1, 0, 0, 0 ]
                            },
                            {
                                "guess_id": "e2a233f8-2e1d-4d74-8ce6-df0f419053f2",
                                "guess": "rajah",
                                "result": [ 1, 0, 0, 0, 0 ]
                            },
                            {
                                "guess_id": "2b84c35e-abc2-4aef-87c9-3a90eb77b976",
                                "guess": "imbue",
                                "result": [ 1, 0, 0, 0, 0 ]
                            },
                            {
                                "guess_id": "fb68b729-5445-4a41-a7a3-c4a3ee305759",
                                "guess": "funky",
                                "result": [ 0, 0, 0, 0, 0 ]
                            },
                            {
                                "guess_id": "fce2024e-a2b3-4d4e-9806-2b7c44d1b5b5",
                                "guess": "kneel",
                                "result": [ 0, 0, 0, 0, 0 ]
                            }
                        ],
                        "guess_durations_ns": [ 1001874500, 1000829600, 1000972000, 1000716200, 1000844600, 1000734400 ]
                    }
                ]
            },
            {
                "player_id": "ad249978-1d49-4e24-bff7-bd1e1c7fcca6",
                "definition": {
                    "name": "schwordler",
                    "description": "the brave"
                },
                "games_played": [
                    {
                        "game_id": "37f77f29-5865-477c-9bef-ca01d22483c2",
                        "guess_results": [
                            {
                                "guess_id": "f8be7f56-8efa-4dfd-b9df-af3765f7c28e",
                                "guess": "crane",
                                "result": [ 0, 2, 0, 0, 0 ]
                            },
                            {
                                "guess_id": "3a804ca6-8f15-4e12-b561-81b82cc4aad2",
                                "guess": "droit",
                                "result": [ 2, 2, 2, 2, 2 ]
                            }
                        ],
                        "correct": true,
                        "guess_durations_ns": [ 1458600, 1386800 ]
                    }
                ]
            }
        ],
        "games": [ { "game_id": "37f77f29-5865-477c-9bef-ca01d22483c2", "answer": "droit" } ],
        "rounds_per_game": 6,
        "letters_per_word": 5
    }
}
```

### Releasing
Any time you add a new tag, a release is automatically built and deployed to github. [goreleaser](https://goreleaser.com/) is awesome.

