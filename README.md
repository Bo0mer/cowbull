# cowbull
The infamous game cows &amp; bulls

## User's guide
### Running it
The game consist ot two parts - the game server and the client content.
If you want to install and run it, just go get it and start the `cowbull`
process in the root of this repo.

```bash
go get github.com/Bo0mer/cowbull
cd $GOPATH/src/github.com/Bo0mer/cowbull
cowbull
```

The server will be available at `http://127.0.0.1:8080/` by default.
If you want to fine-tune your settings, this parameters are supported during
startup as command-line arguments:
```
 -address string
    Server address. (default "127.0.0.1:8080")
 -skip-origin-check
    Skip Origin header check upon WebSocket connection negotiation.
```

Example:
```
cowbull -address "10.244.0.34:6060"
```

### Playing it
The game has 4 modes - you can play against the computer or a real person,
and you can be either a thinker or a guesser.

You select **the role of your opponent** in the combo-box within the game.
If you choose to play against a thinker, you must fill up the number of digits
input field. Otherwise, you'll be prompted to think for a number.

ATM, if you start a game with guessers, you will be prompted to enter a list
of desired opponent names. The list **must be comma separated**. Note also that
each player will be prompted to guess your number once it's time for his turn.
Unfortunately, before the first ask a player won't be notified that she is
participating in a game (PRs are welcome:).


## Developer's guide
### Running the tests
The tests use `ginkgo`, so you'll need to `go get`-it.
```
go get github.com/onsi/ginkgo
cd $GOPATH/src/github.com/Bo0mer/cowbull/
ginkgo -r --race
```

Or, you can ofcourse use go test
```
go test ./...
```
