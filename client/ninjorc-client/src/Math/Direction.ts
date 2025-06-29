export enum Direction {
  UP,
  DOWN,
  LEFT,
  RIGHT,
  NONE,
}

export function getDirectionString(dir: Direction): string {
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