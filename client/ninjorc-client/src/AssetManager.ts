import { AnimatedSprite, Assets } from "pixi.js";
import assetPathsJson from '../public/asset_paths.json' with { type: 'json'};

type AssetPaths = {
  [key: string]: string[];
};

export class AssetManager {
    areSpritesReady: boolean = false;
    static assetPaths: AssetPaths = assetPathsJson;

    static populateSprite(jsonAssetPath: string): AnimatedSprite {
        console.log(jsonAssetPath);
        const sheet = Assets.get(jsonAssetPath);
        const jsonAssetName = this.getNameFromPath(jsonAssetPath);
        // TO-DO: Add error handling
        let s: AnimatedSprite;
        try {
            s = new AnimatedSprite(sheet.animations[jsonAssetName]); 
        } catch (error) {
            console.log(error);
        }
 
        return s;
    }

    static async preloadTextures() {
        for (const [key, spritePaths] of Object.entries(AssetManager.assetPaths)) {
            for (const path of spritePaths) {
                await Assets.load(path);
            }
        }
    }

    private static getNameFromPath(path: string): string {
        return path.split('/')
                .slice(-1)[0]
                .split('.')[0];
    }
}