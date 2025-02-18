import { Application, Assets, AnimatedSprite} from 'pixi.js';
import { PlayerController } from './PlayerController';
import { MAP_HEIGHT, MAP_WIDTH, Orc, Point } from './Orc';
import { Game } from './Game';
import { Connection } from './Connection';
import { AssetManager } from './AssetManager';

const assetPathsAttack = [
    'sprites/orc_attack_up.json',
    'sprites/orc_attack_down.json',
    'sprites/orc_attack_left.json',
    'sprites/orc_attack_right.json',
];

const assetPathsIdle = [
    'sprites/orc_idle_up.json',
    'sprites/orc_idle_down.json',
    'sprites/orc_idle_left.json',
    'sprites/orc_idle_right.json',
];

const assetPathsRun = [
    'sprites/orc_run_up.json',
    'sprites/orc_run_down.json',
    'sprites/orc_run_left.json',
    'sprites/orc_run_right.json',
];

(async () => {
    const app = new Application();
    await app.init({
        resolution: window.devicePixelRatio || 1,
        autoDensity: true,
        backgroundColor: 0x6495ed,
        width: MAP_WIDTH,
        height: MAP_HEIGHT
    });

    document.body.appendChild(app.canvas);
    const manager = new AssetManager();
    const orc = new Orc(1, new Point(100, 100));

    const attackSprite1 = await manager.populateSprite(assetPathsAttack[0]);
    const attackSprite2 = await manager.populateSprite(assetPathsAttack[1]);
    const attackSprite3 = await manager.populateSprite(assetPathsAttack[2]);
    const attackSprite4 = await manager.populateSprite(assetPathsAttack[3]);

    const idleSprite1 = await manager.populateSprite(assetPathsIdle[0]);
    const idleSprite2 = await manager.populateSprite(assetPathsIdle[1]);
    const idleSprite3 = await manager.populateSprite(assetPathsIdle[2]);
    const idleSprite4 = await manager.populateSprite(assetPathsIdle[3]);

    const runSprite1 = await manager.populateSprite(assetPathsRun[0]);
    const runSprite2 = await manager.populateSprite(assetPathsRun[1]);
    const runSprite3 = await manager.populateSprite(assetPathsRun[2]);
    const runSprite4 = await manager.populateSprite(assetPathsRun[3]);
    
    orc.addAttackSprites([attackSprite1, attackSprite2, attackSprite3, attackSprite4]);
    orc.addIdleSprites([idleSprite1, idleSprite2, idleSprite3, idleSprite4]);
    orc.addRunSprites([runSprite1, runSprite2, runSprite3, runSprite4]);
    app.stage.addChild(orc);
    orc.setDefaultSprite();
    const controller = new PlayerController(1, orc);
    controller.setupInputHandler();
    const conn = new Connection();
    const game = new Game(app.ticker, conn, controller);
    game.run();
})();