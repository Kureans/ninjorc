import { Direction } from "./Math/Direction";
import { Action, Orc } from "./Orc";

type TDirection = {
    tag: "D";
    value: Direction;
};

type TAction = {
    tag: "A";
    value: Action;
};

type GameInput = TDirection | TAction;

export class PlayerController {
    id: number;
    orc: Orc;
    hasReceivedInput: boolean;
    keysCurrentlyPressed: number;
    keyMap: Map<Direction, boolean>;
    inputQueue: GameInput[]= [];

    constructor(id: number, orc: Orc) {
        this.id = id;
        this.orc = orc;
        this.hasReceivedInput = false;
        this.keysCurrentlyPressed = 0;
        this.keyMap = new Map<Direction, boolean>();
        this.keyMap.set(Direction.DOWN, false);
        this.keyMap.set(Direction.UP, false);
        this.keyMap.set(Direction.LEFT, false);
        this.keyMap.set(Direction.RIGHT, false);
    }

    setupInputHandler() {
        // TODO: get keyboard events only on app element instead of document
        // const canvas = document.getElementById("app");
        // if (canvas == null) {
        //     throw new Error("Canvas element doesn't exist.")
        // }
        document.addEventListener("keydown", (event) => {
            if (event.key == "ArrowDown") {
                this.orc.Direction = Direction.DOWN;
                if (!this.keyMap.get(Direction.DOWN)) {
                    this.keysCurrentlyPressed++;
                    this.keyMap.set(Direction.DOWN, true);
                }
                console.log("down");
                this.inputQueue.push({tag: "D", value: Direction.DOWN});
            }
            else if (event.key == "ArrowUp") {
                this.orc.Direction = Direction.UP;
                if (!this.keyMap.get(Direction.UP)) {
                    this.keysCurrentlyPressed++;
                    this.keyMap.set(Direction.UP, true);
                }
                console.log("up");
                this.inputQueue.push({tag: "D", value: Direction.UP});
            }
            else if (event.key == "ArrowLeft") {
                this.orc.Direction = Direction.LEFT;
                if (!this.keyMap.get(Direction.LEFT)) {
                    this.keysCurrentlyPressed++;
                    this.keyMap.set(Direction.LEFT, true);
                }
                console.log("left");
                this.inputQueue.push({tag: "D", value: Direction.LEFT});
            }
            else if (event.key == "ArrowRight") {
                this.orc.Direction = Direction.RIGHT;
                if (!this.keyMap.get(Direction.RIGHT)) {
                    this.keysCurrentlyPressed++;
                    this.keyMap.set(Direction.RIGHT, true);
                }
                console.log("right");
                this.inputQueue.push({tag: "D", value: Direction.RIGHT});
            }
            else if (event.key == "Z" || event.key == "z") {
                this.orc.Action = Action.MELEE;
                console.log("swing");
                this.inputQueue.push({tag: "A", value: Action.MELEE});
            }
            else if (event.key == "X" || event.key == "x") {
                this.orc.Action = Action.PROJECTILE;
                console.log("fire");
                this.inputQueue.push({tag: "A", value: Action.PROJECTILE});
            }
        })

        document.addEventListener("keyup", (event) => {
            if (event.key == "ArrowDown") {
                if (this.keyMap.get(Direction.DOWN)) {
                    this.keysCurrentlyPressed--;
                    this.keyMap.set(Direction.DOWN, false);
                }
            }
            else if (event.key == "ArrowUp") {
                if (this.keyMap.get(Direction.UP)) {
                    this.keysCurrentlyPressed--;
                    this.keyMap.set(Direction.UP, false);
                }
            }
            else if (event.key == "ArrowLeft") {
                if (this.keyMap.get(Direction.LEFT)) {
                    this.keysCurrentlyPressed--;
                    this.keyMap.set(Direction.LEFT, false);
                }
            }
            else if (event.key == "ArrowRight") {
                if (this.keyMap.get(Direction.RIGHT)) {
                    this.keysCurrentlyPressed--;
                    this.keyMap.set(Direction.RIGHT, false);
                }
            }

            if (this.keysCurrentlyPressed == 0) {
                this.orc.resetDirection();
            }
        })
    }

    updateOrc(deltaTime: number) {
        this.orc.update(deltaTime);
    }
}