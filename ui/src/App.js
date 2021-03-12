import React from "react"
import "./App.css"
import Chessboard from "chessboardjsx"
import rough from "roughjs/bundled/rough.cjs"

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

class App extends React.Component {
    state = {
        fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
        result: "NotStarted",
    }

    constructor() {
        super()
    }

    componentDidMount() {
        this.playGame()
    }
    playGame = async () => {
        this.setState({ result: "NoResult" })
        let result = "NoResult"

        while (result === "NoResult") {
            let response = await fetch(
                "http://localhost:8080/bestmove?fen=" +
                    encodeURI(this.state.fen)
            ).then((res) => res.json())

            result = response.result
            this.setState({ fen: response.fen, result: response.result })
        }
    }

    render() {
        return (
            <div className="App">
                <Chessboard
                    position={this.state.fen}
                    roughSquare={roughSquare}
                />
            </div>
        )
    }
}

export default App
