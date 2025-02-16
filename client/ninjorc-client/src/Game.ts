import { Ticker } from "pixi.js";
import { Orc } from "./Orc";
import { Connection } from "./Connection";
import { PlayerController } from "./PlayerController";

export class Game {
    gamestate: GameState;
    ticker: Ticker;
    conn: Connection;
    playerController: PlayerController;
    constructor(ticker: Ticker, conn: Connection, controller: PlayerController) {
        this.gamestate = new GameState();
        this.ticker = ticker;
        this.conn = conn;
        this.playerController = controller;
    }

    run() {
        this.ticker.add((time) => {
            console.log("Running Game");

            //handling player's own orc
            this.playerController.updateOrc(time.deltaMS);
            this.playerController.orc.setupSprite();
            //handling other orcs
        })
    }
}

class GameState {
    playerOrc: Orc;
    otherOrcs: Orc[];

    constructor() {

    }
}