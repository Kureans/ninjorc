import { Direction } from "./Math/Direction";
import { Point } from "./Math/Point";
import { Action } from "./Orc";
import { None, Some, Option } from "./types/Option";

export class Connection {
    id!: number;
    socket: WebSocket;
    isConnectionOpen = false;
    gameInitCtxQueue: GameInitContext[] = [];
    deltaQueue: DeltaState[] = [];

    constructor() {
        this.socket = new WebSocket("ws://localhost:8080");
        this.socket.onmessage = (event) => {
            console.log(event.data);
            const packet: Packet = JSON.parse(event.data);
            switch (packet.Type) {
            case 'A':
                console.log(packet.Data); // Game Init Context
                const ctx: GameInitContext = mapJsonToGameInitContext(packet.Data[0]);
                console.log(ctx.clientId)
                console.log(ctx.orcCount)
                console.log(ctx.idToOrcLocations)
                this.id = ctx.clientId;
                this.gameInitCtxQueue.push(ctx);
                break;
            case 'L':
                console.log("Notification that game is starting");
                break;
            case 'G':
                const delta: DeltaState = mapJsonToDeltaState(packet.Data[0]);
                console.log("Delta: ", delta);
                this.deltaQueue.push(delta);
                break;
            default:
                console.log("Something else");
            }
        }
    }

    getNextGameInitContext(): Option<GameInitContext> {
        if (this.gameInitCtxQueue.length == 0) {
            return new None();
        }
        return new Some(this.gameInitCtxQueue.shift());
    }

    checkNextGameInitContext(): boolean {
        return this.gameInitCtxQueue.length == 0;
    }

    getNextDeltaState(): Option<DeltaState> {
        if (this.deltaQueue.length == 0) {
            return new None();
        }
        console.log("getting next delta");
        console.log(this.deltaQueue[0]);
        return new Some(this.deltaQueue.shift());      
    }

    sendPacket(packet: Packet) {
        if (!this.isConnectionOpen) {
            this.socket.addEventListener("open", () => {
                this.socket.send(JSON.stringify(packet));
                this.isConnectionOpen = true;
            }, {once: true});
        } else {
            this.socket.send(JSON.stringify(packet));
        }
    }
}

export type Packet = {
    Id: number;
    Type: string;
    Size: number;
    Data: any[];
};

export type DeltaState = {
    orcDeltas: OrcDelta[]; 
};

//TODO: Separate Orc Sprite and Orc GameInfo
export type OrcDelta = {
  Id: number;
  Health: number;
  Point: Point;
  Direction: Direction;
  Action: Action;
  IsSwinging: boolean;
}

export type GameInitContext = {
    orcCount: number;
    clientId: number;
    idToOrcLocations: IdPointPair[];
};

// might turn this into a generic next time
// but also feel like the benefit of json is faster iteration and dynamic typing so idk if im throwing that away if im doing all this mapping on my own
function mapJsonToGameInitContext(jsonObj: any): GameInitContext {
    const ctx: GameInitContext = {
        orcCount: jsonObj.OrcCount,
        clientId: jsonObj.ClientId,
        idToOrcLocations: []
    };
    for (const [k, v] of Object.entries(jsonObj.IdToOrcLocations)) {
        ctx.idToOrcLocations.push(new IdPointPair(parseInt(k), new Point(v.X, v.Y)));
        console.log(k, v);
    }
    return ctx;
}

function mapJsonToDeltaState(jsonObj: any): DeltaState {
    const deltaState: DeltaState = {
        orcDeltas: []
    };
    const orcArr = jsonObj["Orcs"];
    orcArr.forEach((orc: any) => {
        const delta: OrcDelta = {
            Id: orc.Id,
            Health: orc.Health,
            Point: new Point(orc.Point.X, orc.Point.Y),
            Direction: orc.Direction,
            Action: orc.Action,
            IsSwinging: orc.IsSwinging
        };
        deltaState.orcDeltas.push(delta);
    });
    console.log(deltaState);
    return deltaState;
}  

// {"Id":0,"Type":"G","Size":1,"Data":[{"Orcs":[{"Id":0,"Health":100,"Point":{"X":170,"Y":210},"Direction":0,"Action":0,"IsSwinging":true}]}]}

class IdPointPair {
    id: number;
    point: Point;

    constructor(id: number, point: Point) {
        this.id = id;
        this.point = point;
    }
}