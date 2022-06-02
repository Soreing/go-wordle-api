# Demo Go Wordle API
A small project demonstrating APIs written in Go. The project contains an API with different routes to play wordle, and a game manager process that generates new games every 2 hours.

## How to Use
### Getting the Current game
To play the game, you need to get an ID for the current game. Along with the ID, you are given a unix timestamp in milliseconds for the expiry date of the game.
```
GET /game/current
```
```
Output:
{"Id":40,"ExpiresAt":1654184444241}
```
### Creating a Player
Before you can make guesses, you need a player. You can create a player without any required data.
```
POST /player
```
```
Output:
7F11FBEE6AE9DABA
```
### Making a Guess
Once you have a Game Id, a Player and a Word in mind, you can make a guess. The inputs must be sent as x-www-form-urlencoded field by field. If the game is already solved, or all guesses are used, a 400 status is returned with a message.
```
POST /guess
GameId=40&Player=7F11FBEE6AE9DABA&Word=crane
```
```
Output:
Ok (200)
```
### Get Game Status
With a game Id and a player code, you can get your game state with colored letters as follows. Guesses come as an array of strings where each string contains pairs of characters, the capital being the letter and the lower case is the color.
```
GET /game/:GameId/:PlayerCode
```
```
Output:
{
    "Solved": false,
    "GuessCount": 2,
    "Guesses": [
        "bC bR bA bN bE ",
        "bL yU bN bG gS "
    ]
}
```