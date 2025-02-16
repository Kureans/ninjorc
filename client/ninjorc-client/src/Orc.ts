import { AnimatedSprite, Container, Sprite } from "pixi.js";
import { AssetManager } from "./AssetManager";

const ORC_SPEED = 1;
export const MAP_WIDTH        = 960;
export const MAP_HEIGHT       = 640;

export class Orc extends Container {
    Id: number;
    Health: number;
    lastFacing: Direction = 1;
    Direction: Direction = 0;
    Action: Action = 0;
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

    setupSprite() {
        this.removeChild(this.currentSprite);
        if (this.Action == Action.NONE) {
            this.currentSprite = this.idleSprites[this.lastFacing];
        } else if (this.Action == Action.MELEE) {
            this.currentSprite = this.attackSprites[this.lastFacing];
        }
        this.currentSprite.scale.set(2);
        this.currentSprite.anchor.set(0.5);
        // this.currentSprite.x = this.position.x;
        // this.currentSprite.y = this.position.y;
        this.currentSprite.animationSpeed = 0.2;
        this.currentSprite.loop = true;
        this.currentSprite.onComplete = () => {
            this.currentSprite.currentFrame = 0;
        }
        this.currentSprite.play();
        this.addChild(this.currentSprite);
    }

    addIdleSprites(sprites: AnimatedSprite[]) {
        this.idleSprites = sprites;
    }

    addAttackSprites(sprites: AnimatedSprite[]) {
        this.attackSprites = sprites;
    }

    update(deltaTime: number) {
        this.lastFacing = this.Direction;
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

    resetState() {
        this.Direction = Direction.NONE;
        this.Action = Action.NONE;
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
    NONE,
    UP,
    DOWN,
    LEFT,
    RIGHT,
}

export enum Action {
    NONE,
    MELEE,
    PROJECTILE,
    BLINK
}