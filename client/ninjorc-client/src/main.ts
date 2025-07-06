import { Application } from 'pixi.js';
import { MAP_HEIGHT, MAP_WIDTH } from './Orc';
import { Game } from './Game';
import { Connection, GameInitContext, Packet } from './Connection';
import { AssetManager } from './AssetManager';
import { Option, None } from './types/Option';

const uiRoot = document.getElementById("lobby")!;
const startButton = document.getElementById("start-button");
startButton?.addEventListener("click", onStart);
await AssetManager.preloadTextures(); // loads textures into memory/cache

export function onStart() {
    const conn = new Connection();
    const readyMsg: Packet = {
        "Id": 1,
        "Type": "L",
        "Size": 1,
        "Data": [{"IsReady": true}]
    }
    conn.sendPacket(readyMsg);
    document.body.removeChild(uiRoot);
    startGame(conn);
}

async function startGame(conn: Connection) {
    const app = new Application();
    await app.init({
        resolution: window.devicePixelRatio || 1,
        autoDensity: true,
        backgroundColor: 0x6495ed,
        width: MAP_WIDTH,
        height: MAP_HEIGHT
    });

    document.body.appendChild(app.canvas);
    await wait(conn);
    const gameInitOption = conn.getNextGameInitContext();

    if (gameInitOption instanceof None) {
        console.log("test");
        console.error("Cannot initialise Game!");
        return;
    }
    const gameInitCtx = gameInitOption.value;
    console.log("Context Obj: ", gameInitCtx);
    const game = new Game(app.ticker, conn, app.stage, gameInitCtx);
    game.run();
};

async function wait(conn: Connection) {
    let flag = false;
    while (!flag) {
        if (conn.checkNextGameInitContext()) {
            await sleep(100);
        } else {
            flag = true;
        }
    }
}

function sleep(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms));
}