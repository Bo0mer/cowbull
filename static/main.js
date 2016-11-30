$(function() {

    var $settingsDiv = $('.settingsDiv');
    var $gameDiv = $('.gameDiv');
    var $gameLog = $('.gameLogDiv');

    initView();

    function initView() {
        $('.playButton').click(clickPlay);
        $('.numberInput').keydown(keydownNumber);
    }

    function clickPlay(event) {
        var digits;
        var playerRole;
        var againstAI; 
        var opponents;

        $digitsInput = $('.digitsInput');
        digits = parseInt(cleanInput($digitsInput.val().trim()))

        var opponentName = prompt("Enter name of your opponent");
        var opponentId = getPlayerId(opponentName);
        if (opponentId === "") {
            alert("Player " + opponentName + " not found.");
            return
        }

        switch ($('#opponentSelect').val()) {
        case "ai_thinker":
            playerRole = "guesser"
            againstAI = true;
            opponents = [];
            break;
        case "thinker":
            playerRole = "guesser";
            againstAI = false;
            opponents = [opponentId];
            break;
        }

        beginGame(againstAI, digits, playerRole, opponents);
    }
    
    function keydownNumber(event) {
        if (event.which === 13)  {
            number = cleanInput($('.numberInput').val().trim());
            sendGuess(number);
            $(this).val('');
        }
    }

    function initGameField(playerRole) {
        $settingsDiv.fadeOut();
        $gameDiv.show();

        var $numberInput = $('.numberInput');
        if (playerRole == "guesser") {
            $numberInput.show();
        } else {
            $numberInput.fadeOut();
        }
    }

    function resetGameField() {
        $settingsDiv.show();
        $gameDiv.fadeOut();
        $gameLog.empty();
    }

    function showGuessRequest(digitsCount) {
        var logEntry = "The number has " + digitsCount.toString() + " digits.";
        gameLog(logEntry);
    }

    function showGuessResult(cows, bulls) {
        var logEntry = "There are " + cows + " cows and " + bulls +" bulls.";
        gameLog(logEntry); 
    }

    function showTryRequest(number) {
        var logEntry = "Remote players has a guess " + number.toString();
        gameLog(logEntry);
    }

    function showTryResponse(cows, bulls) {
        var logEntry = "The guess has " + cows + " cows and " + bulls + " bulls.";
        gameLog(logEntry);
    }

    function showPlayers(players) {
        $playersDiv = $('.playersDiv');
        $playersDiv.empty();
        for (var i = 0; i < players.length; i++) {
            var displayName = players[i].name;
            if (displayName === "" || displayName === undefined) {
                displayName = players[i].id;
            }
            $playersDiv.append('<span>'+displayName+'</span><br/>');
        }
    }

    function promptForNumber() {
        return prompt("You have been challenged. Pick your number!");
    }

    function promptForName() {
        return prompt("Please choose your in-game name");
    }

    function showGameEnd(won) {
        if (won) {
            alert("You have just WON!!!");
        } else {
            alert("The remote player guessed your number.");
        }
        resetGameField();
    }

    function showDisconnected() {
        alert("You have been disconnected. Please reload the page to connect again.");
    }

    function gameLog(message) {
        // TODO(ivan): this is a security hole
        $gameLog.append('<span>' + message + '</span><br/>');
    }

    function cleanInput (input) {
        return $('<div/>').text(input).text();
    }

    /////////////// controller ///////////////////

    var socket;
    var inGame = false;
    var currentNumber;
    var currentNumberDigits;

    var connectedPlayers;
    
    initController();

    function initController() {
        socket = new WebSocket("ws://127.0.0.1:8080/websocket");
        socket.onopen = onOpen;
        socket.onerror = onError;
        socket.onmessage = onMessage;
        socket.onclose = onClose;
    }

    function setName(name) {
        var name = {
            name: "name",
            data: JSON.stringify({name: name}),
        };
        socket.send(JSON.stringify(name));    
    }

    function beginGame(againstAI, digits, playerRole, opponents) {
        initGameField(playerRole);
        inGame = true;

        var play = {
            name: "play",
            data: JSON.stringify({
                AI: againstAI,
                digits: digits,
                role: playerRole,
                opponents: opponents,
            }),
        };
        socket.send(JSON.stringify(play))
    }

    function endGame(won) {
        inGame = false;
        showGameEnd(won);
    }

    function sendGuess(number) {
        var guess = {
            name: "guess",
            data: JSON.stringify({number: number}),
        };
        socket.send(JSON.stringify(guess));
    }

    function onOpen(event) {
        console.log("socket opened");
        setName(promptForName());
        var connect = {
            name: "connect",
            data: "",
        };
        socket.send(JSON.stringify(connect));
    }

    function onError(event) {
        console.log("error with socket");
    }

    function onMessage(event) {
        var msg = JSON.parse(event.data);
        switch (msg.name) {
        case "guess":
            console.log("guess message recved");
            handleGuess(msg.data);
            break;
        case "tell":
            console.log("tell message recved");
            handleTell(msg.data)
            break;
        case "think":
            console.log("think message recved");
            handleThink(msg.data);
            break;
        case "try":
            console.log("try message recved");
            handleTry(msg.data);
        case "players":
            console.log("players message recved");
            handlePlayers(msg.data);
        }
    }

    function onClose(event) {
        showDisconnected();
    }

    function handleGuess(data) {
        if (!inGame) {
            return;
        }
        var guess = JSON.parse(data);
        currentNumberDigits = guess.digits;

        showGuessRequest(guess.digits); 
    }

    function handleTell(data) {
        if (!inGame) {
            return;
        }

        var cowsbulls = JSON.parse(data); 
        showGuessResult(cowsbulls.cows, cowsbulls.bulls); 

        if (cowsbulls.bulls === currentNumberDigits) {
            endGame(true);
        }
    }

    function handleThink(data) {
        if (inGame) {
            return;
        } 

        inGame = true;
        initGameField("thinker");

        currentNumber = parseInt(promptForNumber());
        currentNumberDigits = currentNumber.toString().length;

        var think = {
            name: "think",
            data: JSON.stringify({digits: currentNumberDigits}),
        };
        socket.send(JSON.stringify(think));
    }

    function handleTry(data) {
        var tryRequest = JSON.parse(data);
        tryNumber = tryRequest.number;

        showTryRequest(tryNumber);
        var cowsbulls = getCowsBulls(currentNumber, tryNumber);
        showTryResponse(cowsbulls.cows, cowsbulls.bulls);
        
        var tryResponse = {
            name: "try",
            data: JSON.stringify(cowsbulls),
        };
        socket.send(JSON.stringify(tryResponse));

        if (cowsbulls.bulls === currentNumberDigits) {
            endGame(false);
        }
    }

    function handlePlayers(data) {
        var players = JSON.parse(data);
        connectedPlayers = players;
        showPlayers(players);
    }

    function getCowsBulls(number, guess) {
        number = number.toString();
        guess = guess.toString();
        var cowsbulls = {cows: 0, bulls: 0};

        if (number.length != guess.length) { return cowsbulls;  }

        for (var i = 0; i < number.length; i++) {
            if (number.indexOf(guess[i]) == i) { cowsbulls.bulls++;  }
            else if (number.indexOf(guess[i]) != -1) { cowsbulls.cows++;  }
        }
        return cowsbulls;
    }

    function getPlayerId(name) {
        for (var i = 0; i < connectedPlayers.length; i++) {
            if (connectedPlayers[i].name === name) {
                return connectedPlayers[i].id;
            }
        }
        return "";
    }

});
