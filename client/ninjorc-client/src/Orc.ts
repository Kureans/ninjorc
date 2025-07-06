import { AnimatedSprite, Container } from "pixi.js";
import { GameState } from "./Game";
import { Projectile } from "./Projectile";
import { AssetManager } from "./AssetManager";
import { Direction } from "./Math/Direction";
import { Point } from "./Math/Point";
import { DeltaState, OrcDelta } from "./Connection";

const ORC_SPEED = 0.5;
export const MAP_WIDTH = 960;
export const MAP_HEIGHT = 640;

export class Orc extends Container {
  Id: number;
  Health: number;
  lastFacing: Direction = Direction.NONE;
  lastAction: Action = Action.NONE;
  Direction: Direction = Direction.NONE; // Direction & Action properties directly set by PlayerController locally
  Action: Action = Action.NONE;
  IsSwinging: boolean;
  idleSprites: AnimatedSprite[] = [];
  attackSprites: AnimatedSprite[] = [];
  runSprites: AnimatedSprite[] = [];
  fireballSprites: AnimatedSprite[] = [];
  currentSprite!: AnimatedSprite;
  gamestateRef!: GameState;

  constructor(id: number, point: Point) {
    super();
    this.Id = id;
    this.Health = 100;
    this.IsSwinging = false;
    this.position.x = point.x;
    this.position.y = point.y;
  }

  populateSprites() {
    const spriteKeys = ["attack", "idle", "run"];
    for (const key of spriteKeys) {
      console.log(key)
      const spriteArray = []
      for (const path of AssetManager.assetPaths[key]) {
        spriteArray.push(AssetManager.populateSprite(path));
      }
      this.addSprites(key, spriteArray);
    }
  }

  addSprites(key: string, sprites: AnimatedSprite[]) {
    switch (key) {
      case "attack":
        this.addAttackSprites(sprites);
        break;
      case "idle":
        this.addIdleSprites(sprites);
        break;
      case "run":
        this.addRunSprites(sprites);
        break;
      default:
        console.error("Could not add sprites with invalid key ", key);
    }
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

  addRunSprites(sprites: AnimatedSprite[]) {
    this.runSprites = sprites;
    for (let i = 0; i < this.idleSprites.length; i++) {
      this.runSprites[i].anchor.set(0.5);
      this.runSprites[i].scale.set(2);
      this.runSprites[i].animationSpeed = 0.2;
      this.runSprites[i].loop = true;
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
      };
    }
  }

  updateFromServer(delta: OrcDelta) {
    this.Direction = delta.Direction;
    this.Action = delta.Action;
    if (this.lastFacing != this.Direction || this.lastAction != this.Action) {
      // Only update animations when required
      this.updateSprite();
    }
    if (this.Direction != Direction.NONE) {
      this.lastFacing = this.Direction; // Direction.NONE is not a valid direction to face
    }
    this.lastAction = this.Action;
    this.Health = delta.Health;
    this.position.x = delta.Point.x;
    this.position.y = delta.Point.y;
    this.IsSwinging = delta.IsSwinging;
    console.log(`Delta dir: ${delta.Direction}, action: ${delta.Action}`);
    console.log(`Delta values: ${delta.Point.x}, ${delta.Point.y}`)
    console.log(`New position: ${this.position}`);
  }

  update(deltaTime: number) {
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
        this.position.y = Math.max(0, this.position.y - ORC_SPEED * deltaTime);
        break;
      case Direction.DOWN:
        this.position.y = Math.min(
          MAP_HEIGHT,
          this.position.y + ORC_SPEED * deltaTime
        );
        break;
      case Direction.LEFT:
        this.position.x = Math.max(0, this.position.x - ORC_SPEED * deltaTime);
        break;
      case Direction.RIGHT:
        this.position.x = Math.min(
          MAP_WIDTH,
          this.position.x + ORC_SPEED * deltaTime
        );
        break;
    }
  }

  setDefaultSprite() {
    this.currentSprite = this.idleSprites[Direction.UP-1];
    this.currentSprite.play();
    this.addChild(this.currentSprite);
  }

  updateSprite() {
    this.removeChild(this.currentSprite);
    const currentDirection =
      Direction.NONE == this.Direction ? this.lastFacing : this.Direction;
    if (this.Action == Action.NONE) {
      this.currentSprite =
        Direction.NONE == this.Direction
          ? this.idleSprites[currentDirection-1]
          : this.runSprites[currentDirection-1];
    } else if (this.Action == Action.MELEE) {
      this.currentSprite = this.attackSprites[currentDirection-1];
    } else if (this.Action == Action.PROJECTILE) {
      console.log(this.gamestateRef);
      this.gamestateRef.projectiles.push(new Projectile(1, new Point(this.position.x, this.position.y), currentDirection, this.gamestateRef.stageRef));
      this.Action = Action.NONE;
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
