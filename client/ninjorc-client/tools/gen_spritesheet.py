import argparse
import json
import os.path

parser = argparse.ArgumentParser(
    prog='genSpritesheet'
)

parser.add_argument('-n', '--name', type=str)
parser.add_argument('-xs', '--x_start', default=0, type=int)
parser.add_argument('-ys', '--y_start', default=0, type=int)
parser.add_argument('-xi', '--x_interval', default=0, type=int)
parser.add_argument('-yi', '--y_interval', default=0, type=int)
parser.add_argument('-f', '--frame_count', type=int)
parser.add_argument('-p', '--path', type=str)
parser.add_argument('-ip', '--image_path', type=str)
parser.add_argument('-w', '--image_width', type=int)
parser.add_argument('-ht', '--image_height', type=int)

sprite_size = 32
args = parser.parse_args()
# image_name = os.path.basename(args.image_path)

spritesheet_obj = {
    "meta": {
        "image": args.image_path,
        "size": {"w": args.image_width, "h": args.image_height},
        "scale": "1"
    }
}

frames_obj = {}
animations_arr = []

for i in range(args.frame_count):
    frames_obj.update({f"{args.name}_{i}.png": {
        "frame": {
            "x": args.x_start+(i*args.x_interval),
            "y": args.y_start+(i*args.y_interval),
            "w": sprite_size,
            "h": sprite_size
        },
        "rotated": False,
        "trimmed": False,
        "spriteSourceSize": {"x":0,"y":0,"w":sprite_size,"h":sprite_size},
        "sourceSize": {"w":sprite_size,"h":sprite_size}
    }})
    animations_arr.append(f"{args.name}_{i}.png")

spritesheet_obj.update({"frames": frames_obj})
spritesheet_obj.update({"animations": {
    args.name: animations_arr
}})

dest_path = os.path.join(args.path, f"{args.name}.json")
with open(dest_path, "w") as f:
    json.dump(spritesheet_obj, f)