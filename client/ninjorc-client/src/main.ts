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
    const textures = await manager.loadTextures(assetPathsAttack);
    const orc = new Orc(1, new Point(100, 100));
    manager.loadSpritesFromTextures(textures, orc.attackSprites);

    // const attackSprites = await manager.populateSprites(assetPathsAttack);
    // const idleSprites = await manager.populateSprites(assetPathsIdle);
    // orc.addAttackSprites(attackSprites);
    // orc.addIdleSprites(await manager.populateSprites(assetPathsIdle));
    // orc.idleSprites[0].scale = 2;
    // orc.setupSprite();
    // app.stage.addChild(orc);
    // const controller = new PlayerController(1, orc);
    // controller.setupInputHandler();
    // const conn = new Connection();
    // const game = new Game(app.ticker, conn, controller);
    // game.run();
})();