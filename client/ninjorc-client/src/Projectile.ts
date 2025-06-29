import { AnimatedSprite, Container, ContainerChild } from "pixi.js";
import { AssetManager } from "./AssetManager";

const PROJECTILE_SPEED = 0.1;
export const MAP_WIDTH = 960;
export const MAP_HEIGHT = 640;

export class Projectile extends Container {
  Id: number;
  Direction: Direction = Direction.NONE;
  Speed: number;
  currentSprite!: AnimatedSprite;

  constructor(id: number, point: Point, dir: Direction, stage: Container<ContainerChild>) {
    super();
    stage.addChild(this);
    this.Id = id;
    this.Speed = PROJECTILE_SPEED;
    this.position.x = point.x;
    this.position.y = point.y;
    this.Direction = dir;
    const assetPath = "sprites/fireball_right.json";
    this.currentSprite = AssetManager.populateSprite(assetPath);
    console.log(this.currentSprite);
    this.addChild(this.currentSprite);
    this.currentSprite.anchor.set(0.5);
    this.currentSprite.scale.set(2);
    this.currentSprite.animationSpeed = 0.2;
    this.currentSprite.loop = true;
    this.currentSprite.play();
  }

  update(deltaTime: number) {
    switch (this.Direction) {
      case Direction.UP:
        this.position.y = Math.max(0, this.position.y - PROJECTILE_SPEED * deltaTime);
        break;
      case Direction.DOWN:
        this.position.y = Math.min(
          MAP_HEIGHT,
          this.position.y + PROJECTILE_SPEED * deltaTime
        );
        break;
      case Direction.LEFT:
        this.position.x = Math.max(0, this.position.x - PROJECTILE_SPEED * deltaTime);
        break;
      case Direction.RIGHT:
        this.position.x = Math.min(
          MAP_WIDTH,
          this.position.x + PROJECTILE_SPEED * deltaTime
        );
        break;
    }
    console.log(this.position);
  }

//   setDefaultSprite() {
//     this.currentSprite = this.idleSprites[Direction.UP];
//     this.currentSprite.play();
//     this.addChild(this.currentSprite);
//   }

//   updateSprite() {
//     this.removeChild(this.currentSprite);
//     const currentDirection =
//       Direction.NONE == this.Direction ? this.lastFacing : this.Direction;
//     if (this.Action == Action.NONE) {
//       this.currentSprite =
//         Direction.NONE == this.Direction
//           ? this.idleSprites[currentDirection]
//           : this.runSprites[currentDirection];
//     } else if (this.Action == Action.MELEE) {
//       this.currentSprite = this.attackSprites[currentDirection];
//     } else if (this.Action == Action.PROJECTILE) {
//       this.currentSprite = this.fireballSprites[currentDirection];
//     }

//     this.currentSprite.play();
//     this.addChild(this.currentSprite);
//   }
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
    case Action.PROJECTILE:
      return "fireball";
    case Action.BLINK:
      return "blink";
  }
}

export enum Action {
  NONE,
  MELEE,
  PROJECTILE,
  BLINK,
}
