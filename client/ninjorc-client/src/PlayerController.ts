import { Action, Direction, Orc } from "./Orc";

export class PlayerController {
    id: number;
    orc: Orc;
    hasReceivedInput: boolean;

    constructor(id: number, orc: Orc) {
        this.id = id;
        this.orc = orc;
        this.hasReceivedInput = false;
    }

    setupInputHandler() {
        // const canvas = document.getElementById("app");
        // if (canvas == null) {
        //     throw new Error("Canvas element doesn't exist.")
        // }
        document.addEventListener("keydown", (event) => {
            if (event.key == "ArrowDown") {
                this.orc.Direction = Direction.DOWN;
                console.log("down");
            }
            else if (event.key == "ArrowUp") {
                this.orc.Direction = Direction.UP;
                console.log("up");
            }
            else if (event.key == "ArrowLeft") {
                this.orc.Direction = Direction.LEFT;
                console.log("left");
            }
            else if (event.key == "ArrowRight") {
                this.orc.Direction = Direction.RIGHT;
                console.log("right");
            }
            else if (event.key == "Z" || event.key == "z") {
                this.orc.Action = Action.MELEE;
                console.log("swing");
            }
            else if (event.key == "X" || event.key == "x") {
                this.orc.Action = Action.PROJECTILE;
                console.log("fire");
            }
            this.hasReceivedInput = true;
        })
    }

    updateOrc(deltaTime: number) {
        if (this.hasReceivedInput) {
            this.orc.update(deltaTime);
            this.hasReceivedInput = false;
        }
    }
}