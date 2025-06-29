import { Application } from 'pixi.js';
import { MAP_HEIGHT, MAP_WIDTH } from './Orc';
import { Game } from './Game';
import { Connection, Packet } from './Connection';
import { AssetManager } from './AssetManager';

const uiRoot = document.getElementById("lobby")!;
const startButton = document.getElementById("start-button");
startButton?.addEventListener("click", onStart);
await AssetManager.preloadTextures(); // loads textures into memory/cache

export function onStart() {
    const socket = new WebSocket("ws://localhost:8080");
    const readyMsg = {
        "Id": 1,
        "Type": "L",
        "Size": 1,
        "Data": [{"Lobby":{"IsReady": true}}]
    }
    
    socket.onopen = (event) => {
        socket.send(JSON.stringify(readyMsg));
    }
    socket.onmessage = (event) => {
        console.log(event.data);
        const packet: Packet = JSON.parse(event.data);
        switch (packet.Type) {
        case 'L':
            console.log("Notification that game is starting");
            break;
        default:
            console.log("Something else");
        }
    }

    document.body.removeChild(uiRoot);
    startGame();
}

async function startGame() {
    const app = new Application();
    await app.init({
        resolution: window.devicePixelRatio || 1,
        autoDensity: true,
        backgroundColor: 0x6495ed,
        width: MAP_WIDTH,
        height: MAP_HEIGHT
    });

    document.body.appendChild(app.canvas);
    const conn = new Connection();
    const game = new Game(app.ticker, conn, app.stage);
    console.log(game.gamestate);
    game.run();
};