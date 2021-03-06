import React from "react"
import "./App.css"
import Chessboard from "chessboardjsx"
import rough from "roughjs/bundled/rough.cjs"
import Chess from "chess.js"

const roughSquare = ({ squareElement, squareWidth }) => {
    let rc = rough.svg(squareElement)
    const chessSquare = rc.rectangle(0, 0, squareWidth, squareWidth, {
        roughness: 0.2,
        fill: "rgb(236,217,185)",
        bowing: 5,
        fillStyle: "hatched",
        fillWeight: 0.5,
    })
    squareElement.appendChild(chessSquare)
}
const startFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
const baseurl = window.location.host.startsWith("localhost")
    ? "http://localhost:8080/"
    : "https://baidachess.westeurope.cloudapp.azure.com/"

class App extends React.Component {
    constructor() {
        super()

        this.game = new Chess(startFEN)
    }

    state = {
        fen: startFEN,
        result: "NotStarted",
        userBlocked: false,
    }

    playGame = async () => {
        this.setState({ result: "NoResult", userBlocked: true })
        let result = "NoResult"

        while (result === "NoResult") {
            let response = await fetch(
                baseurl +
                    "bestmove?time=" +
                    120 +
                    "&fen=" +
                    encodeURI(this.state.fen)
            ).then((res) => res.json())

            result = response.result
            this.setState({ fen: response.fen, result: response.result })
        }
    }

    onDrop = async ({ sourceSquare, targetSquare }) => {
        if (this.state.userBlocked) {
            return
        }

        // see if the move is legal
        let move = this.game.move({
            from: sourceSquare,
            to: targetSquare,
            promotion: "q",
        })

        // illegal move
        if (move === null) {
            console.log("Illegal move")
            console.log(this.game.board())
            return
        }

        const result = this.game.game_over() ? "GameEnded" : "NoResult"
        const fen = this.game.fen()

        this.setState({
            fen,
            userBlocked: true,
            result,
        })

        if (result === "NoResult") {
            // Play computer move
            let response = await fetch(
                baseurl + "bestmove?time=" + 120 + "&fen=" + encodeURI(fen)
            ).then((res) => res.json())

            this.game = new Chess(response.fen)
            this.setState({
                fen: response.fen,
                result: this.game.game_over() ? "GameEnded" : "NoResult",
                userBlocked: this.game.game_over(),
            })
        }
    }

    render() {
        return (
            <div className="App">
                <Chessboard
                    position={this.state.fen}
                    roughSquare={roughSquare}
                    onDrop={this.onDrop}
                />
                {this.state.result !== "NoResult" && (
                    <button id="start-game" onClick={() => this.playGame()}>
                        START BOT SELF-PLAY
                    </button>
                )}
                {this.state.userBlocked && this.state.result === "NoResult" && (
                    <span>Thinking...</span>
                )}
                {this.state.userBlocked && this.state.result !== "NoResult" && (
                    <span>
                        Game ended, refresh the page to play another one
                    </span>
                )}
            </div>
        )
    }
}

export default App
