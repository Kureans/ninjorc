import { AnimatedSprite, Container, Sprite } from "pixi.js";
import { AssetManager } from "./AssetManager";

const ORC_SPEED = 0.5;
export const MAP_WIDTH        = 960;
export const MAP_HEIGHT       = 640;

export class Orc extends Container {
    Id: number;
    Health: number;
    lastFacing: Direction = Direction.NONE;
    lastAction: Action = Action.NONE;
    Direction: Direction = Direction.NONE;
    Action: Action = Action.NONE;
    IsSwinging: boolean;
    idleSprites: AnimatedSprite[] = [];
    attackSprites: AnimatedSprite[] = [];
    currentSprite!: AnimatedSprite;

    constructor(id: number, point: Point) {
        super();
        this.Id = id;
        this.Health = 100;
        this.IsSwinging = false;
        this.position.x = point.x;
        this.position.y = point.y;
    }

    addIdleSprites(sprites: AnimatedSprite[]) {
        this.idleSprites = sprites;
        for (let i = 0; i < this.idleSprites.length; i++) {
            this.idleSprites[i].anchor.set(0.5);
            this.idleSprites[i].scale.set(2);
            this.idleSprites[i].animationSpeed = 0.2;
            this.idleSprites[i].loop = true;
        }
    }

    addAttackSprites(sprites: AnimatedSprite[]) {
        this.attackSprites = sprites;
        for (let i = 0; i < this.attackSprites.length; i++) {
            this.attackSprites[i].anchor.set(0.5);
            this.attackSprites[i].scale.set(2);
            this.attackSprites[i].animationSpeed = 0.2;
            this.attackSprites[i].loop = false;
            this.attackSprites[i].onComplete = () => {
                this.currentSprite.currentFrame = 0;
                this.IsSwinging = false; 
                this.Action = Action.NONE;
            }
        }
    }

    update(deltaTime: number) {
        // console.log("Orc Direction: ", getDirectionString(this.Direction));
        console.log("Orc Action: ", getActionString(this.Action));
        if (this.lastFacing != this.Direction || this.lastAction != this.Action) {
            // Only update animations when required
            this.updateSprite();
        }
        if (this.Direction != Direction.NONE) {
            this.lastFacing = this.Direction; // Direction.NONE is not a valid direction to face
        }
        this.lastAction = this.Action;
        switch (this.Direction) {
            case Direction.UP:
                this.position.y = Math.max(0, this.position.y-ORC_SPEED*deltaTime);
                break;
            case Direction.DOWN:
                this.position.y = Math.min(MAP_HEIGHT, this.position.y+ORC_SPEED*deltaTime);
                break;
            case Direction.LEFT:
                this.position.x = Math.max(0, this.position.x-ORC_SPEED*deltaTime);
                break;
            case Direction.RIGHT:
                this.position.x = Math.min(MAP_WIDTH, this.position.x+ORC_SPEED*deltaTime);
                break;
        }

    }

    setDefaultSprite() {
        this.currentSprite = this.idleSprites[Direction.UP];
        this.currentSprite.play();
        this.addChild(this.currentSprite);
    }

    updateSprite() {
        this.removeChild(this.currentSprite);
        const currentDirection = (Direction.NONE == this.Direction) ? this.lastFacing : this.Direction;
        if (this.Action == Action.NONE) {
            this.currentSprite = this.idleSprites[currentDirection];
        } else if (this.Action == Action.MELEE) {
            this.currentSprite = this.attackSprites[currentDirection];
        }

        this.currentSprite.play();
        this.addChild(this.currentSprite);
    }

    resetState() {
        this.Direction = Direction.NONE;
        this.Action = Action.NONE;
    }

    resetDirection() {
        this.Direction = Direction.NONE;
    }
}

export class Point {
    x: number;
    y: number;

    constructor(x: number, y: number) {
        this.x = x;
        this.y = y;
    }
}

export enum Direction {
    UP,
    DOWN,
    LEFT,
    RIGHT,
    NONE,
}

function getDirectionString(dir: Direction): string {
    switch (dir) {
        case Direction.UP:
            return "up";
        case Direction.DOWN:
            return "down";
        case Direction.LEFT:
            return "left";
        case Direction.RIGHT:
            return "right";
        case Direction.NONE:
            return "none";
    }
}

function getActionString(action: Action): string {
    switch (action) {
        case Action.NONE:
            return "none";
        case Action.MELEE:
            return "melee";
    }
}

export enum Action {
    NONE,
    MELEE,
    PROJECTILE,
    BLINK
}