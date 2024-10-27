import { AnimatedSprite } from "pixi.js";

const ORC_SPEED = 1;
const MAP_WIDTH        = 1500
const MAP_HEIGHT       = 700

function max(a: number, b: number): number {
	if (a > b) {
		return a
	}
	return b
}

function min(a: number, b: number): number {
	if (a < b) {
		return a
	}
	return b
}

export class Orc {
    Id: number;
    Health: number;
    Direction: Direction = 0;
    Action: Action = 0;
    IsSwinging: boolean;
    sprite: AnimatedSprite;

    constructor(id: number, x: number, y: number, sprite: AnimatedSprite) {
        this.Id = id;
        this.Health = 100;
        this.IsSwinging = false;
        this.sprite = sprite;

        this.sprite.scale.set(2);
        this.sprite.anchor.set(0.5);
        this.sprite.x = x;
        this.sprite.y = y;
        this.sprite.animationSpeed = 0.2;
    }

    update(deltaTime: number) {
        switch (this.Direction) {
            case Direction.UP:
                this.sprite.y = max(0, this.sprite.y-ORC_SPEED*deltaTime);
                break;
            case Direction.DOWN:
                this.sprite.y = min(MAP_HEIGHT, this.sprite.y+ORC_SPEED*deltaTime);
                break;
            case Direction.LEFT:
                this.sprite.x = max(0, this.sprite.x-ORC_SPEED*deltaTime);
                break;
            case Direction.RIGHT:
                this.sprite.x = min(MAP_WIDTH, this.sprite.x+ORC_SPEED*deltaTime);
                break;
        }

        switch (this.Action) {
            case Action.MELEE:
                this.sprite.play();
                break;
        }
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
    RIGHT
}

export enum Action {
    NONE,
    MELEE,
    PROJECTILE,
    BLINK
}