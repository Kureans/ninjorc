import { Container, ContainerChild, Ticker } from "pixi.js";
import { Orc, Point } from "./Orc";
import { Connection } from "./Connection";
import { PlayerController } from "./PlayerController";
import { Projectile } from "./Projectile";

export class Game {
    gamestate: GameState;
    ticker: Ticker;
    conn: Connection;
    playerController: PlayerController;
    constructor(ticker: Ticker, conn: Connection, stage: Container<ContainerChild>) {
        this.gamestate = new GameState(stage);
        this.ticker = ticker;
        this.conn = conn;
        const orc = new Orc(1, new Point(100, 100));
        orc.populateSprites();
        orc.setDefaultSprite();
        orc.gamestateRef = this.gamestate;
        stage.addChild(orc);
        const controller = new PlayerController(1, orc);
        controller.setupInputHandler();
        this.playerController = controller;
    }

    run() {
        this.ticker.add((time) => {
            // console.log("Running Game");

            //handling player's own orc
            this.playerController.updateOrc(time.deltaMS);
            this.gamestate.update(time.deltaMS);
            // this.playerController.orc.setupSprite();
            //handling other orcs
        })
    }
}

export class GameState {
    otherOrcs: Orc[] = [];
    projectiles: Projectile[] = [];
    stageRef: Container<ContainerChild>;

    constructor(stage: Container<ContainerChild>) {
        this.stageRef = stage;
    }

    update(deltaTime: number) {
        for (let i = 0; i < this.otherOrcs.length; i++) {
            this.otherOrcs[i].update(deltaTime);
        }
        for (let i = 0; i < this.projectiles.length; i++) {
            console.log("Projectile " + i + "Updating");
            this.projectiles[i].update(deltaTime);
        }
    }
}