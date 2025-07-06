import { Container, ContainerChild, Ticker } from "pixi.js";
import { Action, Orc } from "./Orc";
import { Connection, DeltaState, GameInitContext, Packet } from "./Connection";
import { PlayerController } from "./PlayerController";
import { Projectile } from "./Projectile";
import { Point } from "./Math/Point";
import { Direction } from "./Math/Direction";
import { None } from "./types/Option";

const CLIENT_UPDATE_FREQUENCY = 5;
const CLIENT_UPDATE_FREQUENCY_DEBUG = 2000;

export class Game {
    gamestate: GameState;
    ticker: Ticker;
    conn: Connection;
    playerController!: PlayerController;

    constructor(
        ticker: Ticker,
        conn: Connection,
        stage: Container<ContainerChild>,
        ctx: GameInitContext) {
        this.gamestate = new GameState(stage);
        this.ticker = ticker;
        this.conn = conn;
        //server needs to tell us by now how many orcs to spawn
        for (let i = 0; i < ctx.orcCount; i++) {
            console.log(ctx.idToOrcLocations);
            const clientId = ctx.idToOrcLocations[i].id;
            const location = ctx.idToOrcLocations[i].point;
            console.log(location);
            const orc = new Orc(clientId, location); 
            orc.populateSprites();
            orc.setDefaultSprite();
            orc.gamestateRef = this.gamestate;
            stage.addChild(orc);
            this.gamestate.otherOrcs.push(orc);
            if (clientId == conn.id) {
                const controller = new PlayerController(conn.id, orc);
                controller.setupInputHandler();
                this.playerController = controller;
                console.log("Setup Player Controller");
            }
        }

        setInterval(() => {
            if (this.playerController.inputQueue.length == 0) {
                return;
            }
            const packet: Packet = {
                Id: this.conn.id,
                Type: "G",
                Size: this.playerController.inputQueue.length,
                Data: this.playerController.inputQueue
            };
            this.conn.sendPacket(packet);
            this.playerController.inputQueue.length = 0;
        }, CLIENT_UPDATE_FREQUENCY);
    }

    run() {
        this.ticker.add((time) => {
            // get game state update
            const deltaStateOption = this.conn.getNextDeltaState();
            if (deltaStateOption instanceof None) {
                return;
            } 
            console.log("Delta Queue Size: ", this.conn.deltaQueue.length);
            //handling player's own orc
            // this.playerController.updateOrc(time.deltaMS);
            this.gamestate.update(time.deltaMS, deltaStateOption.value);

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

    update(deltaTime: number, newState: DeltaState) {
        for (let i = 0; i < this.otherOrcs.length; i++) {
            console.log("Updating from server");
            this.otherOrcs[i].updateFromServer(newState.orcDeltas[i]);
        }
        for (let i = 0; i < this.projectiles.length; i++) {
            console.log("Projectile " + i + "Updating");
            this.projectiles[i].update(deltaTime);
        }
    }
}