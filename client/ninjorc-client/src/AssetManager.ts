import { AnimatedSprite, Assets, Sprite, Texture } from "pixi.js";

export class AssetManager {
    areSpritesReady: boolean = false;
    static spriteMap: Record<string, Texture> = {};

    async populateSprite(jsonAssetPath: string): Promise<AnimatedSprite> {
        const sheet = await Assets.load(jsonAssetPath);
        const jsonAssetName = this.getNameFromPath(jsonAssetPath);
        // TO-DO: Add error handling
        let s: AnimatedSprite;
        try {
            s = new AnimatedSprite(sheet.animations[jsonAssetName]); 
        } catch (error) {
            console.log(error);
        }

        AssetManager.spriteMap[jsonAssetPath] = sheet.animations[jsonAssetName];
 
        return s;
    }

    loadTextures(jsonAssetPaths: string[]): Promise<Record<string, Texture>> {
        let spriteNames: string[] = [];
        jsonAssetPaths.forEach((path) => {
            const name = this.getNameFromPath(path);
            Assets.add({alias: name, src: path});
            spriteNames.push(name);
        });
        return Assets.load(spriteNames);
    }

    loadSpritesFromTextures(textures: Record<string, Texture>, sprites: AnimatedSprite[]) {
        console.log(textures);
        for (const key in textures) {
            Sprite.from(textures[key]);
        }
    }

    private getNameFromPath(path: string): string {
        return path.split('/')
                .slice(-1)[0]
                .split('.')[0];
    }
}